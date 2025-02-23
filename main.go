package main

import (
	"Stant/ECommerce/domain"
	"Stant/ECommerce/stores"
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Database: %s\n", err)
	}
	defer db.Close()
	productStore := stores.NewProductStore(db)
	categoryStore := stores.NewCategoryStore(db)
	sellerStore := stores.NewSellerStore(db)
	userStore := stores.NewUserStore(db)

	loggingMiddleware := LoggingMiddleware(*log.Default())

	serveMux := NewMux(categoryStore, sellerStore, productStore, userStore)

	server := &http.Server{
		Addr:    "localhost:5050",
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
) *http.ServeMux {
	styles := http.FileServer(http.Dir("views/static"))
	serveMux := &http.ServeMux{}
	serveMux.Handle("/static/", http.StripPrefix("/static/", styles))
	serveMux.Handle("/", HandleIndex(categories))
	serveMux.Handle("/category/{id}", HandleCategory(categories, products))
	serveMux.Handle("/seller/{id}", HandleSeller(sellers, products))
	serveMux.Handle("/product/{id}", HandleProduct(products))

	serveMux.Handle("POST /register", HandleRegistration(users))
	serveMux.Handle("/login", http.HandlerFunc(http.NotFound))

	return serveMux
}
