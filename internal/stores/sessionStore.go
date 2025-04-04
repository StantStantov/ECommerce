package stores

import (
	"Stant/ECommerce/internal/domain"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type SessionStore struct {
	db              *sql.DB
	sessionLifetime time.Time
}

func NewSessionStore(db *sql.DB, sessionLifetime time.Time) *SessionStore {
	return &SessionStore{db: db, sessionLifetime: sessionLifetime}
}

const createSession = `
  INSERT INTO market.sessions
  (user_id, session_token, csrf_token, expire_on)
  VALUES
  ($1, $2, $3, $4)
  ;
`

func (s SessionStore) Create(userID, sessionToken, csrfToken string) error {
	if _, err := s.db.Exec(createSession, userID, sessionToken, csrfToken, s.sessionLifetime); err != nil {
		return fmt.Errorf("stores.SessionStore.Create: [%w]", err)
	}
	return nil
}

const readSession = `
  SELECT * FROM market.sessions
  WHERE session_token = $1
  LIMIT 1
  ;
`

func (s SessionStore) Read(sessionToken string) (domain.Session, error) {
	row := s.db.QueryRow(readSession, sessionToken)
	session, err := scanSession(row)
	if err != nil {
		return session, fmt.Errorf("stores.SessionStore.Read: [%w]", err)
	}
	return session, nil
}

const deleteSessionByToken = `
  DELETE FROM market.sessions
  WHERE session_token = $1 
  ;
`

func (s SessionStore) Delete(sessionToken string) error {
	if _, err := s.db.Exec(deleteSessionByToken, sessionToken); err != nil {
		return fmt.Errorf("stores.SessionStore.Delete: [%w]", err)
	}
	return nil
}

const deleteExpiredSessions = `
  DELETE FROM market.sessions
  WHERE expire_on < now() 
  ;
`

func (s SessionStore) DeleteAllExpired() error {
	if _, err := s.db.Exec(deleteExpiredSessions); err != nil {
		return fmt.Errorf("stores.SessionStore.DeleteAllExpired: [%w]", err)
	}
	return nil
}

func scanSession(row sqlRow) (domain.Session, error) {
	var (
		userID       string
		sessionToken string
		csrfToken    string
		expireOn     time.Time
	)
	if err := row.Scan(&userID, &sessionToken, &csrfToken, &expireOn); err != nil {
		return domain.Session{}, fmt.Errorf("stores.scanSession: [%w]", err)
	}
	return domain.NewSession(userID, sessionToken, csrfToken, expireOn), nil
}

func (s SessionStore) StartCleanup(
	logger log.Logger,
	interval time.Duration,
) (chan<- struct{}, <-chan struct{}) {
	quit, done := make(chan struct{}), make(chan struct{})
	go s.clean(logger, interval, quit, done)
	return quit, done
}

func (s SessionStore) StopCleanup(quit chan<- struct{}, done <-chan struct{}) {
	quit <- struct{}{}
	<-done
}

func (s SessionStore) clean(
	logger log.Logger,
	interval time.Duration,
	quit <-chan struct{},
	done chan<- struct{},
) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-quit:
			done <- struct{}{}
			return
		case <-ticker.C:
			if err := s.DeleteAllExpired(); err != nil {
				logger.Printf("Error = %v", fmt.Errorf("stores.SessionStore.clean: [%w]", err))
			}
		}
	}
}
