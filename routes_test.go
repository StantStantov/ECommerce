package main

import (
	"Stant/ECommerce/domain"
	"Stant/ECommerce/stores"
	"Stant/ECommerce/views"
	"bytes"
	"context"
	"database/sql"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestHandlers(t *testing.T) {
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil || db.Ping() != nil {
		t.Fatalf("Database: %s\n", err)
	}

	categoryStore := stores.NewCategoryStore(db)
	productStore := stores.NewProductStore(db)
	server := NewMux(categoryStore, productStore)

	t.Run("Test Index", func(t *testing.T) {
		t.Parallel()
		testIndexHandler(t, server, categoryStore)
	})
	t.Run("Test Category", func(t *testing.T) {
		t.Parallel()
		testCategoryHandler(t, server, categoryStore, productStore)
	})
	t.Run("Test Product", func(t *testing.T) {
		t.Parallel()
		testProductHandler(t, server, productStore)
	})
}

func testIndexHandler(t *testing.T, server *http.ServeMux, store domain.CategoryStore) {
	t.Helper()

	got := httptest.NewRecorder()
	server.ServeHTTP(got, newGetRequest(t, "/", nil))

	want := httptest.NewRecorder()
	wantCategories, _ := store.ReadAll()
	views.Index(wantCategories).Render(context.Background(), want)

	checkResponseStatus(t, got.Code, http.StatusOK)
	checkResponseBody(t, *got.Body, *want.Body)
}

func testCategoryHandler(t *testing.T,
	server *http.ServeMux,
	categoryStore domain.CategoryStore,
	productStore domain.ProductStore,
) {
	t.Helper()

	got := httptest.NewRecorder()
	server.ServeHTTP(got, newGetRequest(t, "/category/1", nil))

	want := httptest.NewRecorder()
	wantCategory, _ := categoryStore.Read(1)
	wantProducts, _ := productStore.ReadAll()
	views.Category(wantCategory.Name(), wantProducts).Render(context.Background(), want)

	checkResponseStatus(t, got.Code, http.StatusOK)
	checkResponseBody(t, *got.Body, *want.Body)
}

func testProductHandler(t *testing.T, server *http.ServeMux, store domain.ProductStore) {
	t.Helper()

	got := httptest.NewRecorder()
	server.ServeHTTP(got, newGetRequest(t, "/product/1", nil))

	want := httptest.NewRecorder()
	wantProduct, _ := store.Read(1)
	views.Product(wantProduct).Render(context.Background(), want)

	checkResponseStatus(t, got.Code, http.StatusOK)
	checkResponseBody(t, *got.Body, *want.Body)
}

func newGetRequest(t testing.TB, url string, body io.Reader) *http.Request {
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
