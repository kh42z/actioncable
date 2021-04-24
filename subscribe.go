package actioncable

import (
	"encoding/json"
)

func formatSubscribeMessage(channel string, ID int) *message {
	data, _ := json.Marshal(command{
		Channel: channel,
		ID:      ID,
	})
	return &message{
		Command:    "subscribe",
		Identifier: string(data),
	}
}
func (ac *Client) Subscribe(channel string, ID int) {
	ac.emit <- formatSubscribeMessage(channel, ID)
}
