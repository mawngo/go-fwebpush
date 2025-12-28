package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"github.com/mawngo/go-fwebpush"
	"os"
	"strings"
	"time"
)

func main() {
	vapid := flag.String("vapid", "", "vapid keypair")
	sub := flag.String("subscription", "", "webpush subscription json")
	subject := flag.String("subject", "example@example.com", "webpush subject")
	msg := flag.String("message", "Test at "+time.Now().Format("15:04:05"), "webpush message")

	if *vapid == "" {
		*vapid = mustReadFile(".vapid.txt")
		if *vapid == "" {
			println("vapid key pair is required")
			return
		}
	}

	if *sub == "" {
		*sub = mustReadFile(".subscription.json")
		if *sub == "" {
			println("subscription is required")
			return
		}
	}

	// Decode subscription.
	s := fwebpush.Subscription{}
	err := json.Unmarshal([]byte(*sub), &s)
	if err != nil {
		panic(err)
	}

	keypair := strings.Split(*vapid, ":")
	pusher, err := fwebpush.NewVAPIDPusher(
		*subject,
		keypair[1],
		keypair[0],
	)

	if err != nil {
		panic(err)
	}

	// Send Notification.
	resp, err := pusher.SendNotificationOptions(context.Background(), []byte(*msg), &s, fwebpush.Options{TTL: 30})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func mustReadFile(filename string) string {
	b, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ""
		}
		panic(err)
	}
	return string(b)
}
