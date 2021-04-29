package actioncable

import (
	"encoding/json"
)

// ChannelHandler is an interface designed to handle events sent by an ActionCable Server on a specific Channel ("ChatroomChannel" for example).
// SubscriptionHandler is triggered by a successful subscribe on a channel
// MessageHandler is triggered by a message received on a channel
type ChannelHandler interface {
	SubscriptionHandler(c *Client, channelId int)
	MessageHandler(c *Client, message []byte, channelId int)
}

// AddChannelHandler adds handlers for channel ("UserChannel" for example)
func (ac *Client) AddChannelHandler(name string, event ChannelHandler) {
	ac.channels[name] = event
}

func (ac *Client) handleEvent(event *event) {
	var i identifier
	err := json.Unmarshal([]byte(event.Identifier), &i)
	if err != nil {
		ac.logger.Println("handleEvent: ", err)
		return
	}
	for name, e := range ac.channels {
		if name == i.Channel {
			e.MessageHandler(ac, event.Message, i.ID)
		}
	}
}
