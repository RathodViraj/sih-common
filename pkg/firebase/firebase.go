package firebase

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var App *firebase.App

func InitFirebase() {
	credFile := os.Getenv("FIREBASE_CREDENTIALS")
	if credFile == "" {
		log.Fatal("FIREBASE_CREDENTIALS not set")
	}

	opt := option.WithCredentialsFile(credFile)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v", err)
	}
	App = app
}
