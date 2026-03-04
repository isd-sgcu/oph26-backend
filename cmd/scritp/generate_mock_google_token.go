package main

import (
	"flag"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	email := flag.String("email", "", "Email to include in the mock Google OAuth token")
	flag.Parse()

	if flag.NArg() > 0 {
		// If email is provided as a positional argument, use it instead of the flag.
		*email = flag.Arg(0)
	}

	if *email == "" {
		panic("Email is required. Use -email flag or provide as a positional argument.")
	}

	mockToken := createMockGoogleOAuthToken(*email)
	println("Email: " + *email)
	println("Mock Google OAuth Token: " + mockToken)
}

func createMockGoogleOAuthToken(email string) string {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
	})
	// Sign the token with a dummy secret since we won't validate it in development.
	// Function that validates the token is at AuthUsecaseImpl.validateGoogleToken, which will parse the claims without validating the signature in development environment.
	tokenString, err := jwtToken.SignedString([]byte("dummy_secret_for_dev"))
	if err != nil {
		panic("Failed to create mock token: " + err.Error())
	}
	return tokenString
}
