package eventsub

import (
	"time"
)

type CreateEventSubSubscriptionRequest struct {
	Subscriptions []*SubscriptionRequest
}

type CreateEventSubSubscriptionResponse struct {
	Data         []Subscription `json:"data"`
	Total        int            `json:"total"`
	TotalCost    int            `json:"total_cost"`
	MaxTotalCost int            `json:"max_total_cost"`
}

var (
	SubscriptionMethodWebhook   SubscriptionMethod = "webhook"
	SubscriptionMethodWebsocket SubscriptionMethod = "websocket"
	SubscriptionMethodConduit   SubscriptionMethod = "conduit"
)

var (
	SubscriptionStatusEnabled                            SubscriptionStatus = "enabled"
	SubscriptionStatusWebhookCallbackVerificationPending SubscriptionStatus = "webhook_callback_verification_pending"
)

type SubscriptionMethod string

type SubscriptionStatus string

type SubscriptionCondition struct{}
type SubscriptionTransport struct {
	Method SubscriptionMethod `json:"method"`

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
	ID          string                `json:"id"`
	Status      SubscriptionStatus    `json:"status"`
	Type        string                `json:"type"`
	Version     string                `json:"version"`
	Condition   SubscriptionCondition `json:"condition"`
	CreatedAt   time.Time             `json:"created_at"`
	Transport   SubscriptionTransport `json:"transport"`
	ConnectedAt time.Time             `json:"connected_at"`
	ConduitID   string                `json:"conduit_id"`
	Cost        int                   `json:"cost"`
}

type SubscriptionRequest struct {
	Type      string                `json:"type"`
	Version   string                `json:"version"`
	Condition SubscriptionCondition `json:"condition"`
	Transport SubscriptionTransport `json:"transport"`
}
