package twitch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

type Client struct {
	Config     *oauth2.Config
	HttpClient *http.Client
}

func (*Client) NewRequest(ctx context.Context, method, path string, data any) (*http.Request, error) {
	uri := url.URL{
		Scheme: "https",
		Host:   "api.twitch.tv",
		Path:   path,
	}

	var r io.Reader
	urlValues, isURLEncoded := data.(url.Values)
	if isURLEncoded {
		uri.RawQuery = urlValues.Encode()
	} else {
		buf, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		r = bytes.NewReader(buf)
	}

	req, err := http.NewRequestWithContext(ctx, method, uri.String(), r)
	if err != nil {
		return nil, err
	}

	if !isURLEncoded {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) Execute(ctx context.Context, r *http.Request, dst any) error {
	r2 := r.Clone(ctx)
	r2.Header.Set("Client-Id", c.Config.ClientID)
	resp, err := c.HttpClient.Do(r2)
	if err != nil {
		return err
	}

	return parseResponse(resp, dst)
}

func parseBadResponse(resp *http.Response) error {
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusBadRequest:
		var badRequestErr BadRequestError
		if err := json.Unmarshal(buf, &badRequestErr); err != nil {
			return err
		}
		return badRequestErr
	case http.StatusUnauthorized:
		return ErrUnauthorized
	default:
		return fmt.Errorf("status code %d", resp.StatusCode)
	}
}

func parseResponse(resp *http.Response, dst any) error {
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return parseBadResponse(resp)
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(buf, dst)
}

func New(cfg *oauth2.Config, initialToken *oauth2.Token) *Client {
	return &Client{
		Config:     cfg,
		HttpClient: cfg.Client(context.Background(), initialToken),
	}
}
