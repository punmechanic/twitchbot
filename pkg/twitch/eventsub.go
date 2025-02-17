package twitch

import (
	"context"
	"fmt"

	"example.com/twitchbot/pkg/twitch/eventsub"
	"golang.org/x/sync/errgroup"
)

func (c *Client) SubscribeEvents(ctx context.Context, sessionID string, events []string) error {
	var (
		maxConcurrency = 10
		grp, reqCtx    = errgroup.WithContext(ctx)
	)
	grp.SetLimit(maxConcurrency)

	for _, event := range events {
		grp.Go(func() error {
			req := eventsub.SubscribeRequest{
				Type:    event,
				Version: "1",
				Condition: eventsub.Condition{
					UserID: "meppermintpocha",
				},
				Transport: eventsub.Transport{
					Method:    eventsub.MethodWebsocket,
					SessionID: sessionID,
				},
			}
			_, err := eventsub.Subscribe(reqCtx, c, &req)
			if err != nil {
				return fmt.Errorf("subscribe: %w", err)
			}
			return nil
		})
	}

	return grp.Wait()
}
