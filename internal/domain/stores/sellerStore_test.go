package stores_test

import (
	"Stant/ECommerce/internal/domain/models"
	"Stant/ECommerce/internal/domain/stores"
	"slices"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestSellerStore(t *testing.T) {
	db, err := stores.NewDBConn()
	if err != nil {
		t.Fatal(err)
	}

	store := stores.NewSellerStore(db)

	t.Run("Test Read", func(t *testing.T) {
		t.Parallel()
		testSellerRead(t, store)
	})
	t.Run("Test ReadAll", func(t *testing.T) {
		t.Parallel()
		testSellerReadAll(t, store)
	})
}

func testSellerRead(t *testing.T, store models.SellerStore) {
	t.Helper()

	got, err := store.Read("f4d234ff-7aa5-4986-954c-8c2cc61ea0fc")
	if err != nil {
		t.Error(err)
	}

	want := models.NewSeller("f4d234ff-7aa5-4986-954c-8c2cc61ea0fc", "Balam Industries")

	checkSeller(t, got, want)
}

func testSellerReadAll(t *testing.T, store models.SellerStore) {
	t.Helper()

	got, err := store.ReadAll()
	if err != nil {
		t.Error(err)
	}

	want := []models.Seller{
		models.NewSeller("f4d234ff-7aa5-4986-954c-8c2cc61ea0fc", "Balam Industries"),
		models.NewSeller("7e13d4e2-408b-494f-a611-1950a3a36616", "Arquebus Corporation"),
	}

	if !slices.EqualFunc(got, want, isEqualSellers) {
		t.Errorf("Incorrect Sellers:\ngot: \n%v\nwant: \n%v\n", got, want)
	}
}

func isEqualSellers(E1, E2 models.Seller) bool {
	if E1.ID() != E2.ID() {
		return false
	}
	if E1.Name() != E2.Name() {
		return false
	}
	return true
}

func checkSeller(t *testing.T, got, want models.Seller) {
	t.Helper()

	if got.ID() != want.ID() {
		t.Errorf("Incorrect Seller ID: got %v, want %v", got.ID(), want.ID())
	}
	if got.Name() != want.Name() {
		t.Errorf("Incorrect Seller Name: got %v, want %v", got.Name(), want.Name())
	}
}
