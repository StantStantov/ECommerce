package main

import (
	"Stant/ECommerce/domain"
	"Stant/ECommerce/views"
	"log"
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

func handleCategory(store domain.ProductStore) http.Handler {
	renderer := views.Category
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			products, err := store.ReadAll()
			if err != nil {
				log.Printf("Category Handler: %s", err)
				http.NotFound(w, r)
				return
			}

			w.WriteHeader(http.StatusOK)
			renderer(r.PathValue("name"), products).Render(r.Context(), w)
		},
	)
}
