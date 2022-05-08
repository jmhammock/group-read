package wsroom

import (
	"sync"
	"testing"

	"github.com/jmhammock/ereader/cmd/server/events"
)

type mockConnection struct{}

func (mc mockConnection) WriteJSON(interface{}) error {
	return nil
}

func (mc mockConnection) ReadJSON(interface{}) error {
	return nil
}

func (mc mockConnection) Close() error {
	return nil
}

var tracker = map[string]*events.Event{}

type mockMember struct {
	id   string
	conn Connectioner
}

func (m *mockMember) Send(e *events.Event) error {
	tracker[m.id] = e
	return nil
}

func (m *mockMember) Read(e *events.Event) error {
	return nil
}

func (m *mockMember) Close() error {
	return nil
}

func (m *mockMember) GetId() string {
	return m.id
}

func TestJoin(t *testing.T) {
	room := NewWSRoom("test_room")
	member := NewMember(&mockConnection{})

	wg := sync.WaitGroup{}
	wg.Add(1)
	var joinEvent *events.Event
	go func() {
		defer wg.Done()
		for e := range room.messageChan {
			joinEvent = e
		}
	}()

	room.Join(member)
	membersCount := room.MembersLen()
	room.Close(member.GetId())

	if membersCount != 1 {
		t.Logf("want 1 member. got %d members\n", room.MembersLen())
		t.Fail()
	}

	if joinEvent == nil {
		t.Log("want join event. got nil")
		t.FailNow()
	}

	if joinEvent.Type != "member.join" {
		t.Logf("want event type member.join. got event type %s\n", joinEvent.Type)
		t.Fail()
	}

	if joinEvent.SenderId != member.GetId() {
		t.Logf("want sender id %s. got sender id %s\n", member.GetId(), joinEvent.SenderId)
		t.Fail()
	}

	if memberId := joinEvent.Data["member_id"]; memberId != member.GetId() {
		t.Logf("want data.member_id %s. got data.member_id %s\n", member.GetId(), memberId)
		t.Fail()
	}
}

func TestLeave(t *testing.T) {
	room := NewWSRoom("test_room")
	member := NewMember(&mockConnection{})

	wg := sync.WaitGroup{}
	wg.Add(1)
	var evs []*events.Event
	go func() {
		defer wg.Done()
		for e := range room.messageChan {
			evs = append(evs, e)
		}
	}()

	room.Join(member)
	room.Leave(member.GetId())
	membersCount := room.MembersLen()
	room.Close(member.GetId())
	wg.Wait()

	if membersCount != 0 {
		t.Logf("want 0 members. got %d members\n", room.MembersLen())
		t.Fail()
	}

	leaveEvent := evs[1]

	if leaveEvent == nil {
		t.Log("want leave event. got nil")
		t.FailNow()
	}

	if leaveEvent.Type != "member.leave" {
		t.Logf("want event type member.leave. got event type %s\n", leaveEvent.Type)
		t.Fail()
	}

	if leaveEvent.SenderId != member.GetId() {
		t.Logf("want sender id %s. got sender id %s\n", member.GetId(), leaveEvent.SenderId)
		t.Fail()
	}

	if memberId := leaveEvent.Data["member_id"]; memberId != member.GetId() {
		t.Logf("want data.member_id %s. got data.member_id %s\n", member.GetId(), memberId)
		t.Fail()
	}
}

func TestBroadcast(t *testing.T) {
	cOne := &mockMember{
		id: "one",
	}
	cTwo := &mockMember{
		id: "two",
	}

	room := NewWSRoom("test")
	defer room.Close("test_user")
	room.Receiver()
	room.Join(cOne)
	room.Join(cTwo)

	testMessage := &events.Event{
		Type:     "message.test",
		SenderId: "test_user",
		Data: map[string]interface{}{
			"message": "hello, world!",
		},
	}

	room.Broadcast(testMessage)

	if tracker[cOne.id].Type != "message.test" {
		t.Logf("want lastEvent.Type to be message.test. got %s\n", tracker[cOne.id].Type)
		t.Fail()
	}

	if tracker[cTwo.id].Type != "message.test" {
		t.Logf("want lastEvent.Type to be message.test. got %s\n", tracker[cTwo.id].Type)
		t.Fail()
	}
}

func TestAddRoom(t *testing.T) {
	manager := NewWSRoomManager()
	room := NewWSRoom("test")
	manager.AddRoom(room)
	if manager.RoomsLen() != 1 {
		t.Logf("want 1 room. got %d rooms\n", manager.RoomsLen())
		t.Fail()
	}
	if _, err := manager.GetRoom("test"); err != nil {
		t.Log(err)
		t.Fail()
	}

	if _, err := manager.GetRoom("id doesn't exist"); err == nil {
		t.Log("want error. got nil")
		t.Fail()
	}
}

func TestRemoveRoom(t *testing.T) {
	manager := NewWSRoomManager()
	room := NewWSRoom("test")
	manager.AddRoom(room)
	wg := sync.WaitGroup{}

	var closeEvent *events.Event
	wg.Add(1)
	go func() {
		defer wg.Done()
		for e := range room.messageChan {
			closeEvent = e
		}
	}()

	manager.RemoveRoom(room.Id, "test user")
	if closeEvent == nil {
		t.Log("want close event. got nil")
		t.FailNow()
	}

	if closeEvent.Type != "room.close" {
		t.Logf("want event type room.close. got event type %s\n", closeEvent.Type)
		t.Fail()
	}

	if closeEvent.SenderId != "test user" {
		t.Logf("want sender id 'test user'. got sender id %s\n", closeEvent.SenderId)
	}

}
