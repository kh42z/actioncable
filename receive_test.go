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
	}{
		"welcome": {
			payload: "{\"type\":\"welcome\"}",
			want:    &Event{Type: "welcome"}},
		"ping": {
			payload: "{\"type\":\"ping\",\"message\":1619511683}",
			want: &Event{
				Type:    "ping",
				Message: []byte("1619511683")}},
		"message": {
			payload: "{\"identifier\":\"{\\\"channel\\\":\\\"UserChannel\\\",\\\"id\\\":7}\",\"message\": \"data\"}",
			want: &Event{
				Identifier: []byte("\"{\\\"channel\\\":\\\"UserChannel\\\",\\\"id\\\":7}\""),
				Message:    []byte("\"data\"")}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rw := &wsMock{
				ReadLimit:   1,
				ReadPayload: []byte(tc.payload),
			}
			b := &Bot{}
			c := NewClient(rw, b)
			if err := c.Run(); err != nil && err.Error() != "Done" {
				t.Fatal("got an unexpected error ", err)
			}
			if !reflect.DeepEqual(b.Event, tc.want) {
				t.Fatalf("expected: %v, got: %v", tc.want, b.Event)
			}
		})
	}
}
