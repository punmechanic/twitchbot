package eventsub

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"example.com/twitchbot/pkg/twitch"
	"golang.org/x/net/websocket"
)

type state struct {
	SessionID string
}

func Dial(ctx context.Context, options ...Option) (*Conn, error) {
	ws, err := websocket.Dial("wss://eventsub.wss.twitch.tv/ws", "wss", "wss://eventsub.wss.twitch.tv/ws")
	if err != nil {
		log.Fatalf("init websocket: %s", err)
	}

	conn := &Conn{
		r:      ws,
		Events: make(chan Notification),
	}

	for _, opt := range options {
		opt(conn)
	}

	return conn, err
}

// Conn is a websocket connection to the eventsub API.
//
// Websocket connections to the Eventsub API are read-only. Receive Events from the Events channel, and then respond to them using the Twitch API.
type Conn struct {
	// Events is a channel that can be read from to poll events from the Conn.
	Events chan Notification

	r io.ReadCloser

	subscriptions []string
	client        *twitch.Client
}

func (c *Conn) serveMessage(msg *message) error {
	switch msg.Metadata.MessageType {
	case "session_welcome":
		var payload welcome
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return err
		}

		err := c.setupSubscriptions(payload.Session.ID)
		if err != nil {
			return fmt.Errorf("subscribe: %w", err)
		}

	case "keepalive":
	case "notification":
		var payload Notification
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return err
		}

		c.Events <- payload
	}

	return nil
}

func (c *Conn) Listen() error {
	dec := json.NewDecoder(c.r)
	for {
		var msg message
		if err := dec.Decode(&msg); err != nil {
			return err
		}

		err := c.serveMessage(&msg)
		if err != nil {
			return err
		}
	}
}

func (c *Conn) setupSubscriptions(sessionID string) error {
	// TODO: impl
	return nil
}

type Option func(*Conn)

func WithSubscriptions(events []string, client *twitch.Client) Option {
	return func(c *Conn) {
		c.subscriptions = events
		c.client = client
	}
}
