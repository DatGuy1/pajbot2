package web

import (
	"strings"
	"time"

	"github.com/pajlada/pajbot2/common"
)

func getUserPayload(channel, username string) interface{} {
	username = strings.ToLower(username)
	if !redis.IsValidUser(channel, username) {
		return newError("user not found")
	}
	u := common.User{
		Name: username,
	}
	u.LoadData(redis.Pool, channel)
	defer u.RedisData.CloseWithoutSaving()
	return user{
		Name:                username,
		DisplayName:         u.DisplayName(),
		Points:              int64(u.RedisData.Points),
		Level:               int64(u.Level),
		TotalMessageCount:   int64(u.RedisData.TotalMessageCount),
		OfflineMessageCount: int64(u.RedisData.OfflineMessageCount),
		OnlineMessageCount:  int64(u.RedisData.OnlineMessageCount),
		LastSeen:            u.RedisData.LastSeen.Format(time.UnixDate),
		LastActive:          u.RedisData.LastActive.Format(time.UnixDate),
	}
}
