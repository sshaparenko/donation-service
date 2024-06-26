package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

/*
GetValue reads .env and loading env varables
*/
func GetValue(key string) string {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file\n")
		log.Fatalf(err.Error())
	}

	return os.Getenv(key)
}
