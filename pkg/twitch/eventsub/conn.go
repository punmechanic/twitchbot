package eventsub

import (
	"context"
	"encoding/json"
	"io"
	"log"

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
		r:             ws,
		Notifications: make(chan Notification),
		SessionID:     make(chan string, 1),
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
	// Notifications is a channel that can be read from to poll events from the Conn.
	Notifications chan Notification
	r             io.ReadCloser
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
		var payload Notification
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return err
		}

		c.Notifications <- payload
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
