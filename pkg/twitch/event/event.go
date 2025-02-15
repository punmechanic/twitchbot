package event

import "time"

type ChannelFollow struct {
	UserID               string    `json:"user_id"`
	UserLogin            string    `json:"user_login"`
	UserName             string    `json:"user_name"`
	BroadcasterUserID    string    `json:"broadcaster_user_id"`
	BroadcasterUserLogin string    `json:"broadcaster_user_login"`
	BroadcasterUserName  string    `json:"broadcaster_user_name"`
	FollowedAt           time.Time `json:"followed_at`
}
