package actioncable

import "testing"

type ChannelMock struct {
	Subscribed bool
	Messaged bool
}

func (c *ChannelMock) OnSubscription(_ int) {
	c.Subscribed = true
}

func (c *ChannelMock) OnMessage(_ []byte, _ int) {
	c.Messaged = true
}

func TestClient_RegisterChannelCallbacks(t *testing.T) {
	rw := &wsMock{
		ReadLimit: 1,
		ReadPayload: []byte("{\"identifier\":\"{\\\"channel\\\":\\\"UserChannel\\\",\\\"id\\\":1}\",\"type\":\"confirm_subscription\"}"),
	}
	c := NewClient(rw)

	channelCallback := &ChannelMock{false, false}
	c.RegisterChannelCallbacks("UserChannel", channelCallback)
	c.Start()
	if channelCallback.Subscribed == false {
		t.Errorf("Expecting UserChannel to be triggered (true) got %t", channelCallback.Subscribed)
	}
}
