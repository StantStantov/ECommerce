package internal 

import (
	"Stant/ECommerce/internal/domain"
	"Stant/ECommerce/internal/views"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/sync/errgroup"
)

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

			categoryChan := make(chan domain.Category)
			defer close(categoryChan)
			productsChan := make(chan []domain.Product)
			defer close(productsChan)

			var eg errgroup.Group
			eg.Go(func() error {
				category, err := categoryStore.Read(id)
				if err != nil {
					return err
				}
				categoryChan <- category
				return nil
			})
			eg.Go(func() error {
				products, err := productStore.ReadAllByFilter(id, 0)
				if err != nil {
					return err
				}
				productsChan <- products
				return nil
			})
			category := <-categoryChan
			products := <-productsChan
			if err := eg.Wait(); err != nil {
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
			if err != nil {
				log.Printf("Seller Handler: %s", err)
				http.NotFound(w, r)
				return
			}

			sellerChan := make(chan domain.Seller)
			defer close(sellerChan)
			productsChan := make(chan []domain.Product)
			defer close(productsChan)

			var eg errgroup.Group
			eg.Go(func() error {
				seller, err := sellerStore.Read(id)
				if err != nil {
					return err
				}
				sellerChan <- seller
				return nil
			})
			eg.Go(func() error {
				products, err := productStore.ReadAllByFilter(0, id)
				if err != nil {
					return err
				}
				productsChan <- products
				return nil
			})
			seller := <-sellerChan
			products := <-productsChan
			if err := eg.Wait(); err != nil {
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

func HandleRegistration(users domain.UserStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseForm(); err != nil {
				log.Printf("Register Handler: %s", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			email := r.FormValue("email")
			firstName := r.FormValue("firstName")
			secondName := r.FormValue("secondName")
			password := r.FormValue("password")

			hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
			if err != nil {
				log.Printf("Register Handler: %s", err)
				http.Error(w, "Fail to hash password", http.StatusInternalServerError)
				return
			}

			exists, err := users.IsExists(email)
			if err != nil {
				log.Printf("Register Handler: %s", err)
				http.Error(w, "SQL Error", http.StatusInternalServerError)
				return
			}
			if exists {
				http.Error(w, "User already exists", http.StatusConflict)
				return
			}
			if err := users.Create(email, firstName, secondName, string(hash)); err != nil {
				log.Printf("Register Handler: %s", err)
				http.Error(w, "SQL Error", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/", http.StatusFound)
			return
		},
	)
}

func HandleLogin(users domain.UserStore, sessions domain.SessionStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseForm(); err != nil {
				log.Printf("Register Handler: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			email := r.FormValue("email")
			password := r.FormValue("password")

			user, err := users.ReadByEmail(email)
			if err != nil {
				log.Printf("Register Handler: [%v]", err)
				http.Error(w, "Failed to get User info by Email", http.StatusInternalServerError)
				return
			}
			if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword()), []byte(password)); err != nil {
				http.Error(w, "Password is incorrect", http.StatusConflict)
				return
			}

			expireOn := time.Now().Add(1 * time.Hour)
			sessionToken, err := generateToken(64)
			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    sessionToken,
				Expires:  expireOn,
				HttpOnly: true,
			})
			csrfToken, err := generateToken(64)
			http.SetCookie(w, &http.Cookie{
				Name:     "csrf_token",
				Value:    csrfToken,
				Expires:  expireOn,
				HttpOnly: false,
			})

			if err := sessions.Create(sessionToken, csrfToken, expireOn); err != nil {
				log.Printf("Register Handler: [%v]", err)
				http.Error(w, "Failed to store session", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/", http.StatusFound)
			return
		},
	)
}

func generateToken(length int) (string, error) {
	token := make([]byte, length)
	if _, err := rand.Read(token); err != nil {
		return "", fmt.Errorf("generateToken: [%w]", err)
	}
	return base64.URLEncoding.EncodeToString(token), nil
}
