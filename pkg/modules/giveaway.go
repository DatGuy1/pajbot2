package modules

import (
	"math/rand"
	"strings"

	"github.com/pajlada/pajbot2/pkg"
)

type giveaway struct {
	botChannel pkg.BotChannel

	server *server

	state string

	entrants []string
}

func newGiveaway() pkg.Module {
	return &giveaway{
		server: &_server,
		state:  "inactive",
	}
}

var giveawaySpec = moduleSpec{
	id:    "giveaway",
	name:  "Giveaway",
	maker: newGiveaway,
}

func (m *giveaway) Initialize(botChannel pkg.BotChannel, settings []byte) error {
	m.botChannel = botChannel

	return nil
}

func (m *giveaway) Disable() error {
	return nil
}

func (m *giveaway) Spec() pkg.ModuleSpec {
	return &giveawaySpec
}

func (m *giveaway) BotChannel() pkg.BotChannel {
	return m.botChannel
}

const forsen25ID = "1788703"
const pajlada25ID = "908917"
const pajaWID = "80481"

func (m giveaway) OnWhisper(bot pkg.Sender, user pkg.User, message pkg.Message) error {
	return nil
}

func (m *giveaway) OnMessage(bot pkg.Sender, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) error {
	// if channel.GetChannel() != "forsen" {
	// 	return nil
	// }

	const giveawayEmote = forsen25ID

	text := message.GetText()

	// Commands
	if user.IsModerator() || user.GetName() == "forsen" || user.GetName() == "pajlada" {
		if strings.HasPrefix(text, "!25start") {
			if m.state == "inactive" {
				m.state = "started"
				m.entrants = []string{}
				bot.Say(channel, "Started giveaway")

				return nil
			} else {
				bot.Say(channel, "Giveaway already started")
				return nil
			}
		}

		if strings.HasPrefix(text, "!25stop") {
			if m.state == "started" {
				m.state = "inactive"
				bot.Say(channel, "Stopped accepting people into the giveaway")

				return nil
			}
		}

		if strings.HasPrefix(text, "!25draw") {
			if len(m.entrants) == 0 {
				bot.Say(channel, "No one has joined the giveaway")
				return nil
			}

			winnerIndex := rand.Intn(len(m.entrants))
			winnerUsername := m.entrants[winnerIndex]
			bot.Say(channel, winnerUsername+" just won the sub emote giveaway PogChamp")

			m.entrants = append(m.entrants[:winnerIndex], m.entrants[winnerIndex+1:]...)

			return nil
		}
	}

	if m.state == "started" {
		enterGiveaway := false

		reader := message.GetTwitchReader()
		for reader.Next() {
			emote := reader.Get()
			if emote.GetID() == giveawayEmote {
				// bot.Say(channel, fmt.Sprintf("%#v", emote))
				// User wants to enter the giveaway
				enterGiveaway = true
				break
			}
		}

		if enterGiveaway {
			for _, entrant := range m.entrants {
				if entrant == user.GetName() {
					// User has already joined
					return nil
				}

			}
			m.entrants = append(m.entrants, user.GetName())

			bot.Say(channel, "@"+user.GetName()+", you have been entered into the sub emote giveaway")
			return nil
		}
	}

	return nil
}
