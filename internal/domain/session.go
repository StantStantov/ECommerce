package domain

import "time"

type Session struct {
	userID       int32
	sessionToken string
	csrfToken    string
	expireOn     time.Time
}

func NewSession(userID int32, sessionToken, csrfToken string, expireOn time.Time) Session {
	return Session{
		userID:       userID,
		sessionToken: sessionToken,
		csrfToken:    csrfToken,
		expireOn:     expireOn,
	}
}

func (s Session) UserID() int32 {
	return s.userID
}

func (s Session) SessionToken() string {
	return s.sessionToken
}

func (s Session) CsrfToken() string {
	return s.csrfToken
}

func (s Session) ExpireOn() time.Time {
	return s.expireOn
}
