package models

/*
Responce struct is a representation of an API endpoint responce
*/
type Responce[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}
