# Go Webpush

Web Push API Encryption with VAPID. Sacrifice security and features for raw speed.

Inspired by https://github.com/SherClockHolmes/webpush-go

## Installation

Require go 1.27+

```shell
go get -u github.com/mawngo/go-fwebpush/v2
```

## Optimization

The library applies various optimizations to improve performance, including:

- Reducing allocation (pre-allocating buffer, string joining, ...)
- Optimized encoding
- Caching keys and token (see bellow)
- No padding by default

### Optional Optimization

- `WithVAPIDTokenTTL` **(Enabled by default)** Caching jwt token + local public key and curve. When preparing requests,
  this option improving performance by ~2x
- `WithLocalSecretTTL` Reuse the local public key and ikm if available, or else generate new ones and set
  to `Subscription`, user can save those keys for reuse later. When preparing requests, this option improving 
  performance by ~1.5x if enabled alone and ~14x if enabled along with the `WithVAPIDTokenTTL` above (this huge 
  difference is due to some design decisions that favor simplicity in implementation).

## Usage

For a full example, refer to the code in the [example](example/) directory.

```go
package main

import (
  "context"
  "encoding/json"
  "github.com/mawngo/go-fwebpush/v2"
)

func main() {
  // Decode subscription.
  s := fwebpush.Subscription{}
  err := json.Unmarshal([]byte("<YOUR_SUBSCRIPTION>"), &s)
  if err != nil {
    panic(err)
  }

  pusher, err := fwebpush.NewVAPIDPusher(
    "example@example.com",
    "<YOUR_VAPID_PUBLIC_KEY>",
    "<YOUR_VAPID_PRIVATE_KEY>",
  )

  if err != nil {
    panic(err)
  }

  // Send Notification.
  resp, err := pusher.SendNotificationOptions(
    context.Background(),
    []byte("Test"),
    &s,
    fwebpush.Options{TTL: 30},
  )
  if err != nil {
    // TODO: Handle error
  }
  defer resp.Body.Close()
}
```

### Generating VAPID Keys

Use the helper method `GenerateVAPIDKeys` to generate the VAPID key pair.

```golang
privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
if err != nil {
// TODO: Handle error
}
```

### Dependencies

This library only depends on `golang.org/x/crypto`.