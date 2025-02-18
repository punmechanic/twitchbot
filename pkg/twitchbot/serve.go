package twitchbot

import (
	"context"
	"errors"
	"fmt"
	"log"

	"example.com/twitchbot/pkg/twitch"
	"example.com/twitchbot/pkg/twitch/eventsub"
	"example.com/twitchbot/pkg/twitch/subscriptions"
	"github.com/urfave/cli/v3"
	"golang.org/x/oauth2"
)

// serve runs the Twitch bot.
func serve(ctx context.Context, c *cli.Command) error {
	var (
		cfg                     = initTwitchConfig([]string{"user:read:chat"})
		broadcasterIDs []string = c.StringSlice("broadcaster-ids")
	)

	token, usedKeyringToken, err := fetchTokenWithFallback(ctx, cfg)
	defer saveTokenInKeyring(token)

	userInfo, err := fetchUserInfo(ctx, cfg, token)
	if err != nil {
		return fmt.Errorf("fetch user info: %w", err)
	}

	// Our websocket is useless without having a valid token for the Twitch API, so wait to have one before we continue.
	client := twitch.New(cfg, token)
	conn, err := eventsub.Dial(ctx)
	if err != nil {
		return fmt.Errorf("init websocket: %w", err)
	}

	listenErrCh := make(chan error, 1)
	go func() {
		err = conn.Listen()
		if err != nil {
			listenErrCh <- err
		}
		close(listenErrCh)
	}()

	id := <-conn.SessionID

	var reqs []*twitch.SubscribeRequest
	for _, broadcasterID := range broadcasterIDs {
		reqs = append(reqs, &twitch.SubscribeRequest{
			Type: subscriptions.ChannelChatMessage,
			Condition: eventsub.Condition{
				UserID:            userInfo.Sub,
				BroadcasterUserID: broadcasterID,
			},
			Transport: eventsub.Transport{
				Method:    eventsub.MethodWebsocket,
				SessionID: id,
			},
		})
	}
	err = client.SubscribeEvents(ctx, reqs)

	if err != nil {
		var retrieveErr *oauth2.RetrieveError
		if errors.As(err, &retrieveErr) && usedKeyringToken {
			// If we are here, it means that the token in the keyring has expired.  We will need to re-subscribe.
			// And, since fetching a new token might take longer than the 10 seconds twitch gives us before killing
			// our websocket, we will also need to re-Dial.
			// TODO: implement this
			panic("not yet implemented")
		}

		return fmt.Errorf("setup events: %s", err)
	}

	log.Println("listening for events")
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-listenErrCh:
			return fmt.Errorf("listener error: %w", err)
		case ev := <-conn.ChannelFollowed:
			log.Printf("follow: %#v", ev)
		case ev := <-conn.ChannelChatMessage:
			log.Printf("[%s] %s: %s", ev.BroadcasterUserName, ev.ChatterUserName, ev.Message.Text)
		}
	}
}
