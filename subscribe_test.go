package actioncable

import "testing"

func TestClient_Subscribe(t *testing.T) {
	tests := map[string]struct {
		name    string
		id      int
		content string
		want    string
	}{
		"UserChannel 1": {name: "UserChannel",
			id:   1,
			want: "{\"command\":\"subscribe\",\"identifier\":\"{\\\"channel\\\":\\\"UserChannel\\\",\\\"id\\\":1}\"}"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rw := &wsMock{
				WriteLimit: 1,
				NoRead:     true,
			}
			c := NewClient(rw)
			go c.Run()
			c.Subscribe("UserChannel", 1)
			if string(rw.WritePayload) != tc.want {
				t.Errorf("expecting payload to be [%s] got [%s]", tc.want, rw.WritePayload)
			}
			rw.CancelRead()
		})
	}
}
