package helper

import (
	"fmt"
	"log"
	"sirloinapi/config"

	"firebase.google.com/go/v4/messaging"
)

func PushNotification(msg, token string) error {
	client, ctx, err := config.InitFCMClient()
	if err != nil {
		log.Println("error initializing FCM client: " + err.Error())
	}
	// Define the message to be sent
	message := messaging.Message{
		Data: map[string]string{
			"message": msg,
		},
		Token: token,
	}

	// Send the message to the device
	response, err := client.Send(ctx, &message)
	log.Println(response)
	if err != nil {
		fmt.Printf("error sending message: %v\n", err)
		return err
	}

	return nil
}
