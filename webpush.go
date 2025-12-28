package fwebpush

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"golang.org/x/crypto/hkdf"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var ErrEncryption = errors.New("encryption error")

var (
	nonceInfo                = []byte("Content-Encoding: nonce\x00")
	contentEncryptionKeyInfo = []byte("Content-Encoding: aes128gcm\x00")
	webpushInfo              = []byte("WebPush: info\x00")
)

const (
	hkdfLen    = 60
	saltLen    = 16
	gcmTagSize = 16
	rsLen      = 4
)

type VAPIDPusher struct {
	client                   *http.Client
	mailtoSub                string        // Sub in VAPID JWT token.
	vapidPublicKeyHeaderPart string        // VAPID public key passed in the VAPID Authorization header (format: `, k=<key`).
	vapidPrivateKey          []byte        // VAPID private key, used to sign VAPID JWT token.
	vapidTokenTTL            time.Duration // Optional, expiration for VAPID JWT token.
	// vapidTTLBuffer additional duration added to expiration.
	// The key will expire later than configured expiration this amount of duration,
	// while the validation of the key will expire sooner than configured expiration this amount of duration,
	// thus make the actual expiration time equal to configured expiration.
	// It is recommended to set this field to at least 30 minutes.
	vapidTTLBuffer   time.Duration
	localSecretTTLFn func() time.Duration // Optional, enable local public key and secret reuse.
	base64Encoding   Base64Encoding
	randReader       io.Reader

	mu    sync.RWMutex
	cache map[string]*reusableKey // Cache of VAPID JWT token by audience.
}

func NewVAPIDPusher(
	subject string,
	vapidPublicKey string,
	vapidPrivateKey string,
	options ...VAPIDPusherOption,
) (*VAPIDPusher, error) {
	c := &VAPIDPusher{
		mailtoSub:      "mailto:" + subject,
		vapidTokenTTL:  12 * time.Hour,
		cache:          make(map[string]*reusableKey),
		base64Encoding: base64.RawURLEncoding,
		vapidTTLBuffer: 1 * time.Hour,
		randReader:     rand.Reader,
	}
	for _, opt := range options {
		opt(c)
	}

	// Decode the VAPID private key.
	var err error
	c.vapidPrivateKey, err = c.decodeBase64(vapidPrivateKey)
	if err != nil {
		return nil, err
	}
	// Decode the VAPID public key.
	vapidPublicKeyBytes, err := c.decodeBase64(vapidPublicKey)
	if err != nil {
		return nil, err
	}
	c.vapidPublicKeyHeaderPart = ", k=" + c.base64Encoding.EncodeToString(vapidPublicKeyBytes)

	if c.client == nil {
		c.client = &http.Client{
			Timeout: 1 * time.Minute,
		}
	}
	return c, nil
}

type Base64Encoding interface {
	DecodeString(string) ([]byte, error)
	EncodeToString([]byte) string
}

// Options are config and extra params needed to send a notification.
type Options struct {
	Topic   string  // Set the Topic header to collapse pending messages.
	TTL     int     // Set the TTL on the endpoint POST request.
	Urgency Urgency // Set the Urgency header.
}

// Keys are the base64 encoded values from PushSubscription.getKey().
type Keys struct {
	Auth   string `json:"auth"`
	P256dh string `json:"p256dh"`
}

// Subscription represents a PushSubscription object from the Push API.
type Subscription struct {
	Endpoint string    `json:"endpoint"`
	Keys     Keys      `json:"keys"`
	LocalKey *LocalKey `json:"lk"`
}

type LocalKey struct {
	// Public generated public key.
	Public string `json:"p"`
	// Secret generated secret.
	Secret string `json:"s"`
	// At creation timestamp, used for checking expiration.
	At int64 `json:"a"`
}

func (p *VAPIDPusher) IsVapidTokenCachingEnabled() bool {
	return p.vapidTokenTTL > 0
}

func (p *VAPIDPusher) IsLocalSecretCachingEnabled() bool {
	return p.localSecretTTLFn != nil
}

func (p *VAPIDPusher) SendNotification(ctx context.Context, message []byte, sub *Subscription) (*http.Response, error) {
	return p.SendNotificationOptions(ctx, message, sub, Options{})
}

// SendNotificationOptions sends a push notification to a subscription's endpoint.
// Message Encryption for Web Push, and VAPID protocols.
// FOR MORE INFORMATION SEE RFC8291: https://datatracker.ietf.org/doc/rfc8291.
func (p *VAPIDPusher) SendNotificationOptions(ctx context.Context, message []byte, sub *Subscription, options Options) (*http.Response, error) {
	req, err := p.PrepareNotificationRequest(ctx, message, sub, options)
	if err != nil {
		return nil, err
	}
	return p.client.Do(req)
}

// PrepareNotificationRequest prepare a push notification request to a subscription's endpoint.
// Message Encryption for Web Push, and VAPID protocols.
// FOR MORE INFORMATION SEE RFC8291: https://datatracker.ietf.org/doc/rfc8291.
func (p *VAPIDPusher) PrepareNotificationRequest(ctx context.Context, message []byte, sub *Subscription, options Options) (*http.Request, error) {
	authSecret, err := p.decodeBase64(sub.Keys.Auth)
	if err != nil {
		return nil, errors.Join(ErrEncryption, err)
	}
	dh, err := p.decodeBase64(sub.Keys.P256dh)
	if err != nil {
		return nil, errors.Join(ErrEncryption, err)
	}

	// GENERATE VAPID TOKEN AND LOCAL KEYPAIR.
	keys, err := p.getCachedKeys(sub.Endpoint)
	if err != nil {
		return nil, err
	}

	// GENERATE SHARED AND PUBLIC KEY.
	var localPublicKeyBytes []byte
	var sharedECDHSecret []byte
	if p.localSecretTTLFn != nil && sub.LocalKey != nil && sub.LocalKey.At > time.Now().Add(-p.localSecretTTLFn()).UnixMilli() {
		// Use publicKey and secret from LocalKey.
		localPublicKeyBytes, err = p.base64Encoding.DecodeString(sub.LocalKey.Public)
		if err != nil {
			return nil, errors.Join(ErrEncryption, err)
		}
		sharedECDHSecret, err = p.base64Encoding.DecodeString(sub.LocalKey.Secret)
		if err != nil {
			return nil, errors.Join(ErrEncryption, err)
		}
	} else {
		localPublicKeyBytes = keys.localPublicKeyBytes
		// Derive ECDH shared secret.
		dhPublicKey, err := keys.curve.NewPublicKey(dh)
		if err != nil {
			return nil, errors.Join(ErrEncryption, err)
		}
		sharedECDHSecret, err = keys.localPrivateKey.ECDH(dhPublicKey)
		if err != nil {
			return nil, errors.Join(ErrEncryption, err)
		}
		// Update LocalKey if enabled.
		if p.localSecretTTLFn != nil {
			sub.LocalKey = &LocalKey{
				Public: p.base64Encoding.EncodeToString(localPublicKeyBytes),
				Secret: p.base64Encoding.EncodeToString(sharedECDHSecret),
				At:     time.Now().UnixMilli(),
			}
		}
	}

	// GENERATE PAYLOAD.
	// Pre-alloc everything.
	prkLen := len(webpushInfo) + len(dh) + len(localPublicKeyBytes)
	dataLen := len(message) + 1
	cipherTextLen := dataLen + gcmTagSize
	recordLen := saltLen + rsLen + 1 + len(localPublicKeyBytes) + cipherTextLen
	// buf is just for bulk allocation and re-slicing, does not use it directly.
	buf := make([]byte, prkLen+hkdfLen+cipherTextLen+saltLen+recordLen)
	i := 0
	prkInfo := buf[i:prkLen:prkLen]
	i += prkLen
	bufHKDF := buf[i : i+hkdfLen : i+hkdfLen]
	i += hkdfLen
	// Use data slice for both data and cipher text.
	data := buf[i : i+dataLen : i+cipherTextLen]
	i += cipherTextLen
	salt := buf[i : i+saltLen : i+saltLen]
	i += saltLen
	record := buf[i : i+recordLen : i+recordLen]

	// Start generating payload
	hash := sha256.New
	// ikm.
	copy(prkInfo, webpushInfo)
	copy(prkInfo[len(webpushInfo):], dh)
	copy(prkInfo[len(webpushInfo)+len(dh):], localPublicKeyBytes)
	prkHKDF := hkdf.New(hash, sharedECDHSecret, authSecret, prkInfo)
	ikm, err := getHKDFKey(prkHKDF, bufHKDF[0:32:32])
	if err != nil {
		return nil, errors.Join(ErrEncryption, err)
	}

	err = p.genSalt(salt)
	if err != nil {
		return nil, errors.Join(ErrEncryption, err)
	}
	// Derive Content Encryption Key.
	contentHKDF := hkdf.New(hash, ikm, salt, contentEncryptionKeyInfo)
	contentEncryptionKey, err := getHKDFKey(contentHKDF, bufHKDF[32:48:48])
	if err != nil {
		return nil, errors.Join(ErrEncryption, err)
	}
	// Derive the Nonce.
	nonceHKDF := hkdf.New(hash, ikm, salt, nonceInfo)
	nonce, err := getHKDFKey(nonceHKDF, bufHKDF[48:60:60])
	if err != nil {
		return nil, errors.Join(ErrEncryption, err)
	}
	// Cipher.
	c, err := aes.NewCipher(contentEncryptionKey)
	if err != nil {
		return nil, errors.Join(ErrEncryption, err)
	}
	gcm, err := cipher.NewGCMWithTagSize(c, gcmTagSize)
	if err != nil {
		return nil, errors.Join(ErrEncryption, err)
	}

	// Data.
	copy(data, message)
	// End padding.
	data[len(message)] = 2
	// Compose the ciphertext.
	ciphertext := gcm.Seal(data[:0], nonce, data, nil)
	rs := make([]byte, rsLen)
	binary.BigEndian.PutUint32(rs, uint32(len(message)*8))

	// Encryption Content-Coding Header.
	copy(record, salt)
	copy(record[len(salt):], rs)
	w := len(salt) + len(rs)
	record[w] = byte(len(localPublicKeyBytes))
	copy(record[w+1:], localPublicKeyBytes)
	copy(record[w+1+len(localPublicKeyBytes):], ciphertext)

	// SEND REQUEST.
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, sub.Endpoint, bytes.NewReader(record))
	if err != nil {
		return nil, err
	}
	req.Header["Content-Encoding"] = []string{"aes128gcm"}
	req.Header["Content-Length"] = []string{strconv.Itoa(len(ciphertext))}
	req.Header["Content-Type"] = []string{"application/octet-stream"}
	req.Header["TTL"] = []string{strconv.Itoa(options.TTL)}
	if options.Urgency != UrgencyUnset && isValidUrgency(options.Urgency) {
		req.Header["Urgency"] = []string{string(options.Urgency)}
	}
	if options.Topic != "" {
		req.Header["Topic"] = []string{options.Topic}
	}
	req.Header["Authorization"] = []string{keys.vapid}
	return req, nil
}

// GenVAPIDAuthHeader generate the web push vapid auth header.
// Should only be used for debug/logging.
func (p *VAPIDPusher) GenVAPIDAuthHeader(subscriptionEndpoint string) (string, error) {
	keys, err := p.getCachedKeys(subscriptionEndpoint)
	if err != nil {
		return "", err
	}
	return keys.vapid, nil
}

// genSalt generates a salt of 16 bytes.
func (p *VAPIDPusher) genSalt(salt []byte) error {
	_, err := io.ReadFull(p.randReader, salt)
	if err != nil {
		return err
	}
	return nil
}

func (p *VAPIDPusher) decodeBase64(key string) ([]byte, error) {
	b, err := p.base64Encoding.DecodeString(key)
	if err == nil {
		return b, nil
	}
	return decodeBas64Safe(key)
}

// decodeBas64Safe decodes a base64 subscription key.
func decodeBas64Safe(key string) ([]byte, error) {
	b, err := base64.RawURLEncoding.DecodeString(key)
	if err == nil {
		return b, nil
	}
	b, err = base64.URLEncoding.DecodeString(key)
	if err == nil {
		return b, nil
	}
	b, err = base64.RawStdEncoding.DecodeString(key)
	if err == nil {
		return b, nil
	}
	return base64.StdEncoding.DecodeString(key)
}

// Returns a key of length "length" given a hkdf function.
func getHKDFKey(hkdf io.Reader, dst []byte) ([]byte, error) {
	n, err := io.ReadFull(hkdf, dst)
	if n != len(dst) || err != nil {
		return dst, err
	}
	return dst, nil
}
