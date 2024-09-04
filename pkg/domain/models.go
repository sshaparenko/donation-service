package domain

import "github.com/google/uuid"

type User struct {
	ID             uuid.UUID `db:"id"`
	FirstName      string    `db:"first_name"`
	LastName       string    `db:"last_name"`
	Password       string    `db:"password"`
	Email          string    `db:"email"`
	FkOrganization string    `db:"fk_org"`
}

type Organization struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
