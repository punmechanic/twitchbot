package twitch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Authorization interface {
	Apply(r *http.Request) error
}

type Client struct {
	Authorization Authorization
	HttpClient    *http.Client
}

func (*Client) NewRequest(ctx context.Context, method, path string, data any) (*http.Request, error) {
	uri := url.URL{
		Scheme: "https",
		Host:   "api.twitch.tv",
		Path:   path,
	}
	buf, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return http.NewRequestWithContext(ctx, method, uri.String(), bytes.NewReader(buf))
}

func (c *Client) Execute(ctx context.Context, r *http.Request, dst any) error {
	r2 := r.Clone(ctx)
	if err := c.Authorization.Apply(r2); err != nil {
		return fmt.Errorf("authorization error: %w", err)
	}

	resp, err := c.HttpClient.Do(r2)
	if err != nil {
		return err
	}

	return c.parseResponse(resp, dst)
}

func (c *Client) parseResponse(resp *http.Response, dst any) error {
	if resp.StatusCode >= 400 {
		log.Printf("status code %d", resp.StatusCode)
		// impl
		panic("not yet implemented")
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, dst)
}

func New(auth Authorization) *Client {
	return &Client{
		Authorization: auth,
		HttpClient:    http.DefaultClient,
	}
}
