package config

import (
	"context"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func InitializeFirebase(ctx context.Context) (*auth.Client, error) {

	err := godotenv.Load() // Add this line to load the .env file
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	credentialsFile := os.Getenv("FIREBASE_CREDENTIALS_FILE")
	if credentialsFile == "" {
		return nil, fmt.Errorf("FIREBASE_CREDENTIALS_FILE environment variable not set")
	}

	opt := option.WithCredentialsFile(credentialsFile)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %v", err)
	}

	client, err := app.Auth(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting auth client: %v", err)
	}

	log.Println("Firebase app initialized")
	return client, nil
}
