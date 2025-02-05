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
	mux := http.NewServeMux()
	mux.Handle("/", handleIndex())

	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	got := httptest.NewRecorder()
	mux.ServeHTTP(got, request)

	want := httptest.NewRecorder()
	views.Index().Render(context.Background(), want)

	checkResponseStatus(t, got.Code, http.StatusOK)
	checkResponseBody(t, *got.Body, *want.Body)
}

func TestCategoryHandler(t *testing.T) {
	products := []domain.Product{
		domain.NewProduct(1, "Huawei"),
		domain.NewProduct(2, "Lenovo"),
		domain.NewProduct(3, "ThinkPad"),
	}
	store := InMemoryStore{products}
	mux := http.NewServeMux()
	mux.Handle("/category/{name}", handleCategory(store))

	request, _ := http.NewRequest(http.MethodGet, "/category/Laptops", nil)
	got := httptest.NewRecorder()
	mux.ServeHTTP(got, request)

	want := httptest.NewRecorder()
	test, _ := store.ReadAll()
	views.Category("Laptops", test).Render(context.Background(), want)

	checkResponseStatus(t, got.Code, http.StatusOK)
	checkResponseBody(t, *got.Body, *want.Body)
}

func TestProductHandler(t *testing.T) {
	products := []domain.Product{
		domain.NewProduct(1, "Huawei"),
		domain.NewProduct(2, "Lenovo"),
		domain.NewProduct(3, "ThinkPad"),
	}
	store := InMemoryStore{products}
	mux := http.NewServeMux()
	mux.Handle("/product/{productID}", handleProduct(store))

	request, _ := http.NewRequest(http.MethodGet, "/product/2", nil)
	got := httptest.NewRecorder()
	mux.ServeHTTP(got, request)

	want := httptest.NewRecorder()
	views.Product(products[2]).Render(context.Background(), want)

	checkResponseStatus(t, got.Code, http.StatusOK)
	checkResponseBody(t, *got.Body, *want.Body)
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
