package main

import (
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
	productStore := newSQLProductStore(db)

	loggingMiddleware := LoggingMiddleware(*log.Default())

	styles := http.FileServer(http.Dir("views/static"))
	serveMux := &http.ServeMux{}
	serveMux.Handle("/static/", http.StripPrefix("/static/", styles))
	serveMux.Handle("/", handleIndex())
	serveMux.Handle("/category/{name}", handleCategory(productStore))
	serveMux.Handle("/product/{id}", handleProduct(productStore))

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
