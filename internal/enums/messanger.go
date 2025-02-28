package enums

type Messenger string

const (
	Telegram   Messenger = "tg"
	MatterMost Messenger = "mm"
)

func Messengers() []Messenger {
	return []Messenger{
		Telegram, MatterMost,
	}
}
