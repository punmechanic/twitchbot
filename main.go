package main

import (
	"context"
	"log"
	"net/http"

	"example.com/twitchbot/pkg/twitch"
	"example.com/twitchbot/pkg/twitch/eventsub"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()

	// TODO: Try to add Public flow using localhost
	// Twitch seems to let us do localhost during test but I don't know if they would allow it in production...
	//
	// If not, that means spinning up complex infrastructure where we have an Oauth2 flow that the twitch bot (conduit, I guess)
	// uses to interact w/ twitch and we have to develop a client the end-user can use to interact with ours, and that client would
	// use the public flow
	//
	// tbh the latter is likely more approachable for most twitch users.
	cfg := oauth2.Config{
		ClientID: "",
		Scopes:   []string{},
	}

	client := twitch.New(&oauth2Authorization{
		ClientID:    cfg.ClientID,
		TokenSource: cfg.TokenSource(context.Background(), nil),
	})

	conn, err := eventsub.Dial(ctx)
	if err != nil {
		log.Fatalf("init websocket: %s", err)
	}

	go func() {
		err = conn.Listen()
		if err != nil {
			log.Fatalf("closed with error: %s", err)
		}
	}()

	err = client.SubscribeEvents(ctx, <-conn.SessionID, []string{"channel_follow"})
	if err != nil {
		log.Fatalf("setup events: %s", err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case ev := <-conn.ChannelFollowed:
			log.Printf("follow: %#v", ev)
		}
	}
}

type oauth2Authorization struct {
	ClientID    string
	TokenSource oauth2.TokenSource
}

func (a *oauth2Authorization) Apply(r *http.Request) error {
	return nil
}
