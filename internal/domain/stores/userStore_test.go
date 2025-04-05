package stores_test

import (
	"Stant/ECommerce/internal/domain/models"
	"Stant/ECommerce/internal/domain/stores"
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
}

func testUserCreate(t *testing.T, store models.UserStore) {
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

func testUserIsExists(t *testing.T, store models.UserStore) {
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

func testUserRead(t *testing.T, store models.UserStore) {
	t.Helper()

	id := "ad43dfbf-1152-478c-a595-e3ebe5ad0085"
	got, err := store.Read(id)
	if err != nil {
		t.Fatalf("Failed to confirm existence of user in DB: [%v]", err)
	}

	email := "readME@test.com"
	firstName := "read"
	secondName := "ME"
	hashedPassword := "$2a$10$sgEy3LehHNpbZ7NjqDhMiejJ8gaQTcykfv1VFJL42aPN8pZJL45EW"
	want := models.NewUser(id, email, firstName, secondName, hashedPassword)

	checkUser(t, got, want)
}

func checkUser(t *testing.T, got, want models.User) {
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
