package twitchbot

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

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

func serveCallback(cfg *oauth2.Config, ch <-chan tokenFetchJob) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Remove panics
		var job tokenFetchJob
		select {
		case job = <-ch:
			break
		case <-r.Context().Done():
			// timeout
			panic("timeout")
		}

		// TODO: Verify state
		// Client side check so I don't care to do it right now
		token, err := cfg.Exchange(r.Context(), r.FormValue("code"), oauth2.S256ChallengeOption(job.Verifier))
		if err != nil {
			panic(fmt.Sprintf("token exchange: %s", err))
		}

		job.C <- token
		fmt.Fprintln(w, "You can close this window now.")
	})
}

// runServerForCallback runs a HTTP server listening for a valid OAuth2 callback.
func runServerForCallback(ctx context.Context, cfg *oauth2.Config, ch <-chan tokenFetchJob) error {
	mux := http.NewServeMux()
	mux.Handle("/oauth2/{provider}/callback", serveCallback(cfg, ch))
	server := http.Server{Handler: mux, Addr: ":8080"}

	go func() {
		<-ctx.Done()
		server.Shutdown(context.Background())
	}()

	return server.ListenAndServe()
}

func initTwitchConfig() *oauth2.Config {
	// Twitch seems to let us do localhost during test but I don't know if they would allow it in production...
	//
	// If not, that means spinning up complex infrastructure where we have an Oauth2 flow that the twitch bot (conduit,
	// I guess) uses to interact w/ twitch and we have to develop a client the end-user can use to interact with ours,
	// and that client would use the public flow
	//
	// tbh the latter is likely more approachable for most twitch users.
	return &oauth2.Config{
		ClientID:     os.Getenv("TWITCH_CLIENT_ID"),
		ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
		// From https://id.twitch.tv/oauth2/.well-known/openid-configuration
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://id.twitch.tv/oauth2/authorize",
			TokenURL: "https://id.twitch.tv/oauth2/token",
		},
		RedirectURL: "http://localhost:8080/oauth2/twitch/callback",
		Scopes:      []string{"openid"},
	}
}

func fetchTokenFromTwitch(ctx context.Context, cfg *oauth2.Config) (*oauth2.Token, error) {
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
