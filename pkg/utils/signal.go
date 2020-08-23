package utils

import (
	"os"
	"os/signal"
	"syscall"
)

func HandleSignals(handler func(os.Signal)) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		var sig = <-signals
		handler(sig)
	}()
}

