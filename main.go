package main

import (
	"context"
	"encoding/json"
	"log"

	"example.com/twitchbot/pkg/twitch"
	"example.com/twitchbot/pkg/twitch/events"
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

	go loop(conn.Events)

	err = conn.Listen()
	if err != nil {
		log.Fatalf("closed with error: %s", err)
	}
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
