package domain

import "time"

type Session struct {
	sessionToken string
	csrfToken    string
	expireOn     time.Time
}

func NewSession(sessionToken, csrfToken string, expireOn time.Time) Session {
	return Session{
		sessionToken: sessionToken,
		csrfToken:    csrfToken,
		expireOn:     expireOn,
	}
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
