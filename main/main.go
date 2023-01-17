package main

import (
	"AOI/kit"
	"AOI/quadtree"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	quadtree.Start()
	quadtree.StartEnterTest()

	kit.CvsStart()

	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, syscall.SIGINT, syscall.SIGTERM)

	select {
	case _ = <-chSignal:
		signal.Reset(syscall.SIGINT, syscall.SIGTERM)
	}
}
