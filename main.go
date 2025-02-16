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
	client := twitch.New(&oauth2Authorization{
		ClientID:    "..",
		TokenSource: nil,
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
