package twitch

import (
	"context"
	"net/url"
	"time"
)

type UsersRequest struct {
	ID    []string
	Login []string
}

type User struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
	// ViewCount is deprecated.
	ViewCount int       `json:"view_count"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at`
}

type UsersResponse struct {
	Data []User `json:"data"`
}

func (c *Client) Users(ctx context.Context, r *UsersRequest) (*UsersResponse, error) {
	values := url.Values{}
	for _, id := range r.ID {
		values.Add("id", id)
	}
	for _, login := range r.Login {
		values.Add("login", login)
	}

	req, err := c.NewRequest(ctx, "GET", "/helix/users", values)
	if err != nil {
		return nil, err
	}

	var resp UsersResponse
	err = c.Execute(ctx, req, &resp)
	return &resp, err
}
