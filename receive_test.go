package actioncable

import (
	"reflect"
	"testing"
)

type Bot struct {
	Event *Event
}

func (b *Bot) Handle(_ *Client, e *Event) {
	b.Event = e
}

func TestClient_receive(t *testing.T) {
	tests := map[string]struct {
		payload string
		want    *Event
		error	string
	}{
		"welcome": {
			payload: "{\"type\":\"welcome\"}",
			want:    &Event{Type: "welcome"},
			error:   "Done"},
		"ping": {
			payload: "{\"type\":\"ping\",\"message\":1619511683}",
			want: &Event{
				Type:    "ping",
				Message: []byte("1619511683")},
			error: "Done",
		},
		"disconnect": {
			payload: "{\"type\":\"disconnect\"}",
			want:    nil,
			error: "disconnect"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rw := &wsMock{
				ReadLimit:   1,
				ReadPayload: []byte(tc.payload),
			}
			b := &Bot{}
			c := NewClient(rw, b)
			if err := c.Run(); err != nil && err.Error() != tc.error {
				t.Fatal("got an unexpected error ", err)
			}
			if !reflect.DeepEqual(b.Event, tc.want) {
				t.Fatalf("expected: %v, got: %v", tc.want, b.Event)
			}
		})
	}
}
