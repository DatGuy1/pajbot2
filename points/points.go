package points

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/plog"
)

var log = plog.GetLogger()

// GivePoints xD
func GivePoints(b *bot.Bot, user *common.User, args []string) error {
	if len(args) < 2 {
		return errors.New("not enough args xD")
	}
	targetUserName := args[0]
	pts, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return err
	}
	if pts > 1000 || pts < -1000 {
		return fmt.Errorf("you can only give 1000 points at a time pajaHop")
	}
	if !b.Redis.IsValidUser(b.Channel.Name, targetUserName) {
		return errors.New("invalid user")
	}
	targetUser := common.RedisUser{}
	log.Debug("Load redis user...")
	err = common.LoadRedisUser(b.Redis.Pool, b.Channel.Name, targetUserName, &targetUser)
	if err != nil {
		log.Error(err)
		targetUser.CloseWithoutSaving()
		return err
	}
	oldUser := targetUser
	log.Debug(oldUser.Points)
	log.Debug(targetUser.Points)
	targetUser.Points += int(pts)
	log.Debug(targetUser.Points)
	targetUser.Close(&oldUser)
	return nil
}
