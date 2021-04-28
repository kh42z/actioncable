package actioncable

import (
	"testing"
)

func TestClient_receive(t *testing.T) {
	tests := map[string]struct {
		payload string
		want    string
	}{
		"welcome": {
			payload: "{\"type\":\"welcome\"}",
			want:    "Done"},
		"ping": {
			payload: "{\"type\":\"ping\",\"message\":1619511683}",
			want:    "Done",
		},
		"disconnect": {
			payload: "{\"type\":\"disconnect\"}",
			want:    "disconnect"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rw := &wsReadMock{
				Welcomed:    false,
				ReadLimit:   1,
				ReadPayload: []byte(tc.payload),
			}
			c := NewClient(rw)
			if err := c.Run(); err != nil && err.Error() != tc.want {
				t.Fatal("got an unexpected error ", err)
			}
		})
	}
}
