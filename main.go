package main

import (
	"encoding/json"
	"io"
	"log"

	"example.com/twitchbot/pkg/twitch/eventsub"
	"golang.org/x/net/websocket"
)

func main() {
	conn, err := websocket.Dial("wss://eventsub.wss.twitch.tv/ws", "wss", "wss://eventsub.wss.twitch.tv/ws")
	if err != nil {
		log.Fatalf("init websocket: %s", err)
	}
	defer conn.Close()

	err = Loop(conn, conn)
	if err != nil {
		log.Fatalf("loop: %s", err)
	}
}

func Loop(w io.Writer, r io.Reader) error {
	// enc := json.NewEncoder(conn)
	dec := json.NewDecoder(r)
	for {
		var msg eventsub.Message
		if err := dec.Decode(&msg); err != nil {
			return err
		}

		log.Printf("%#v", msg)
	}
}
