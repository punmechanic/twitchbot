package twitch

import (
	"context"
	"log"

	"example.com/twitchbot/pkg/twitch/eventsub"
)

func (c *Client) SubscribeEvents(ctx context.Context, sessionID string, events []string) error {
	var req eventsub.CreateEventSubSubscriptionRequest
	for _, event := range events {
		req.Subscriptions = append(req.Subscriptions, &eventsub.SubscriptionRequest{
			Type:    event,
			Version: "1",
			Transport: eventsub.SubscriptionTransport{
				Method:    eventsub.SubscriptionMethodWebhook,
				SessionID: sessionID,
			},
		})
	}

	resp, err := eventsub.Subscribe(ctx, c, &req)
	if err != nil {
		return err
	}

	log.Printf("%#v", resp)
	// TODO: impl
	return nil
}
