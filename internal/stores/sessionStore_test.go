package stores_test

import (
	"Stant/ECommerce/internal/domain"
	"Stant/ECommerce/internal/security"
	"Stant/ECommerce/internal/stores"
	"testing"
	"time"
)

func TestSessionStore(t *testing.T) {
	db, err := stores.NewDBConn()
	if err != nil {
		t.Skipf("Failed to connect to DB: [%v]", err)
	}

	store := stores.NewSessionStore(db, time.Now().Add(1*time.Hour))

	t.Run("Test Create", func(t *testing.T) {
		t.Parallel()
		testSessionCreate(t, store)
	})
	t.Run("Test Read", func(t *testing.T) {
		t.Parallel()
		testSessionRead(t, store)
	})
	t.Run("Test Delete", func(t *testing.T) {
		t.Parallel()
		testSessionDelete(t, store)
	})
}

func testSessionCreate(t *testing.T, store domain.SessionStore) {
	t.Helper()

	sessionCookie, err := security.NewSessionCookie()
	if err != nil {
		t.Errorf("Failed to create Session Cookie: [%v]", err)
	}
	csrfCookie, err := security.NewCsrfCookie()
	if err != nil {
		t.Errorf("Failed to create CSRF Cookie: [%v]", err)
	}

	userID := "02f95483-7934-41b7-af1e-40eaf67817fc"
	if err := store.Create(userID, sessionCookie.Value, csrfCookie.Value); err != nil {
		t.Fatalf("Failed to create session in DB: [%v]", err)
	}
}

func testSessionRead(t *testing.T, store domain.SessionStore) {
	t.Helper()

	userID := "ad43dfbf-1152-478c-a595-e3ebe5ad0085"
	sessionToken := "2-bJbG-BU5h1fKovzqoEnwOxDsz9bm1-8vVRHYav5Z29DcaDUchc0LNufSGCEjKFsXGNtn0ZF0FdcXi9_npSGg=="
	csrfToken := "DpwoY8fzNfVyBnJDl9mEclJoZcWW8kxtZIo-CdMMvGnGfwzrrqwogUyVnUZknwazD_MXxEop5ewgxp2S-wTtig=="
	expireOn := time.Date(2025, 3, 13, 14, 33, 57, 0, time.Now().Location())

	got, err := store.Read(sessionToken)
	if err != nil {
		t.Fatalf("Failed to read session from DB: [%v]", err)
	}

	want := domain.NewSession(userID, sessionToken, csrfToken, expireOn)

	checkSession(t, got, want)
}

func testSessionDelete(t *testing.T, store domain.SessionStore) {
	t.Helper()

	sessionToken := "lv1qhEGQgUn1Z3tFqKdCuvq_--W2ptuGp0oV5wRajtlC0sPN9xqsAxEZ6w2RGd-JX7nrk4_rO51tJXhoSONgmw=="

	if err := store.Delete(sessionToken); err != nil {
		t.Fatalf("Failed to delete session from DB: [%v]", err)
	}
}

func checkSession(t *testing.T, got, want domain.Session) {
	t.Helper()

	if got.SessionToken() != want.SessionToken() {
		t.Errorf("Incorrect Session Token: got %v, want %v", got.SessionToken(), want.SessionToken())
	}
	if got.CsrfToken() != want.CsrfToken() {
		t.Errorf("Incorrect CSRF Token: got %v, want %v", got.CsrfToken(), want.CsrfToken())
	}
	if got.ExpireOn() != want.ExpireOn() {
		t.Errorf("Incorrect Expiration Date: got %v, want %v", got.ExpireOn(), want.ExpireOn())
	}
}
