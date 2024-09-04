package handlers

import (
	"errors"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/render"
	"github.com/gorilla/sessions"
	"github.com/sshaparenko/donation-service/pkg/domain"
	"github.com/sshaparenko/donation-service/pkg/services"
	"golang.org/x/oauth2"
)

var mutex sync.Mutex
var ch = make(chan *sessions.Session, 5)
var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

func Register(w http.ResponseWriter, r *http.Request) {
	registerRequset := r.Context().Value("request").(domain.Register)
	err := services.Register(&registerRequset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	render.Render(w, r, &domain.Responce[string]{ //nolint
		HTTPStatusCode: http.StatusCreated,
		Success:        true,
		Message:        "New user was created!",
		Data:           "",
	})
}

func SignIn(w http.ResponseWriter, r *http.Request) {

}

func SignInWithProvider(w http.ResponseWriter, r *http.Request) {
	verifier := oauth2.GenerateVerifier()
	url := generateUrl(r, verifier)
	go createSession(w, r, verifier)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func generateUrl(r *http.Request, verifier string) string {
	conf := r.Context().Value("conf").(*oauth2.Config)
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
	return url
}

func createSession(w http.ResponseWriter, r *http.Request, verifier string) {
	mutex.Lock()
	defer mutex.Unlock()

	session, _ := store.Get(r, "session")
	session.Values["verifier"] = verifier
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   120,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	session.Save(r, w) //nolint
	session.IsNew = false
	ch <- session
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {

	code, err := getAuthorizationCode(r)
	if err != nil {
		renderErrorResponce(w, r, err, domain.InvalidGrant)
		return
	}

	verifier, err := getVerifier()
	if err != nil {
		renderErrorResponce(w, r, err, domain.InvalidGrant)
		return
	}
	tok := getToken(r, code, verifier)
	// client := conf.Client(r.Context(), tok)
	render.Render(w, r, &domain.OAuthResponse{ //nolint
		AccessToken:  tok.AccessToken,
		TokenType:    tok.TokenType,
		ExpiresIn:    tok.Expiry,
		RefreshToken: tok.RefreshToken,
		Scope:        "",
	})
}

func getAuthorizationCode(r *http.Request) (string, error) {
	var code string = r.URL.Query().Get("code")
	if code == "" {
		return "", errors.New("authorization code is empty")
	}
	return code, nil
}

func getVerifier() (string, error) {
	session := <-ch
	if session.IsNew {
		return "", errors.New("verifier string is empty")
	}
	verifier := session.Values["verifier"].(string)
	return verifier, nil
}

func getToken(r *http.Request, code string, verifier string) *oauth2.Token {
	conf := r.Context().Value("conf").(*oauth2.Config)
	tok, err := conf.Exchange(r.Context(), code, oauth2.VerifierOption(verifier))
	if err != nil {
		log.Fatal(err)
	}
	return tok
}

func renderErrorResponce(w http.ResponseWriter, r *http.Request, err error, errorType domain.OAuthErrorResponseType) {
	render.Render(w, r, &domain.OAuthErrorResponse{ //nolint
		Error:       errorType,
		Description: err.Error(),
		Uri:         "",
	})
}
