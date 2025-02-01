package main

import (
	"Stant/ECommerce/views"
	"net/http"
)

func handleIndex() http.Handler {
	renderer := views.Index()
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			renderer.Render(r.Context(), w)
		},
	)
}

func handleCategory() http.Handler {
	renderer := views.Category
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			renderer(r.PathValue("name")).Render(r.Context(), w)
		},
	)
}
