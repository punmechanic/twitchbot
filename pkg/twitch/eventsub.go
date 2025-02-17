package twitch

import (
	"context"
	"fmt"

	"example.com/twitchbot/pkg/twitch/eventsub"
	"example.com/twitchbot/pkg/twitch/eventsub/subscriptions"
	"golang.org/x/sync/errgroup"
)

type SubscribeRequest struct {
	Type      subscriptions.Type
	Condition eventsub.Condition
	Transport eventsub.Transport
}

func (c *Client) SubscribeEvents(ctx context.Context, events []*SubscribeRequest) error {
	var (
		maxConcurrency = 10
		grp, reqCtx    = errgroup.WithContext(ctx)
	)
	grp.SetLimit(maxConcurrency)

	for _, event := range events {
		grp.Go(func() error {
			_, err := eventsub.Subscribe(reqCtx, c, &eventsub.SubscribeRequest{
				Type:      event.Type.Name,
				Version:   event.Type.Version,
				Transport: event.Transport,
				Condition: event.Condition,
			})
			if err != nil {
				return fmt.Errorf("subscribe: %w", err)
			}
			return nil
		})
	}

	return grp.Wait()
}
