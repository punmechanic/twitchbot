package twitchbot

import (
	"context"
	"io"
	"log"
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
