package main

import (
	"Stant/ECommerce/domain"
	"Stant/ECommerce/views"
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	categories := []domain.Category{
		domain.NewCategory(1, "Electronics"),
		domain.NewCategory(2, "Laptops"),
		domain.NewCategory(3, "Phones"),
	}
	store := mockCategoryStore{categories}
	mux := http.NewServeMux()
	mux.Handle("/", handleIndex(store))

	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	got := httptest.NewRecorder()
	mux.ServeHTTP(got, request)

	want := httptest.NewRecorder()
	wantCategories, _ := store.ReadAll()
	views.Index(wantCategories).Render(context.Background(), want)

	checkResponseStatus(t, got.Code, http.StatusOK)
	checkResponseBody(t, *got.Body, *want.Body)
}

func TestCategoryHandler(t *testing.T) {
	categories := []domain.Category{
		domain.NewCategory(0, "Electronics"),
		domain.NewCategory(1, "Laptops"),
		domain.NewCategory(2, "Phones"),
	}
	categoryStore := mockCategoryStore{categories}
	products := []domain.Product{
		domain.NewProduct(0, "MacBook", 0, 1, 100),
		domain.NewProduct(1, "ThinkPad", 1, 1, 100),
		domain.NewProduct(2, "Foundation", 2, 1, 100),
	}
	productStore := mockProductStore{products}
	mux := http.NewServeMux()
	mux.Handle("/category/{id}", handleCategory(categoryStore, productStore))

	request, _ := http.NewRequest(http.MethodGet, "/category/1", nil)
	got := httptest.NewRecorder()
	mux.ServeHTTP(got, request)

	want := httptest.NewRecorder()
	wantProducts, _ := productStore.ReadAll()
	views.Category(categories[1].Name(), wantProducts).Render(context.Background(), want)

	checkResponseStatus(t, got.Code, http.StatusOK)
	checkResponseBody(t, *got.Body, *want.Body)
}

func TestProductHandler(t *testing.T) {
	products := []domain.Product{
		domain.NewProduct(0, "MacBook", 0, 1, 100),
		domain.NewProduct(1, "ThinkPad", 1, 1, 100),
		domain.NewProduct(2, "Foundation", 2, 1, 100),
	}
	store := mockProductStore{products}
	mux := http.NewServeMux()
	mux.Handle("/product/{id}", handleProduct(store))

	request, _ := http.NewRequest(http.MethodGet, "/product/2", nil)
	got := httptest.NewRecorder()
	mux.ServeHTTP(got, request)

	want := httptest.NewRecorder()
	views.Product(products[2]).Render(context.Background(), want)

	checkResponseStatus(t, got.Code, http.StatusOK)
	checkResponseBody(t, *got.Body, *want.Body)
}

type mockProductStore struct {
	db []domain.Product
}

func newmockProductStore(db []domain.Product) *mockProductStore {
	return &mockProductStore{db: db}
}

func (st mockProductStore) Read(id int) (domain.Product, error) {
	return st.db[id], nil
}

func (st mockProductStore) ReadAll() ([]domain.Product, error) {
	return st.db, nil
}

func (st mockProductStore) ReadAllByFilter(categoryID int) ([]domain.Product, error) {
	products := []domain.Product{}
	for _, product := range st.db {
		if product.CategoryID() == categoryID {
			products = append(products, product)
		}
	}
	return products, nil
}

type mockCategoryStore struct {
	db []domain.Category
}

func newMockCategoryStore(db []domain.Category) *mockCategoryStore {
	return &mockCategoryStore{db: db}
}

func (st mockCategoryStore) Read(id int) (domain.Category, error) {
	return st.db[id], nil
}

func (st mockCategoryStore) ReadAll() ([]domain.Category, error) {
	return st.db, nil
}

func checkResponseStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("Did not get correct HTTP Body: got: %d, want: %d", got, want)
	}
}

func checkResponseBody(t testing.TB, got, want bytes.Buffer) {
	t.Helper()
	if got.String() != want.String() {
		t.Errorf("Did not get correct HTTP Body:\ngot: \n%q\nwant: \n%q\n", got.String(), want.String())
	}
}
