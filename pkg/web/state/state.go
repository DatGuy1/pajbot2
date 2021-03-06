package state

import (
	"database/sql"
	"fmt"
	"net/http"
	"sync"

	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/utils"
)

var (
	sqlClient       *sql.DB
	twitchUserStore pkg.UserStore
	pubSub          pkg.PubSub

	mutex = &sync.RWMutex{}

	sessionStore = &SessionStore{}
)

func StoreSQL(sql_ *sql.DB) {
	mutex.Lock()
	sqlClient = sql_
	mutex.Unlock()
}

func StoreTwitchUserStore(twitchUserStore_ pkg.UserStore) {
	mutex.Lock()
	twitchUserStore = twitchUserStore_
	mutex.Unlock()
}

func StorePubSub(pubSub_ pkg.PubSub) {
	mutex.Lock()
	pubSub = pubSub_
	mutex.Unlock()
}

type Session struct {
	ID     string
	UserID uint64

	TwitchUserID   string
	TwitchUserName string
}

type State struct {
	SQL             *sql.DB
	TwitchUserStore pkg.UserStore
	PubSub          pkg.PubSub
	Session         *Session
	SessionID       *string
}

func (s *State) CreateSession(userID int64) (sessionID string, err error) {
	// TODO: Switch this out for a proper randomizer or something KKona
	sessionID, err = utils.GenerateRandomString(64)
	if err != nil {
		return
	}

	const queryF = `
INSERT INTO
	UserSession
(id, user_id)
	VALUES (?, ?)`

	// TODO: Make sure the exec didn't error
	_, err = sqlClient.Exec(queryF, sessionID, userID)
	if err != nil {
		return
	}

	return
}

func Context(w http.ResponseWriter, r *http.Request) State {
	mutex.RLock()
	state := State{
		SQL:             sqlClient,
		TwitchUserStore: twitchUserStore,
		PubSub:          pubSub,
	}
	mutex.RUnlock()

	// Authorization via header api key (not implemented yet)
	credentials := r.Header.Get("Authorization")
	if credentials != "" {
		fmt.Println("Credentials:", credentials)
	}

	state.SessionID = getCookie(r, SessionIDCookie)
	if state.SessionID != nil && *state.SessionID != "" {
		state.Session = sessionStore.Get(state.SQL, *state.SessionID)
	}

	return state
}
