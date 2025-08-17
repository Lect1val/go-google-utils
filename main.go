package main

import (
	"log"

	"github.com/Lect1val/go-google-utils/auth"
)

func main() {
	if err := auth.GenerateTokenInteractive("token.json"); err != nil {
		log.Fatal(err)
	}
}
