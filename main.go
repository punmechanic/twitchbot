package main

import (
	"context"
	"log"

	"example.com/twitchbot/pkg/twitch"
	"example.com/twitchbot/pkg/twitch/eventsub"
)

func main() {
	var client twitch.Client
	conn, err := eventsub.Dial(
		context.Background(),
		eventsub.WithSubscriptions([]string{"channel_follow"}, &client),
	)
	if err != nil {
		log.Fatalf("init websocket: %s", err)
	}

	go func() {
		for event := range conn.Events {
			log.Printf("%#v", event)
		}
	}()

	err = conn.Listen()
	if err != nil {
		log.Fatalf("closed with error: %s", err)
	}
}
