package message

import (
	"bytes"
	"encoding/binary"
	"errors"
)

const (
	MSG_HEADER_SIZE     = 16
	MSG_HEADER_PROLOGUE = 0x4853 // "HS"
	MSG_MAX_SIZE        = 128
)

type MessageType uint8

const (
	Initialize MessageType = iota
	InitializeResponse
	FatalError
	Error
	AsyncLock
	AsyncLockResponse
	Data
	DataEnd
	DeviceClearComplete
	DeviceClearAcknowledge
	AsyncRemoteLocalControl
	AsyncRemoteLocalResponse
	Trigger
	Interrupted
	AsyncInterrupted
	AsyncMaximumMessageSize
	AsyncMaximumMessageSizeResponse
	AsyncInitialize
	AsyncInitializeResponse
	AsyncDeviceClear
	AsyncServiceRequest
	AsyncStatusQuery
	AsyncStatusResponse
	AsyncDeviceClearAcknowledge
	AsyncLockInfo
	AsyncLockInfoResponse
)

type Header struct {
	Prologue     uint16
	Message_type MessageType
	Control_code uint8 // TOOD: specific type
	Parameter    uint32
}

type Message struct {
	Header  *Header
	Payload []byte
}

func NewMessage(message_type MessageType, control_code uint8, parameter uint32, payload []byte) *Message {
	msg := &Message{}
	msg.Header = &Header{}
	msg.Header.Prologue = MSG_HEADER_PROLOGUE
	msg.Header.Message_type = message_type
	msg.Header.Control_code = control_code
	msg.Header.Parameter = parameter
	msg.Payload = payload
	return msg
}

func (m *Message) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.BigEndian, m.Header.Prologue)
	binary.Write(buf, binary.BigEndian, m.Header.Message_type)
	binary.Write(buf, binary.BigEndian, m.Header.Control_code)
	binary.Write(buf, binary.BigEndian, m.Header.Parameter)

	payloadLen := uint64(len(m.Payload))
	binary.Write(buf, binary.BigEndian, payloadLen)

	buf.Write(m.Payload)

	return buf.Bytes(), nil
}

func Decode(data []byte) (*Message, error) {
	r := bytes.NewReader(data)
	var prologue uint16
	var messageType uint8
	var controlCode uint8
	var parameter uint32

	binary.Read(r, binary.BigEndian, &prologue)
	if prologue != MSG_HEADER_PROLOGUE {
		return nil, errors.New("invalid prologue")
	}
	binary.Read(r, binary.BigEndian, &messageType)
	binary.Read(r, binary.BigEndian, &controlCode)
	binary.Read(r, binary.BigEndian, &parameter)

	var payloadLen uint64
	binary.Read(r, binary.BigEndian, &payloadLen)

	payload := make([]byte, payloadLen)
	r.Read(payload)
	return NewMessage(MessageType(messageType), controlCode, parameter, payload), nil
}
