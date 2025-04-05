package stores_test

import (
	"Stant/ECommerce/internal/domain/models"
	"Stant/ECommerce/internal/domain/stores"
	"slices"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestProductStore(t *testing.T) {
	db, err := stores.NewDBConn()
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

func testRead(t *testing.T, store models.ProductStore) {
	t.Helper()

	id := "961dd7d0-95b3-492e-833e-c33875a64d0f"
	got, err := store.Read(id)
	if err != nil {
		t.Error(err)
	}

	want := models.NewProduct(id, "HD-033M VERRILL",
		models.NewSeller("f4d234ff-7aa5-4986-954c-8c2cc61ea0fc", "Balam Industries"),
		models.NewCategory("c735f60a-bebf-4d2f-a016-190a883eb99f", "Head"),
		205000)

	checkProduct(t, got, want)
}

func testReadAll(t *testing.T, store models.ProductStore) {
	t.Helper()

	got, err := store.ReadAll()
	if err != nil {
		t.Error(err)
	}

	want := []models.Product{
		models.NewProduct("961dd7d0-95b3-492e-833e-c33875a64d0f", "HD-033M VERRILL",
			models.NewSeller("f4d234ff-7aa5-4986-954c-8c2cc61ea0fc", "Balam Industries"),
			models.NewCategory("c735f60a-bebf-4d2f-a016-190a883eb99f", "Head"),
			205000),
		models.NewProduct("02cab72f-e225-4c3d-b725-faaa5d66ca74", "HD-011 MELANDER",
			models.NewSeller("f4d234ff-7aa5-4986-954c-8c2cc61ea0fc", "Balam Industries"),
			models.NewCategory("c735f60a-bebf-4d2f-a016-190a883eb99f", "Head"),
			75000),
		models.NewProduct("29f9978e-77f1-4e03-a054-9244d6bb00d0", "VP-44D",
			models.NewSeller("7e13d4e2-408b-494f-a611-1950a3a36616", "Arquebus Corporation"),
			models.NewCategory("c735f60a-bebf-4d2f-a016-190a883eb99f", "Head"),
			231000),
		models.NewProduct("55a47ff6-d8e2-491b-9378-8d536c1f4e44", "VP-44S",
			models.NewSeller("7e13d4e2-408b-494f-a611-1950a3a36616", "Arquebus Corporation"),
			models.NewCategory("c735f60a-bebf-4d2f-a016-190a883eb99f", "Head"),
			124000),
	}

	checkProducts(t, got, want)
}

func testReadAllByFilter(t *testing.T, store models.ProductStore) {
	t.Helper()

	got, err := store.ReadAllByFilter("c735f60a-bebf-4d2f-a016-190a883eb99f", "7e13d4e2-408b-494f-a611-1950a3a36616")
	if err != nil {
		t.Error(err)
	}

	want := []models.Product{
		models.NewProduct("29f9978e-77f1-4e03-a054-9244d6bb00d0", "VP-44D",
			models.NewSeller("7e13d4e2-408b-494f-a611-1950a3a36616", "Arquebus Corporation"),
			models.NewCategory("c735f60a-bebf-4d2f-a016-190a883eb99f", "Head"),
			231000),
		models.NewProduct("55a47ff6-d8e2-491b-9378-8d536c1f4e44", "VP-44S",
			models.NewSeller("7e13d4e2-408b-494f-a611-1950a3a36616", "Arquebus Corporation"),
			models.NewCategory("c735f60a-bebf-4d2f-a016-190a883eb99f", "Head"),
			124000),
	}

	checkProducts(t, got, want)
}

func checkProducts(t *testing.T, got, want []models.Product) {
	t.Helper()

	if !slices.EqualFunc(got, want, isEqualProducts) {
		t.Errorf("Incorrect Products:\ngot: \n%v\nwant: \n%v\n", got, want)
	}
}

func isEqualProducts(E1, E2 models.Product) bool {
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

func checkProduct(t *testing.T, got, want models.Product) {
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
