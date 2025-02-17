package twitchbot

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"strings"

	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

func fetchTokenFromKeyring() (*oauth2.Token, error) {
	token, err := keyring.Get("twitchbot", "twitch:tokens")
	if err != nil {
		return nil, err
	}

	var r io.Reader = strings.NewReader(token)
	r = base64.NewDecoder(base64.RawStdEncoding, r)
	de := json.NewDecoder(r)
	var tok *oauth2.Token
	err = de.Decode(&tok)
	return tok, err
}

func saveTokenInKeyring(token *oauth2.Token) error {
	var buf strings.Builder
	enc := json.NewEncoder(base64.NewEncoder(base64.RawStdEncoding, &buf))
	err := enc.Encode(token)
	if err != nil {
		return err
	}

	return keyring.Set("twitchbot", "twitch:tokens", buf.String())
}
