package main

import (
	"os"
	"os/signal"
	"syscall"
)

func monitorShutdown(exit chan int) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)
	for _ = range ch {
		exit <- 1
		break
	}
}
