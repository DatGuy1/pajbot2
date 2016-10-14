package modules

import (
	"fmt"
	"strings"
	"time"

	"github.com/dankeroni/gotwitch"
	"github.com/pajlada/pajbot2/apirequest"
	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/command"
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/common/basemodule"
	"github.com/pajlada/pajbot2/helper"
)

/*
Debug xD
*/
type Debug struct {
	basemodule.BaseModule

	commandHandler command.Handler
}

// Ensure the module implements the interface properly
var _ Module = (*Debug)(nil)

// NewDebug xD
func NewDebug() *Debug {
	m := Debug{
		BaseModule: basemodule.NewBaseModule(),
	}
	m.ID = "debug"
	return &m
}

func cmdDebug(b *bot.Bot, msg *common.Msg, action *bot.Action) {
	m := helper.GetTriggersN(msg.Text, 1)

	if len(m) == 0 {
		// missing argument to !test
		return
	}

	switch m[0] {
	case "say":
		b.Say(strings.Join(m[1:], " "))
	case "lasttweet":
		if len(m) > 1 {
			tweet := b.Twitter.LastTweetString(m[1])
			b.Sayf("last tweet from %s ", tweet)
		} else {
			b.Say("Usage: !test lasttweet pajlada")
		}
	case "follow":
		if len(m) > 1 {
			b.Twitter.Follow(m[1])
			b.Sayf("now streaming %s's timeline", m[1])
		} else {
			b.Say("Usage: !test follow pajlada")
		}
	case "unfollow":
		b.Say("not implemented yet")
	case "api":
		if len(m) > 1 {
			apirequest.Twitch.GetStream(m[1],
				func(stream gotwitch.Stream) {
					b.Sayf("Stream info: %+v", stream)
				},
				func(statusCode int, statusMessage, errorMessage string) {
					b.Sayf("ERROR: %d", statusCode)
					b.Say(statusMessage)
					b.Say(errorMessage)
				}, func(err error) {
					b.Say("Internal error")
				})
		} else {
			b.Say("Usage: !test api pajlada")
		}
	case "resub":
		testMessage := `@badges=staff/1,broadcaster/1,turbo/1;color=#008000;display-name=TWITCH_UserName;emotes=;mod=0;msg-id=resub;msg-param-months=6;room-id=1337;subscriber=1;system-msg=TWITCH_UserName\shas\ssubscribed\sfor\s6\smonths!;login=twitch_username;turbo=1;user-id=1337;user-type=staff :tmi.twitch.tv USERNOTICE #%s :Great stream -- keep it up!`
		b.RawRead <- fmt.Sprintf(testMessage, b.Channel.Name)

	case "whisper":
		log.Debugf("WHISPER %s", msg.User.Name)
		b.Whisper(msg.User.Name, "TEST WHISPER")

	default:
		b.Sayf("Unhandled action %s", m[0])
		return
	}
}

func cmdMyInfo(b *bot.Bot, msg *common.Msg, action *bot.Action) {
	b.Mentionf(msg.User, "ID: %d, username: %s, type: %s, level: %d",
		msg.User.ID, msg.User.DisplayName, msg.User.Type, msg.User.Level)
}

// Init xD
func (module *Debug) Init(bot *bot.Bot) (string, bool) {
	module.SetDefaults("debug")
	module.EnabledDefault = true
	module.ParseState(bot.Redis, bot.Channel.Name)

	myInfoCommand := command.NewFuncCommand([]string{"myinfo"})
	myInfoCommand.Function = cmdMyInfo
	myInfoCommand.Cooldown = time.Second * 2

	module.commandHandler.AddCommand(myInfoCommand)

	return "test", true
}

// DeInit xD
func (module *Debug) DeInit(b *bot.Bot) {

}

// Check xD
func (module *Debug) Check(b *bot.Bot, msg *common.Msg, action *bot.Action) error {
	return module.commandHandler.Check(b, msg, action)
}
