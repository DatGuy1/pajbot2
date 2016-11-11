package redismanager

import (
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/common/config"
	"github.com/pajlada/pajbot2/helper"
)

// RedisManager keeps the pool of redis connections
type RedisManager struct {
	Pool *redis.Pool
}

// Init connects to redis and returns redis client
func Init(config *config.Config) *RedisManager {
	r := &RedisManager{}
	pool := redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", config.RedisHost)
		if err != nil {
			return nil, err
		}
		if config.RedisDatabase >= 0 {
			_, err = c.Do("SELECT", config.RedisDatabase)
			if err != nil {
				return nil, err
			}
		}
		return c, err
	}, 69)
	// forsenGASM
	r.Pool = pool
	return r
}

// UpdateGlobalUser sets global values for a user (in other words, values that transcend channels)
// Only globally banned users and admins have a level in global redis
func (r *RedisManager) UpdateGlobalUser(channel string, user *common.User, u *common.GlobalUser) {
	/*
		conn := r.Pool.Get()
		conn.Send("HSET", "global:last_active", user.Name, time.Now().Unix())

		// Don't update the channel if the channel is empty (i.e. a whisper)
		if channel != "" {
			conn.Send("HSET", "global:channel", user.Name, channel)
		}
		conn.Flush()
		conn.Close()
	*/
}

// GetGlobalUser fills in the user and u objects with values from the users global values
func (r *RedisManager) GetGlobalUser(channel string, user *common.User, u *common.GlobalUser) {
	/*
		conn := r.Pool.Get()
		defer conn.Close()
		exist, err := conn.Do("HEXISTS", "global:last_active", user.Name)
		e, _ := redis.Bool(exist, err)
		if e {
			conn.Send("HGET", "global:level", user.Name)
			conn.Send("HGET", "global:channel", user.Name)
			conn.Flush()
			// can this be done in a loop somehow?
			// Level
			res, err := conn.Receive()
			level, _ := redis.Int(res, err) // will be 0 unless user is admin or globally banned
			if level > user.Level {
				// XXX: Should this set u.Level instead? not sure!
				user.Level = level
			}
			// Channel
			res, err = conn.Receive()
			u.Channel, _ = redis.String(res, err)
		} else {
			r.UpdateGlobalUser(u.Channel, user, u)
			// r.GetGlobalUser(channel, user, u)
		}
		r.UpdateGlobalUser(channel, user, u)
	*/
}

// IsValidUser checks if the user is in the database
func (r *RedisManager) IsValidUser(channel string, _user string) bool {
	conn := r.Pool.Get()
	defer conn.Close()
	user := strings.ToLower(_user)
	res, err := redis.Bool(conn.Do("HEXISTS", channel+":users:last_seen", user))
	if err != nil {
		log.Error(err)
	}
	return res
}

func (r *RedisManager) newUser(channel string, user *common.User) {
	conn := r.Pool.Get()
	defer conn.Close()
	conn.Send("HSET", channel+":users:last_seen", user.Name, time.Now().Unix())

	// Why is this called?
	conn.Send("HSET", channel+":users:level", user.Name, r.getLevel(createLevel(0, 1), user))

	conn.Flush()
}

// SetLevel sets the users level in the given channel
func (r *RedisManager) SetLevel(channel string, user *common.User, level int) {
	conn := r.Pool.Get()
	defer conn.Close()
	conn.Send("HSET", channel+":users:level", user.Name, createLevel(uint32(level), 0)) // XXX: Make sure the flags are right
	conn.Flush()
}

// ResetLevel resets a users level to 0/default XXX
func (r *RedisManager) ResetLevel(channel string, user *common.User) {
	conn := r.Pool.Get()
	defer conn.Close()
	conn.Send("HSET", channel+":users:level", user.Name, r.getLevel(createLevel(0, 1), user))
	conn.Flush()
}

// Flags for user level values
const (
	// Specifies whether the level was set automatically or not.
	// If the level was set automatically, that means it will always be
	// re-evaluated
	LevelFlagAutomatic = 1 << iota
)

/*
First 32 bits specifies the level Value
Last 32 bits specifies the level Flags
*/
func (r *RedisManager) getLevel(levelCombined uint64, user *common.User) int {
	// Split up the values from levelCombined
	levelFlags, levelValue := helper.SplitUint64(levelCombined)

	// Returns the users Global Level if its higher than the Channel Level
	if user.Level > int(levelValue) {
		return user.Level
	}

	// If the level was set manually, return that level
	if !helper.CheckFlag(levelFlags, LevelFlagAutomatic) {
		return int(levelValue)
	}

	// Otherwise, return the automatic level that's appropriate for the user
	if user.ChannelOwner {
		return 1500
	} else if user.Mod {
		return 500
	} else if user.Sub {
		return 250
	}

	// Normal user
	return 100
}

func createLevel(levelValue uint32, levelFlags uint32) uint64 {
	// XXX: Make sure this is correct
	return helper.CombineUint32(levelFlags, levelValue)
}
