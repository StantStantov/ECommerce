package stores_test

import (
	"Stant/ECommerce/internal/domain"
	"Stant/ECommerce/internal/stores"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestUserStore(t *testing.T) {
	db, err := stores.NewDBConn()
	if err != nil {
		t.Skipf("Failed to connect to DB: [%v]", err)
	}

	store := stores.NewUserStore(db)

	t.Run("Test Create", func(t *testing.T) {
		t.Parallel()
		testUserCreate(t, store)
	})
	t.Run("Test IsExists", func(t *testing.T) {
		t.Parallel()
		testUserIsExists(t, store)
	})
	t.Run("Test Read", func(t *testing.T) {
		t.Parallel()
    testUserRead(t, store)
	})
	// t.Run("Test ReadAll", func(t *testing.T) {
	// 	t.Parallel()
	// })
}

func testUserCreate(t *testing.T, store domain.UserStore) {
	t.Helper()

	email := "simple@test.com"
	firstName := "test"
	secondName := "test"
	password := "12345"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		t.Fatalf("Failed to Hash password: [%v]", err)
	}

	if err := store.Create(email, firstName, secondName, string(hash)); err != nil {
		t.Fatalf("Failed to Hash password: [%v]", err)
	}
}

func testUserIsExists(t *testing.T, store domain.UserStore) {
	t.Helper()

	email := "readME@test.com"
	exist, err := store.IsExists(email)
	if err != nil {
		t.Fatalf("Failed to confirm existence of user in DB: [%v]", err)
	}
	if !exist {
		t.Fatalf("Failed to confirm existence of user in DB: got %v, want %v", exist, true)
	}
}

func testUserRead(t *testing.T, store domain.UserStore) {
	t.Helper()

	got, err := store.Read(1)
	if err != nil {
		t.Fatalf("Failed to confirm existence of user in DB: [%v]", err)
	}

	email := "readME@test.com"
	firstName := "read"
	secondName := "ME"
	hashedPassword := "$2a$10$sgEy3LehHNpbZ7NjqDhMiejJ8gaQTcykfv1VFJL42aPN8pZJL45EW"
	want := domain.NewUser(1, email, firstName, secondName, hashedPassword)

	checkUser(t, got, want)
}

func checkUser(t *testing.T, got, want domain.User) {
	t.Helper()

	if got.Email() != want.Email() {
		t.Errorf("Incorrect Email: got %v, want %v", got.Email(), want.Email())
	}
	if got.FirstName() != want.FirstName() {
		t.Errorf("Incorrect First Name: got %v, want %v", got.FirstName(), want.FirstName())
	}
	if got.SecondName() != want.SecondName() {
		t.Errorf("Incorrect Second Name: got %v, want %v", got.SecondName(), want.SecondName())
	}
	if got.HashedPassword() != want.HashedPassword() {
		t.Errorf("Incorrect Hashed Password: got %v, want %v", got.HashedPassword(), want.HashedPassword())
	}
}
