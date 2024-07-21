package main

import (
	"bytes"
	"io"
	"strings"
)

type DSNMessage struct {
	header    *Header
	questions []*Question
	answers   []*Answer
}

func BuildMessage(b *bytes.Buffer) *DSNMessage {

	header, err := NewHeader(b)
	if err != nil {
		header.rescode = 255
	} else {
		if header.opcode == 0 {
			header.rescode = 0
		} else {
			header.rescode = 4
		}
	}

	questions := []*Question{}
	answers := []*Answer{}

	for i := 0; i < int(header.qdcount); i++ {
		question, err := NewQuestion(b)

		if err != nil {

			if err == io.EOF {
				first := questions[0]
				sg := strings.Split(first.name, ".")
				sg[0] = ""
				question.name += strings.Join(sg, ".")[1:]
				question.qClass = first.qClass
				question.qType = first.qType

			} else {
				header.rescode = 255
			}
		}

		questions = append(questions, question)
	}

	for _, question := range questions {

		answer, err := NewAnswer(question.name, "8.8.8.8")
		if err != nil {
			header.rescode = 255
		}
		answers = append(answers, answer)
	}

	header.qdcount = uint16(len(questions))
	header.answers = uint16(len(answers))

	return &DSNMessage{header, questions, answers}
}
func BuildMessageFrom(b *bytes.Buffer) *DSNMessage {

	header, err := NewHeader(b)
	if err != nil {
		header.rescode = 255
	} else {
		if header.opcode == 0 {
			header.rescode = 0
		} else {
			header.rescode = 4
		}
	}

	questions := []*Question{}
	answers := []*Answer{}

	for i := 0; i < int(header.qdcount); i++ {
		question, err := NewQuestion(b)
		if err != nil {
			header.rescode = 255
		}
		questions = append(questions, question)
	}

	for _, question := range questions {

		answer, err := NewAnswer(question.name, "8.8.8.8")

		if err != nil {
			header.rescode = 255
		}

		err = answer.build(b)

		if err != nil {
			header.rescode = 255
		}

		answers = append(answers, answer)
	}

	header.qdcount = uint16(len(questions))
	header.answers = uint16(len(answers))

	return &DSNMessage{header, questions, answers}
}

func (ms *DSNMessage) pack() ([]byte, error) {

	questionsStream := []byte{}
	answersStream := []byte{}

	for _, question := range ms.questions {
		packedQuestion, err := question.pack()
		if err != nil {
			return nil, err
		}

		questionsStream = append(questionsStream, packedQuestion...)
	}

	for _, answer := range ms.answers {

		packedAnswer, err := answer.pack()
		if err != nil {
			return nil, err
		}
		answersStream = append(answersStream, packedAnswer...)

	}

	response, err := ms.header.pack()
	if err != nil {
		return nil, err
	}

	response = append(response, questionsStream...)
	response = append(response, answersStream...)

	return response, err
}
