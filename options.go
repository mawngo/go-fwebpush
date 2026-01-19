package fwebpush

import (
	"io"
	"net/http"
	"time"
)

// VAPIDPusherOption modify VAPIDPusher configs.
type VAPIDPusherOption = func(pusher *VAPIDPusher)

// WithClient set the client for VAPIDPusher.
func WithClient(client *http.Client) VAPIDPusherOption {
	return func(pusher *VAPIDPusher) {
		pusher.client = client
	}
}

// WithVAPIDTokenTTL configure vapid token caching and expiration.
// Token cache invalidation will be set to exp.
func WithVAPIDTokenTTL(exp time.Duration) VAPIDPusherOption {
	return func(pusher *VAPIDPusher) {
		pusher.vapidTokenTTL = exp
	}
}

// WithVAPIDTokenTTLExt additional duration added to expiration.
// The key will expire later than configured expiration this amount of duration,
// while the validation of the key will expire sooner than configured expiration this amount of duration,
// thus making the actual expiration time equal to configured expiration.
//
// Default 30 minutes.
func WithVAPIDTokenTTLExt(exp time.Duration) VAPIDPusherOption {
	return func(pusher *VAPIDPusher) {
		pusher.vapidTTLBuffer = exp
	}
}

// WithRandReader allow switching randReader implementation.
func WithRandReader(rand io.Reader) VAPIDPusherOption {
	return func(pusher *VAPIDPusher) {
		pusher.randReader = rand
	}
}

// WithLocalSecretTTL configure reusing of the local secret and public key.
// Set to 0 to disable.
// When enabled, the pusher will check the LocalKey of the Subscription and generate if not have one or expired.
// You can save the generated LocalKey with the Subscription to reuse later.
func WithLocalSecretTTL(exp time.Duration) VAPIDPusherOption {
	return func(pusher *VAPIDPusher) {
		pusher.localSecretTTLFn = func() time.Duration {
			return exp
		}
	}
}

// WithLocalSecretTTLFn configure reusing of the local secret and public key.
// Set to nil to disable.
//
// See [WithLocalSecretTTL].
func WithLocalSecretTTLFn(fn func() time.Duration) VAPIDPusherOption {
	return func(pusher *VAPIDPusher) {
		pusher.localSecretTTLFn = fn
	}
}

// WithRecordSize configure padding of the message payload.
// Payload that has length exceed the configured size will not be padded.
// The maximum accepted value is [MaxRecordSize].
// The default value is 0 (disabled).
func WithRecordSize(size int) VAPIDPusherOption {
	return func(pusher *VAPIDPusher) {
		pusher.recordSize = min(max(size, 0), MaxRecordSize)
	}
}

// WithMaxRecordSize configure the maximum message payload size.
// If the payload exceeds this size, the push will fail with [ErrMaxSizeExceeded]
// The default value is [MaxRecordSize].
// The minimum accepted value is 103, which includes: Absent header (86 octets), padding (minimum 1 octet),
// and expansion for AEAD_AES_128_GCM (16 octets)
//
// Can be set to 0 to disable max record size validation.
func WithMaxRecordSize(size int) VAPIDPusherOption {
	return func(pusher *VAPIDPusher) {
		if size <= 0 {
			pusher.maxRecordSize = 0
			return
		}
		pusher.maxRecordSize = min(max(size, 103), MaxRecordSize)
	}
}
