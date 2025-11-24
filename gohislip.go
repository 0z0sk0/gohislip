package gohislip

import (
	"errors"

	"github.com/0z0sk0/gohislip/internal/resource"
	"github.com/0z0sk0/gohislip/internal/session"
	"github.com/0z0sk0/gohislip/internal/utils"
)

type Hislip interface {
	Connect(address string) error
	Write(data string) error
	Query(data string) (string, error)
}

type HislipResource struct {
	*resource.Resource
}

func NewHislipResource() *HislipResource {
	return &HislipResource{
		Resource: resource.NewResource(),
	}
}

func (res *HislipResource) Connect(address string) error {
	s, err := session.NewSession()
	if err != nil {
		return err
	}

	final_address, err := utils.SplitAddress(address)
	if err != nil {
		s.Free()
		return err
	}

	// Create for sync channel
	sync_conn, sessionId, err := res.SetupSync(final_address)
	if err != nil {
		s.Free()
		return err
	}

	// Create for async channel
	async_conn, err := res.SetupAsync(final_address, sessionId)
	if err != nil {
		s.Free()
		return err
	}

	s.SyncConn = sync_conn
	s.AsyncConn = async_conn
	s.SessionId = sessionId

	res.Session = s
	return nil
}

func (res *HislipResource) Write(data string) error {
	if res.Session == nil {
		return errors.New("you are not connected")
	}

	bytes := []byte(data)
	if err := res.SyncSend(bytes); err != nil {
		return err
	}

	return nil
}

func (res *HislipResource) Query(query string) (string, error) {
	if res.Session == nil {
		return "", errors.New("you are not connected")
	}

	bytes := []byte(query)
	if err := res.SyncSend(bytes); err != nil {
		return "", err
	}

	data, err := res.SyncReceive()
	if err != nil {
		return "", err
	}

	return string(data), nil
}
