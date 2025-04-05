package models

import "time"

type Session struct {
	userID       string
	sessionToken string
	csrfToken    string
	expireOn     time.Time
}

func NewSession(userID string, sessionToken, csrfToken string, expireOn time.Time) Session {
	return Session{
		userID:       userID,
		sessionToken: sessionToken,
		csrfToken:    csrfToken,
		expireOn:     expireOn,
	}
}

func (s Session) UserID() string {
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
