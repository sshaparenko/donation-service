package domain

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
)

type OAuthErrorResponseType string

const (
	InvalidRequest       OAuthErrorResponseType = "invalid_request"
	InvalidClient        OAuthErrorResponseType = "invalid_client"
	InvalidGrant         OAuthErrorResponseType = "invalid_grant"
	UnauthorizedClient   OAuthErrorResponseType = "unauthorized_client"
	UnsupportedGrantType OAuthErrorResponseType = "unsupported_grant_type"
	InvalidScope         OAuthErrorResponseType = "invalid_scope"
)

/*
Responce struct is a representation of an API endpoint responce
*/
type Responce[T any] struct {
	HTTPStatusCode int    `json:"-"`
	Success        bool   `json:"success"`
	Message        string `json:"message"`
	Data           T      `json:"data"`
}

func (rs *Responce[T]) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, rs.HTTPStatusCode)
	return nil
}

type ErrorResponse struct {
	ErrorMessage string `json:"error_message"`
	Field        string `json:"field"`
}

type OAuthResponse struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    time.Time `json:"expires_in"`
	RefreshToken string    `json:"refresh_token"`
	Scope        string    `json:"scope"`
}

func (rs *OAuthResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, http.StatusOK)
	return nil
}

type OAuthErrorResponse struct {
	Error       OAuthErrorResponseType `json:"error"`
	Description string                 `json:"error_description"`
	Uri         string                 `json:"error_uri"`
}

func (rs *OAuthErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	switch rs.Error {
	case "unauthorized_client":
		render.Status(r, http.StatusUnauthorized)
	default:
		render.Status(r, http.StatusBadRequest)
	}
	return nil
}
