package actioncable

import "testing"

func TestClient_SendMessage(t *testing.T) {
	tests := map[string]struct {
		name    string
		id      int
		content string
		want    string
	}{
		"plain": {name: "UserChannel",
			id:      1,
			content: "Hello", want: "{\"command\":\"message\",\"data\":\"Hello\",\"identifier\":\"{\\\"channel\\\":\\\"UserChannel\\\",\\\"id\\\":1}\"}"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rw := &wsMock{
				WriteLimit: 1,
				NoRead:     true,
			}
			c := NewClient(rw)
			go c.Start()
			c.SendMessage(tc.name, tc.id, tc.content)
			rw.CancelRead()
			if string(rw.WritePayload) != tc.want {
				t.Errorf("expecting payload to be [%s] got [%s]", tc.want, rw.WritePayload)
			}
		})
	}
}
