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
	"fmt"
	"golang.org/x/crypto/hkdf"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const MaxRecordSize = 4096

var ErrMaxSizeExceeded = errors.New("message too large")
var ErrEncryption = errors.New("encryption error")

var (
	nonceInfo                = []byte("Content-Encoding: nonce\x00")
	contentEncryptionKeyInfo = []byte("Content-Encoding: aes128gcm\x00")
	webpushInfo              = []byte("WebPush: info\x00")
)

// Pre-allocated byte buffer format
//   - [hkdf] authSecret (16)
//   - [hkdf] sharedECDHSecret (32)
//   - [hkdf] prkHKDF (60)
//   - [prk] webpushInfo (14)
//   - [prk] dh (65)
//   - [prk] localPublicKey (65)
//   - [record] salt (16)
//   - [record] rs (4)
//   - [record] localPublicKeyLen (1)
//   - [record] localPublicKey (65)
//   - [record] data
//   - [record] padding delimiter (1)
//   - [record] padding
//   - [record] gcmTag (16)
const (
	authSecretLen       = 16
	sharedECDHSecretLen = 32
	hkdfLen             = 60

	webPushInfoLen    = 14
	p256dhLen         = 65
	localPublicKeyLen = 65

	saltLen   = 16
	rsLen     = 4
	gcmTagLen = 16
	headerLen = saltLen + rsLen + 1 + localPublicKeyLen
)

const (
	authSecretOffset       = 0
	sharedECDHSecretOffset = authSecretLen
	hkdfOffset             = sharedECDHSecretOffset + sharedECDHSecretLen

	prkOffset            = hkdfOffset + hkdfLen
	dhOffset             = prkOffset + webPushInfoLen
	localPublicKeyOffset = dhOffset + p256dhLen

	recordOffset = localPublicKeyOffset + localPublicKeyLen
	rsOffset     = recordOffset + saltLen
	keyOffset    = rsOffset + rsLen
	dataOffset   = keyOffset + 1 + localPublicKeyLen
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
	// thus making the actual expiration time equal to configured expiration.
	// It is recommended to set this field to at least 30 minutes.
	vapidTTLBuffer   time.Duration
	localSecretTTLFn func() time.Duration // Optional, enable local public key and secret reuse.
	randReader       io.Reader
	recordSize       int
	maxRecordSize    int

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
		vapidTTLBuffer: 1 * time.Hour,
		randReader:     rand.Reader,
		maxRecordSize:  MaxRecordSize,
	}
	for _, opt := range options {
		opt(c)
	}

	// Decode the VAPID private key.
	var err error
	c.vapidPrivateKey, err = decodeBase64(vapidPrivateKey)
	if err != nil {
		return nil, err
	}
	// Decode the VAPID public key.
	vapidPublicKeyBytes, err := decodeBase64(vapidPublicKey)
	if err != nil {
		return nil, err
	}
	c.vapidPublicKeyHeaderPart = ", k=" + encodeBase64String(vapidPublicKeyBytes)

	if c.client == nil {
		c.client = &http.Client{
			Timeout: 1 * time.Minute,
		}
	}
	return c, nil
}

// Options are config and extra params needed to send a notification.
type Options struct {
	Topic      string  // Set the Topic header to collapse a pending message.
	TTL        int     // Set the TTL on the endpoint POST request.
	Urgency    Urgency // Set the Urgency header.
	RecordSize int     // Set the target record size for padding.
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

// IsVapidTokenCachingEnabled returns whether the VAPID token caching feature is enabled.
func (p *VAPIDPusher) IsVapidTokenCachingEnabled() bool {
	return p.vapidTokenTTL > 0
}

// IsLocalSecretCachingEnabled returns whether the local secret caching feature is enabled.
func (p *VAPIDPusher) IsLocalSecretCachingEnabled() bool {
	return p.localSecretTTLFn != nil
}

// SendNotification ends a push notification to a subscription's endpoint.
// Message Encryption for Web Push, and VAPID protocols.
// FOR MORE INFORMATION SEE RFC8291: https://datatracker.ietf.org/doc/rfc8291.
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
// The request can then be sent using any http client or [VAPIDPusher.ExecuteRequest].
//
// It is recommended to use [VAPIDPusher.SendNotification] directly instead.
func (p *VAPIDPusher) PrepareNotificationRequest(ctx context.Context, message []byte, sub *Subscription, options Options) (*http.Request, error) {
	// GENERATE VAPID TOKEN AND LOCAL KEYPAIR.
	keys, err := p.getCachedKeys(sub.Endpoint)
	if err != nil {
		return nil, err
	}

	// Pre-alloc everything.
	dataLen := len(message) + 1
	cipherTextLen := dataLen + gcmTagLen
	recordLen := headerLen + cipherTextLen
	if recordLen > p.maxRecordSize {
		return nil, ErrMaxSizeExceeded
	}

	// Calculate padded size.
	recordSize := p.recordSize
	if options.RecordSize > 0 {
		recordSize = options.RecordSize
	}
	if recordLen < recordSize {
		padLen := recordSize - recordLen
		recordLen = recordSize
		cipherTextLen += padLen
		dataLen += padLen
	}
	buf := make([]byte, dataOffset+cipherTextLen)

	// Decode auth and P256dh into a pre allocated buffer.
	authSecret := buf[authSecretOffset : authSecretOffset+authSecretLen : authSecretOffset+authSecretLen]
	dh := buf[dhOffset : dhOffset+p256dhLen : dhOffset+p256dhLen]
	if err := decodeBase64Buff(sub.Keys.Auth, authSecret); err != nil {
		return nil, errors.Join(ErrEncryption, err)
	}
	if err := decodeBase64Buff(sub.Keys.P256dh, dh); err != nil {
		return nil, errors.Join(ErrEncryption, err)
	}

	// GENERATE SHARED AND PUBLIC KEY.
	localPublicKeyBytes := buf[localPublicKeyOffset : localPublicKeyOffset+localPublicKeyLen : localPublicKeyOffset+localPublicKeyLen]
	var sharedECDHSecret []byte

	var now time.Time
	isLocalSecretCacheEnabled := p.localSecretTTLFn != nil && sub.LocalKey != nil
	if isLocalSecretCacheEnabled {
		now = time.Now()
	}
	if isLocalSecretCacheEnabled && sub.LocalKey.At > now.Add(-p.localSecretTTLFn()).UnixMilli() {
		// Use publicKey and secret from LocalKey.
		if err = decodeBase64Buff(sub.LocalKey.Public, localPublicKeyBytes); err != nil {
			return nil, errors.Join(ErrEncryption, err)
		}
		sharedECDHSecret = buf[sharedECDHSecretOffset : sharedECDHSecretOffset+sharedECDHSecretLen : sharedECDHSecretOffset+sharedECDHSecretLen]
		if err = decodeBase64Buff(sub.LocalKey.Secret, sharedECDHSecret); err != nil {
			return nil, errors.Join(ErrEncryption, err)
		}
	} else {
		// We need to copy instead of re-assign, as the localPublicKeyBytes is actually a required part
		// of prkInfo
		copy(localPublicKeyBytes, keys.localPublicKeyBytes)
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
				Public: encodeBase64String(localPublicKeyBytes),
				Secret: encodeBase64String(sharedECDHSecret),
				At:     now.UnixMilli(),
			}
		}
	}

	// GENERATE PAYLOAD.
	bufHKDF := buf[hkdfOffset:prkOffset:prkOffset]
	prkInfo := buf[prkOffset:recordOffset:recordOffset]

	record := buf[recordOffset : recordOffset+recordLen : recordOffset+recordLen]
	salt := buf[recordOffset : recordOffset+saltLen : recordOffset+saltLen]
	rs := buf[rsOffset : rsOffset+rsLen : rsOffset+rsLen]
	// Use data slice for both data and cipher text.
	data := buf[dataOffset : dataOffset+dataLen : dataOffset+cipherTextLen]

	// Start generating payload
	hash := sha256.New
	// ikm.
	copy(prkInfo, webpushInfo)
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
	gcm, err := cipher.NewGCMWithTagSize(c, gcmTagLen)
	if err != nil {
		return nil, errors.Join(ErrEncryption, err)
	}

	// Data.
	copy(data, message)
	// End padding.
	data[len(message)] = 2
	// Compose the ciphertext.
	ciphertext := gcm.Seal(data[:0], nonce, data, nil)
	// From the spec, rs must greater than: plaintext data + padding delimiter + padding + gcmTag,
	// which equal to computed cipherTextLen.
	// Most of the lib I found just use 4096 here, as it is the payload limit.
	binary.BigEndian.PutUint32(rs, MaxRecordSize)

	// Encryption Content-Coding Header.
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

// ExecuteRequest send an [http.Request] using the underlying client,
// usually the request prepared by [VAPIDPusher.PrepareNotificationRequest].
// Useful when you want to measure the request preparation and the request execution time separately.
//
// It is recommended to use [VAPIDPusher.SendNotification] directly instead.
func (p *VAPIDPusher) ExecuteRequest(req *http.Request) (*http.Response, error) {
	return p.client.Do(req)
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

func decodeBase64(key string) ([]byte, error) {
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

func decodeBase64Buff(key string, buff []byte) error {
	expectedLen := cap(buff)
	src := []byte(key)
	if base64.RawURLEncoding.DecodedLen(len(src)) == expectedLen {
		n, err := base64.RawURLEncoding.Decode(buff, src)
		if err == nil && n == expectedLen {
			return nil
		}
	}
	if base64.URLEncoding.DecodedLen(len(src)) == expectedLen {
		n, err := base64.URLEncoding.Decode(buff, src)
		if err == nil && n == expectedLen {
			return nil
		}
	}
	if base64.RawStdEncoding.DecodedLen(len(src)) == expectedLen {
		n, err := base64.RawStdEncoding.Decode(buff, src)
		if err == nil && n == expectedLen {
			return nil
		}
	}
	if base64.StdEncoding.DecodedLen(len(src)) == expectedLen {
		n, err := base64.StdEncoding.Decode(buff, src)
		if err == nil && n == expectedLen {
			return nil
		}
	}
	return fmt.Errorf("invalid base64 data length")
}

func encodeBase64String(src []byte) string {
	return base64.RawURLEncoding.EncodeToString(src)
}

// getHKDFKey Returns a key of length "length" given a hkdf function.
func getHKDFKey(hkdf io.Reader, dst []byte) ([]byte, error) {
	n, err := io.ReadFull(hkdf, dst)
	if n != len(dst) || err != nil {
		return dst, err
	}
	return dst, nil
}
