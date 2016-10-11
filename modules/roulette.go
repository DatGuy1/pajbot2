package modules

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/command"
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/common/basemodule"
	"github.com/pajlada/pajbot2/helper"
)

// Roulette module
type Roulette struct {
	basemodule.BaseModule
	commandHandler command.Handler

	// Settings
	WinMessage         StringSetting
	LoseMessage        StringSetting
	WinPercentage      IntSetting
	OnlineCooldown     IntSetting
	OnlineUserCooldown IntSetting
	MinimumBet         IntSetting
	MaximumBet         IntSetting

	// TODO: Implement
	CanWhisper BoolSetting
	// TODO: Implement
	OutputType OptionSetting
	// TODO: Implement
	MinimumShowPoints IntSetting
	// TODO: Implement
	OnlyAllowRouletteAfterSub IntSetting
}

var _ Module = (*Roulette)(nil)

// NewRoulette xD
func NewRoulette() *Roulette {
	return &Roulette{
		BaseModule: basemodule.BaseModule{
			ID: "roulette",
		},
		WinMessage: StringSetting{
			Label:        `Win message | %d = Points bet`,
			DefaultValue: `$(source.name) won %d points in roulette and now has $(source.points) points! FeelsGoodMan`,
			MinLength:    10,
			MaxLength:    400,
		},
		LoseMessage: StringSetting{
			Label:        `Lose message | %d = Points bet`,
			DefaultValue: `$(source.name) lost %d points in roulette and now has $(source.points) points! FeelsBadMan`,
			MinLength:    10,
			MaxLength:    400,
		},
		WinPercentage: IntSetting{
			Label:        `Win % chance: 0-100`,
			DefaultValue: 50,
			MinValue:     0,
			MaxValue:     100,
		},
		OnlineCooldown: IntSetting{
			Label:        `Online global cooldown (seconds)`,
			DefaultValue: 0,
			MinValue:     0,
			MaxValue:     600,
		},
		OnlineUserCooldown: IntSetting{
			Label:        `Online user cooldown (seconds)`,
			DefaultValue: 60,
			MinValue:     0,
			MaxValue:     600,
		},
		MinimumBet: IntSetting{
			Label:        `Minimum roulette bet`,
			DefaultValue: 1,
			MinValue:     1,
			MaxValue:     50000,
		},
		MaximumBet: IntSetting{
			Label:        `Maximum roulette bet`,
			DefaultValue: -1,
			MinValue:     -1,
			MaxValue:     50000000,
		},
		CanWhisper: BoolSetting{
			Label:        `Can roulette in whisper`,
			DefaultValue: false,
		},
		OutputType: OptionSetting{
			Label:        `Result output type`,
			DefaultValue: 0,
			Options: []string{
				`Show results in chat`,
				`Show results in whispers`,
				`Show results in chat if it's over X points, else whisper results`,
			},
		},
		MinimumShowPoints: IntSetting{
			Label:        `Minimum points you need to bet to have the results show up in chat (with option 3)`,
			DefaultValue: 100,
			MinValue:     1,
			MaxValue:     150000,
		},
		OnlyAllowRouletteAfterSub: IntSetting{
			Label:        `Only allow roulette after subbing (-1 = off, rest = time in seconds people can roulette after a sub occurs)`,
			DefaultValue: -1,
			MinValue:     -1,
			MaxValue:     600,
		},
	}
}

// Check roulette
func (module *Roulette) cmdRoulette(b *bot.Bot, msg *common.Msg, action *bot.Action) {
	const usageMessage = `Usage: !roul 123 OR !roul all`
	const betTooSmallMessage = `You can't roulette 0 or less points OMGScoots`
	const betTooBigMessage = `You can't roulette more points than you have. You have %d points to use`
	const betBelowLimitMessage = `You need to roulette at least %d points`
	const betAboveLimitMessage = `You can roulette at most %d points`

	minBet := module.MinimumBet.Int()
	maxBet := module.MaximumBet.Int()
	bet := 0

	args := helper.GetTriggersN(msg.Text, 1)
	user := &msg.User
	if user.Points == 0 {
		b.Mention(msg.User, "You don't have enough points to roulette")
		return
	}

	if len(args) < 1 {
		b.Mention(msg.User, usageMessage)
		return
	}

	if args[0] == "all" || args[0] == "allin" {
		bet = user.Points
	} else {
		_bet, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			b.Mention(msg.User, usageMessage)
		}
		bet = int(_bet)
	}

	if bet <= 0 {
		b.Mention(msg.User, betTooSmallMessage)
		return
	}

	if bet < minBet {
		b.Mentionf(msg.User, betBelowLimitMessage, minBet)
		return
	}

	if maxBet > 0 && bet > maxBet {
		b.Mentionf(msg.User, betAboveLimitMessage, maxBet)
		return
	}

	if bet > user.Points {
		b.Mentionf(msg.User, betTooSmallMessage, user.Points)
		return
	}

	module.runRoulette(b, msg, bet)
}

func (module *Roulette) runRoulette(b *bot.Bot, msg *common.Msg, points int) {
	user := &msg.User
	won := module.WinPercentage.Int() >= rand.Intn(101)
	if won {
		user.Points += points
		b.SayFormat(module.WinMessage.String(), msg, points)
	} else {
		user.Points -= points
		b.SayFormat(module.LoseMessage.String(), msg, points)
	}

}

// Init xD
func (module *Roulette) Init(bot *bot.Bot) (string, bool) {
	module.SetDefaults("roulette")
	module.EnabledDefault = true
	module.ParseState(bot.Redis, bot.Channel.Name)

	rouletteCommand := command.NewFuncCommand([]string{"roul"})
	rouletteCommand.Function = module.cmdRoulette
	rouletteCommand.Cooldown = time.Second * time.Duration(module.OnlineCooldown.Int())
	rouletteCommand.UserCooldown = time.Second * time.Duration(module.OnlineUserCooldown.Int())

	module.commandHandler.AddCommand(rouletteCommand)

	return "roulette", isModuleEnabled(bot, "roulette", true)
}

// DeInit xD
func (module *Roulette) DeInit(b *bot.Bot) {

}

// Check xD
func (module *Roulette) Check(b *bot.Bot, msg *common.Msg, action *bot.Action) error {
	return module.commandHandler.Check(b, msg, action)
}
