package main

import "github.com/mawngo/go-fwebpush"

func main() {
	privateKey, publicKey, err := fwebpush.GenerateVAPIDKeys()
	if err != nil {
		panic(err)
	}
	println(publicKey)
	println(privateKey)
}
