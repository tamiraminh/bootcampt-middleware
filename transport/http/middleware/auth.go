package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/oauth"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/spf13/viper"
)

type Authentication struct {
	db *infras.MySQLConn
}

const (
	HeaderAuthorization = "Authorization"
)

func ProvideAuthentication(db *infras.MySQLConn) *Authentication {
	return &Authentication{
		db: db,
	}
}

func (a *Authentication) ClientCredential(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get(HeaderAuthorization)
		token := oauth.New(a.db.Read, oauth.Config{})

		parseToken, err := token.ParseWithAccessToken(accessToken)
		if err != nil {
			response.WithMessage(w, http.StatusUnauthorized, err.Error())
			return
		}

		if !parseToken.VerifyExpireIn() {
			response.WithMessage(w, http.StatusUnauthorized, err.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *Authentication) ClientCredentialWithQueryParameter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		token := params.Get("token")
		tokenType := params.Get("token_type")
		accessToken := tokenType + " " + token

		auth := oauth.New(a.db.Read, oauth.Config{})
		parseToken, err := auth.ParseWithAccessToken(accessToken)
		if err != nil {
			response.WithMessage(w, http.StatusUnauthorized, err.Error())
			return
		}

		if !parseToken.VerifyExpireIn() {
			response.WithMessage(w, http.StatusUnauthorized, err.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *Authentication) Password(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get(HeaderAuthorization)
		token := oauth.New(a.db.Read, oauth.Config{})

		parseToken, err := token.ParseWithAccessToken(accessToken)
		if err != nil {
			response.WithMessage(w, http.StatusUnauthorized, err.Error())
			return
		}

		if !parseToken.VerifyExpireIn() {
			response.WithMessage(w, http.StatusUnauthorized, err.Error())
			return
		}

		if !parseToken.VerifyUserLoggedIn() {
			response.WithMessage(w, http.StatusUnauthorized, oauth.ErrorInvalidPassword)
			return
		}

		next.ServeHTTP(w, r)
	})
}


func (a *Authentication) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		apiKey := viper.Get("API_KEY")

		xApiKey := r.Header.Get("X-API-KEY")

		if apiKey != xApiKey {
			fmt.Println("API KEY didn't match")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *Authentication) CustomMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		log.Println("custom middleware is executed")
		next.ServeHTTP(w, r)
	})
}