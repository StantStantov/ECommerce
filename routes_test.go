package main

import (
	"Stant/ECommerce/domain"
	"Stant/ECommerce/stores"
	"Stant/ECommerce/views"
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestHandlers(t *testing.T) {
	db, err := stores.NewDBConn()
	if err != nil || db.Ping() != nil {
		t.Fatalf("Database: %s\n", err)
	}

	categoryStore := stores.NewCategoryStore(db)
	productStore := stores.NewProductStore(db)
	sellerStore := stores.NewSellerStore(db)
	userStore := stores.NewUserStore(db)
	server := NewMux(categoryStore, sellerStore, productStore, userStore)

	t.Run("Test Index", func(t *testing.T) {
		t.Parallel()
		testIndexHandler(t, server, categoryStore)
	})
	t.Run("Test Category", func(t *testing.T) {
		t.Parallel()
		testCategoryHandler(t, server, categoryStore, productStore)
	})
	t.Run("Test Seller", func(t *testing.T) {
		t.Parallel()
		testSellerHandler(t, server, sellerStore, productStore)
	})
	t.Run("Test Product", func(t *testing.T) {
		t.Parallel()
		testProductHandler(t, server, productStore)
	})
	t.Run("Test Registration", func(t *testing.T) {
		t.Parallel()
		testRegisterHandler(t, server, userStore)
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
	wantProducts, _ := productStore.ReadAllByFilter(1, 0)
	views.Category(wantCategory.Name(), wantProducts).Render(context.Background(), want)

	checkResponseStatus(t, got.Code, http.StatusOK)
	checkResponseBody(t, *got.Body, *want.Body)
}

func testSellerHandler(t *testing.T, server *http.ServeMux, sellerStore domain.SellerStore, productStore domain.ProductStore) {
	t.Helper()

	got := httptest.NewRecorder()
	server.ServeHTTP(got, newGetRequest(t, "/seller/1", nil))

	want := httptest.NewRecorder()
	wantSeller, _ := sellerStore.Read(1)
	wantProduct, _ := productStore.ReadAllByFilter(0, 1)
	views.Seller(wantSeller.Name(), wantProduct).Render(context.Background(), want)

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

func testRegisterHandler(t *testing.T, server *http.ServeMux, users domain.UserStore) {
	t.Helper()

	email := "user@test.com"
	content := fmt.Sprintf("email=%s&firstName=test&secondName=test&password=12345", email)
	body := strings.NewReader(content)
	got := httptest.NewRecorder()
	server.ServeHTTP(got, newPostRequest(t, "/register", body))

	checkResponseStatus(t, got.Code, http.StatusFound)
	exists, err := users.IsExists(email)
	checkError(t, err, nil)
	if !exists {
		t.Errorf("Did not get correct User exist status: got: %v, want: %v", got, true)
	}
}

func BenchmarkHandlers(t *testing.B) {
	db, err := sql.Open("pgx", os.Getenv("TEST_DATABASE_URL"))
	if err != nil || db.Ping() != nil {
		t.Fatalf("Database: %s\n", err)
	}

	categoryStore := stores.NewCategoryStore(db)
	productStore := stores.NewProductStore(db)
	sellerStore := stores.NewSellerStore(db)
	userStore := stores.NewUserStore(db)
	server := NewMux(categoryStore, sellerStore, productStore, userStore)

	t.Run("Index", func(t *testing.B) {
		benchmarkIndex(t, server)
	})
	t.Run("Category", func(t *testing.B) {
		benchmarkCategory(t, server)
	})
	t.Run("Seller", func(t *testing.B) {
		benchmarkSeller(t, server)
	})
	t.Run("Product", func(t *testing.B) {
		benchmarkProduct(t, server)
	})
}

func benchmarkIndex(t *testing.B, server *http.ServeMux) {
	t.Helper()
	for range t.N {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetRequest(t, "/", nil))
	}
}

func benchmarkCategory(t *testing.B, server *http.ServeMux) {
	t.Helper()
	for range t.N {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetRequest(t, "/category/1", nil))
	}
}

func benchmarkSeller(t *testing.B, server *http.ServeMux) {
	t.Helper()
	for range t.N {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetRequest(t, "/seller/1", nil))
	}
}

func benchmarkProduct(t *testing.B, server *http.ServeMux) {
	t.Helper()
	for range t.N {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetRequest(t, "/product/1", nil))
	}
}

func newGetRequest(t testing.TB, url string, body io.Reader) *http.Request {
	request, err := http.NewRequest(http.MethodGet, url, body)
	checkError(t, err, nil)
	return request
}

func newPostRequest(t testing.TB, url string, body io.Reader) *http.Request {
	request, err := http.NewRequest(http.MethodPost, url, body)
	checkError(t, err, nil)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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
