package modules

import (
	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/command"
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/common/basemodule"
)

/*
Uptime xD
*/
type Uptime struct {
	basemodule.BaseModule

	commandHandler command.Handler
}

// Ensure the module implements the interface properly
var _ Module = (*Uptime)(nil)

// NewUptime xD
func NewUptime() *Uptime {
	m := Uptime{
		BaseModule: basemodule.NewBaseModule(),
	}
	m.ID = "uptime"
	m.EnabledDefault = true
	return &m
}

// Init xD
func (module *Uptime) Init(bot *bot.Bot) (string, bool) {
	module.ParseState(bot.Redis, bot.Channel.Name)

	uptimeCommand := &command.FuncCommand{
		BaseCommand: command.BaseCommand{
			Triggers: []string{
				"uptime",
				"downtime",
			},
		},
		Function: module.cmdUptime,
	}
	module.commandHandler.AddCommand(uptimeCommand)

	return "uptime", isModuleEnabled(bot, "uptime", true)
}

// DeInit xD
func (module *Uptime) DeInit(b *bot.Bot) {

}

// Check xD
func (module *Uptime) Check(b *bot.Bot, msg *common.Msg, action *bot.Action) error {
	return module.commandHandler.Check(b, msg, action)
}

func (module *Uptime) cmdUptime(b *bot.Bot, msg *common.Msg, action *bot.Action) {
	if b.Channel.Online() {
		b.Mentionf(msg.User, "%s has been online for %s", b.Channel.Name, b.Channel.UptimeString())
	} else {
		b.Mentionf(msg.User, "%s has been offline for %s", b.Channel.Name, b.Channel.DowntimeString())
	}
}
