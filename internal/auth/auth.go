package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
	"net/http"
	"time"
	"log"
	"fmt"
	"strings"
)


func GetBearerToken(headers http.Header) (string, error) {

	header := headers.Get("Authorization")

	if header == "" {
		return "", fmt.Errorf("The header was not found")
	}
	
	return strings.Split(header, " ")[1], nil

}


func HashPassword(password string) (string, error) {

    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err

}

func CheckPasswordHash(password, hash string) bool {

    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil

}


func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error){

	claims := &jwt.RegisteredClaims{
		Issuer:		"chirpy",
		IssuedAt:	jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt:	jwt.NewNumericDate(time.Now().Add(expiresIn).UTC()),
		Subject:	userID.String(),	
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	claims := &jwt.RegisteredClaims{}
	
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(tokenSecret), nil
	})

	if err != nil {
	    log.Println(err)
	    return uuid.UUID{}, nil
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
	    return uuid.UUID{}, fmt.Errorf("token is expired")
	}
	
	uuidStr, err := token.Claims.GetSubject()

	if err != nil {
		log.Println(err)
		return uuid.UUID{}, nil
	}

	id, err := uuid.Parse(uuidStr)

	if err != nil {
		log.Println(err)
		return uuid.UUID{}, nil 
	}
	
	return id, err
}
