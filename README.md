# Go Webpush

Web Push API Encryption with VAPID. Sacrifice security and features for raw speed.

Inspired by https://github.com/SherClockHolmes/webpush-go

## Installation

Require go 1.25+

```shell
go get -u github.com/mawngo/go-fwebpush
```

## Optimization

The library applies various optimizations to improve performance, including:

- Reducing allocation (pre-allocating buffer, string joining, ...)
- Avoid header canonicalization (which, in turn, further reducing allocation)
- Optimized encoding
- Caching keys and token (see bellow)

### Optional Optimization

- `WithVAPIDTokenTTL` **(Enabled by default)** Caching jwt token + local public key and curve, improving performance
  by ~2.5x
- `WithLocalSecretTTL` Reuse the local public key and secret if available, or else generate new ones and set
  to `Subscription`, user can save those keys for reuse later. Improving performance by 1.5x if enabled alone, and 15x
  if enabled along with the `WithVAPIDTokenTTL` above (this huge different is due to some design decision that favor
  simplicity in implementation).
- `WithBase64Encoding` Allow user to opt in more performance base64 implementation, for
  example: https://github.com/cristalhq/base64. However, be aware that benchmark shows no difference when switched to
  that
  library.
- `WithRandReader` Allow user to opt in faster rand.Reader.

### Generating VAPID Keys

Use the helper method `GenerateVAPIDKeys` to generate the VAPID key pair.

```golang
privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
if err != nil {
// TODO: Handle error
}
```