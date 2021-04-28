package actioncable

import "testing"

func TestClient_Subscribe(t *testing.T) {
	tests := map[string]struct {
		name    string
		id      int
		content string
		want    string
	}{
		"subscribe": {
			name: "UserChannel",
			id:   1,
			want: "{\"command\":\"subscribe\",\"identifier\":\"{\\\"channel\\\":\\\"UserChannel\\\",\\\"id\\\":1}\"}"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rw := &wsWriteMock{
				WriteLimit: 1,
				Over:       false,
			}
			c := NewClient(rw)
			go c.Run()
			c.Subscribe("UserChannel", 1)
			rw.WaitForWrite()
			if string(rw.WritePayload) != tc.want {
				t.Errorf("expecting payload to be [%s] got [%s]", tc.want, rw.WritePayload)
			}
			rw.CancelRead()
		})
	}
}

func TestClient_UnSubscribe(t *testing.T) {
	tests := map[string]struct {
		name    string
		id      int
		content string
		want    string
	}{
		"unsubscribe": {
			name: "UserChannel",
			id:   1,
			want: "{\"command\":\"unsubscribe\",\"identifier\":\"{\\\"channel\\\":\\\"UserChannel\\\",\\\"id\\\":1}\"}"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rw := &wsWriteMock{
				WriteLimit: 1,
				Over:       false,
			}
			c := NewClient(rw)
			go c.Run()
			c.Unsubscribe("UserChannel", 1)
			rw.WaitForWrite()
			if string(rw.WritePayload) != tc.want {
				t.Errorf("expecting payload to be [%s] got [%s]", tc.want, rw.WritePayload)
			}
			rw.CancelRead()
		})
	}
}
