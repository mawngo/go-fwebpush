package fwebpush

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/hkdf"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func getURLEncodedTestSubscription() Subscription {
	return Subscription{
		Endpoint: "https://updates.push.services.mozilla.com/wpush/v2/gAAAAA",
		Keys: Keys{
			P256dh: "BNNL5ZaTfK81qhXOx23-wewhigUeFb632jN6LvRWCFH1ubQr77FE_9qV1FuojuRmHP42zmf34rXgW80OvUVDgTk",
			Auth:   "zqbxT6JKstKSY9JKibZLSQ",
		},
	}
}

func getStandardEncodedTestSubscription() Subscription {
	return Subscription{
		Endpoint: "https://updates.push.services.mozilla.com/wpush/v2/gAAAAA",
		Keys: Keys{
			P256dh: "BNNL5ZaTfK81qhXOx23+wewhigUeFb632jN6LvRWCFH1ubQr77FE/9qV1FuojuRmHP42zmf34rXgW80OvUVDgTk=",
			Auth:   "zqbxT6JKstKSY9JKibZLSQ==",
		},
	}
}

var subs = []string{
	`{
  		"endpoint": "https://fcm.googleapis.com/fcm/send/d4Kva_Hbz0o:APA91bETic_l7GTsOG7W18SMgG-U8dU4azrmxyLU9jFlZ73W1x0hZstgh6ODOXU_ogs2TaYKnbgg2zidzYVLO0aQi_CN1N3_r96gSr5NF27bVLldgfKBozCqu94ubxmPb2hBzOrgOOzG",
  		"expirationTime": null,
  		"keys": {
  		  "p256dh": "BLNORfMiAA0TJ6unnAKaGcvo8KLQocmbez5dRNRYka42-12CjM8YBgBoPrT1jJDBPnjKyhAzB1Bif9cBtKrtiDU",
  		  "auth": "bWqqGJUm3wHSM8XHfV-gOg"
  		}
	}`,
	`{
		"endpoint": "https://updates.push.services.mozilla.com/wpush/v2/gAAAAABmh8irktRlLlzMTleHFEIb8vg_dxkb7IK3IW_ahHcAkBNIKji45qsyLqe3dlEu6mAM14b7zZIy-UUzOPVsHH2RdUglL63r-hjg4M1IGOpGSRFKEzdSVq01lhbWUnqlb1cEaARK1DC9pm9KM2jyogGMh9kF7vKMgufwjBferT06WZ_GNdY",
		"expirationTime": null,
		"keys": {
			"auth": "gOb0-WdnfhYc3rq3VffHQA",
			"p256dh": "BH_Z4gt6Wgp57RCrZ_rLmMplfId6cpAg-LT3OQAtD8YAXazG4m5yZq8xTx2at_6qHDyvUEgPdMUkeW1nV4Wv-1I"
		}
	}`,
	`{
		"endpoint": "https://updates.push.services.mozilla.com/wpush/v2/gAAAAABmgoYCadcGmB9z0nbU-HBbl-virVetbHe7aucmxH7kgE7Y2fSZZq6FqX3JSy2iZmSAPYN4tIxWyAsdN_tWWJTCymWIZkw2N9L3bTYMtjXyKD4DC84qEI9ebMMyRd5_Dz4MRGjMDKCX0phvUUbUPio6_60S7rL3jPBJf1fPDD3pFe5UjN4",
		"expirationTime": null,
		"keys": {
			"auth": "z2xnt5uzJMo6h8MkbsTsHA",
			"p256dh": "BJ-jM5fUDa5xTrNBHv0PbU4TjIZ9DtX_csQNgrdKufRez6xtgwFlWqLZdmiPxK5YKNkH19BxVWyrQsZg0fO6x7E"
		}
	}`,
	`{
  		"endpoint": "https://fcm.googleapis.com/fcm/send/eKAWKNUIYFw:APA91bHkYaziMvso61arnA20A8j83Mv7uv8ud-lIaMoyCEY5UGILMlTc3O-A8r2mBiGQfMZRlZDWMNFOH4EH-oNrWPLGv130bWsRhUgUupwlWoxwxjWuI6YxNoR3c-wsrMFA7CoWwq5E",
  		"expirationTime": null,
  		"keys": {
  		  "p256dh": "BOzEjUBLHGnCmVNnAQ_1Nlfe3Q43N4vE-u-GzIQe2ZjqZ-brL4w8HSA63p1s3qjrD7y1xlpMD9T-kZUNNgBE6XM",
  		  "auth": "UXIaCOt5-izK-hxTRZDQfg"
  		}
	}`,
	`{
  		"endpoint": "https://fcm.googleapis.com/fcm/send/eKgOwo4F0LE:APA91bFMTsfBvshuyJJAL_em_UjS0Xr0ESwiQ5GPFH4uZRuKbtlKxIkAbtpJ--JrNE_u1m0PT-maaOQ1TW_cRJkjNUKOXNyEy2HwgHkYOH7B7o3p8pDZ_W-YKA95zNFkbCuUNqLaorSy",
  		"expirationTime": null,
  		"keys": {
  		  "p256dh": "BIEjQrKvMsFemU-O44mNiC25ia_AtVq7nDWTdV_-1oz09uricpf1vLMnKqgTzncqi7ap8QlLdUvaXKXatKpPVk4",
  		  "auth": "6EedK_-W3uaRObiylehcxw"
  		}
	}`,
	`{
  		"endpoint": "https://fcm.googleapis.com/fcm/send/ds3U0FZvBZI:APA91bEE0qMd4CJ-gnPN-FeQnIagFRB1XiAHygwdckvCAY2VpQo8OcHpXRhpq5qrojGwIIU-GB-6a6WBPUazekKZhjyUr1p7m1KOdjI4AwdbMBPgTanOP--BMul4NyYUzM6VByL0U56S",
  		"expirationTime": null,
  		"keys": {
  		  "p256dh": "BHO6fU-oWDehyIbYk-Scndqu-2qmX4FvkUK7h6bjsMdjqpk9Sv8HZqL93VQz0a6XVvzv3-GfJth1UaG5-QPnmAo",
  		  "auth": "9jZHfLkD29kmNtDSl33fLA"
  		}
	}`,
	`{
		"endpoint": "https://updates.push.services.mozilla.com/wpush/v2/gAAAAABl98DYNXaQ_S-Oz4qlgqZppVvv2-pIc1n8frNu_4QybKMKxlTmHidleQOZsC8spKAtukeXdTm5d1f9haBvGeWyLfGbJPyWT9wcXmOIbuhteKi73tYbVTj_LCEl47WM0jk0OggIgAntMduBHzntgFzRrbI11L9G9KR-Uq3h4ZNczpH_xWY",
		"expirationTime": null,
		"keys": {
			"auth": "GOIZSVgaUcEWN8Eugn70Tw",
			"p256dh": "BCpAtBJMJ8McmJ-d6FKJvn3vLe7R5amXO_UmHzisGdKh2WP8GNLpvEDMUsr0nwmlHkKGIsM-Qzsvx3r63JgPe-E"
		}
	}`,
}
var message = []byte("Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged")

func BenchmarkDefaultConfig(b *testing.B) {
	benchEachSub(b, func(b *testing.B, pusher *VAPIDPusher, sub Subscription, i int) {
		b.Run(fmt.Sprintf("run_%d", i), func(b *testing.B) {
			for b.Loop() {
				sub := sub
				_, err := pusher.PrepareNotificationRequest(context.Background(), message, &sub, Options{})
				if err != nil {
					b.Fatal(err)
					return
				}
			}
		})
	})
}

func BenchmarkOldImpl(b *testing.B) {
	for i, v := range subs {
		s := Subscription{}
		err := json.Unmarshal([]byte(v), &s)
		if err != nil {
			b.Fatal(err)
			return
		}
		b.Run(fmt.Sprintf("run_%d", i), func(b *testing.B) {
			for b.Loop() {
				sub := s
				_, err := goWebpushImpl(
					context.Background(),
					"example@example.com",
					"BDUGjzk8wKOBI96Ip6xG3PVNPfK3RcSCIjhwxY6irbQwNpE5f-1mfBq2rcxhrexjQ5alPA5aiST_PuERnhoaiUM",
					"qtFPvMd1wkVVbPzdRU1TdnXzCV8F4YIWRze7BQGQuy0",
					&sub,
					Options{},
				)
				if err != nil {
					b.Fatal(err)
					return
				}
			}
		})
	}
}

func BenchmarkNoCaching(b *testing.B) {
	benchEachSub(b, func(b *testing.B, pusher *VAPIDPusher, sub Subscription, i int) {
		b.Run(fmt.Sprintf("run_%d", i), func(b *testing.B) {
			for b.Loop() {
				sub := sub
				_, err := pusher.PrepareNotificationRequest(context.Background(), message, &sub, Options{})
				if err != nil {
					b.Fatal(err)
					return
				}
			}
		})
	}, WithVAPIDTokenTTL(0))
}

func BenchmarkVapidAndLocalSecretCachingExpired(b *testing.B) {
	benchEachSub(b, func(b *testing.B, pusher *VAPIDPusher, sub Subscription, i int) {
		_, err := pusher.PrepareNotificationRequest(context.Background(), message, &sub, Options{})
		if err != nil {
			b.Fatal(err)
			return
		}
		if sub.LocalKey == nil {
			b.Fatal("local key not generated")
		}

		b.Run(fmt.Sprintf("run_%d", i), func(b *testing.B) {
			for b.Loop() {
				sub := sub
				_, err := pusher.PrepareNotificationRequest(context.Background(), message, &sub, Options{})
				if err != nil {
					b.Fatal(err)
					return
				}
			}
		})
	}, WithVAPIDTokenTTL(time.Hour), WithLocalSecretTTL(1)) // 1 nanosecond - always expire, if you were not going to run this on a super computer.
}

func BenchmarkVAPIDCaching(b *testing.B) {
	benchEachSub(b, func(b *testing.B, pusher *VAPIDPusher, sub Subscription, i int) {
		b.Run(fmt.Sprintf("run_%d", i), func(b *testing.B) {
			for b.Loop() {
				sub := sub
				_, err := pusher.PrepareNotificationRequest(context.Background(), message, &sub, Options{})
				if err != nil {
					b.Fatal(err)
					return
				}
			}
		})
	}, WithVAPIDTokenTTL(time.Hour))
}

func BenchmarkLocalSecretCaching(b *testing.B) {
	benchEachSub(b, func(b *testing.B, pusher *VAPIDPusher, sub Subscription, i int) {
		_, err := pusher.PrepareNotificationRequest(context.Background(), message, &sub, Options{})
		if err != nil {
			b.Fatal(err)
			return
		}
		if sub.LocalKey == nil {
			b.Fatal("local key not generated")
		}

		b.Run(fmt.Sprintf("run_%d", i), func(b *testing.B) {
			for b.Loop() {
				sub := sub
				_, err := pusher.PrepareNotificationRequest(context.Background(), message, &sub, Options{})
				if err != nil {
					b.Fatal(err)
					return
				}
			}
		})
	}, WithVAPIDTokenTTL(0), WithLocalSecretTTL(time.Hour))
}

func BenchmarkVapidAndLocalSecretCaching(b *testing.B) {
	benchEachSub(b, func(b *testing.B, pusher *VAPIDPusher, sub Subscription, i int) {
		_, err := pusher.PrepareNotificationRequest(context.Background(), message, &sub, Options{})
		if err != nil {
			b.Fatal(err)
			return
		}
		if sub.LocalKey == nil {
			b.Fatal("local key not generated")
		}

		b.Run(fmt.Sprintf("run_%d", i), func(b *testing.B) {
			for b.Loop() {
				sub := sub
				_, err := pusher.PrepareNotificationRequest(context.Background(), message, &sub, Options{})
				if err != nil {
					b.Fatal(err)
					return
				}
			}
		})
	}, WithVAPIDTokenTTL(time.Hour), WithLocalSecretTTL(time.Hour))
}

func BenchmarkVapidAndLocalSecretCachingCacheInit(b *testing.B) {
	benchEachSub(b, func(b *testing.B, pusher *VAPIDPusher, sub Subscription, i int) {
		b.Run(fmt.Sprintf("run_%d", i), func(b *testing.B) {
			for b.Loop() {
				sub := sub
				_, err := pusher.PrepareNotificationRequest(context.Background(), message, &sub, Options{})
				if err != nil {
					b.Fatal(err)
					return
				}
			}
		})
	}, WithVAPIDTokenTTL(time.Hour), WithLocalSecretTTL(time.Hour))
}

func benchEachSub(b *testing.B, bench func(b *testing.B, pusher *VAPIDPusher, sub Subscription, i int), options ...VAPIDPusherOption) {
	pusher, err := NewVAPIDPusher(
		"example@example.com",
		"BDUGjzk8wKOBI96Ip6xG3PVNPfK3RcSCIjhwxY6irbQwNpE5f-1mfBq2rcxhrexjQ5alPA5aiST_PuERnhoaiUM",
		"qtFPvMd1wkVVbPzdRU1TdnXzCV8F4YIWRze7BQGQuy0",
		options...,
	)
	if err != nil {
		b.Fatal(err)
		return
	}
	for i, v := range subs {
		s := Subscription{}
		err := json.Unmarshal([]byte(v), &s)
		if err != nil {
			b.Fatal(err)
			return
		}
		bench(b, pusher, s, i)
	}
}

// goWebpushImpl Implement of https://github.com/SherClockHolmes/webpush-go/blob/master/webpush.go at 2024-05-07.
// We removed the padding logic for better comparison.
//
//nolint:staticcheck
func goWebpushImpl(ctx context.Context, sub, pub, priv string, s *Subscription, options Options) (*http.Request, error) {
	// Authentication secret (auth_secret).
	authSecret, err := base64.RawURLEncoding.DecodeString(s.Keys.Auth)
	if err != nil {
		return nil, err
	}

	// dh (Diffie Hellman).
	dh, err := base64.RawURLEncoding.DecodeString(s.Keys.P256dh)
	if err != nil {
		return nil, err
	}

	// Generate 16 byte salt.
	salt := make([]byte, 16)
	_, err = io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, err
	}

	// Create the ecdh_secret shared key pair.
	curve := elliptic.P256()

	// Application server key pairs (single use).
	localPrivateKey, x, y, err := elliptic.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}

	localPublicKey := elliptic.Marshal(curve, x, y)

	// Combine application keys with receiver's EC public key.
	sharedX, sharedY := elliptic.Unmarshal(curve, dh)
	if sharedX == nil {
		return nil, errors.New("unmarshal Error: Public key is not a valid point on the curve")
	}

	// Derive ECDH shared secret.
	sx, sy := curve.ScalarMult(sharedX, sharedY, localPrivateKey)
	if !curve.IsOnCurve(sx, sy) {
		return nil, errors.New("encryption error: ECDH shared secret isn't on curve")
	}
	mlen := curve.Params().BitSize / 8
	sharedECDHSecret := make([]byte, mlen)
	sx.FillBytes(sharedECDHSecret)

	hash := sha256.New

	// ikm.
	prkInfoBuf := bytes.NewBufferString("WebPush: info\x00")
	prkInfoBuf.Write(dh)
	prkInfoBuf.Write(localPublicKey)

	buf := make([]byte, 60)
	prkHKDF := hkdf.New(hash, sharedECDHSecret, authSecret, prkInfoBuf.Bytes())
	ikm, err := getHKDFKey(prkHKDF, buf[0:32])
	if err != nil {
		return nil, err
	}

	// Derive Content Encryption Key.
	contentEncryptionKeyInfo := []byte("Content-Encoding: aes128gcm\x00")
	contentHKDF := hkdf.New(hash, ikm, salt, contentEncryptionKeyInfo)
	contentEncryptionKey, err := getHKDFKey(contentHKDF, buf[32:48])
	if err != nil {
		return nil, err
	}

	// Derive the Nonce.
	nonceInfo := []byte("Content-Encoding: nonce\x00")
	nonceHKDF := hkdf.New(hash, ikm, salt, nonceInfo)
	nonce, err := getHKDFKey(nonceHKDF, buf[48:60])
	if err != nil {
		return nil, err
	}

	// Cipher.
	c, err := aes.NewCipher(contentEncryptionKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	// Encryption Content-Coding Header.
	recordBuf := bytes.NewBuffer(salt)

	rs := make([]byte, 4)
	binary.BigEndian.PutUint32(rs, uint32(len(message)*8))

	recordBuf.Write(rs)
	recordBuf.Write([]byte{byte(len(localPublicKey))})
	recordBuf.Write(localPublicKey)
	// Data.
	dataBuf := bytes.NewBuffer(message)
	// Padding ending delimeter.
	dataBuf.WriteString("\x02")
	// Compose the ciphertext.
	ciphertext := gcm.Seal([]byte{}, nonce, dataBuf.Bytes(), nil)
	recordBuf.Write(ciphertext)

	// POST request.
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.Endpoint, recordBuf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Encoding", "aes128gcm")
	req.Header.Set("Content-Length", strconv.Itoa(len(ciphertext)))
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("TTL", strconv.Itoa(options.TTL))

	// Сheck the optional headers.
	if len(options.Topic) > 0 {
		req.Header.Set("Topic", options.Topic)
	}

	if isValidUrgency(options.Urgency) {
		req.Header.Set("Urgency", string(options.Urgency))
	}

	// Get VAPID Authorization header
	vapidAuthHeader, err := getVAPIDAuthorizationHeader(
		s.Endpoint,
		sub,
		pub,
		priv,
		time.Now().Add(time.Hour*12),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", vapidAuthHeader)
	return req, nil
}

func getVAPIDAuthorizationHeader(
	endpoint,
	subscriber,
	vapidPublicKey,
	vapidPrivateKey string,
	expiration time.Time,
) (string, error) {
	// Create the JWT token
	subURL, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"aud": fmt.Sprintf("%s://%s", subURL.Scheme, subURL.Host),
		"exp": expiration.Unix(),
		"sub": fmt.Sprintf("mailto:%s", subscriber),
	})

	// Decode the VAPID private key
	decodedVapidPrivateKey, err := base64.RawURLEncoding.DecodeString(vapidPrivateKey)
	if err != nil {
		return "", err
	}

	privKey := generateVAPIDHeaderKeys(decodedVapidPrivateKey)

	// Sign token with private key
	jwtString, err := token.SignedString(privKey)
	if err != nil {
		return "", err
	}

	// Decode the VAPID public key
	pubKey, err := base64.RawURLEncoding.DecodeString(vapidPublicKey)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"vapid t=%s, k=%s",
		jwtString,
		base64.RawURLEncoding.EncodeToString(pubKey),
	), nil
}

// func TestSendNotificationToURLEncodedSubscription(t *testing.T) {
//	pusher, err := NewVAPIDPusher(
//		"<EMAIL@EXAMPLE.COM>",
//		"test-public",
//		"test-private",
//	)
//	if err != nil {
//		t.Fatal(err)
//	}
//	resp, err := pusher.SendNotification(context.Background(), []byte("Test"), getURLEncodedTestSubscription())
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusCreated {
//		t.Fatalf(
//			"Incorreect status code, expected=%d, got=%d",
//			resp.StatusCode,
//			http.StatusCreated,
//		)
//	}
//}
//
// func TestSendNotificationToStandardEncodedSubscription(t *testing.T) {
//	pusher, err := NewVAPIDPusher(
//		"<EMAIL@EXAMPLE.COM>",
//		"test-public",
//		"test-private",
//	)
//	if err != nil {
//		t.Fatal(err)
//	}
//	resp, err := pusher.SendNotification(context.Background(), []byte("Test"), getStandardEncodedTestSubscription())
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusCreated {
//		t.Fatalf(
//			"Incorreect status code, expected=%d, got=%d",
//			resp.StatusCode,
//			http.StatusCreated,
//		)
//	}
//}
