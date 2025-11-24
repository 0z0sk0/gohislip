package resource

import (
	"bytes"
	"errors"
	"net"
	"time"

	"github.com/0z0sk0/gohislip/internal/message"
	"github.com/0z0sk0/gohislip/internal/session"
	"github.com/0z0sk0/gohislip/internal/utils"
)

const (
	HISLIP_VERSION_MAJOR = 0
	HISLIP_VERSION_MINOR = 0
	HISLIP_VENDOR_ID     = 42
)

type Hislip interface {
	Connect(address string) error
	Write(data string) error
	Query(data string) (string, error)
}

type Resource struct {
	Session *session.Session
}

func NewResource() *Resource {
	return &Resource{}
}

func (res *Resource) SetupSync(address string) (net.Conn, uint32, error) {
	sync_conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, 0, err
	}

	sync_conn.SetWriteDeadline(time.Now().Add(time.Second * 20))
	sync_conn.SetReadDeadline(time.Now().Add(time.Second * 20))

	version := (HISLIP_VERSION_MAJOR << 8) + HISLIP_VERSION_MINOR
	parameter := (version << 16) + HISLIP_VENDOR_ID

	// TODO: control codes
	msg := message.NewMessage(message.Initialize, 0, uint32(parameter), nil)
	bytes, err := msg.Encode()
	if err != nil {
		return nil, 0, err
	}

	_, err = sync_conn.Write(bytes)
	if err != nil {
		return nil, 0, err
	}

	buf := make([]byte, 128)
	sync_conn.Read(buf)

	received, err := message.Decode(buf)
	if err != nil {
		return nil, 0, err
	}

	if received.Header.Message_type != message.InitializeResponse {
		return nil, 0, errors.New("wrong response to Initialize")
	}

	// Recognize received session id
	sessionId := received.Header.Parameter
	sessionId &= 0xFFFF

	return sync_conn, sessionId, err
}

func (res *Resource) SetupAsync(address string, sessionId uint32) (net.Conn, error) {
	async_conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	// TODO: control codes
	msg := message.NewMessage(message.AsyncInitialize, 0, sessionId, nil)
	bytes, err := msg.Encode()
	if err != nil {
		return nil, err
	}

	_, err = async_conn.Write(bytes)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 128)
	async_conn.Read(buf)

	received, err := message.Decode(buf)
	if err != nil {
		return nil, err
	}

	if received.Header.Message_type != message.AsyncInitializeResponse {
		return nil, errors.New("wrong response to AsyncInitialize")
	}

	return async_conn, err
}

func (res *Resource) SyncSend(data []byte) error {
	conn := res.Session.SyncConn
	payloadMax := message.MSG_MAX_SIZE - message.MSG_HEADER_SIZE
	length := len(data)

	if payloadMax <= 0 {
		return errors.New("invalid payload max size")
	}

	messageCount := length / payloadMax
	messageBytesRemaining := messageCount * payloadMax
	if length%payloadMax != 0 {
		messageBytesRemaining = length % payloadMax
	}

	offset := 0

	for messageBytesRemaining > 0 {
		chunkSize := payloadMax
		if messageBytesRemaining < chunkSize {
			chunkSize = int(messageBytesRemaining)
		}

		parameter := res.Session.MessageId

		res.Session.MessageId += 2

		if offset+chunkSize > length {
			return errors.New("internal chunking error: slice out of range")
		}
		payload := data[offset : offset+chunkSize]

		msgType := message.Data
		if uint64(offset+chunkSize) == uint64(length) {
			msgType = message.DataEnd
		}

		// TODO: control codes
		msg := message.NewMessage(msgType, 0, uint32(parameter), payload)
		bytes, err := msg.Encode()
		if err != nil {
			return err
		}

		if err := utils.WriteAll(conn, bytes); err != nil {
			return err
		}

		offset += chunkSize
		messageBytesRemaining -= chunkSize
	}

	return nil
}

func (res *Resource) SyncReceive() ([]byte, error) {
	conn := res.Session.SyncConn
	var payloadBuffer bytes.Buffer

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			return nil, err
		}
		buf = buf[:n]

		received, err := message.Decode(buf)
		if err != nil {
			return nil, err
		}

		payloadBuffer.Write([]byte(received.Payload))

		if received.Header.Message_type == message.DataEnd {
			break
		}
	}

	return payloadBuffer.Bytes(), nil
}
