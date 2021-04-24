package actioncable

import (
	"encoding/json"
)

func formatSubscribeMessage(channel string, ID int) *Message {
	data, _ := json.Marshal(Command{
		Channel: channel,
		ID:      ID,
	})
	return &Message{
		Command:    "subscribe",
		Identifier: string(data),
	}
}
func (ac *Client) Subscribe(channel string, ID int) {
	ac.emit <- formatSubscribeMessage(channel, ID)
}
