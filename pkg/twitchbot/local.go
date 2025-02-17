package twitchbot

import (
	"context"
	"fmt"
	"log"

	"example.com/twitchbot/pkg/twitch"
	"example.com/twitchbot/pkg/twitch/eventsub"
	"golang.org/x/oauth2"
)

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
