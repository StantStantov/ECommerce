package main

import (
	"Stant/ECommerce/domain"
	"Stant/ECommerce/stores"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	db, err := stores.NewDBConn()
	if err != nil {
		log.Fatalf("Database: %s\n", err)
	}

	defer db.Close()
	productStore := stores.NewProductStore(db)
	categoryStore := stores.NewCategoryStore(db)
	sellerStore := stores.NewSellerStore(db)
	userStore := stores.NewUserStore(db)
	sessionStore := stores.NewSessionStore(db)
	defer sessionStore.StopCleanup(sessionStore.StartCleanup(*log.Default(), time.Minute*10))

	loggingMiddleware := LoggingMiddleware(*log.Default())

	serveMux := NewMux(categoryStore, sellerStore, productStore, userStore, sessionStore)

	server := &http.Server{
		Addr:    ":8080",
		Handler: loggingMiddleware(serveMux),
	}

	go func() {
		log.Println("Server started listening")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	log.Println("Server stopped listening")
}

func NewMux(categories domain.CategoryStore,
	sellers domain.SellerStore,
	products domain.ProductStore,
	users domain.UserStore,
	sessions domain.SessionStore,
) *http.ServeMux {
	styles := http.FileServer(http.Dir("views/static"))
	serveMux := &http.ServeMux{}
	serveMux.Handle("/static/", http.StripPrefix("/static/", styles))
	serveMux.Handle("/", HandleIndex(categories))
	serveMux.Handle("/category/{id}", HandleCategory(categories, products))
	serveMux.Handle("/seller/{id}", HandleSeller(sellers, products))
	serveMux.Handle("/product/{id}", HandleProduct(products))

	serveMux.Handle("POST /register", HandleRegistration(users))
	serveMux.Handle("POST /login", HandleLogin(users, sessions))

	return serveMux
}
