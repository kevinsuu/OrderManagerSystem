package firebase

import (
	"context"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/db"
	"google.golang.org/api/option"
)

type Firebase struct {
	App      *firebase.App
	Auth     *auth.Client
	Database *db.Client
}

func InitFirebase(ctx context.Context, credFile, projectID string) (*Firebase, error) {
	opt := option.WithCredentialsFile(credFile)
	config := &firebase.Config{
		ProjectID:   projectID,
		DatabaseURL: "https://order-manager-system-a6931-default-rtdb.firebaseio.com",
	}

	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v\n", err)
		return nil, err
	}

	auth, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("Error initializing Firebase auth: %v\n", err)
		return nil, err
	}

	database, err := app.Database(ctx)
	if err != nil {
		log.Fatalf("Error initializing Realtime Database: %v\n", err)
		return nil, err
	}

	return &Firebase{
		App:      app,
		Auth:     auth,
		Database: database,
	}, nil
}

func InitializeApp() (*firebase.App, error) {
	// 從環境變量獲取憑證文件路徑
	credentialsFile := os.Getenv("FIREBASE_CREDENTIALS")
	if credentialsFile == "" {
		return nil, fmt.Errorf("FIREBASE_CREDENTIALS environment variable not set")
	}

	conf := &firebase.Config{
		DatabaseURL: "https://order-manager-system-a6931-default-rtdb.firebaseio.com",
	}

	opt := option.WithCredentialsFile(credentialsFile)
	app, err := firebase.NewApp(context.Background(), conf, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	return app, nil
}
