package main

import (
	"log"

	"github.com/Lect1val/go-google-utils/auth"
)

func main() {
	// if err := auth.GenerateTokenInteractive("token.json"); err != nil {
	// 	log.Fatal(err)
	// }

	if err := auth.GenerateCalendarTokenInteractive("calendar-token.json"); err != nil {
		log.Fatal(err)
	}

	// logger, _ := zap.NewProduction()
	// defer logger.Sync()

	// es := email.NewEmailService(logger)

	// if err := es.SendIndividualEmail("lectival857@gmail.com", "Test", "text/plain", "This is a test message"); err != nil {
	// 	log.Fatal(err)
	// }
}
