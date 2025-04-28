package auth

import (
    "testing"
    "time"
    // "fmt"

    // "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {

	t.Run("[Test 1] MakeJWT function is testing:", func(t *testing.T) {

		userID := uuid.New()
		secret := "secret"
		expiresIn := 5*time.Second

		_, err := MakeJWT(userID, secret, expiresIn)

		if err != nil {
			t.Fatalf("Error trying to build the JWT %v", err)
		}
		
		t.Log("Success!")
	})

	t.Run("[Test 2] ValidateJWT function is testing", func(t *testing.T){
		userID := uuid.New()
		secret := "secret"
		expiresIn := 5*time.Second

		token, err := MakeJWT(userID, secret, expiresIn)

		if err != nil {
			t.Fatalf("Error trying to build the JWT %v", err)
		}

		t.Log("Testgin signing")

		uuid, err := ValidateJWT(token, "secret")

		if err != nil {
			t.Fatalf("Error trying to validate %v", err)
		}

		t.Logf("Success! %v", uuid)
	})

	t.Run("[Test 3] ValidateJWT function is testing", func(t *testing.T){
		userID := uuid.New()
		secret := "secret"
		expiresIn := 2*time.Second

		token, err := MakeJWT(userID, secret, expiresIn)
		
		if err != nil {
		    t.Fatalf("Error trying to build the JWT %v", err)
		}

		t.Log("Testing expiring time")
		time.Sleep(3*time.Second)
		uuid, err := ValidateJWT(token, "secret")

		if err != nil {
		    t.Fatalf("Error trying to validate %v", err)
		}
		
		t.Logf("Success! %v", uuid)
	})
}
