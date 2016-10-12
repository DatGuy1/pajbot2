package common

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// User xD
type User struct {
	ID                  int
	Name                string
	DisplayName         string
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
	LastSeen            time.Time // should this be time.Time or int/float?
	LastActive          time.Time
}

const noPing = string("\u05C4")

// NameNoPing xD
func (u *User) NameNoPing() string {
	return string(u.DisplayName[0]) + noPing + u.DisplayName[1:]
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
