# Example

All commands below should be run from this directory (example).

## Start the server

You can supply VAPID keypair using `-vapid=<priv>:<pub>`, or specify in `.vapid` file in current working
directory.

```bash
go run cmd/server/main.go
```

This command will generate a `.vapid` file in the current working directory if not supplied.

Go to `http://localhost:8080`, allow notification permission and copy the logged subscription from the console.

## Test send a notification

Create a `.subscription.json` file using the subscription you had from the first section (the server should
automatically generate this file already after you subscribed)

```bash
go run cmd/send/main.go
```
