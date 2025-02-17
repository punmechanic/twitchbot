package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"example.com/twitchbot/pkg/twitch"
	"example.com/twitchbot/pkg/twitch/eventsub"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
)

func init() {
	browser.Stderr = io.Discard
}

type tokenFetchJob struct {
	C               chan *oauth2.Token
	Verifier, State string
}

// runServerForCallback runs a HTTP server listening for a valid OAuth2 callback.
func runServerForCallback(ctx context.Context, cfg *oauth2.Config, ch <-chan tokenFetchJob) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/oauth2/{provider}/callback", func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-ch:
			log.Println(r.PathValue("provider"))
		case <-r.Context().Done():
			// timeout
		}
	})
	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	go func() {
		<-ctx.Done()
		server.Shutdown(context.Background())
	}()

	return server.ListenAndServe()
}

func fetchInitialToken(ctx context.Context, cfg *oauth2.Config) (*oauth2.Token, error) {
	var (
		// jobQueue contains a queue of all jobs waiting for a token.
		//
		// The capacity must be set to 1 to prevent a deadlock.
		jobQueue = make(chan tokenFetchJob, 1)
		// tokCh is our channel which will receive a token.
		//
		// The capacity must be set to 1 to prevent a deadlock.
		tokCh     = make(chan *oauth2.Token, 1)
		httpErrCh = make(chan error)
	)

	go func() {
		err := runServerForCallback(ctx, cfg, jobQueue)
		if err != nil {
			httpErrCh <- err
		}
		close(httpErrCh)
	}()

	state := "opaque string goes here"
	verifier := oauth2.GenerateVerifier()
	err := browser.OpenURL(cfg.AuthCodeURL(state, oauth2.VerifierOption(verifier)))
	if err != nil {
		return nil, err
	}

	// Queue our job.
	jobQueue <- tokenFetchJob{
		Verifier: verifier,
		State:    state,
		C:        tokCh,
	}
	// Wait for our token.
	select {
	case tok := <-tokCh:
		return tok, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// runLocal attempts to run the twitch bot locally.
func runLocal(ctx context.Context) error {
	// Twitch seems to let us do localhost during test but I don't know if they would allow it in production...
	//
	// If not, that means spinning up complex infrastructure where we have an Oauth2 flow that the twitch bot (conduit,
	// I guess) uses to interact w/ twitch and we have to develop a client the end-user can use to interact with ours,
	// and that client would use the public flow
	//
	// tbh the latter is likely more approachable for most twitch users.
	cfg := oauth2.Config{
		ClientID: "??",
		// From https://id.twitch.tv/oauth2/.well-known/openid-configuration
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://id.twitch.tv/oauth2/authorize",
			TokenURL: "https://id.twitch.tv/oauth2/token",
		},
		RedirectURL: "http://localhost:8080/oauth2/twitch/callback",
		Scopes:      []string{"openid"},
	}

	tok, err := fetchInitialToken(ctx, &cfg)
	if err != nil {
		return fmt.Errorf("fetch initial token: %w", ctx.Err())
	}

	// Our websocket is useless without having a valid token for the Twitch API, so wait to have one before we continue.
	client := twitch.New(&cfg, tok)

	conn, err := eventsub.Dial(ctx)
	if err != nil {
		return fmt.Errorf("init websocket: %s", err)
	}

	listenErrCh := make(chan error, 1)
	go func() {
		err = conn.Listen()
		if err != nil {
			listenErrCh <- err
		}
		close(listenErrCh)
	}()

	err = client.SubscribeEvents(ctx, <-conn.SessionID, []string{"channel_follow"})
	if err != nil {
		return fmt.Errorf("setup events: %s", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-listenErrCh:
			return fmt.Errorf("listener error: %w", err)
		case ev := <-conn.ChannelFollowed:
			log.Printf("follow: %#v", ev)
		}
	}
}

func main() {
	ctx := context.Background()

	err := runLocal(ctx)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
