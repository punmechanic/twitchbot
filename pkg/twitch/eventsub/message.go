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

type Subscription struct {
	ID        string            `json:"id"`
	Status    string            `json:"status"`
	Type      string            `json:"type"`
	Version   string            `json:"version"`
	Cost      string            `json:"cost"`
	Condition map[string]string `json:"condition"`
	Transport struct {
		Method    string `json:"string"`
		SessionID string `json:"session_id"`
	} `json:"transport"`
	CreatedAt time.Time `json:"created_at"`
}

type Notification struct {
	Subscription Subscription    `json:"subscription"`
	Event        json.RawMessage `json:"event"`
}
