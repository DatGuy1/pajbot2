package common

import (
	"encoding/json"
	"time"

	"github.com/garyburd/redigo/redis"
)

// RedisUser :P
type RedisUser struct {
	Name        string    `json:"Name"`
	DisplayName string    `json:"DisplayName,omitempty"`
	Mod         bool      `json:"Mod,omitempty"`
	Sub         bool      `json:"Sub,omitempty"`
	FixedLevel  int       `json:"FixedLevel,omitempty"`
	LastSeen    time.Time `json:"LastSeen,omitempty"`
	LastActive  time.Time `json:"LastActive,omitempty"`

	// Extended
	// TotalMessageCount (STREAMER:users:total_message_count) ZSET
	TotalMessageCount int `json:"-"`
	// OnlineMessageCount (STREAMER:users:online_message_count) ZSET
	OnlineMessageCount int `json:"-"`
	// OfflineMessageCount (STREAMER:users:offline_message_count) ZSET
	OfflineMessageCount int `json:"-"`
	// Points (STREAMER:users:points) ZSET
	Points int `json:"-"`

	// Set to true if base object was not set
	New bool `json:"-"`

	pool    *redis.Pool
	channel string
}

// LoadRedisUserConn XD
func LoadRedisUserConn(conn redis.Conn, channel string, name string, u *RedisUser) error {
	u.channel = channel
	u.Name = name

	// Load base JSON data
	jsonString, err := redis.Bytes(conn.Do("HGET", channel+":users", u.Name))
	if err != nil {
		return err
	}

	// Load extended data
	conn.Send("ZSCORE", channel+":users:points", u.Name)
	conn.Send("ZSCORE", channel+":users:total_message_count", u.Name)
	conn.Send("ZSCORE", channel+":users:online_message_count", u.Name)
	conn.Send("ZSCORE", channel+":users:offline_message_count", u.Name)

	conn.Flush()

	u.Points, _ = redis.Int(conn.Receive())
	u.TotalMessageCount, _ = redis.Int(conn.Receive())
	u.OnlineMessageCount, _ = redis.Int(conn.Receive())
	u.OfflineMessageCount, _ = redis.Int(conn.Receive())

	err = json.Unmarshal(jsonString, u)
	if err != nil {
		return err
	}

	// Set user level to 1500 if this is the broadcaster
	if u.Name == u.channel && u.FixedLevel < 1500 {
		u.FixedLevel = 1500
	}

	return nil
}

// LoadRedisUser takes a username as argument and returns a filled RedisUser object
func LoadRedisUser(pool *redis.Pool, channel string, name string, u *RedisUser) error {
	u.pool = pool
	conn := pool.Get()
	defer conn.Close()

	return LoadRedisUserConn(conn, channel, name, u)
}

// Close xD
func (u *RedisUser) Close(oldUser *RedisUser) {
	err := u.save(oldUser)
	if err != nil {
		log.Error(err)
	}
}

// CloseWithoutSaving xD
func (u *RedisUser) CloseWithoutSaving() {
	// do nothing
}

// save saves the RedisUser object into the database
func (u *RedisUser) save(oldUser *RedisUser) error {
	conn := u.pool.Get()
	defer conn.Close()

	bytes, err := json.Marshal(u)
	if err != nil {
		return err
	}

	conn.Send("HSET", u.channel+":users", u.Name, bytes)
	conn.Send("HSET", u.channel+":users:last_seen", u.Name, time.Now().Unix())
	conn.Send("HSET", u.channel+":users:last_active", u.Name, time.Now().Unix())

	// Update total message count if needed
	if u.TotalMessageCount != oldUser.TotalMessageCount {
		diff := u.TotalMessageCount - oldUser.TotalMessageCount
		conn.Send("ZINCRBY", u.channel+":users:total_message_count", diff, u.Name)
	}

	// Update online message count if needed
	if u.OnlineMessageCount != oldUser.OnlineMessageCount {
		diff := u.OnlineMessageCount - oldUser.OnlineMessageCount
		conn.Send("ZINCRBY", u.channel+":users:online_message_count", diff, u.Name)
	}

	// Update offline message count if needed
	if u.OfflineMessageCount != oldUser.OfflineMessageCount {
		diff := u.OfflineMessageCount - oldUser.OfflineMessageCount
		conn.Send("ZINCRBY", u.channel+":users:offline_message_count", diff, u.Name)
	}

	// Update points if needed
	if u.Points != oldUser.Points {
		diff := u.Points - oldUser.Points
		log.Debugf("Increase points by %d", diff)
		conn.Send("ZINCRBY", u.channel+":users:points", diff, u.Name)
	}

	conn.Flush()

	return nil
}
