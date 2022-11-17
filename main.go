package main

import (
	"fmt"
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
	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}
