package modules

import (
	"encoding/json"
	"reflect"

	"github.com/fatih/structs"
	"github.com/garyburd/redigo/redis"
	"github.com/pajlada/pajbot2/common/basemodule"
	"github.com/pajlada/pajbot2/redismanager"
)

type baseSetting struct {
	Valid bool
}

// IntSetting xD
type IntSetting struct {
	baseSetting

	Label string
	Value int

	// Constraints
	MinValue int
	MaxValue int
}

// Int returns the int value of an int setting
// If the Value is nil then we will just return the DefaultValue int
func (s *IntSetting) Int() int {
	return s.Value
}

// StringSetting xD
type StringSetting struct {
	baseSetting

	// Label is used to describe the setting
	Label string

	Value string

	// Constraints
	MinLength int
	MaxLength int
}

func setString(m interface{}, field string, value string) {
	v := reflect.ValueOf(m).Elem().FieldByName(field).FieldByName("Value")
	if v.IsValid() {
		v.SetString(value)
	}
}

func setInt(m interface{}, field string, value int) {
	v := reflect.ValueOf(m).Elem().FieldByName(field).FieldByName("Value")
	if v.IsValid() {
		v.SetInt(int64(value))
	}
}
func setBool(m interface{}, field string, value bool) {
	v := reflect.ValueOf(m).Elem().FieldByName(field).FieldByName("Value")
	if v.IsValid() {
		v.SetBool(value)
	}
}

// String returns the string value of a string setting
// If the Value is nil then we will just return the DefaultValue string
func (s *StringSetting) String() string {
	return s.Value
}

// BoolSetting xD
type BoolSetting struct {
	baseSetting

	// Label is used to describe the setting
	Label string

	Value bool
}

// Bool returns the string value of a string setting
// If the Value is nil then we will just return the DefaultValue string
func (s *BoolSetting) Bool() bool {
	return s.Value
}

// OptionSetting xD
type OptionSetting struct {
	baseSetting

	// Label is used to describe the setting
	Label string

	Value   int
	Options []string
}

// Int returns the index of the currently selected index
// If the Index is nil then we will just return the DefaultValue string
func (s *OptionSetting) Int() int {
	return s.Value
}

// String returns the index of the currently selected index
// If the Index is nil then we will just return the DefaultValue string
func (s *OptionSetting) String() string {
	index := s.Int()

	// XXX(pajlada): Make bounds check here
	return s.Options[index]
}

func loadJSONData(r *redismanager.RedisManager, channelName string, id string) map[string]interface{} {
	conn := r.Pool.Get()
	defer conn.Close()
	ret := make(map[string]interface{})
	data, err := redis.Bytes(conn.Do("HGET", channelName+":modules:settings", id))
	if err != nil {
		return nil
	}

	err = json.Unmarshal(data, &ret)
	if err != nil {
		log.Error(err)
		return nil
	}

	return ret
}

// LoadSettings xD
func LoadSettings(r *redismanager.RedisManager, streamer string, m interface{}) {
	var settings map[string]interface{}

	for _, f := range structs.Fields(m) {
		if !f.IsExported() {
			// Ignore unexported objects
			continue
		}

		name := f.Name()
		val := f.Value()
		switch val.(type) {
		case basemodule.BaseModule:
			// Load JSON settings from redis into "settings" variable
			settings = loadJSONData(r, streamer, val.(basemodule.BaseModule).ID)

		case IntSetting:
			if data, ok := settings[name]; ok {
				setInt(m, name, data.(int))
			}

		case StringSetting:
			if data, ok := settings[name]; ok {
				setString(m, name, data.(string))
			}

		case OptionSetting:
			if data, ok := settings[name]; ok {
				setInt(m, name, data.(int))
			}

		case BoolSetting:
			if data, ok := settings[name]; ok {
				setBool(m, name, data.(bool))
			}
		}
	}
}
