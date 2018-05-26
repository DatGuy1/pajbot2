package pkg

import (
	twitch "github.com/gempir/go-twitch-irc"
)

type Module interface {
	Register() error
	OnMessage(channel string, user twitch.User, message twitch.Message) error
}