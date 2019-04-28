package main

import (
	"os"
	"os/signal"
)

func onTerminate(callback func()) {
	// Setup termination signal listener
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	go func() {
		<-stop
		callback()
	}()
}
