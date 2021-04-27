package actioncable

import (
	"reflect"
	"testing"
)

type ChannelHandlerStub struct {
	Subscribed bool
	Messaged   bool
}

func (c *ChannelHandlerStub) OnSubscription(_ *Client, id int) {
	c.Subscribed = true
}

func (c *ChannelHandlerStub) OnMessage(_ *Client, _ []byte, _ int) {
	c.Messaged = true
}

func TestClient_RegisterChannelCallbacks(t *testing.T) {
	tests := map[string]struct {
		payload string
		want    *ChannelHandlerStub
	}{
		"onSubscription": {payload: "{\"identifier\":\"{\\\"channel\\\":\\\"UserChannel\\\",\\\"id\\\":1}\",\"type\":\"confirm_subscription\"}",
			want: &ChannelHandlerStub{Subscribed: true, Messaged: false}},
		"onMessage": {payload: "{\"identifier\":\"{\\\"channel\\\":\\\"UserChannel\\\",\\\"id\\\":7}\",\"message\":{\"content\":\"Some content?\"}}",
			want: &ChannelHandlerStub{Subscribed: false, Messaged: true}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rw := &wsMock{
				ReadLimit:   1,
				ReadPayload: []byte(tc.payload),
			}
			c := NewClient(rw)
			channelHS := &ChannelHandlerStub{false, false}
			c.AddChannelHandler("UserChannel", tc.want)
			if err := c.Run(); err != nil && err.Error() != "Done" {
				t.Errorf("unexpect error %s", err)
			}
			if reflect.DeepEqual(tc.want, channelHS) {
				t.Errorf("expecting onMessage to be %t got %t", channelHS, tc.want)
			}
		})
	}
}
