package eventsub

import (
	"context"
	"net/http"
)

type Requestor interface {
	Execute(ctx context.Context, req *http.Request, data any) error
	NewRequest(ctx context.Context, method, path string, data any) (*http.Request, error)
}

type SubscribeRequest struct {
	Type      string    `json:"type"`
	Version   string    `json:"version"`
	Condition Condition `json:"condition"`
	Transport Transport `json:"transport"`
}

type SubscribeResponse struct {
	Data         []Subscription `json:"data"`
	Total        int            `json:"total"`
	TotalCost    int            `json:"total_cost"`
	MaxTotalCost int            `json:"max_total_cost"`
}

func Subscribe(ctx context.Context, requestor Requestor, req *SubscribeRequest) (*SubscribeResponse, error) {
	var data SubscribeResponse
	r, err := requestor.NewRequest(ctx, "POST", "/helix/eventsub/subscriptions", req)
	if err != nil {
		return nil, err
	}

	err = requestor.Execute(ctx, r, &data)
	return &data, err
}
