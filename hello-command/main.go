package main

import (
	"log"
	"net"
)

func main() {
	conn, err := net.ListenPacket("udp", ":0")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	dst, err := net.ResolveUDPAddr("udp", "192.168.10.1:8889")
	if err != nil {
		log.Fatal(err)
	}

	// The connection can write data to the desired address.
	_, err = conn.WriteTo([]byte("command"), dst)
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.WriteTo([]byte("takeoff"), dst)
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.WriteTo([]byte("land"), dst)
	if err != nil {
		log.Fatal(err)
	}

}
