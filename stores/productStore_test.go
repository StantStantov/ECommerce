package stores_test

import (
	"Stant/ECommerce/domain"
	"Stant/ECommerce/stores"
	"database/sql"
	"os"
	"slices"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestProductStore(t *testing.T) {
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}

	store := stores.NewProductStore(db)

	t.Run("Test Read", func(t *testing.T) {
		t.Parallel()
		testRead(t, store)
	})
	t.Run("Test ReadAll", func(t *testing.T) {
		t.Parallel()
		testReadAll(t, store)
	})
	t.Run("Test ReadAllByFilter", func(t *testing.T) {
		t.Parallel()
		testReadAllByFilter(t, store)
	})
}

func testRead(t *testing.T, store domain.ProductStore) {
	t.Helper()

	got, err := store.Read(1)
	if err != nil {
		t.Error(err)
	}

	want := domain.NewProduct(1, "ThinkPad", "Lenovo", "Laptops", 10000)

	checkProduct(t, got, want)
}

func testReadAll(t *testing.T, store domain.ProductStore) {
	t.Helper()

	got, err := store.ReadAll()
	if err != nil {
		t.Error(err)
	}

	want := []domain.Product{
		domain.NewProduct(1, "ThinkPad", "Lenovo", "Laptops", 10000),
	}

	checkProducts(t, got, want)
}

func testReadAllByFilter(t *testing.T, store domain.ProductStore) {
	t.Helper()

	got, err := store.ReadAllByFilter(1)
	if err != nil {
		t.Error(err)
	}

	want := []domain.Product{
		domain.NewProduct(1, "ThinkPad", "Lenovo", "Laptops", 10000),
	}

	checkProducts(t, got, want)
}

func checkProducts(t *testing.T, got, want []domain.Product) {
	t.Helper()

	if !slices.EqualFunc(got, want, isEqualProducts) {
		t.Errorf("Incorrect Products:\ngot: \n%v\nwant: \n%v\n", got, want)
	}
}

func isEqualProducts(E1, E2 domain.Product) bool {
	if E1.ID() != E2.ID() {
		return false
	}
	if E1.Name() != E2.Name() {
		return false
	}
	if E1.Seller() != E2.Seller() {
		return false
	}
	if E1.Category() != E2.Category() {
		return false
	}
	if E1.Price() != E2.Price() {
		return false
	}
	return true
}

func checkProduct(t *testing.T, got, want domain.Product) {
	t.Helper()

	if got.ID() != want.ID() {
		t.Errorf("Incorrect Product ID: got %v, want %v", got.ID(), want.ID())
	}
	if got.Name() != want.Name() {
		t.Errorf("Incorrect Product Name: got %v, want %v", got.Name(), want.Name())
	}
	if got.Seller() != want.Seller() {
		t.Errorf("Incorrect Product Seller: got %v, want %v", got.Seller(), want.Seller())
	}
	if got.Category() != want.Category() {
		t.Errorf("Incorrect Product Category: got %v, want %v", got.Category(), want.Category())
	}
	if got.Price() != want.Price() {
		t.Errorf("Incorrect Product Price: got %v, want %v", got.Price(), want.Price())
	}
}
