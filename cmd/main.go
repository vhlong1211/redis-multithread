package main

import (
	"TCPServer/IO_MultiplexingServer"
	"TCPServer/Multithread_server"
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

	//This is single thread
	//go IO_MultiplexingServer.RunIoMultiplexingServer(&wg)

	//This is multithread
	s := Multithread_server.NewServer()
	go s.StartSingleListener(&wg)
	//go s.StartMultiListeners(&wg)
	go IO_MultiplexingServer.WaitForSignal(&wg, signals)
	wg.Wait()
}
