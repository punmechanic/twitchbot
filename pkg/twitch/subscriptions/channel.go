package subscriptions

var (
	ChannelFollow = Type{
		Version: "2",
		Name:    "channel.follow",
	}

	// ChannelChatMessage requires user:read:chat scope from the chatting user.
	ChannelChatMessage = Type{
		Version: "1",
		Name:    "channel.chat.message",
	}
)
