package auth

import (
    "testing"
    "time"
	"net/http"
    "github.com/google/uuid"
)

func TestGetBearerToken(t *testing.T) {
    t.Run("[Test 1] Valid Bearer header", func(t *testing.T) {
        headers := http.Header{}
        headers.Set("Bearer", "some-token-value")

        token, err := GetBearerToken(headers)

        if err != nil {
            t.Fatalf("expected no error, got %v", err)
        }

        if token != "some-token-value" {
            t.Fatalf("expected token to be 'some-token-value', got %v", token)
        }

        t.Log("[Test 1] status: SUCCESS")
    })

    t.Run("[Test 2] Missing Bearer header", func(t *testing.T) {
        headers := http.Header{} // No header set

        token, err := GetBearerToken(headers)

        if err == nil {
            t.Fatalf("expected an error for missing Bearer header, got nil")
        }

        if token != "" {
            t.Fatalf("expected token to be empty string, got %v", token)
        }

        t.Log("[Test 2] status: SUCCESS")
    })
}

func TestJWT(t *testing.T) {

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
