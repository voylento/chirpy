package auth

import (
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPasswordHash(t *testing.T) {
	password1 := "Testity123!1"
	password2 := "TestityTestTest123!1"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name			string
		password	string
		hash			string
		wantErr		bool
	}{
		{
			name:				"Correct Password",
			password:		password1,
			hash:				hash1,
			wantErr:		false,
		},
		{
			name:				"Incorrect Password",
			password:		"TestityTest321!1",
			hash:				hash1,
			wantErr:		true,
		},
		{
			name:				"Password doesn't match different hash",
			password:		password1,
			hash:				hash2,
			wantErr:		true,
		},
		{
			name:				"Empty password",
			password:		"",
			hash:				hash1,
			wantErr:		true,
		},
		{
			name:				"Invalid hash",
			password: 	password1,
			hash:				"invalidHash",
			wantErr:		true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestMakeJWT_Success(t *testing.T) {
	userID 			:= uuid.New()
	tokenSecret := "testing-secret"
	expiresIn		:= time.Hour

	tokenStr, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Expected no error in call to MakeJWT, got %v", err)
	}

	if tokenStr == "" {
		t.Fatalf("Expected non-empty token string")
	}

	// Basic check that it looks like a JWT (has 2 dots)
	dotCount := 0
	for _, char := range tokenStr {
		if char == '.' {
			dotCount++
		}
	}

	if dotCount != 2 {
		t.Fatalf("Expected JWT format with 2 dots, got %d dots", dotCount)
	}
}

func TestValidateJWT_Success(t *testing.T) {
	userID 			:= uuid.New()
	tokenSecret := "testing-secret"
	expiresIn		:= time.Hour

	tokenStr, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Expected no error in call to MakeJWT, got %v", err)
	}

	parsedUserID, err := ValidateJWT(tokenStr, tokenSecret)
	if err != nil {
		t.Fatalf("Expected no error in call to ValidateJWT, got %v", err)
	}

	if parsedUserID != userID {
		t.Fatalf("Expected parsedUserID %v to equal userID %v", parsedUserID, userID)
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret"
	expiresIn := -time.Hour // Token expired 1 hour ago

	// Create an already expired token
	tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	// Try to validate expired token
	_, err = ValidateJWT(tokenString, tokenSecret)

	if err == nil {
		t.Fatal("Expected error when validating expired token, got nil")
	}
}

func TestValidateJWT_MalformedToken(t *testing.T) {
	tokenSecret := "test-secret"
	malformedToken := "this.is.not.a.valid.jwt"

	// Try to validate malformed token
	_, err := ValidateJWT(malformedToken, tokenSecret)

	if err == nil {
		t.Fatal("Expected error when validating malformed token, got nil")
	}
}

func TestExtractBearerToken(t *testing.T) {
	tests := []struct {
		name		  string
		headerVal	string
		expectErr	bool
		expectVal	string
	}{
		{
			"Valid token",
			"Bearer abc.def.ghi",
			false,
			"abc.def.ghi",
		},
		{
			"Missing Bearer Prefix",
			"abc.def.ghi",
			true,
			"",
		},
		{
			"Empty token",
			"Bearer ",
			true,
			"",
		},
		{
			"Empty header",
			"",
			true,
			"",
		},
	}
	
	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		if tt.headerVal != "" {
			req.Header.Set("Authorization", tt.headerVal)
		}

		got, err := GetBearerToken(req.Header)
		if (err != nil) != tt.expectErr {
			t.Errorf("%s: Expected error = %v, got %v", tt.name, tt.expectErr, err)
		}
		if got != tt.expectVal {
			t.Errorf("%s: Expected token value = %v, got %v", tt.name, tt.expectVal, got)
		}
	}
}
