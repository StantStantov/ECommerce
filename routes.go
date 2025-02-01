package main

import (
	"Stant/ECommerce/views"
	"database/sql"
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

func handleCategory(db *sql.DB) http.Handler {
	q := "SELECT name FROM laptops"
	renderer := views.Category
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			rows, err := db.Query(q)
			if err != nil {
				log.Printf("Category: %s\n", err)
				http.NotFound(w, r)
				return
			}
			products := []string{}
			for rows.Next() {
				var product string
				if err := rows.Scan(&product); err != nil {
					log.Printf("Category: %s\n", err)
					http.NotFound(w, r)
					return
				}
				products = append(products, product)
			}

			w.WriteHeader(http.StatusOK)
			renderer(r.PathValue("name"), products).Render(r.Context(), w)
		},
	)
}
