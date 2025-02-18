package twitchbot

import (
	"context"
	"errors"
	"fmt"
	"log"

	"example.com/twitchbot/pkg/twitch"
	"example.com/twitchbot/pkg/twitch/eventsub"
	"example.com/twitchbot/pkg/twitch/subscriptions"
	"golang.org/x/oauth2"
)

// runLocal attempts to run the twitch bot locally.
func runLocal(ctx context.Context) error {
	var (
		cfg                       = initTwitchConfig([]string{"user:read:chat"})
		usedKeyringToken          = true
		userID, broadcasterUserID string
	)

	token, err := fetchTokenFromKeyring()
	if err != nil {
		usedKeyringToken = false
		token, err = fetchTokenFromTwitch(ctx, cfg)
		if err != nil {
			return err
		}
	}

	// Our websocket is useless without having a valid token for the Twitch API, so wait to have one before we continue.
	client := twitch.New(cfg, token)
	// We need to get the user IDs!
	users, err := client.Users(ctx, &twitch.UsersRequest{
		Login: []string{"punmechanic", "piratesoftware"},
	})
	if err != nil {
		return fmt.Errorf("fetch users: %w", err)
	}

	for _, user := range users.Data {
		if user.Login == "punmechanic" {
			userID = user.ID
		} else if user.Login == "piratesoftware" {
			broadcasterUserID = user.ID
		}
	}

	conn, err := eventsub.Dial(ctx)
	if err != nil {
		return fmt.Errorf("init websocket: %s", err)
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
	err = client.SubscribeEvents(ctx, []*twitch.SubscribeRequest{
		{
			Type: subscriptions.ChannelChatMessage,
			Condition: eventsub.Condition{
				UserID:            userID,
				BroadcasterUserID: broadcasterUserID,
			},
			Transport: eventsub.Transport{
				Method:    eventsub.MethodWebsocket,
				SessionID: id,
			},
		},
	})

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

	// If we got here, the token we used was valid, and we should store it.
	// TODO: If the token changes while the bot is running, we will still be storing the old value. We should make sure
	// to catch any new tokens if the token is refreshed. We can probably do this by exposing a channel on twitch.Client
	if err := saveTokenInKeyring(token); err != nil {
		// Print a warning but don't die
		log.Printf("could not save token to OS keyring: %s", err)
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
