package actioncable

import (
	"log"
	"testing"
)

type ChannelMock struct {
	Subscribed bool
	Messaged bool
}

func (c *ChannelMock) OnSubscription(_ int) {
	c.Subscribed = true
}

func (c *ChannelMock) OnMessage(_ []byte, _ int) {
	log.Println("###############")
	c.Messaged = true
}

func TestClient_RegisterChannelCallbacks(t *testing.T) {
	tests := map[string]struct {
		payload string
		want  *ChannelMock
	}{
		"onSubscription":       {payload: "{\"identifier\":\"{\\\"channel\\\":\\\"UserChannel\\\",\\\"id\\\":1}\",\"type\":\"confirm_subscription\"}",
			want: &ChannelMock{Subscribed: true, Messaged: false}},
		"onMessage":       {payload: "{\"identifier\":\"{\\\"channel\\\":\\\"UserChannel\\\",\\\"id\\\":7}\",\"message\":{\"content\":\"Some content?\"}}",
			want: &ChannelMock{Subscribed: false, Messaged: true}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rw := &wsMock{
				ReadLimit: 1,
				ReadPayload: []byte(tc.payload),
			}
			c := NewClient(rw)
			channelCallback := &ChannelMock{false, false}
			c.RegisterChannelCallbacks("UserChannel", channelCallback)
			c.Start()
			if tc.want.Messaged != channelCallback.Messaged {
				t.Errorf("expecting onMessage to be %t got %t", tc.want.Messaged, channelCallback.Messaged)
			}
			if tc.want.Subscribed != channelCallback.Subscribed {
				t.Errorf("expecting onSubscription to be %t got %t", tc.want.Subscribed, channelCallback.Subscribed)
			}
		})
	}
}

