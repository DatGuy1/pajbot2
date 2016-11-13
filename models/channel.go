package models

import (
	"database/sql"
	"time"

	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/sqlmanager"
)

const channelQ = "SELECT id, name, nickname, enabled, bot_id FROM pb_channel"

// Channel contains data about the channel
type Channel struct {
	// ID in the database
	ID int

	// Name of the channel, i.e. forsenlol
	Name string

	// Nickname of the channel, i.e. Forsen (could be used as an alias for website)
	Nickname string

	// Enabled decides whether we join this channel or not
	Enabled bool

	// Channel ID (fetched from the twitch API)
	TwitchChannelID int64

	// XXX: this should probably we renamed to BotAcountID instead of naming it BotID or bot_id everywhere
	BotID int

	Emotes common.ExtensionEmotes

	online     bool
	start      time.Time
	end        *time.Time
	updateTime time.Time
}

func (c *Channel) updateStatus() {
	const statusCacheTime = -20 * time.Second
	if time.Now().Add(statusCacheTime).After(c.updateTime) {
		conn := db.Pool.Get()
		defer conn.Close()

		stream, err := GetLastStream(c.Name)
		if err != nil {
			return
		}

		c.updateTime = time.Now()

		c.online = stream.End == nil
		c.start = stream.Start
		c.end = stream.End
	}
}

// Online return online status of channel
func (c *Channel) Online() bool {
	c.updateStatus()

	return c.online
}

// Start return time when the stream started
func (c *Channel) Start() time.Time {
	c.updateStatus()

	return c.start
}

// End return time when the stream started
func (c *Channel) End() *time.Time {
	c.updateStatus()

	return c.end
}

// Uptime returns a time.Time object for how long the stream has been online/offline
func (c *Channel) Uptime() time.Duration {
	if c.Online() {
		return time.Since(c.Start())
	}

	end := c.End()
	if end != nil {
		return time.Since(*end)
	}

	return 5 * time.Second
}

// ChannelSQLWrapper contains data about the channel that's stored in MySQL
type ChannelSQLWrapper struct {
	ID              int
	Name            string
	Nickname        sql.NullString
	Enabled         int
	TwitchChannelID sql.NullInt64
	BotID           int
}

/*
FetchFromWrapper is the linker function that links the values of
the ChannelSQLWrapper and the Channel struct together.
This needs to be kept up to date when the table structure of pb_channel
is changed.
*/
func (c *Channel) FetchFromWrapper(w ChannelSQLWrapper) {
	c.ID = w.ID
	c.Name = w.Name
	if w.Nickname.Valid {
		c.Nickname = w.Nickname.String
	}
	if w.Enabled != 0 {
		c.Enabled = true
	}
	if w.TwitchChannelID.Valid {
		c.TwitchChannelID = w.TwitchChannelID.Int64
	}
	c.BotID = w.BotID
}

// Scannable means sql.Rows or sql.Row
type Scannable interface {
	Scan(dest ...interface{}) error
}

// FetchFromSQL populates the given object with data from SQL based on the
// given argument
func (c *Channel) FetchFromSQL(row Scannable) error {
	w := ChannelSQLWrapper{}

	err := row.Scan(&w.ID, &w.Name, &w.Nickname, &w.Enabled, &w.BotID)

	if err != nil {
		log.Error(err)
		return err
	}

	c.FetchFromWrapper(w)

	return nil
}

// InsertNewToSQL inserts the given channel to SQL
func (c *Channel) InsertNewToSQL(sql *sqlmanager.SQLManager) error {
	const queryF = `INSERT INTO pb_channel (name, bot_id) VALUES (?, ?)`

	stmt, err := sql.Session.Prepare(queryF)
	if err != nil {
		// XXX
		log.Error(err)
		return err
	}

	_, err = stmt.Exec(c.Name, c.BotID)

	if err != nil {
		// XXX
		log.Error(err)
		return err
	}
	return nil
}

// SQLSetEnabled updates the enabled state of the given channel
func (c *Channel) SQLSetEnabled(sql *sqlmanager.SQLManager, enabled int) error {
	const queryF = `UPDATE pb_channel SET enabled=? WHERE id=?`

	stmt, err := sql.Session.Prepare(queryF)
	if err != nil {
		// XXX
		log.Fatal(err)
		return err
	}

	_, err = stmt.Exec(enabled, c.ID)

	if err != nil {
		// XXX
		log.Fatal(err)
		return err
	}
	return nil
}

// SQLSetBotID updates the enabled state of the given channel
func (c *Channel) SQLSetBotID(sql *sqlmanager.SQLManager, botID int) error {
	const queryF = `UPDATE pb_channel SET bot_id=? WHERE id=?`

	stmt, err := sql.Session.Prepare(queryF)
	if err != nil {
		// XXX
		log.Fatal(err)
		return err
	}

	_, err = stmt.Exec(botID, c.ID)

	if err != nil {
		// XXX
		log.Fatal(err)
		return err
	}
	return nil
}

func getChannelsByQuery(session *sql.DB, query string, dest ...interface{}) ([]Channel, error) {
	stmt, err := session.Prepare(query)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(dest...)
	if err != nil {
		return nil, err
	}

	var channels []Channel

	for rows.Next() {
		c := Channel{}
		err = c.FetchFromSQL(rows)
		if err != nil {
			log.Error(err)
		} else {
			// The channel was fetched properly
			channels = append(channels, c)
		}
	}

	return channels, nil

}

// GetChannelsByName xD
func GetChannelsByName(session *sql.DB, name string) ([]Channel, error) {
	const queryF = channelQ + " WHERE name=?"

	return getChannelsByQuery(session, queryF, name)
}

// GetChannelsByBotID xD
func GetChannelsByBotID(session *sql.DB, botID int) ([]Channel, error) {
	const queryF = channelQ + " WHERE bot_id=?"

	return getChannelsByQuery(session, queryF, botID)
}
