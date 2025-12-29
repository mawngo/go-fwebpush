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
			println("VAPID key pair is required")
			return
		}
	}

	if *sub == "" {
		*sub = mustReadFile(".subscription.json")
		if *sub == "" {
			println("Subscription is required")
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
		fwebpush.WithLocalSecretTTL(12*time.Hour),
	)

	if err != nil {
		panic(err)
	}

	// Send Notification.
	start := time.Now()
	resp, err := pusher.SendNotificationOptions(context.Background(), []byte(*msg), &s, fwebpush.Options{TTL: 30})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		println("Notification sent took", time.Since(start).String())
		file, err := os.Create(".subscription.json")
		if err != nil {
			println("Error creating file:", err)
			return
		}
		defer file.Close()
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		err = encoder.Encode(s)
		if err != nil {
			println("Error encoding json to file:", err)
			return
		}
		return
	}
	println("error status", resp.StatusCode)
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
