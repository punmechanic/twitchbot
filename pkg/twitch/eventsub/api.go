package eventsub

import (
	"context"
	"net/http"
)

type Requestor interface {
	Execute(ctx context.Context, req *http.Request, data any) error
	NewRequest(ctx context.Context, method, path string, data any) (*http.Request, error)
}

func Subscribe(ctx context.Context, requestor Requestor, req *CreateEventSubSubscriptionRequest) (*CreateEventSubSubscriptionResponse, error) {
	var data CreateEventSubSubscriptionResponse
	r, err := requestor.NewRequest(ctx, "POST", "/helix/eventsub/subscriptions", req)
	if err != nil {
		return nil, err
	}

	err = requestor.Execute(ctx, r, &data)
	return &data, err
}
