package internal

import (
	"Stant/ECommerce/internal/domain/models"
	"Stant/ECommerce/internal/middleware"
	"Stant/ECommerce/internal/security"
	"Stant/ECommerce/internal/views"
	"log"
	"net/http"

	"golang.org/x/sync/errgroup"
)

func NewMux(categories models.CategoryStore,
	sellers models.SellerStore,
	products models.ProductStore,
	users models.UserStore,
	sessions models.SessionStore,
) *http.ServeMux {
	styles := http.FileServer(http.Dir("web/static"))
	serveMux := &http.ServeMux{}
	serveMux.Handle("/static/", http.StripPrefix("/static/", styles))

	checkSession := middleware.CheckSessionMiddleware(sessions)
	serveMux.Handle("/", checkSession(HandleIndex(categories, users)))
	serveMux.Handle("/category/{id}", checkSession(HandleCategory(categories, products, users)))
	serveMux.Handle("/seller/{id}", checkSession(HandleSeller(sellers, products, users)))
	serveMux.Handle("/product/{id}", checkSession(HandleProduct(products, users)))
	serveMux.Handle("/search/", checkSession(HandleSearch(products, users)))

	serveMux.Handle("GET /register", HandleRegistrationPage())
	serveMux.Handle("POST /register", HandleRegistration(users))
	serveMux.Handle("GET /login", HandleLoginPage())
	serveMux.Handle("POST /login", HandleLogin(users, sessions))

	return serveMux
}

func HandleIndex(categories models.CategoryStore, users models.UserStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			categories, err := categories.ReadAll()
			if err != nil {
				log.Printf("internal.HandleIndex: [%v]", err)
				http.NotFound(w, r)
				return
			}

			user := models.User{}
			userId, ok := middleware.GetUserId(r.Context())
			if ok {
				user, err = users.Read(userId)
				if err != nil {
					log.Printf("internal.HandleIndex: [%v]", err)
					http.Error(w, "Internal Error", http.StatusInternalServerError)
					return
				}
			}

			w.WriteHeader(http.StatusOK)
			views.RenderIndexPage(categories, user, w, r.Context())
		},
	)
}

func HandleCategory(categories models.CategoryStore, products models.ProductStore, users models.UserStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id")

			categoryChan := make(chan models.Category)
			defer close(categoryChan)
			productsChan := make(chan []models.Product)
			defer close(productsChan)

			var eg errgroup.Group
			eg.Go(func() error {
				category, err := categories.Read(id)
				if err != nil {
					return err
				}
				categoryChan <- category
				return nil
			})
			eg.Go(func() error {
				products, err := products.ReadAllByFilter(id, "00000000-0000-0000-0000-000000000000")
				if err != nil {
					return err
				}
				productsChan <- products
				return nil
			})
			category := <-categoryChan
			products := <-productsChan
			if err := eg.Wait(); err != nil {
				log.Printf("internal.HandleCategory: [%v]", err)
				http.Error(w, "Internal Error", http.StatusInternalServerError)
				return
			}

			user := models.User{}
			userId, ok := middleware.GetUserId(r.Context())
			if ok {
				user, err := users.Read(userId)
				if err != nil {
					log.Printf("internal.HandleCategory: [%v]", err)
					http.Error(w, "Internal Error", http.StatusInternalServerError)
					return
				}
				user = user
			}

			w.WriteHeader(http.StatusOK)
			views.RenderCategoryPage(category, products, user, w, r.Context())
		},
	)
}

func HandleSeller(sellers models.SellerStore, products models.ProductStore, users models.UserStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id")

			sellerChan := make(chan models.Seller)
			defer close(sellerChan)
			productsChan := make(chan []models.Product)
			defer close(productsChan)

			var eg errgroup.Group
			eg.Go(func() error {
				seller, err := sellers.Read(id)
				if err != nil {
					return err
				}
				sellerChan <- seller
				return nil
			})
			eg.Go(func() error {
				products, err := products.ReadAllByFilter("00000000-0000-0000-0000-000000000000", id)
				if err != nil {
					return err
				}
				productsChan <- products
				return nil
			})
			seller := <-sellerChan
			products := <-productsChan
			if err := eg.Wait(); err != nil {
				log.Printf("internal.HandleSeller: [%v]", err)
				http.Error(w, "Internal Error", http.StatusInternalServerError)
				return
			}

			user := models.User{}
			userId, ok := middleware.GetUserId(r.Context())
			if ok {
				user, err := users.Read(userId)
				if err != nil {
					log.Printf("internal.HandleSeller: [%v]", err)
					http.Error(w, "Internal Error", http.StatusInternalServerError)
					return
				}
				user = user
			}

			w.WriteHeader(http.StatusOK)
			views.RenderSellerPage(seller, products, user, w, r.Context())
		},
	)
}

func HandleProduct(products models.ProductStore, users models.UserStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id")

			product, err := products.Read(id)
			if err != nil {
				log.Printf("internal.HandleProduct: [%v]", err)
				http.Error(w, "Internal Error", http.StatusInternalServerError)
				return
			}

			user := models.User{}
			userId, ok := middleware.GetUserId(r.Context())
			if ok {
				user, err = users.Read(userId)
				if err != nil {
					log.Printf("internal.HandleProduct: [%v]", err)
					http.Error(w, "Internal Error", http.StatusInternalServerError)
					return
				}
			}

			w.WriteHeader(http.StatusOK)
			views.RenderProductPage(product, user, w, r.Context())
		},
	)
}

func HandleSearch(products models.ProductStore, users models.UserStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseForm(); err != nil {
				log.Printf("internal.HandleSearch: [%v]", err)
				http.Error(w, "Internal Error", http.StatusInternalServerError)
				return
			}

			query := r.FormValue("text")
			log.Printf("%s", query)

			http.Redirect(w, r, "/", http.StatusFound)
		},
	)
}

func HandleRegistrationPage() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			views.RenderRegistrationPage(w, r.Context())
		},
	)
}

func HandleRegistration(users models.UserStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseForm(); err != nil {
				log.Printf("internal.HandleRegistration: [%v]", err)
				http.Error(w, "Internal Error", http.StatusInternalServerError)
				return
			}

			email := r.FormValue("email")
			firstName := r.FormValue("firstName")
			secondName := r.FormValue("secondName")
			password := r.FormValue("password")

			exist, err := users.IsExists(email)
			if err != nil {
				log.Printf("internal.HandleRegistration: [%v]", err)
				http.Error(w, "Internal Error", http.StatusInternalServerError)
				return
			}
			if exist {
				http.Error(w, "User already exists", http.StatusConflict)
				return
			}

			hash, err := security.HashPassword(password)
			if err != nil {
				log.Printf("internal.HandleRegistration: [%v]", err)
				http.Error(w, "Internal Error", http.StatusInternalServerError)
				return
			}

			if err := users.Create(email, firstName, secondName, string(hash)); err != nil {
				log.Printf("internal.HandleRegistration: [%v]", err)
				http.Error(w, "Internal Error", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/", http.StatusFound)
			return
		},
	)
}

func HandleLoginPage() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			views.RenderLoginPage(w, r.Context())
		},
	)
}

func HandleLogin(users models.UserStore, sessions models.SessionStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseForm(); err != nil {
				log.Printf("internal.HandleLogin: [%v]", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			email := r.FormValue("email")
			password := r.FormValue("password")

			user, err := users.ReadByEmail(email)
			if err != nil {
				log.Printf("internal.HandleLogin: [%v]", err)
				http.Error(w, "Internal Error", http.StatusInternalServerError)
				return
			}
			if !security.IsCorrectPassword(password, user.HashedPassword()) {
				http.Error(w, "Password is incorrect", http.StatusConflict)
				return
			}

			sessionCookie, err := security.NewSessionCookie()
			if err != nil {
				log.Printf("internal.HandleLogin: [%v]", err)
				http.Error(w, "Internal Error", http.StatusInternalServerError)
				return
			}
			csrfCookie, err := security.NewCsrfCookie()
			if err != nil {
				log.Printf("internal.HandleLogin: [%v]", err)
				http.Error(w, "Internal Error", http.StatusInternalServerError)
				return
			}

			if err := sessions.Create(user.ID(), sessionCookie.Value, csrfCookie.Value); err != nil {
				log.Printf("internal.HandleLogin: [%v]", err)
				http.Error(w, "Internal Error", http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, sessionCookie)
			http.SetCookie(w, csrfCookie)

			http.Redirect(w, r, "/", http.StatusFound)
			return
		},
	)
}
