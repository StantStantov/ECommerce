package main

import (
	"Stant/ECommerce/domain"
	"Stant/ECommerce/views"
	"log"
	"net/http"
	"strconv"
)

func HandleIndex(store domain.CategoryStore) http.Handler {
	renderer := views.Index
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			categories, err := store.ReadAll()
			if err != nil {
				log.Printf("Category Handler: %s", err)
				http.NotFound(w, r)
				return
			}

			w.WriteHeader(http.StatusOK)
			renderer(categories).Render(r.Context(), w)
		},
	)
}

func HandleCategory(categoryStore domain.CategoryStore, productStore domain.ProductStore) http.Handler {
	renderer := views.Category
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.Atoi(r.PathValue("id"))
			if err != nil {
				log.Printf("Category Handler: %s", err)
				http.NotFound(w, r)
				return
			}

			category, err := categoryStore.Read(id)
			if err != nil {
				log.Printf("Category Handler: %s", err)
				http.NotFound(w, r)
				return
			}

			products, err := productStore.ReadAllByFilter(id, 0)
			if err != nil {
				log.Printf("Category Handler: %s", err)
				http.NotFound(w, r)
				return
			}

			w.WriteHeader(http.StatusOK)
			renderer(category.Name(), products).Render(r.Context(), w)
		},
	)
}

func HandleSeller(sellerStore domain.SellerStore, productStore domain.ProductStore) http.Handler {
	renderer := views.Seller
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.Atoi(r.PathValue("id"))
			println(r.URL.Query().Encode())
			if err != nil {
				log.Printf("Seller Handler: %s", err)
				http.NotFound(w, r)
				return
			}

			seller, err := sellerStore.Read(id)
			if err != nil {
				log.Printf("Seller Handler: %s", err)
				http.NotFound(w, r)
				return
			}

			products, err := productStore.ReadAllByFilter(0, id)
			if err != nil {
				log.Printf("Seller Handler: %s", err)
				http.NotFound(w, r)
				return
			}

			w.WriteHeader(http.StatusOK)
			renderer(seller.Name(), products).Render(r.Context(), w)
		},
	)
}

func HandleProduct(store domain.ProductStore) http.Handler {
	renderer := views.Product
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.Atoi(r.PathValue("id"))
			if err != nil {
				log.Printf("Product Handler: %s", err)
				http.NotFound(w, r)
				return
			}

			product, err := store.Read(id)
			if err != nil {
				log.Printf("Product Handler: %s", err)
				http.NotFound(w, r)
				return
			}

			w.WriteHeader(http.StatusOK)
			renderer(product).Render(r.Context(), w)
		},
	)
}
