package actioncable

import (
	"testing"
)

func TestClient_receive(t *testing.T) {
	tests := map[string]struct {
		payload string
		welcomed bool
		want    string
	}{
		"unexpected type": {
			payload: "{\"type\":\"unexpected\"}",
			welcomed: false,
			want:    "Done"},
		"ping": {
			payload: "{\"type\":\"ping\",\"message\":1619511683}",
			welcomed: false,
			want:    "Done",
		},
		"disconnect": {
			payload: "{\"type\":\"disconnect\"}",
			welcomed: false,
			want:    "disconnect"},
		"unexpected first message": {
			payload: "{\"type\":\"hello\"}",
			welcomed: true,
			want:    "expecting welcome type message"},
		"invalid identifier": {
			payload: "{\"identifier\":\"{\\\"invalid\\\":\\\"ChatChannel\\\",\\\"id\\\":7}\",\"message\":{\"action\":\"message\",\"content\":\"Test\"}}",
			welcomed: false,
			want:    "Done"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rw := &wsReadMock{
				Welcomed:    tc.welcomed,
				ReadLimit:   1,
				ReadPayload: []byte(tc.payload),
			}
			c := NewClient(rw)
			if err := c.Run(); err != nil && err.Error() != tc.want {
				t.Fatalf("Expected error <%s>, got <%s>", tc.want, err.Error())
			}
		})
	}
}
