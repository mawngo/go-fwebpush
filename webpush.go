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
	"strings"
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
// Key buffer:
//   - [hkdf] authSecret (16)
//   - [hkdf] sharedECDHSecret (32)
//   - [hkdf] prkHKDF (60)
//   - [prk] webpushInfo (14)
//   - [prk] dh (65)
//   - [prk] localPublicKey (65)
//
// Record buffer:
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
	sharedECDHSecretOffset = authSecretLen
	hkdfOffset             = sharedECDHSecretOffset + sharedECDHSecretLen

	prkOffset          = hkdfOffset + hkdfLen
	dhOffset           = prkOffset + webPushInfoLen
	prkPublicKeyOffset = dhOffset + p256dhLen

	rsOffset             = saltLen
	keyOffset            = rsOffset + rsLen
	localPublicKeyOffset = keyOffset + 1
	dataOffset           = localPublicKeyOffset + localPublicKeyLen
)

type VAPIDPusher struct {
	client                   *http.Client
	subject                  string        // Sub in VAPID JWT token.
	vapidPublicKeyHeaderPart string        // VAPID public key passed in the VAPID Authorization header (format: `, k=<key`).
	vapidPrivateKey          []byte        // VAPID private key, used to sign VAPID JWT token.
	vapidTokenTTL            time.Duration // Optional, expiration for VAPID JWT token.
	vapidTTLBuffer           time.Duration
	localSecretTTLFn         func() time.Duration // Optional, enable reuse of the local public key and secret.
	randReader               io.Reader
	recordSize               int
	maxRecordSize            int
	keyBufPool               sync.Pool

	mu    sync.RWMutex
	cache map[string]reusableKey // Cache of VAPID JWT token by audience.
}

func NewVAPIDPusher(
	subject string,
	vapidPublicKey string,
	vapidPrivateKey string,
	options ...VAPIDPusherOption,
) (*VAPIDPusher, error) {
	c := &VAPIDPusher{
		vapidTokenTTL:  1 * time.Hour,
		cache:          make(map[string]reusableKey),
		vapidTTLBuffer: 10 * time.Minute,
		randReader:     rand.Reader,
		maxRecordSize:  MaxRecordSize,
	}
	for _, opt := range options {
		opt(c)
	}

	if c.vapidTokenTTL+c.vapidTTLBuffer > 24*time.Hour {
		return nil, errors.New("total VAPID token must be less than 24 hours")
	}

	if !strings.HasPrefix(subject, "mailto:") && !strings.HasPrefix(subject, "https:") {
		subject = "mailto:" + subject
	}
	c.subject = subject

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
	// Deprecated: switched to IKM caching.
	Secret string `json:"s,omitempty"`
	// IKM generated ikm.
	IKM string `json:"m,omitempty"`
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
	now := time.Now()
	// GENERATE VAPID TOKEN AND LOCAL KEYPAIR.
	keys, err := p.getCachedKeys(sub.Endpoint, now)
	if err != nil {
		return nil, err
	}

	// Pre-alloc for record.
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
	record := make([]byte, recordLen)

	// Pre-alloc for keys.
	// This buffer can be pooled, reduce allocations.
	// Not sure about the performance impact through.
	keyBuf := make([]byte, prkPublicKeyOffset+localPublicKeyLen)
	hash := sha256.New

	// GENERATE IKM AND PUBLIC KEY.
	localPublicKeyBytes := record[localPublicKeyOffset : localPublicKeyOffset+localPublicKeyLen : localPublicKeyOffset+localPublicKeyLen]
	ikm := keyBuf[hkdfOffset : hkdfOffset+32 : hkdfOffset+32]
	if p.localSecretTTLFn != nil && sub.LocalKey != nil && sub.LocalKey.At > now.Add(-p.localSecretTTLFn()).UnixMilli() && sub.LocalKey.IKM != "" {
		// Use publicKey and ikm from LocalKey.
		if err = decodeBase64Buff(sub.LocalKey.Public, localPublicKeyBytes); err != nil {
			return nil, errors.Join(ErrEncryption, err)
		}
		if err = decodeBase64Buff(sub.LocalKey.IKM, ikm); err != nil {
			return nil, errors.Join(ErrEncryption, err)
		}
	} else {
		// We need to copy instead of re-assign, as the localPublicKeyBytes is actually a required part
		// of the record.
		copy(localPublicKeyBytes, keys.localPublicKeyBytes)
		// Derive ECDH shared secret.
		// Decode auth and P256dh into a pre allocated buffer.
		authSecret := keyBuf[:authSecretLen:authSecretLen]
		dh := keyBuf[dhOffset : dhOffset+p256dhLen : dhOffset+p256dhLen]
		if err := decodeBase64Buff(sub.Keys.Auth, authSecret); err != nil {
			return nil, errors.Join(ErrEncryption, err)
		}
		if err := decodeBase64Buff(sub.Keys.P256dh, dh); err != nil {
			return nil, errors.Join(ErrEncryption, err)
		}
		dhPublicKey, err := keys.curve.NewPublicKey(dh)
		if err != nil {
			return nil, errors.Join(ErrEncryption, err)
		}
		sharedECDHSecret, err := keys.localPrivateKey.ECDH(dhPublicKey)
		if err != nil {
			return nil, errors.Join(ErrEncryption, err)
		}

		// ikm.
		copy(keyBuf[prkOffset:prkOffset+webPushInfoLen:prkOffset+webPushInfoLen], webpushInfo)
		copy(keyBuf[prkPublicKeyOffset:prkPublicKeyOffset+localPublicKeyLen:prkPublicKeyOffset+localPublicKeyLen], localPublicKeyBytes)
		prkInfo := keyBuf[prkOffset:]
		prkHKDF := hkdf.New(hash, sharedECDHSecret, authSecret, prkInfo)
		ikm, err = getHKDFKey(prkHKDF, ikm)
		if err != nil {
			return nil, errors.Join(ErrEncryption, err)
		}

		// Update LocalKey if enabled.
		if p.localSecretTTLFn != nil {
			sub.LocalKey = &LocalKey{
				Public: encodeBase64String(localPublicKeyBytes),
				IKM:    encodeBase64String(ikm),
				At:     now.UnixMilli(),
			}
		}
	}

	// GENERATE PAYLOAD.
	bufHKDF := keyBuf[hkdfOffset:prkOffset:prkOffset]
	salt := record[:saltLen:saltLen]
	rs := record[rsOffset : rsOffset+rsLen : rsOffset+rsLen]
	// Use data slice for both data and cipher text.
	data := record[dataOffset : dataOffset+dataLen : dataOffset+cipherTextLen]

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
	record[keyOffset] = byte(len(localPublicKeyBytes))
	copy(record[dataOffset:], ciphertext)

	// PREPARE REQUEST.
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
	keys, err := p.getCachedKeys(subscriptionEndpoint, time.Now())
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

func decodeBase64Buff(key string, buff []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("base64 data overflow")
		}
	}()
	expectedLen := cap(buff)
	src := []byte(key)
	if n, err := base64.RawURLEncoding.Decode(buff, src); err == nil && n == expectedLen {
		return nil
	}
	if n, err := base64.URLEncoding.Decode(buff, src); err == nil && n == expectedLen {
		return nil
	}
	if n, err := base64.RawStdEncoding.Decode(buff, src); err == nil && n == expectedLen {
		return nil
	}
	if n, err := base64.StdEncoding.Decode(buff, src); err == nil && n == expectedLen {
		return nil
	}
	return fmt.Errorf("invalid base64 data length")
}

func encodeBase64String(src []byte) string {
	return base64.RawURLEncoding.EncodeToString(src)
}

// getHKDFKey Returns a key of length "length" given a hkdf function.
func getHKDFKey(hkdf io.Reader, dst []byte) ([]byte, error) {
	n, err := io.ReadFull(hkdf, dst)
	if err != nil {
		return nil, err
	}
	if n != len(dst) {
		return nil, fmt.Errorf("hkdf key length mismatch")
	}
	return dst, nil
}
