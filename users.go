package goodidea

import (
	"time"

	"github.com/rs/xid"
)

// users within the system
//
//	id - unique ID for the user, an xid prefixed with 'u'
//	name - the username for the given user
//	salt - the salt used befored hashing the passwd
//	passwd - the user's password used for logging in
//	createdAt - the date only, day the user signed up
//	admin - if the user is an admin user or not
type user struct {
	ID        string    `json:"id" validate:"required,u20"`
	Name      string    `json:"name" validate:"required,max=48"`
	Salt      string    `json:"-" validate:"omitempty,len=8"`
	Passwd    string    `json:"passwd" validate:"required,sha256"`
	CreatedAt time.Time `json:"createdAt"`
	Admin     bool      `json:"admin"`
}

// sessions are for user logins
type session struct {
	SessionID xid.ID    `json:"sessionId"`
	UserID    string    `json:"userId" validate:"required,u20"`
	CreatedAt time.Time `json:"createdAt"`
	Valid     bool      `json:"valid"`
}
