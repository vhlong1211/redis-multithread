package main

import (
	"TCPServer/IO_MultiplexingServer"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	var signals = make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	var wg sync.WaitGroup
	wg.Add(2)

	go IO_MultiplexingServer.RunIoMultiplexingServer(&wg)
	go IO_MultiplexingServer.WaitForSignal(&wg, signals)
	wg.Wait()
}
