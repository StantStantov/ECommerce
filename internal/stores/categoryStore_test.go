package stores_test

import (
	"Stant/ECommerce/internal/domain"
	"Stant/ECommerce/internal/stores"
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

	id := "c735f60a-bebf-4d2f-a016-190a883eb99f"
	got, err := store.Read(id)
	if err != nil {
		t.Error(err)
	}

	want := domain.NewCategory(id, "Head")

	checkCategory(t, got, want)
}

func testCategoryReadAll(t *testing.T, store domain.CategoryStore) {
	t.Helper()

	got, err := store.ReadAll()
	if err != nil {
		t.Error(err)
	}

	want := []domain.Category{
		domain.NewCategory("c735f60a-bebf-4d2f-a016-190a883eb99f", "Head"),
		domain.NewCategory("7670dd24-fffd-4ede-8fad-17613ec6f2ba", "Core"),
		domain.NewCategory("70b0d225-f526-4c8b-aafd-cdea3f2977d2", "Arms"),
		domain.NewCategory("10021a86-d948-4c54-bdf2-00df93a22add", "Legs"),
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
