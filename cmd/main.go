package main

import (
	"Stant/ECommerce/internal"
	"Stant/ECommerce/internal/middleware"
	"Stant/ECommerce/internal/stores"
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
		log.Fatalf("Main: [%v]\n", err)
	}

	defer db.Close()
	productStore := stores.NewProductStore(db)
	categoryStore := stores.NewCategoryStore(db)
	sellerStore := stores.NewSellerStore(db)
	userStore := stores.NewUserStore(db)
	sessionStore := stores.NewSessionStore(db, time.Now().Add(1*time.Hour))
	defer sessionStore.StopCleanup(sessionStore.StartCleanup(*log.Default(), time.Minute*10))

	logHandlers := middleware.LoggingMiddleware(*log.Default())
	serveMux := internal.NewMux(categoryStore, sellerStore, productStore, userStore, sessionStore)
	server := &http.Server{
		Addr:    ":8080",
		Handler: logHandlers(serveMux),
	}

	go func() {
		log.Println("Server started listening")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Main: [%v]\n", err)
		}
	}()

	<-ctx.Done()
	log.Println("Server stopped listening")
}
