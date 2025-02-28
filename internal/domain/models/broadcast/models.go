package broadcast

type MessangerOptions struct {
	Channel string
	ChatID  string
	Color   string
}

func (o *MessangerOptions) GetChannel() string {
	return o.Channel
}

func (o *MessangerOptions) GetChatID() string {
	return o.ChatID
}

type MessageOptions struct {
	Color string
}

func (o *MessageOptions) GetColor() string {
	return o.Color
}
