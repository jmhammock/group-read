package events

type Event struct {
	Type     string                 `json:"type"`
	SenderId string                 `json:"sender_id"`
	Data     map[string]interface{} `json:"data"`
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
