package main

import (
	"github.com/mawngo/go-fwebpush"
	"os"
)

func main() {
	privateKey, publicKey, err := fwebpush.GenerateVAPIDKeys()
	if err != nil {
		panic(err)
	}
	println(publicKey)
	println(privateKey)

	fi, err := os.Create(".vapid.txt")
	if err != nil {
		println("Error creating file:", err)
	}
	defer fi.Close()
	_, err = fi.WriteString(privateKey + ":" + publicKey)
	if err != nil {
		println("Error writing file:", err)
	}
}
