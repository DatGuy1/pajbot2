package bot

import (
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/format"
)

// Format formats the given line xD
func (bot *Bot) Format(line string, msg *common.Msg) string {
	// catch all errors until we have proper error handling
	defer func() {
		if r := recover(); r != nil {
			log.Error(r)
		}
	}()
	fmtline, rawCommands := format.ParseLine(line)
	for i := range rawCommands {
		bot.ExecCommand(&rawCommands[i], msg)
	}
	return format.RunCommands(fmtline, rawCommands)
}

/*
ExecCommand xD
*/
func (bot *Bot) ExecCommand(cmd *format.Command, msg *common.Msg) {
	switch cmd.C {
	case "source", "sender":
		log.Debugf("Parsing for %s", msg.User.Name)
		cmd.Outcome = format.ParseUser(&msg.User, cmd.SubC)
	case "user":
		if msg.Args != nil {
			user := common.User{
				Name: msg.Args[0],
			}
			if bot.Redis.IsValidUser(msg.Channel, msg.Args[0]) {
				log.Debugf("xD: %s", user.Name)
				user.LoadData(bot.Redis.Pool, bot.Channel.Name)
				defer user.RedisData.CloseWithoutSaving()
				cmd.Outcome = format.ParseUser(&user, cmd.SubC)
				return
			}
		}
		cmd.Outcome = format.ParseUser(&msg.User, cmd.SubC)
	case "lasttweet":
		cmd.Outcome = bot.Twitter.LastTweetString(cmd.SubC[0])
	}

}
