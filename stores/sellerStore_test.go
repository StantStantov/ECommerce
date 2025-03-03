package stores_test

import (
	"Stant/ECommerce/domain"
	"Stant/ECommerce/stores"
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

func testSellerRead(t *testing.T, store domain.SellerStore) {
	t.Helper()

	got, err := store.Read(1)
	if err != nil {
		t.Error(err)
	}

	want := domain.NewSeller(1, "HUAWEI")

	checkSeller(t, got, want)
}

func testSellerReadAll(t *testing.T, store domain.SellerStore) {
	t.Helper()

	got, err := store.ReadAll()
	if err != nil {
		t.Error(err)
	}

	want := []domain.Seller{
		domain.NewSeller(1, "HUAWEI"),
		domain.NewSeller(2, "Lenovo"),
		domain.NewSeller(3, "Apple"),
	}

	if !slices.EqualFunc(got, want, isEqualSellers) {
		t.Errorf("Incorrect Sellers:\ngot: \n%v\nwant: \n%v\n", got, want)
	}
}

func isEqualSellers(E1, E2 domain.Seller) bool {
	if E1.ID() != E2.ID() {
		return false
	}
	if E1.Name() != E2.Name() {
		return false
	}
	return true
}

func checkSeller(t *testing.T, got, want domain.Seller) {
	t.Helper()

	if got.ID() != want.ID() {
		t.Errorf("Incorrect Seller ID: got %v, want %v", got.ID(), want.ID())
	}
	if got.Name() != want.Name() {
		t.Errorf("Incorrect Seller Name: got %v, want %v", got.Name(), want.Name())
	}
}
