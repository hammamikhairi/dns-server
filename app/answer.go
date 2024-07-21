package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
)

type Answer struct {
	name  string
	aType uint16
	class uint16
	ttl   uint32
	rDLen uint16
	rData string
}

func NewAnswer(name string, value string) (*Answer, error) {

	answer := &Answer{
		name:  name,
		aType: 1,
		class: 1,
		ttl:   60,
		rDLen: 4,
		rData: value,
	}

	return answer, nil
}

func (a *Answer) build(b *bytes.Buffer) error {

	domainName, err := readDomainName(b)
	if err != nil {
		return err
	}
	a.name = domainName

	var (
		aType, class uint16
		ttl          uint32
		// rDLen        uint16
	)

	err = binary.Read(b, binary.BigEndian, &aType)
	if err != nil {
		return err
	}
	err = binary.Read(b, binary.BigEndian, &class)
	if err != nil {
		return err
	}
	err = binary.Read(b, binary.BigEndian, &ttl)
	if err != nil {
		return err
	}

	// err = binary.Read(b, binary.BigEndian, &rDLen)
	// if err != nil {
	// 	return err
	// }

	b.Next(2)

	var segments []string
	for i := 0; i < 4; i++ {
		b, err := b.ReadByte()
		if err != nil {
			return err
		}
		segments = append(segments, strconv.Itoa(int(b)))
	}

	a.aType = aType
	a.class = class
	a.ttl = ttl
	a.rDLen = 4
	a.rData = fmt.Sprintf("%s.%s.%s.%s", segments[0], segments[1], segments[2], segments[3])
	return nil
}

func (a *Answer) pack() ([]byte, error) {

	buff := new(bytes.Buffer)

	labels := bytes.Split([]byte(a.name), []byte("."))
	for _, label := range labels {
		buff.WriteByte(byte(len(label)))
		buff.Write(label)
	}

	buff.WriteByte(0)

	if err := binary.Write(
		buff,
		binary.BigEndian,
		a.aType,
	); err != nil {
		return nil, err
	}

	if err := binary.Write(
		buff,
		binary.BigEndian,
		a.class,
	); err != nil {
		return nil, err
	}
	if err := binary.Write(
		buff,
		binary.BigEndian,
		a.ttl,
	); err != nil {
		return nil, err
	}
	if err := binary.Write(
		buff,
		binary.BigEndian,
		a.rDLen,
	); err != nil {
		return nil, err
	}

	segments := bytes.Split([]byte(a.rData), []byte("."))
	for _, segment := range segments {
		seg, _ := strconv.Atoi(string(segment))
		bb := byte(seg)
		buff.WriteByte(bb)
	}

	buff.WriteByte(0)

	return buff.Bytes(), nil
}
