package eventsub

import (
	"encoding/json"
	"time"
)

type metadata struct {
	MessageID        string    `json:"message_id"`
	MessageType      string    `json:"message_type"`
	MessageTimestamp time.Time `json:"message_timestamp"`
}

type message struct {
	Metadata metadata        `json:"metadata"`
	Payload  json.RawMessage `json:"payload"`
}

type welcome struct {
	Session struct {
		ID                      string    `json:"id"`
		Status                  string    `json:"status"`
		ConnectedAt             time.Time `json:"connected_at"`
		KeepaliveTimeoutSeconds int       `json:"keepalive_timeout_seconds"`
		ReconnectURL            *string   `json:"reconnect_url"`
	} `json:"session"`
}

type keepalive struct{}

type notification struct {
	Subscription struct{} `json:"subscription"`
	Event        json.RawMessage
}
