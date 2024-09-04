package routes

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sshaparenko/donation-service/pkg/database"
	"github.com/sshaparenko/donation-service/pkg/handlers"
	"github.com/sshaparenko/donation-service/pkg/middleware"
)

func SetupRotes(r *chi.Mux) {

	r.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Post("/register", handlers.Register)
			r.Get("/donation", handlers.GetAllDonations)
		})

		r.Group(func(r chi.Router) {

			// r.Use(oauth.Authorize(os.Getenv("AUTH_SECRET"), nil))
			r.Use(middleware.OauthMiddleware)
			r.Get("/db", func(w http.ResponseWriter, r *http.Request) {
				stats := database.DB.Stats()
				log.Printf("Open Connections: %d\n", stats.OpenConnections)
				log.Printf("In Use Connections: %d\n", stats.InUse)
				log.Printf("Idle Connections: %d\n", stats.Idle)
			})
			r.Get("/auth/{provider}", handlers.SignInWithProvider)
			r.Get("/auth/{provider}/callback", handlers.CallbackHandler)
		})
	})
}
