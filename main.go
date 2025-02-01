package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	serveMux := &http.ServeMux{}
	serveMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	serveMux.Handle("/", handleIndex())
	serveMux.Handle("/category/{name}", handleCategory())

	server := &http.Server{
		Addr:    "localhost:5050",
		Handler: serveMux,
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
