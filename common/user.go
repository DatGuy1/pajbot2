package common

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/garyburd/redigo/redis"
)

// User defines a user as parsed directly from an IRC message
type User struct {
	// I don't know what this ID is. User ID from TMI?
	ID int

	// Name in full lower case, also known as "login name"
	Name string

	// Name in variable case, also known as "international name" or "nick name"
	IRCDisplayName string `json:"DisplayName"`

	Mod          bool
	Sub          bool
	Turbo        bool
	ChannelOwner bool
	Type         string // admin , staff etc
	Level        int

	// Data stores
	RedisData RedisUser
}

const noPing = string("\u05C4")

// DisplayName xD
func (u *User) DisplayName() string {
	if u.IRCDisplayName == "" {
		if u.RedisData.DisplayName == "" {
			return u.Name
		}

		return u.RedisData.DisplayName
	}

	return u.IRCDisplayName
}

// NameNoPing xD
func (u *User) NameNoPing() string {
	displayName := u.DisplayName()
	return string(displayName[0]) + noPing + displayName[1:]
}

// LoadData loads the redis user associated with this IRC user
// Name is required to be set
func (u *User) LoadData(pool *redis.Pool, channel string) {
	LoadRedisUser(pool, channel, u.Name, &u.RedisData)

	if u.IRCDisplayName != "" {
		u.RedisData.DisplayName = u.IRCDisplayName
	}
}

// LoadDataConn loads the redis user associated with this IRC user
// Name is required to be set
func (u *User) LoadDataConn(conn redis.Conn, channel string) {
	LoadRedisUserConn(conn, channel, u.Name, &u.RedisData)
}

// GetPoints returns a point amount relative to the user with the given arg.
// Available args:
// "all" or "allin": Return all of users points
// "50%": Return 50% of users points
// "13": Return 13 if the user has 13 points, otherwise return error
func (u *User) GetPoints(arg string) (int, error) {
	if arg == "all" || arg == "allin" {
		// Return all points
		return u.RedisData.Points, nil
	}

	var bet int

	if strings.HasSuffix(arg, "%") {
		_bet, err := strconv.ParseFloat(arg[0:len(arg)-1], 64)
		if err != nil {
			log.Error(err)
			return 0, fmt.Errorf(`Invalid argument to GetPoints (using %%)`)
		}
		fpoints := float64(u.RedisData.Points)
		bet = int(fpoints * (_bet / 100))
	} else {
		_bet, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("Invalid argument to GetPoints")
		}
		bet = int(_bet)
	}

	/*
		if bet > u.RedisData.Points {
			return 0, fmt.Errorf("User does not have this many points")
		}
	*/

	return bet, nil
}

// IncrLines increases lines for the user KKona
func (u *User) IncrLines(channelOnline bool) {
	if channelOnline {
		u.RedisData.OnlineMessageCount++
	} else {
		u.RedisData.OfflineMessageCount++
	}
	u.RedisData.TotalMessageCount++
}
