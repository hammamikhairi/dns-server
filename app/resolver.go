package main

import (
	"bytes"
	"fmt"
	"net"
	"time"
)

func resolve(address string, message *DSNMessage) ([]byte, error) {

	message.answers = []*Answer{}

	for i, question := range message.questions {

		msg := DSNMessage{
			header:    message.header.copy(),
			questions: []*Question{question},
			answers:   []*Answer{},
		}

		msg.header.qdcount = 1
		msg.header.answers = 0
		msg.header.opcode = 0
		msg.header.query = false

		packed, _ := msg.pack()
		response, err := UDPRequest(address, packed)
		if err != nil {
			return nil, err
		}

		buff := bytes.NewBuffer(response)
		dnsMsg := BuildMessageFrom(buff)
		if i > 0 {
			question.name = question.name[3:]
		}

		message.answers = append(message.answers, dnsMsg.answers...)
	}

	message.header.query = true
	message.header.answers = uint16(len(message.questions))
	message.header.qdcount = uint16(len(message.questions))

	return message.pack()
}

func UDPRequest(address string, data []byte) ([]byte, error) {
	// Resolve the UDP address
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, fmt.Errorf("error resolving address: %w", err)
	}

	// Create a UDP connection
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("error dialing UDP address: %w", err)
	}
	defer conn.Close()

	// Send the data
	_, err = conn.Write(data)
	if err != nil {
		return nil, fmt.Errorf("error writing data to UDP connection: %w", err)
	}

	// Set a read deadline for the response
	conn.SetReadDeadline(time.Now().Add(30000 * time.Second))

	// Read the response
	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return nil, fmt.Errorf("error reading response from UDP connection: %w", err)
	}

	// Return the response
	return buffer[:n], nil
}
