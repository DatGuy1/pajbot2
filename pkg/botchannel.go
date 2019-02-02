package pkg

type BotChannel interface {
	DatabaseID() int64
	Channel() Channel
	ChannelID() string
	ChannelName() string

	EnableModule(string) error
	DisableModule(string) error

	Stream() Stream

	Say(string)
}
