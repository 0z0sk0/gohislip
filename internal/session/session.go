package session

import (
	"errors"
	"net"
)

const MAX_SESSIONS int = 256

var sessions []Session

type Session struct {
	allocated bool
	SyncConn  net.Conn
	AsyncConn net.Conn
	SessionId uint32
	MessageId uint32
}

func NewSession() (*Session, error) {
	if len(sessions) >= MAX_SESSIONS {
		return nil, errors.New("not available space to session")
	}
	session := &Session{
		allocated: true,
		SessionId: uint32(len(sessions)),
	}
	sessions = append(sessions, *session)
	return session, nil
}

func (s *Session) Free() error {
	if int(s.SessionId) >= MAX_SESSIONS {
		return errors.New("invalid session id")
	}

	if s.allocated {
		return errors.New("already free")
	}

	s.allocated = false
	return nil
}
