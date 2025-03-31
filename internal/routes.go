package internal

import (
	"Stant/ECommerce/internal/domain"
	"Stant/ECommerce/internal/middleware"
	"Stant/ECommerce/internal/security"
	templates "Stant/ECommerce/web"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/sync/errgroup"
)

func NewMux(categories domain.CategoryStore,
	sellers domain.SellerStore,
	products domain.ProductStore,
	users domain.UserStore,
	sessions domain.SessionStore,
) *http.ServeMux {
	styles := http.FileServer(http.Dir("web/static"))
	serveMux := &http.ServeMux{}
	serveMux.Handle("/static/", http.StripPrefix("/static/", styles))

	checkSession := middleware.CheckSessionMiddleware(sessions)
	serveMux.Handle("/", checkSession(HandleIndex(categories, users)))
	serveMux.Handle("/category/{id}", checkSession(HandleCategory(categories, products, users)))
	serveMux.Handle("/seller/{id}", checkSession(HandleSeller(sellers, products, users)))
	serveMux.Handle("/product/{id}", checkSession(HandleProduct(products, users)))

	serveMux.Handle("GET /register", HandleRegistrationPage())
	serveMux.Handle("POST /register", HandleRegistration(users))
	serveMux.Handle("GET /login", HandleLoginPage())
	serveMux.Handle("POST /login", HandleLogin(users, sessions))

	return serveMux
}

func HandleIndex(categories domain.CategoryStore, users domain.UserStore) http.Handler {
	renderer := templates.Index
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			categories, err := categories.ReadAll()
			if err != nil {
				log.Printf("internal.HandleIndex: [%v]", err)
				http.NotFound(w, r)
				return
			}

			userId, ok := middleware.GetUserId(r.Context())
			if ok {
				user, _ := users.Read(userId)
				log.Printf("Session: [%v]", user.FirstName())
			}

			w.WriteHeader(http.StatusOK)
			renderer(categories).Render(r.Context(), w)
		},
	)
}

func HandleCategory(categories domain.CategoryStore, products domain.ProductStore, users domain.UserStore) http.Handler {
	renderer := templates.Category
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.Atoi(r.PathValue("id"))
			if err != nil {
				log.Printf("internal.HandleCategory: [%v]", err)
				http.Error(w, "Coudn't parse ID", http.StatusInternalServerError)
				return
			}

			categoryChan := make(chan domain.Category)
			defer close(categoryChan)
			productsChan := make(chan []domain.Product)
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
				products, err := products.ReadAllByFilter(id, 0)
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

			userId, ok := middleware.GetUserId(r.Context())
			if ok {
				user, _ := users.Read(userId)
				log.Printf("Session: [%v]", user.FirstName())
			}

			w.WriteHeader(http.StatusOK)
			renderer(category.Name(), products).Render(r.Context(), w)
		},
	)
}

func HandleSeller(sellers domain.SellerStore, products domain.ProductStore, users domain.UserStore) http.Handler {
	renderer := templates.Seller
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.Atoi(r.PathValue("id"))
			if err != nil {
				log.Printf("internal.HandleSeller: [%v]", err)
				http.Error(w, "Coudn't parse ID", http.StatusInternalServerError)
				return
			}

			sellerChan := make(chan domain.Seller)
			defer close(sellerChan)
			productsChan := make(chan []domain.Product)
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
				products, err := products.ReadAllByFilter(0, id)
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

			userId, ok := middleware.GetUserId(r.Context())
			if ok {
				user, _ := users.Read(userId)
				log.Printf("Session: [%v]", user.FirstName())
			}

			w.WriteHeader(http.StatusOK)
			renderer(seller.Name(), products).Render(r.Context(), w)
		},
	)
}

func HandleProduct(products domain.ProductStore, users domain.UserStore) http.Handler {
	renderer := templates.Product
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.Atoi(r.PathValue("id"))
			if err != nil {
				log.Printf("internal.HandleProduct: [%v]", err)
				http.Error(w, "Coudn't parse ID", http.StatusInternalServerError)
				return
			}

			product, err := products.Read(id)
			if err != nil {
				log.Printf("internal.HandleProduct: [%v]", err)
				http.Error(w, "Internal Error", http.StatusInternalServerError)
				return
			}

			userId, ok := middleware.GetUserId(r.Context())
			if ok {
				user, _ := users.Read(userId)
				log.Printf("Session: [%v]", user.FirstName())
			}

			w.WriteHeader(http.StatusOK)
			renderer(product).Render(r.Context(), w)
		},
	)
}

func HandleRegistrationPage() http.Handler {
	renderer := templates.Registration()
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			renderer.Render(r.Context(), w)
		},
	)
}

func HandleRegistration(users domain.UserStore) http.Handler {
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
	renderer := templates.Login()
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			renderer.Render(r.Context(), w)
		},
	)
}

func HandleLogin(users domain.UserStore, sessions domain.SessionStore) http.Handler {
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
