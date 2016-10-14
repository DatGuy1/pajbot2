package modules

import (
	"fmt"
	"strings"

	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/common/basemodule"
	"github.com/pajlada/pajbot2/points"
)

// Points module
type Points struct {
	basemodule.BaseModule
}

var _ Module = (*Points)(nil)

// NewPoints xD
func NewPoints() *Points {
	return &Points{
		BaseModule: basemodule.BaseModule{
			ID: "points",
		},
	}
}

// Init xD
func (module *Points) Init(bot *bot.Bot) (string, bool) {
	module.SetDefaults("points")
	module.EnabledDefault = true
	module.ParseState(bot.Redis, bot.Channel.Name)

	return "points", isModuleEnabled(bot, "points", true)
}

// DeInit xD
func (module *Points) DeInit(b *bot.Bot) {

}

// Check xD
func (module *Points) Check(b *bot.Bot, msg *common.Msg, action *bot.Action) error {
	if !strings.HasPrefix(msg.Text, "!") {
		return nil
	}
	m := strings.ToLower(msg.Text)
	spl := strings.Split(m, " ")
	trigger := spl[0]
	var args []string
	if len(spl) > 1 {
		args = spl[1:]
	}

	// using pts to not trigger other bots
	switch trigger {
	case "!givepts":
		err := points.GivePoints(b, &msg.User, args)
		if err != nil {
			b.Say(fmt.Sprint(err))
		}
	case "!pts":
		msg.Args = args
		b.SaySafe(b.Format("$(user.name) has $(user.points) points KKaper", msg))
	case "!resetpts":
		msg.User.Points = 0
		b.Redis.SetPoints(msg.Channel, &msg.User)
	}

	return nil
}
