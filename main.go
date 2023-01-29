package main

import (
	quadtree2 "AOI/internal/aoi"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	quadtree2.Start()
	quadtree2.StartEnterTest()
	//kit.TestPlayer(2)
	quadtree2.CvsStart()


	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, syscall.SIGINT, syscall.SIGTERM)

	select {
	case _ = <-chSignal:
		signal.Reset(syscall.SIGINT, syscall.SIGTERM)
	}
}
