package internal

import (
	"Stant/ECommerce/internal/domain"
	"Stant/ECommerce/internal/stores"
	templates "Stant/ECommerce/web"
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
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestHandlers(t *testing.T) {
	db, err := stores.NewDBConn()
	if err != nil {
		t.Fatalf("Main: [%v]\n", err)
	}

	categoryStore := stores.NewCategoryStore(db)
	productStore := stores.NewProductStore(db)
	sellerStore := stores.NewSellerStore(db)
	userStore := stores.NewUserStore(db)
	sessionStore := stores.NewSessionStore(db, time.Now().Add(1*time.Hour))
	server := NewMux(categoryStore, sellerStore, productStore, userStore, sessionStore)

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
	t.Run("Test Login", func(t *testing.T) {
		t.Parallel()
		testLoginHandler(t, server)
	})
}

func testIndexHandler(t *testing.T, server *http.ServeMux, categories domain.CategoryStore) {
	t.Helper()

	got := httptest.NewRecorder()
	server.ServeHTTP(got, newGetRequest(t, "/", nil))

	want := httptest.NewRecorder()
	wantCategories, _ := categories.ReadAll()
	templates.Index(wantCategories).Render(context.Background(), want)

	checkResponseStatus(t, got.Code, http.StatusOK)
	checkResponseBody(t, *got.Body, *want.Body)
}

func testCategoryHandler(t *testing.T,
	server *http.ServeMux,
	categories domain.CategoryStore,
	products domain.ProductStore,
) {
	t.Helper()

	got := httptest.NewRecorder()
	server.ServeHTTP(got, newGetRequest(t, "/category/1", nil))

	want := httptest.NewRecorder()
	wantCategory, _ := categories.Read(1)
	wantProducts, _ := products.ReadAllByFilter(1, 0)
	templates.Category(wantCategory.Name(), wantProducts).Render(context.Background(), want)

	checkResponseStatus(t, got.Code, http.StatusOK)
	checkResponseBody(t, *got.Body, *want.Body)
}

func testSellerHandler(t *testing.T, server *http.ServeMux,
	sellers domain.SellerStore,
	products domain.ProductStore,
) {
	t.Helper()

	got := httptest.NewRecorder()
	server.ServeHTTP(got, newGetRequest(t, "/seller/1", nil))

	want := httptest.NewRecorder()
	wantSeller, _ := sellers.Read(1)
	wantProduct, _ := products.ReadAllByFilter(0, 1)
	templates.Seller(wantSeller.Name(), wantProduct).Render(context.Background(), want)

	checkResponseStatus(t, got.Code, http.StatusOK)
	checkResponseBody(t, *got.Body, *want.Body)
}

func testProductHandler(t *testing.T, server *http.ServeMux, products domain.ProductStore) {
	t.Helper()

	got := httptest.NewRecorder()
	server.ServeHTTP(got, newGetRequest(t, "/product/1", nil))

	want := httptest.NewRecorder()
	wantProduct, _ := products.Read(1)
	templates.Product(wantProduct).Render(context.Background(), want)

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

func testLoginHandler(t *testing.T, server *http.ServeMux) {
	t.Helper()

	email := "readME@test.com"
	password := "12345"
	content := fmt.Sprintf("email=%s&password=%s", email, password)
	body := strings.NewReader(content)
	got := httptest.NewRecorder()
	server.ServeHTTP(got, newPostRequest(t, "/login", body))

	checkResponseStatus(t, got.Code, http.StatusFound)
	cookies := got.Result().Cookies()
	if len(cookies) != 2 {
		t.Fatalf("Did not get correct Cookies amount: got: %v, want: %v", len(cookies), 2)
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
	sessionStore := stores.NewSessionStore(db, time.Now().Add(1*time.Hour))
	server := NewMux(categoryStore, sellerStore, productStore, userStore, sessionStore)

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
