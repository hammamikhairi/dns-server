package main

import (
	"bytes"
	"encoding/binary"
)

const (
	HEADER_SIZE = 12 // 12 bytes
)

type Header struct {
	id                   uint16
	query                bool
	opcode               uint8
	authoritative_answer bool
	truncation           bool
	recursion_desired    bool
	recursion_available  bool
	z                    uint8
	rescode              uint8

	qdcount               uint16
	answers               uint16
	authoritative_entries uint16
	resource_entries      uint16
}

func NewHeader(b *bytes.Buffer) (*Header, error) {
	header := &Header{}
	err := header.readHeader(b)
	if err != nil {
		return nil, err
	}
	return header, nil
}

func (d *Header) copy() *Header {
	return &Header{
		d.id,
		d.query,
		d.opcode,
		d.authoritative_answer,
		d.truncation,
		d.recursion_desired,
		d.recursion_available,
		d.z,
		d.rescode,

		d.qdcount,
		d.answers,
		d.authoritative_entries,
		d.resource_entries,
	}
}

func (d *Header) readHeader(b *bytes.Buffer) error {
	if err := binary.Read(b, binary.BigEndian, &d.id); err != nil {
		return err
	}

	var flags uint16
	if err := binary.Read(b, binary.BigEndian, &flags); err != nil {
		return err
	}
	d.query = (flags & 0x8000) == 0
	d.opcode = uint8((flags >> 11) & 0xF)
	d.authoritative_answer = (flags & 0x0400) != 0
	d.truncation = (flags & 0x0200) != 0
	d.recursion_desired = (flags & 0x0100) != 0
	d.recursion_available = (flags & 0x0080) != 0
	d.z = uint8((flags >> 4) & 0x7)
	d.rescode = uint8(flags & 0xF)

	if err := binary.Read(b, binary.BigEndian, &d.qdcount); err != nil {
		return err
	}
	if err := binary.Read(b, binary.BigEndian, &d.answers); err != nil {
		return err
	}
	if err := binary.Read(b, binary.BigEndian, &d.authoritative_entries); err != nil {
		return err
	}
	if err := binary.Read(b, binary.BigEndian, &d.resource_entries); err != nil {
		return err
	}
	return nil
}

// thanks @XyaelD <3
func (h *Header) pack() ([]byte, error) {

	buff := new(bytes.Buffer)

	if err := binary.Write(buff, binary.BigEndian, &h.id); err != nil {
		return nil, err
	}

	headerFlags := uint16(0) // 0000-0000-0000-0000

	if h.query {
		headerFlags |= 0x8000 // filp 14th bit --> 1000-0000-0000-0000
	}

	headerFlags |= uint16(h.opcode) << 11

	if h.authoritative_answer {
		headerFlags |= 0x0400
	}

	if h.truncation {
		headerFlags |= 0x0200
	}

	if h.recursion_desired {
		headerFlags |= 0x0100
	}

	if h.recursion_available {
		headerFlags |= 0x0080
	}

	headerFlags |= uint16(h.z) << 4
	headerFlags |= uint16(h.rescode)

	if err := binary.Write(buff, binary.BigEndian, headerFlags); err != nil {
		return nil, err
	}
	if err := binary.Write(buff, binary.BigEndian, &h.qdcount); err != nil {
		return nil, err
	}
	if err := binary.Write(buff, binary.BigEndian, &h.answers); err != nil {
		return nil, err
	}
	if err := binary.Write(buff, binary.BigEndian, &h.authoritative_entries); err != nil {
		return nil, err
	}
	if err := binary.Write(buff, binary.BigEndian, &h.resource_entries); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
