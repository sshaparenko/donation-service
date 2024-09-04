package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/sshaparenko/donation-service/pkg/domain"
	"github.com/sshaparenko/donation-service/pkg/services"
)

func GetAllDonations(w http.ResponseWriter, r *http.Request) {
	orgs, err := services.GetAllOrganizations()
	if err != nil {
		msg := fmt.Sprintf("Error has occured while reading organizations data: %s", err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
	}

	render.Render(w, r, &domain.Responce[[]*domain.Organization]{ //nolint
		HTTPStatusCode: http.StatusOK,
		Success:        true,
		Message:        "Organizations data returned successfully",
		Data:           orgs,
	})
}
