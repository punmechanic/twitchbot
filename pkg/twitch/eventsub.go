package twitch

import (
	"context"
	"log"

	"example.com/twitchbot/pkg/twitch/eventsub"
)

func (c *Client) SubscribeEvents(ctx context.Context, sessionID string, events []string) error {
	var req eventsub.SubscribeRequest
	for _, event := range events {
		req.Subscriptions = append(req.Subscriptions, &eventsub.SubscriptionDefinition{
			Type:    event,
			Version: "1",
			Transport: eventsub.Transport{
				Method:    eventsub.MethodWebhook,
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
