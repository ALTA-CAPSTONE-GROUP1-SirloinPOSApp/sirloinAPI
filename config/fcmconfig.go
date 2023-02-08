package config

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

func InitFCMClient() (*messaging.Client, context.Context, error) {
	ctx := context.Background()

	// Use option.WithCredentialsFile to specify the path to your service account key file.
	opt := option.WithCredentialsJSON([]byte(FIREBASECREDENTIALS))

	// Initialize Firebase app
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		fmt.Printf("error initializing app: %v\n", err)
		return nil, nil, err
	}

	// Get Firebase messaging client
	client, err := app.Messaging(ctx)
	if err != nil {
		fmt.Printf("error getting messaging client: %v\n", err)
		return nil, nil, err
	}
	return client, ctx, err
}
