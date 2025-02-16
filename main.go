package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"example.com/twitchbot/pkg/twitch"
	"example.com/twitchbot/pkg/twitch/events"
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

	err = setupEvents(ctx, client, <-conn.SessionID)
	if err != nil {
		log.Fatalf("setup events: %s", err)
	}

	loop(conn.Notifications)
}

func setupEvents(ctx context.Context, client *twitch.Client, sessionID string) error {
	return client.SubscribeEvents(ctx, sessionID, []string{"channel_follow"})
}

func loop(events <-chan eventsub.Notification) error {
	for event := range events {
		err := doEvent(event)
		if err != nil {
			return err
		}
	}
	return nil
}

func doEvent(event eventsub.Notification) error {
	var followEvent *events.ChannelFollow

	switch event.Subscription.Type {
	case "channel_follow":
		err := json.Unmarshal(event.Event, &followEvent)
		if err != nil {
			return err
		}
		log.Printf("follow: %#v", followEvent)
	}
	return nil
}

type oauth2Authorization struct {
	ClientID    string
	TokenSource oauth2.TokenSource
}

func (a *oauth2Authorization) Apply(r *http.Request) error {
	return nil
}
