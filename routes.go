package main

import (
	"Stant/ECommerce/views"
	"context"
	"net/http"
)

func handleIndex() http.Handler {
	renderer := views.Hello()
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			renderer.Render(context.Background(), w)
		},
	)
}
