package main

import (
	"encoding/base64"
	"time"
	"os"
	"log"
	"github.com/dgrijalva/jwt-go"
)

const usage = `
	Usage: make_token.go <username> <project>
	Set base64-encoded token to OPTIMUS_JWT_KEY_BASE64 environment variable
`

func main() {
	key_base64 := os.Getenv("OPTIMUS_JWT_KEY_BASE64")
	key, err := base64.StdEncoding.DecodeString(key_base64)
	if err != nil || len(key_base64) != 8 {
		log.Print(usage)
		log.Fatal(err)
	}

    args := os.Args[1:]
    if len(args) != 2 {
    	log.Fatal(usage)
    }

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"optimus_project": args[1],
		"sub": args[0],
		"nbf": time.Now().Unix(),
		"exp": time.Now().AddDate(0, 6, 0).Unix(), // 6 months
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(key)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Token:", tokenString)
}
