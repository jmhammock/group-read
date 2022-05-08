package wsroom

import (
	"github.com/google/uuid"
	"github.com/jmhammock/ereader/cmd/server/events"
)

type Client interface {
	Send(*events.Event) error
	Read(*events.Event) error
	Close() error
	GetId() string
}

type Connectioner interface {
	WriteJSON(interface{}) error
	ReadJSON(interface{}) error
	Close() error
}

type Member struct {
	id   string
	conn Connectioner
}

func NewMember(conn Connectioner) *Member {
	return &Member{
		id:   uuid.NewString(),
		conn: conn,
	}
}

func (m *Member) Send(e *events.Event) error {
	return m.conn.WriteJSON(e)
}

func (m *Member) Read(e *events.Event) error {
	return m.conn.ReadJSON(e)
}

func (m *Member) Close() error {
	return m.conn.Close()
}

func (m *Member) GetId() string {
	return m.id
}
