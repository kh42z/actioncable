package actioncable

import (
	"testing"
)

type ChannelHandlerStub struct {
	Subscribed bool
	Messaged   bool
}

func (c *ChannelHandlerStub) SubscriptionHandler(_ *Client, _ int) {
	c.Subscribed = true
}

func (c *ChannelHandlerStub) MessageHandler(_ *Client, _ []byte, _ int) {
	c.Messaged = true
}

func TestClientRegisterChannelCallbacks(t *testing.T) {
	tests := map[string]struct {
		payload string
		want    *ChannelHandlerStub
		error   string
	}{
		"onSubscription": {
			payload: "{\"identifier\":\"{\\\"channel\\\":\\\"UserChannel\\\",\\\"id\\\":1}\",\"type\":\"confirm_subscription\"}",
			want:    &ChannelHandlerStub{Subscribed: true, Messaged: false},
			error:   "Done"},
		"onMessage": {
			payload: "{\"identifier\":\"{\\\"channel\\\":\\\"UserChannel\\\",\\\"id\\\":7}\",\"message\":{\"content\":\"Some content?\"}}",
			want:    &ChannelHandlerStub{Subscribed: false, Messaged: true},
			error:   "Done"},
		"invalidMessage": {
			payload: "{\"identifier\":\"{\\\"channel\\\": 4,\\\"id\\\":7}\",\"message\":{\"action\":\"message\",\"content\":\"Test\"}}",
			want:    &ChannelHandlerStub{Subscribed: false, Messaged: false},
			error:   "Done"},
		"invalidSubscription": {
			payload: "{\"identifier\":\"{\\\"channel\\\": 3,\\\"id\\\":1}\",\"type\":\"confirm_subscription\"}",
			want:    &ChannelHandlerStub{Subscribed: false, Messaged: false},
			error:   "json: cannot unmarshal number into Go struct field identifier.channel of type string"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rw := &wsReadMock{
				Welcomed:    false,
				ReadLimit:   1,
				ReadPayload: []byte(tc.payload),
			}
			c := NewClient(rw)
			channelHS := &ChannelHandlerStub{
				Subscribed: false,
				Messaged:   false}
			c.AddChannelHandler("UserChannel", channelHS)
			if err := c.Run(); err != nil && err.Error() != tc.error {
				t.Errorf("unexpect error %s", err)
			}
			if tc.want.Messaged != channelHS.Messaged {
				t.Errorf("expecting Messaged to be %t got %t", tc.want.Messaged, channelHS.Messaged)
			}
			if tc.want.Subscribed != channelHS.Subscribed {
				t.Errorf("expecting Subscribed to be %t got %t", tc.want.Subscribed, channelHS.Subscribed)
			}
		})
	}
}
