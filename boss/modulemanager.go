package boss

import (
	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/modules"
)

func modulesUnload(b *bot.Bot) {
	// De-init all already-loaded modules
	for _, module := range b.EnabledModules {
		module.DeInit(b)
	}

	b.EnabledModules = nil
}

func modulesInit(b *bot.Bot) {
	// TODO(pajlada): Select which modules should be loaded
	//                via a redis json list or something
	if b.Channel.Name == "nymn_hs" {
		b.AllModules = []bot.Module{
			&modules.Admin{},
			&modules.Debug{},
			&modules.Emotes{},
			modules.NewDebug(),
		}
	} else {
		b.AllModules = []bot.Module{
			&modules.Admin{},
			&modules.Banphrase{},
			&modules.Bingo{},
			&modules.Command{},
			&modules.Emotes{},
			&modules.Points{},
			&modules.Pyramid{},
			&modules.Raffle{},
			&modules.SubAnnounce{},
			&modules.Test{},
			&modules.Top{},
			modules.NewDebug(),
			modules.NewRoulette(),
		}
	}
}

func modulesLoad(b *bot.Bot) {
	// Initialize all loaded modules
	for _, module := range b.AllModules {
		module.Init(b)
	}

	b.EnabledModules = nil

	for _, module := range b.AllModules {
		state := module.GetState()
		if state.IsEnabled() {
			b.EnabledModules = append(b.EnabledModules, module)
		}
	}
}

// modulesReload unloads all loaded modules, then reloads all modules
// that should be enabled
func modulesReload(b *bot.Bot) {
	modulesUnload(b)
	modulesLoad(b)
}
