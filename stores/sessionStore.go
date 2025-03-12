package stores

import (
	"Stant/ECommerce/domain"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type SessionStore struct {
	db *sql.DB
}

func NewSessionStore(db *sql.DB) *SessionStore {
	return &SessionStore{db: db}
}

const createSession = `
  INSERT INTO sessions
  (session_token, csrf_token, expire_on)
  VALUES
  ($1, $2, $3)
  ;
`

func (s SessionStore) Create(sessionToken, csrfToken string, expireOn time.Time) error {
	if _, err := s.db.Exec(createSession, sessionToken, csrfToken, expireOn); err != nil {
		return fmt.Errorf("SessionStore.Create: [%w]", err)
	}
	return nil
}

const readSession = `
  SELECT * FROM sessions
  WHERE sessionToken = $1
  LIMIT 1
  ;
`

func (s SessionStore) Read(sessionToken string) (domain.Session, error) {
	row := s.db.QueryRow(readSession, sessionToken)
	session, err := scanSession(row)
	if err != nil {
		return session, fmt.Errorf("SessionStore.Read: [%w]", err)
	}
	return session, nil
}

const deleteSessionByToken = `
  DELETE FROM sessions
  WHERE sessionToken = $1 
  LIMIT 1
  ;
`

func (s SessionStore) Delete(sessionToken string) error {
	if _, err := s.db.Exec(deleteSessionByToken, sessionToken); err != nil {
		return fmt.Errorf("SessionStore.Delete: [%w]", err)
	}
	return nil
}

func scanSession(row sqlRow) (domain.Session, error) {
	var (
		sessionToken string
		csrfToken    string
		expireOn     time.Time
	)
	if err := row.Scan(sessionToken, csrfToken, expireOn); err != nil {
		return domain.Session{}, fmt.Errorf("ScanSession: [%w]", err)
	}
	return domain.NewSession(sessionToken, csrfToken, expireOn), nil
}

const deleteExpiredSessions = `
  DELETE FROM sessions
  WHERE expire_on < now() 
  ;
`

func (s SessionStore) DeleteAllExpired() error {
	if _, err := s.db.Exec(deleteExpiredSessions); err != nil {
		return fmt.Errorf("SessionStore.DeleteAllExpired: [%w]", err)
	}
	return nil
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
				logger.Printf("Error = %v", fmt.Errorf("SessionStore.clean: [%w]", err))
			}
		}
	}
}
