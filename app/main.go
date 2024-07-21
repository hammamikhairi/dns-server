package main

import (
	"bytes"
	"flag"
	"fmt"
	"net"
)

var addy *string
var port *string

func init() {
	addy = flag.String("resolver", "127.0.0.1", "IP address to send the UDP request to")
	port = flag.String("port", "2053", "IP address to send the UDP request to")
	flag.Parse()
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:"+*port)
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	buf := make([]byte, 512)
	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		receivedData := string(buf[:size])
		b := bytes.NewBufferString(receivedData)

		message := BuildMessage(b)

		response := []byte{}
		if *addy != "127.0.0.1" {
			response, err = resolve(*addy, message)
		} else {
			response, err = message.pack()
		}

		if err != nil {

			panic(err)
		}

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
