package main

import (
	"fmt"
	"github.com/go-chi/cors"
	"log"
	"net/http"
	"os"
	"spreewill-core/pkg/db"
	"spreewill-core/pkg/services/auth"
	"spreewill-core/pkg/services/customer"
	"spreewill-core/pkg/services/vendorx"
	"spreewill-core/pkg/session"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cognitoClient := auth.Init()
	db.Init()

	session := session.CreateSession()

	vendorService := vendorx.NewVendorService(session)
	customerService := customer.NewCustomerService(session)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// pass database to context too
	r.Use(middleware.Logger, middleware.WithValue("CognitoClient", cognitoClient))

	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/signup", auth.SignUp)
		r.Post("/verify", auth.VerifyUser)
		r.Post("/signin", auth.SignIn)
	})

	r.Route("/api/vendor", func(r chi.Router) {
		r.Post("/", vendorService.CreateVendor)

		r.Group(func(r chi.Router) {
			r.Use(ValidateToken)

			r.Get("/{id}", vendorService.GetVendor)
			r.Get("/all", vendorService.GetVendors)
			r.Put("/", vendorService.UpdateVendor)
			r.Delete("/{id}", vendorService.DeleteVendor)
		})
	})

	r.Route("/api/customer", func(r chi.Router) {
		r.Post("/", customerService.CreateCustomer)

		r.Group(func(r chi.Router) {
			r.Use(ValidateToken)

			r.Get("/{id}", customerService.GetCustomer)
			r.Get("/all", customerService.GetCustomers)
			r.Put("/", customerService.UpdateCustomer)
			r.Delete("/{id}", customerService.DeleteCustomer)
		})
	})

	port := os.Getenv("PORT")
	log.Printf("server started @ %s...", port)
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), r)
}
