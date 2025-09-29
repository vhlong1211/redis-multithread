package main

import (
	"log"
	"net"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listening at port 3000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("handle conn from port  = ", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	//receive req and reply response
	for {
		var buf []byte = make([]byte, 1000)
		_, err := conn.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		//pretend to process req
		time.Sleep(time.Second * 1)
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nHello, world\r\n"))
	}

	//conn.Close()
}
