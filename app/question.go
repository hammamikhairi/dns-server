package main

import (
	"bytes"
	"encoding/binary"
)

type Question struct {
	name   string
	qType  uint16
	qClass uint16
}

func NewQuestion(b *bytes.Buffer) (*Question, error) {

	question := &Question{}
	err := question.readQuestion(b)

	return question, err
}

func (q *Question) readQuestion(b *bytes.Buffer) error {

	// start with the name
	domainName, err := readDomainName(b)
	if err != nil {
		return err
	}
	q.name = domainName

	qType := uint16(0)
	err = binary.Read(b, binary.BigEndian, &qType)
	if err != nil {
		return err
	}

	qClass := uint16(0)
	err = binary.Read(b, binary.BigEndian, &qClass)
	if err != nil {
		return err
	}

	q.qType, q.qClass = qType, qClass

	return nil
}

func (q *Question) pack() ([]byte, error) {

	buff := new(bytes.Buffer)

	labels := bytes.Split([]byte(q.name), []byte("."))
	for _, label := range labels {
		buff.WriteByte(byte(len(label)))
		buff.Write(label)
	}

	buff.WriteByte(0)

	// Write the type
	if err := binary.Write(buff, binary.BigEndian, q.qType); err != nil {
		return nil, err
	}

	// Write the class
	if err := binary.Write(buff, binary.BigEndian, q.qClass); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func readDomainName(b *bytes.Buffer) (string, error) {
	domain := new(bytes.Buffer)

	for {
		length, err := b.ReadByte()
		if err != nil {
			break
		}

		if length == 0 {
			break
		}

		label := make([]byte, length)
		_, err = b.Read(label)
		if err != nil {
			return "", err
		}

		if domain.Len() > 0 {
			domain.WriteByte('.')
		}

		domain.Write(label)
	}

	return domain.String(), nil
}
