package common

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Calls the `.Close` method on parameter when there is SIGINT or SIGTERM signal
func Teardown(service Service) {

	cleanupSignal := make(chan os.Signal, 1)
	signal.Notify(cleanupSignal, syscall.SIGINT, syscall.SIGTERM)

	for range cleanupSignal {

		service.Close()

	}

}

func LogFatalOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
