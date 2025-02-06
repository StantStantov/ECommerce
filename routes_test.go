package main

import (
	"Stant/ECommerce/domain"
	"Stant/ECommerce/views"
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	store := newMockCategoryStore()
	mux := newMockServer(store, nil)

	got := httptest.NewRecorder()
	mux.ServeHTTP(got, newMockGetRequest(t, "/", nil))

	want := httptest.NewRecorder()
	wantCategories, _ := store.ReadAll()
	views.Index(wantCategories).Render(context.Background(), want)

	checkResponseStatus(t, got.Code, http.StatusOK)
	checkResponseBody(t, *got.Body, *want.Body)
}

func TestCategoryHandler(t *testing.T) {
	categoryStore := newMockCategoryStore()
	productStore := newMockProductStore()
	mux := newMockServer(categoryStore, productStore)

	got := httptest.NewRecorder()
	mux.ServeHTTP(got, newMockGetRequest(t, "/category/1", nil))

	want := httptest.NewRecorder()
	wantCategory, _ := categoryStore.Read(1)
	wantProducts, _ := productStore.ReadAll()
	views.Category(wantCategory.Name(), wantProducts).Render(context.Background(), want)

	checkResponseStatus(t, got.Code, http.StatusOK)
	checkResponseBody(t, *got.Body, *want.Body)
}

func TestProductHandler(t *testing.T) {
	store := newMockProductStore()
	mux := newMockServer(nil, store)

	got := httptest.NewRecorder()
	mux.ServeHTTP(got, newMockGetRequest(t, "/product/2", nil))

	want := httptest.NewRecorder()
	wantProduct, _ := store.Read(2)
	views.Product(wantProduct).Render(context.Background(), want)

	checkResponseStatus(t, got.Code, http.StatusOK)
	checkResponseBody(t, *got.Body, *want.Body)
}

type mockProductStore struct {
	db []domain.Product
}

func newMockProductStore() *mockProductStore {
	return &mockProductStore{
		db: []domain.Product{
			domain.NewProduct(0, "MacBook", 0, 1, 100),
			domain.NewProduct(1, "ThinkPad", 1, 1, 100),
			domain.NewProduct(2, "Foundation", 2, 1, 100),
		},
	}
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

func newMockCategoryStore() *mockCategoryStore {
	return &mockCategoryStore{
		db: []domain.Category{
			domain.NewCategory(1, "Electronics"),
			domain.NewCategory(2, "Laptops"),
			domain.NewCategory(3, "Phones"),
		},
	}
}

func (st mockCategoryStore) Read(id int) (domain.Category, error) {
	return st.db[id], nil
}

func (st mockCategoryStore) ReadAll() ([]domain.Category, error) {
	return st.db, nil
}

func newMockServer(categoryStore *mockCategoryStore, productStore *mockProductStore) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", handleIndex(categoryStore))
	mux.Handle("/category/{id}", handleCategory(categoryStore, productStore))
	mux.Handle("/product/{id}", handleProduct(productStore))
	return mux
}

func newMockGetRequest(t testing.TB, url string, body io.Reader) *http.Request {
	request, err := http.NewRequest(http.MethodGet, url, body)
	checkError(t, err, nil)
	return request
}

func checkError(t testing.TB, got, want error) {
	t.Helper()
	if got != want {
		t.Errorf("Did not get correct Error: got: %q, want: %q", got, want)
	}
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
