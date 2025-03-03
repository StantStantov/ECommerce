package stores_test

import (
	"Stant/ECommerce/domain"
	"Stant/ECommerce/stores"
	"slices"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestCategoryStore(t *testing.T) {
	db, err := stores.NewDBConn()
	if err != nil {
		t.Fatal(err)
	}

	store := stores.NewCategoryStore(db)

	t.Run("Test Read", func(t *testing.T) {
		t.Parallel()
		testCategoryRead(t, store)
	})
	t.Run("Test ReadAll", func(t *testing.T) {
		t.Parallel()
		testCategoryReadAll(t, store)
	})
}

func testCategoryRead(t *testing.T, store domain.CategoryStore) {
	t.Helper()

	got, err := store.Read(1)
	if err != nil {
		t.Error(err)
	}

	want := domain.NewCategory(1, "Laptops")

	checkCategory(t, got, want)
}

func testCategoryReadAll(t *testing.T, store domain.CategoryStore) {
	t.Helper()

	got, err := store.ReadAll()
	if err != nil {
		t.Error(err)
	}

	want := []domain.Category{
		domain.NewCategory(1, "Laptops"),
		domain.NewCategory(2, "Phones"),
		domain.NewCategory(3, "Electronics"),
	}

	if !slices.EqualFunc(got, want, isEqualCategories) {
		t.Errorf("Incorrect Categories:\ngot: \n%v\nwant: \n%v\n", got, want)
	}
}

func isEqualCategories(E1, E2 domain.Category) bool {
	if E1.ID() != E2.ID() {
		return false
	}
	if E1.Name() != E2.Name() {
		return false
	}
	return true
}

func checkCategory(t *testing.T, got, want domain.Category) {
	t.Helper()

	if got.ID() != want.ID() {
		t.Errorf("Incorrect Category ID: got %v, want %v", got.ID(), want.ID())
	}
	if got.Name() != want.Name() {
		t.Errorf("Incorrect Category Name: got %v, want %v", got.Name(), want.Name())
	}
}
