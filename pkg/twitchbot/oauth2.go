package twitchbot

import (
	"context"
	"fmt"
	"io"
	"net/http"

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
