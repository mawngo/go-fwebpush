# Example

All commands bellow should be run from this directory (example).

## Access index.html

Start the server:

You can supply VAPID keypair using `-vapid=<priv>:<pub>`, or specify in `.vapid` file in current working
directory.

```bash
go run cmd/server/main.go
```

Go to `http://localhost:8080` and copy the logged subscription and generated VAPID key from the console.

## Test send a notification

Create a `.subscription.json` file using the subscription you had from the first section (the first command should
automatically generate this file already after you subscribed)

You can supply VAPID keypair using `-vapid=<priv>:<pub>`, or specify in `.vapid` file in current working
directory (the first command should automatically generate this file already)

```bash
go run cmd/send/main.go
```
