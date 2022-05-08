package events

type Event struct {
	Type         string                 `json:"type"`
	SenderId     string                 `json:"sender_id"`
	RecipientIds []string               `json:"recipient_ids"`
	Data         map[string]interface{} `json:"data"`
}

func NewHighlightEvent(cfiRange string) *Event {
	return &Event{
		Type: "highlight",
		Data: map[string]interface{}{
			"cfiRange": cfiRange,
		},
	}
}

type PageDirection int64

const (
	Left PageDirection = iota
	Right
)

func (pd PageDirection) String() string {
	var s string

	switch pd {
	case Left:
		s = "Left"
	case Right:
		s = "Right"
	}
	return s
}

func NewTurnPageEvent(pd PageDirection) *Event {
	return &Event{
		Type: "turnPage",
		Data: map[string]interface{}{
			"direction": pd.String(),
		},
	}
}

func NewToPageEvent(p int) *Event {
	return &Event{
		Type: "toPage",
		Data: map[string]interface{}{
			"pageNumber": p,
		},
	}
}

func NewJoinEvent(id string) *Event {
	return &Event{
		Type:     "member.join",
		SenderId: id,
		Data: map[string]interface{}{
			"member_id": id,
		},
	}
}

func NewLeaveEvent(id string) *Event {
	return &Event{
		Type:     "member.leave",
		SenderId: id,
		Data: map[string]interface{}{
			"member_id": id,
		},
	}
}

func NewWSRoomCloseEvent(id string) *Event {
	return &Event{
		Type:     "room.close",
		SenderId: id,
		Data: map[string]interface{}{
			"closed_by": id,
		},
	}
}
