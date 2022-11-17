package main

import (
	"fmt"
	"net/http"
	"os"
	"spreewill-core/pkg/util"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

func ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := util.GetHeaderToken(w, r)

		if accessToken == "" {
			util.SendError(w, http.StatusBadRequest, "invalid authorization header")
			return
		}

		pubKeyURL := "https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json"
		formattedURL := fmt.Sprintf(pubKeyURL, os.Getenv("AWS_DEFAULT_REGION"), os.Getenv("COGNITO_USER_POOL_ID"))

		keySet, err := jwk.Fetch(r.Context(), formattedURL)
		if err != nil {
			util.SendError(w, http.StatusInternalServerError, "could not retrieve key set")
			return
		}

		_, err = jwt.Parse([]byte(accessToken), jwt.WithKeySet(keySet), jwt.WithValidate(true))
		if err != nil {
			util.SendError(w, http.StatusUnauthorized, "")
			return
		}

		next.ServeHTTP(w, r)
	})
}
