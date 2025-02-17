package eventsub

import (
	"time"
)

type Method string

var (
	MethodWebhook   Method = "webhook"
	MethodWebsocket Method = "websocket"
	MethodConduit   Method = "conduit"
)

var (
	StatusEnabled                            Status = "enabled"
	StatusWebhookCallbackVerificationPending Status = "webhook_callback_verification_pending"
)

type Status string

// Condition is a list of all conditions for all eventsub evnt types.
//
// Check the documentation for which conditions are valid with which eventsub event types.
// https://dev.twitch.tv/docs/eventsub/eventsub-reference/
type Condition struct {
	UserID            string `json:"user_id,omitempty"`
	BroadcasterUserID string `json:"broadcaster_user_id,omitempty"`
	ModeratorUserID   string `json:"moderator_user_id,omitempty"`
}

type Transport struct {
	Method Method `json:"method"`

	// Callback is the callback URL where the notifications are sent. The URL must use the HTTPS protocol and port 443. See Processing an event. Specify this field only if method is set to webhook.
	//
	// Redirects are not followed.
	Callback *string `json:"callback,omitempty"`

	// Secret is the secret used to verify the signature. The secret must be an ASCII string that’s a minimum of 10 characters long and a maximum of 100 characters long. For information about how the secret is used, see Verifying the event message. Specify this field only if method is set to webhook.
	Secret *string `json:"secret,omitempty"`

	// SessionID is the ID that identifies the WebSocket to send notifications to. When you connect to EventSub using WebSockets, the server returns the ID in the Welcome message. Specify this field only if method is set to websocket.
	SessionID string `json:"session_id,omitempty"`

	// ConduitID is the ID that identifies the conduit to send notifications to. When you create a conduit, the server returns the conduit ID. Specify this field only if method is set to conduit.
	ConduitID string `json:"conduit_id,omitempty"`
}

type Subscription struct {
	ID          string    `json:"id"`
	Status      Status    `json:"status"`
	Type        string    `json:"type"`
	Version     string    `json:"version"`
	Condition   Condition `json:"condition"`
	CreatedAt   time.Time `json:"created_at"`
	Transport   Transport `json:"transport"`
	ConnectedAt time.Time `json:"connected_at"`
	ConduitID   string    `json:"conduit_id"`
	Cost        int       `json:"cost"`
}
