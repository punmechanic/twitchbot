package eventsub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"

	"example.com/twitchbot/pkg/twitch/events"
	"example.com/twitchbot/pkg/twitch/eventsub/subscriptions"
	"golang.org/x/net/websocket"
)

type state struct {
	SessionID string
}

func Dial(ctx context.Context) (*Conn, error) {
	ws, err := websocket.Dial("wss://eventsub.wss.twitch.tv/ws", "wss", "wss://eventsub.wss.twitch.tv/ws")
	if err != nil {
		log.Fatalf("init websocket: %s", err)
	}

	conn := &Conn{
		r:                  ws,
		SessionID:          make(chan string, 1),
		ChannelFollowed:    make(chan events.ChannelFollow),
		ChannelChatMessage: make(chan events.ChannelChatMessage),
	}

	return conn, err
}

// Conn is a websocket connection to the eventsub API.
//
// Websocket connections to the Eventsub API are read-only. Receive notifications from the Notifications channel, and then respond to them using the Twitch API.
type Conn struct {
	// SessionID is a channel that can be interrogated to retrieve the session id of the Conn.
	//
	// This will receive something every time the Conn receives a session_welcome event, which means it may occur if the session reconnects.
	//
	// You should capture values from this channel and present them to the Twitch API to register subscriptions; without subscriptions, the Conn will be disconnected after 10 seconds.
	SessionID chan string
	r         io.ReadCloser

	ChannelFollowed    chan events.ChannelFollow
	ChannelChatMessage chan events.ChannelChatMessage
}

func (c *Conn) processNotification(notification *Notification) error {
	// Parse the notification data, and then pass it to the appropriate channels.
	var (
		channelFollow      events.ChannelFollow
		channelChatMessage events.ChannelChatMessage
		err                error
	)

	typ := notification.Subscription.Type
	switch typ {
	case subscriptions.ChannelFollow.Name:
		err = json.Unmarshal(notification.Event, &channelFollow)
		if err == nil {
			c.ChannelFollowed <- channelFollow
		}
	case subscriptions.ChannelChatMessage.Name:
		err := json.Unmarshal(notification.Event, &channelChatMessage)
		if err == nil {
			c.ChannelChatMessage <- channelChatMessage
		}
	default:
		return fmt.Errorf("unknown subscription type %q", typ)
	}
	if err != nil {
		return &UnmarshalEventErr{JSON: string(notification.Event), Type: typ, Cause: err}
	}

	return nil
}

func (c *Conn) serveMessage(_ context.Context, msg *message) error {
	switch msg.Metadata.MessageType {
	case "session_welcome":
		var payload welcome
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return err
		}

		c.SessionID <- payload.Session.ID
	case "keepalive":
	case "notification":
		var (
			payload           Notification
			unmarshalEventErr *UnmarshalEventErr
		)

		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return err
		}

		err := c.processNotification(&payload)
		if errors.As(err, &unmarshalEventErr) {
			// Silently drop, but don't kill the bot.
			log.Printf("%s", err)
		} else if err != nil {
			return err
		}
	}

	return nil
}

func (c *Conn) Listen() error {
	// Placeholder context that isn't used for anything.
	// One might use this to enforce timeouts on message handling.
	ctx := context.Background()
	dec := json.NewDecoder(c.r)
	for {
		var msg message
		if err := dec.Decode(&msg); err != nil {
			return err
		}

		err := c.serveMessage(ctx, &msg)
		if err != nil {
			return err
		}
	}
}

type UnmarshalEventErr struct {
	JSON  string
	Type  string
	Cause error
}

func (e *UnmarshalEventErr) Error() string {
	return fmt.Sprintf("unmarshal event %s, JSON: %s: %s", e.Type, e.JSON, e.Cause)
}
