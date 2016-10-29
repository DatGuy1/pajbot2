package common

import (
	"encoding/json"
	"time"

	"github.com/garyburd/redigo/redis"
)

type redisUser struct {
	Name        string    `json:"Name"`
	DisplayName string    `json:"DisplayName,omitempty"`
	Mod         bool      `json:"Mod,omitempty"`
	Sub         bool      `json:"Sub,omitempty"`
	FixedLevel  int       `json:"FixedLevel,omitempty"`
	LastSeen    time.Time `json:"LastSeen,omitempty"`
	LastActive  time.Time `json:"LastActive,omitempty"`
}

// LoadRedisUser takes a username as argument and returns a filled redisUser object
func LoadRedisUser(pool *redis.Pool, channel string, name string, u *redisUser) error {
	conn := pool.Get()
	defer conn.Close()
	jsonString, err := redis.Bytes(conn.Do("HGET", channel+":users"))
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonString, u)
	if err != nil {
		return err
	}

	return nil
}

// Save saves the redisUser object into the database
func (u *redisUser) Save(pool *redis.Pool, channel string) error {
	conn := pool.Get()
	defer conn.Close()

	bytes, err := json.Marshal(u)
	if err != nil {
		return err
	}

	_, err = conn.Do("HSET", bytes)
	if err != nil {
		return err
	}

	// Saved successfully
	return nil
}
