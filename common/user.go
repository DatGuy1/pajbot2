package common

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/garyburd/redigo/redis"
)

// User defines a user as parsed directly from an IRC message
// TODO: the TotalMessageCount etc should probably be moved
// - OnlineMessageCount: MOVE
// - OfflineMessageCount: MOVE
// We should instead contain an object with RedisUser
type User struct {
	// I don't know what this ID is. User ID from TMI?
	ID int

	// Name in full lower case, also known as "login name"
	Name string

	// Name in variable case, also known as "international name" or "nick name"
	DisplayName string

	Mod                 bool
	Sub                 bool
	Turbo               bool
	ChannelOwner        bool
	Type                string // admin , staff etc
	Level               int
	TotalMessageCount   int
	OnlineMessageCount  int
	OfflineMessageCount int
	Points              int

	// Data stores
	RedisUser redisUser
}

const noPing = string("\u05C4")

// NameNoPing xD
func (u *User) NameNoPing() string {
	return string(u.DisplayName[0]) + noPing + u.DisplayName[1:]
}

// LoadRedisUser loads the redis user associated with this IRC user
func (u *User) LoadRedisUser(pool *redis.Pool, channel string) {
	LoadRedisUser(pool, channel, u.Name, &u.RedisUser)
}

// GetPoints returns a point amount relative to the user with the given arg.
// Available args:
// "all" or "allin": Return all of users points
// "50%": Return 50% of users points
// "13": Return 13 if the user has 13 points, otherwise return error
func (u *User) GetPoints(arg string) (int, error) {
	if arg == "all" || arg == "allin" {
		// Return all points
		return u.Points, nil
	}

	var bet int

	if strings.HasSuffix(arg, "%") {
		_bet, err := strconv.ParseFloat(arg[0:len(arg)-1], 64)
		if err != nil {
			log.Error(err)
			return 0, fmt.Errorf(`Invalid argument to GetPoints (using %%)`)
		}
		fpoints := float64(u.Points)
		bet = int(fpoints * (_bet / 100))
	} else {
		_bet, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("Invalid argument to GetPoints")
		}
		bet = int(_bet)
	}

	/*
		if bet > u.Points {
			return 0, fmt.Errorf("User does not have this many points")
		}
	*/

	return bet, nil
}
