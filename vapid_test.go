package fwebpush

import (
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	jwt2 "github.com/mawngo/go-fwebpush/internal/jwt"
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
				WithVAPIDTokenTTLExt(1*time.Hour),
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
			publicKey, err := func() (*ecdsa.PublicKey, error) {
				// To decode the token, it needs the VAPID public key
				b64 := base64.RawURLEncoding
				decodedVapidPrivateKey, err := b64.DecodeString(vapidPrivateKey)
				if err != nil {
					return nil, fmt.Errorf("could not decode VAPID private key: %w", err)
				}
				privKey, err := generateVAPIDHeaderKeys(decodedVapidPrivateKey)
				if err != nil {
					return nil, err
				}
				return privKey.Public().(*ecdsa.PublicKey), nil
			}()
			if err != nil {
				t.Fatal(err)
			}

			// Validate the token in the Authorization header
			tokenString := getTokenFromAuthorizationHeader(keys.vapid, t)
			verifier, err := jwt2.NewVerifierES(jwt2.ES256, publicKey)
			if err != nil {
				t.Fatal(err)
			}
			token, err := jwt2.Parse([]byte(tokenString), verifier)
			if err != nil {
				t.Fatal(err)
			}
			if ok := token.Header().Algorithm == jwt2.ES256; !ok {
				t.Fatal("Wrong validation method need ECDSA!")
			}

			// Check the claims on the token.
			var claims jwt2.RegisteredClaims
			if err := json.Unmarshal(token.Claims(), &claims); err == nil {
				expectedSub := fmt.Sprintf("mailto:%s", sub)
				if expectedSub != claims.Subject {
					t.Fatalf(
						"Incorreect mailto, expected=%s, got=%s",
						expectedSub,
						claims.Subject,
					)
				}

				if claims.ExpiresAt < now.Add(13*time.Hour).Unix() {
					t.Fatalf(
						"Incorreect exp, expected>%d, got=%d",
						now.Add(13*time.Hour).Unix(),
						claims.ExpiresAt,
					)
				}

				if claims.ExpiresAt > now.Add(14*time.Hour).Unix() {
					t.Fatalf(
						"Incorreect exp, expected<%d, got=%d",
						now.Add(14*time.Hour).Unix(),
						claims.ExpiresAt,
					)
				}

				if claims.Audience == "" {
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
