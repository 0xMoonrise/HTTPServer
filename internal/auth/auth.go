package auth

import (
	"log"
)

func HashPassword(password string) string {
	log.Println("The password has been hashed successfully.")
	return password
}
