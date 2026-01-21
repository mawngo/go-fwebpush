package fwebpush

import (
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt"
	"strings"
	"testing"
	"time"
)

func TestVAPID(t *testing.T) {
	now := time.Now()
	sub := "test@test.com"
	subscriptions := []Subscription{getStandardEncodedTestSubscription(), getURLEncodedTestSubscription()}
	for _, s := range subscriptions {
		t.Run("vapid", func(t *testing.T) {
			// Generate vapid keys
			vapidPrivateKey, vapidPublicKey, err := GenerateVAPIDKeys()
			if err != nil {
				t.Fatal(err)
			}

			p, err := NewVAPIDPusher(
				sub,
				vapidPublicKey,
				vapidPrivateKey,
				WithVAPIDTokenTTL(12*time.Hour),
			)
			if err != nil {
				t.Fatal(err)
			}

			// Get authentication header
			keys, err := p.getCachedKeys(s.Endpoint, time.Now())
			if err != nil {
				t.Fatal(err)
			}
			vapid := keys.vapid

			// Validate the token in the Authorization header
			tokenString := getTokenFromAuthorizationHeader(keys.vapid, t)

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
				if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
					t.Fatal("Wrong validation method need ECDSA!")
				}

				// To decode the token it needs the VAPID public key
				b64 := base64.RawURLEncoding
				decodedVapidPrivateKey, err := b64.DecodeString(vapidPrivateKey)
				if err != nil {
					t.Fatal("Could not decode VAPID private key")
				}

				privKey := generateVAPIDHeaderKeys(decodedVapidPrivateKey)
				return privKey.Public(), nil
			})

			// Check the claims on the token.
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				expectedSub := fmt.Sprintf("mailto:%s", sub)
				if expectedSub != claims["sub"] {
					t.Fatalf(
						"Incorreect mailto, expected=%s, got=%s",
						expectedSub,
						claims["sub"],
					)
				}

				if int64(claims["exp"].(float64)) < now.Add(13*time.Hour).Unix() {
					t.Fatalf(
						"Incorreect exp, expected>%d, got=%s",
						now.Add(13*time.Hour).Unix(),
						claims["exp"],
					)
				}

				if int64(claims["exp"].(float64)) > now.Add(14*time.Hour).Unix() {
					t.Fatalf(
						"Incorreect exp, expected<%d, got=%s",
						now.Add(14*time.Hour).Unix(),
						claims["exp"],
					)
				}

				if claims["aud"] == "" {
					t.Fatal("Audience should not be empty")
				}
			} else {
				t.Fatal(err)
			}

			regenerate, err := p.getCachedKeys(s.Endpoint, time.Now())
			if err != nil {
				t.Fatal(err)
			}
			if regenerate.vapid != vapid {
				t.Fatal("regeneration does not reuse vapid header")
			}
		})
	}
}

func TestVAPIDCacheExpired(t *testing.T) {
	s := getStandardEncodedTestSubscription()
	sub := "test@test.com"

	// Generate vapid keys
	vapidPrivateKey, vapidPublicKey, err := GenerateVAPIDKeys()
	if err != nil {
		t.Fatal(err)
	}

	p, err := NewVAPIDPusher(
		sub,
		vapidPublicKey,
		vapidPrivateKey,
		// Always expire.
		WithVAPIDTokenTTL(0),
	)
	if err != nil {
		t.Fatal(err)
	}

	// Get authentication header
	keys, err := p.getCachedKeys(s.Endpoint, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	vapid := keys.vapid
	regenerate, err := p.getCachedKeys(s.Endpoint, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if regenerate.vapid == vapid {
		t.Fatal("regeneration reuse expired token")
	}
}

func TestVAPIDKeys(t *testing.T) {
	privateKey, publicKey, err := GenerateVAPIDKeys()
	if err != nil {
		t.Fatal(err)
	}

	if len(privateKey) != 43 {
		t.Fatal("Generated incorrect VAPID private key")
	}

	if len(publicKey) != 87 {
		t.Fatal("Generated incorrect VAPID public key")
	}
}

// Helper function for extracting the token from the Authorization header.
func getTokenFromAuthorizationHeader(tokenHeader string, t *testing.T) string {
	hsplit := strings.Split(tokenHeader, " ")
	if len(hsplit) < 3 {
		t.Fatal("Failed to auth split header")
	}

	tsplit := strings.Split(hsplit[1], "=")
	if len(tsplit) < 2 {
		t.Fatal("Failed to t split header on =")
	}

	return tsplit[1][:len(tsplit[1])-1]
}
