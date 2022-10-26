package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"spreewill-core/pkg/auth"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	pubKeyURL := "https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json"
	formattedURL := fmt.Sprintf(pubKeyURL, os.Getenv("AWS_DEFAULT_REGION"), os.Getenv("COGNITO_USER_POOL_ID"))

	pubKey, err := jwk.Fetch(context.Background(), formattedURL)
	if err != nil {
		log.Fatalln(err)
	}

	tokenAuth = jwtauth.New("HS256", jwt.WithKeySet(pubKey), nil)
	jwt.Parse()

	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": 123})
	fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)
}

func main() {
	cognitoClient := auth.Init()

	r := chi.NewRouter()

	r.Use(middleware.Logger, middleware.WithValue("CognitoClient", cognitoClient))

	r.Route("/auth", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Post("/signup", auth.SignUp)
		r.Post("/verify", auth.VerifyUser)
		r.Post("/signin", auth.SignIn)
	})

	r.NotFoundHandler()
	port := os.Getenv("PORT")

	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}
