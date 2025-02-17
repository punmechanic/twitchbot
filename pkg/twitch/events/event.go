package events

import "time"

type ChannelFollow struct {
	UserID               string    `json:"user_id"`
	UserLogin            string    `json:"user_login"`
	UserName             string    `json:"user_name"`
	BroadcasterUserID    string    `json:"broadcaster_user_id"`
	BroadcasterUserLogin string    `json:"broadcaster_user_login"`
	BroadcasterUserName  string    `json:"broadcaster_user_name"`
	FollowedAt           time.Time `json:"followed_at"`
}

type ChannelChatMessage struct {
	BroadcasterUserID           string  `json:"broadcaster_user_id"`
	BroadcasterUserLogin        string  `json:"broadcaster_user_login"`
	BroadcasterUserName         string  `json:"broadcaster_user_name"`
	ChatterUserID               string  `json:"chatter_user_id"`
	ChatterUserLogin            string  `json:"chatter_user_login"`
	ChatterUserName             string  `json:"chatter_user_name"`
	MessageID                   string  `json:"message_id"`
	Message                     Message `json:"message"`
	Color                       string  `json:"color"`
	Badges                      []Badge `json:"badges"`
	MessageType                 string  `json:"message_type"`
	Cheer                       *string `json:"cheer"`
	Reply                       *Reply  `json:"reply"`
	ChannelPointsCustomRewardID *string `json:"channel_points_custom_reward_id"`
}
