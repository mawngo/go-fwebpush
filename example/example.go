package main

import (
	"context"
	"encoding/json"
	"github.com/mawngo/go-fwebpush/v2"
)

func main() {
	// Decode subscription.
	s := fwebpush.Subscription{}
	err := json.Unmarshal([]byte("<YOUR_SUBSCRIPTION>"), &s)
	if err != nil {
		panic(err)
	}

	pusher, err := fwebpush.NewVAPIDPusher(
		"example@example.com",
		"<YOUR_VAPID_PUBLIC_KEY>",
		"<YOUR_VAPID_PRIVATE_KEY>",
	)

	if err != nil {
		panic(err)
	}

	// Send Notification.
	resp, err := pusher.SendNotificationOptions(
		context.Background(),
		[]byte("Test"),
		&s,
		fwebpush.Options{TTL: 30},
	)
	if err != nil {
		// TODO: Handle error
	}
	defer resp.Body.Close()
}
