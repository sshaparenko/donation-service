package services

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sshaparenko/donation-service/pkg/database"
	"github.com/sshaparenko/donation-service/pkg/domain"
	"golang.org/x/crypto/bcrypt"
)

func Register(request *domain.Register) error {
	var uuid uuid.UUID = uuid.New()
	query := "INSERT INTO SERVICE_USER VALUES (:id, :first_name, :last_name, :password, :email, :fk_org)"

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	result, err := database.DB.NamedExec(query, domain.User{
		ID:             uuid,
		FirstName:      request.FirstName,
		LastName:       request.LastName,
		Password:       string(hashedPass),
		Email:          request.Email,
		FkOrganization: "0",
	})
	fmt.Print(result)
	if err != nil {
		return err
	}
	return nil
}

func GetUserByUsername(username string) (*domain.User, error) {
	var usr domain.User
	var query string = "SELECT * FROM SERVICE_USER WHERE USERNAME = $1"

	result := database.DB.QueryRow(query, username)
	if result.Err() != nil {
		return nil, result.Err()
	}
	err := result.Scan(&usr)
	if err != nil {
		return nil, err
	}
	return &usr, nil
}

func GetAllOrganizations() ([]*domain.Organization, error) {
	var organizations []*domain.Organization
	var query string = "SELECT * FROM ORGANIZATION"

	rows, err := database.DB.Queryx(query)
	if err != nil {
		return organizations, err
	}

	parseRows(rows, &organizations)

	return organizations, nil
}

func parseRows[T any](rows *sqlx.Rows, slice *[]*T) {
	for rows.Next() {
		var model T
		err := rows.StructScan(&model)
		if err != nil {
			log.Fatalln(err)
		}
		*slice = append(*slice, &model)
	}
}
