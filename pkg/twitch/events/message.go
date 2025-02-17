package events

type Message struct {
	Text      string            `json:"text"`
	Fragments []MessageFragment `json:"fragments"`
}

type MessageFragment struct {
	Type      string   `json:"type"`
	Text      string   `json:"text"`
	Cheermote *string  `json:"cheermote"`
	Emote     *Emote   `json:"emote"`
	Mention   *Mention `json:"mention"`
}

type Mention struct{}

type Badge struct {
	SetID string `json:"set_id"`
	ID    string `json:"id"`
	Info  string `json:"info"`
}

type Reply struct{}

type Emote struct{}
