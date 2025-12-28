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

// WithBase64Encoding allow switching base64 implementation.
// Must be url-safe, no padding encoding.
func WithBase64Encoding(enc Base64Encoding) VAPIDPusherOption {
	return func(pusher *VAPIDPusher) {
		pusher.base64Encoding = enc
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
		pusher.localSecretTTL = exp
	}
}
