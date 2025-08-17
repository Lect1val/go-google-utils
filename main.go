package main

import (
	"log"

	"github.com/Lect1val/go-google-utils/email"
	"go.uber.org/zap"
)

func main() {
	// if err := auth.GenerateTokenInteractive("token.json"); err != nil {
	// 	log.Fatal(err)
	// }

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	es := email.NewEmailService(logger)
	if err := es.SendIndividualEmail("mock@example.com", "Test", "text/plain", "This is a test message"); err != nil {
		log.Fatal(err)
	}
}
