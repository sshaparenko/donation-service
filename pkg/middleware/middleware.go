package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/sshaparenko/donation-service/pkg/domain"
	"golang.org/x/oauth2"
)

func JwtMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string = w.Header().Get("Authorization")
		fmt.Print(token)
	})
}

func JsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Method == http.MethodDelete {
			next.ServeHTTP(w, r)
			return
		}

		var registerRequest domain.Register
		// get the content type header from request
		ct := r.Header.Get("Content-Type")
		if ct != "" {
			normalizeHeader(ct, w)
		}
		// enforce a maximum read of 1MBfrom the responce body
		// a request body larger then that will now result in
		// Decode() returning a "http: request body too large" error
		r.Body = http.MaxBytesReader(w, r.Body, 1048576)

		// setup the decoder and call DisallowUnknownFields
		// This will cause Decode() to return "json: unknown field ..." error
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		err := dec.Decode(&registerRequest)
		if err != nil {
			processError(w, err)
			return
		}
		err = dec.Decode(&struct{}{})
		if !errors.Is(err, io.EOF) {
			msg := "Request body must contain a single JSON object"
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), "request", registerRequest) //nolint
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func normalizeHeader(ct string, w http.ResponseWriter) {
	// trim white space
	// get to lowercase
	// parse and normalize the header to remove
	// any additional parameters
	// (like charset or boundary information)
	mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
	if mediaType != "application/json" {
		msg := "Content-Type header is not application/json"
		http.Error(w, msg, http.StatusUnsupportedMediaType)
		return
	}
}

func processError(w http.ResponseWriter, err error) {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError

	switch {
	// Catch any JSON syntax errors and send an error message
	// witch interpolates the location of the problem to make it
	// easier for the client to fix
	case errors.As(err, &syntaxError):
		msg := fmt.Sprintf("Request body contains badly-fromed JSON (at position %d)", &syntaxError.Offset)
		http.Error(w, msg, http.StatusBadRequest)
	// In some circumstances Decode()
	// could return io.ErrUnexpectedEOF
	// There is an open issue regarding this at
	// https://github.com/golang/go/issues/25956
	case errors.Is(err, io.ErrUnexpectedEOF):
		msg := "Request body contains badly-formated JSON"
		http.Error(w, msg, http.StatusBadRequest)
	// Catch any type errors, like assigning int value
	// to a string field of our struct. We can enterpolate the
	// relevant field name and position into an error
	case errors.As(err, &unmarshalTypeError):
		msg := fmt.Sprintf("Request body contains invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
		http.Error(w, msg, http.StatusBadRequest)
	// Catch error caused by any unexpected fields in the request
	// body.
	case strings.HasPrefix(err.Error(), "json: unknown field"):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
		msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
		http.Error(w, msg, http.StatusBadRequest)
	// An io.EOF caused if body is empty
	case errors.Is(err, io.EOF):
		msg := "Request body must not be empty"
		http.Error(w, msg, http.StatusBadRequest)
	// Catch error if a request body is too large
	case err.Error() == "http: request body is too large":
		msg := "Request body must not be larger then 1MB"
		http.Error(w, msg, http.StatusBadRequest)
	// Othewise default to loggin the error and sending a 500 Ineral Server Error
	default:
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func OauthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conf := &oauth2.Config{
			ClientID:     os.Getenv("OAUTH_ID"),
			ClientSecret: os.Getenv("OAUTH_SECRET"),
			Scopes:       []string{"email", "openid", "profile"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://accounts.google.com/o/oauth2/auth",
				TokenURL: "https://oauth2.googleapis.com/token",
			},
			RedirectURL: "http://localhost:8080/api/v1/auth/google/callback",
		}

		ctx := context.WithValue(r.Context(), "conf", conf) //nolint
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
