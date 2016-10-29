package common

// this should make things easier with redis

import (
	"time"

	"github.com/pajlada/pajbot2/plog"
)

var log = plog.GetLogger()

// BuildTime is the time when the binary was built
// filled in with ./build.sh (ldflags)
var BuildTime string

// GlobalUser will only be used by boss to check if user is admin
// and to decide what channel to send the message to if its a whisper
type GlobalUser struct {
	LastActive time.Time
	Channel    string
	Level      int
}

// MsgType specifies the message's type, for example PRIVMSG or WHISPER
type MsgType uint32

// Various message types which describe what sort of message they are
const (
	MsgPrivmsg MsgType = iota + 1
	MsgWhisper
	MsgSub
	MsgThrowAway
	MsgUnknown
	MsgUserNotice
	MsgReSub
	MsgNotice
	MsgRoomState
	MsgSubsOn
	MsgSubsOff
	MsgSlowOn
	MsgSlowOff
	MsgR9kOn
	MsgR9kOff
	MsgHostOn
	MsgHostOff
	MsgTimeoutSuccess
)

/*
Msg contains all the information about an IRC message.
This included already-parsed ircv3-tags and the User object
*/
type Msg struct {
	// User who sent the message, or user who the message is about
	User User

	// Text contains the text of the message, i.e. the text someone sent in chat or the resub message
	Text string

	// Channel the message was sent in
	Channel string

	// Type of message (PRIVMSG, WHISPER, SUB, RESUB?)
	Type MsgType

	// If the message is a /me message
	Me bool

	// List of emotes contained in the message
	// XXX(pajlada): Not sure if they're sorted
	Emotes []Emote

	// All remaining unparsed ircv3-tags
	Tags map[string]string

	// Used to send along arguments about the user (what direction?)
	// used in bot.Format
	Args []string
}
