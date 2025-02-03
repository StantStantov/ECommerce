package main

import (
	"Stant/ECommerce/domain"
	views "Stant/ECommerce/views/templates"
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
		domain.NewProduct("Huawei"), domain.NewProduct("Lenovo"), domain.NewProduct("ThinkPad"),
	}
	store := InMemoryStore{products}
	mux := http.NewServeMux()
	mux.Handle("/category/{name}", handleCategory(store))

	request, _ := http.NewRequest(http.MethodGet, "/category/Laptops", nil)
	got := httptest.NewRecorder()
	mux.ServeHTTP(got, request)

	want := httptest.NewRecorder()
	views.Category("Laptops", products).Render(context.Background(), want)

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
