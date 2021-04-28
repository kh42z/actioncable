package actioncable

import (
	"encoding/json"
)

func formatSubscribeMessage(channel string, ID int, cmd string) *message {
	data, _ := json.Marshal(identifier{
		Channel: channel,
		ID:      ID,
	})
	return &message{
		Command:    cmd,
		Identifier: string(data),
	}
}
func (ac *Client) Subscribe(channel string, ID int) {
	ac.emit <- formatSubscribeMessage(channel, ID, "subscribe")
}

func (ac *Client) Unsubscribe(channel string, ID int) {
	ac.emit <- formatSubscribeMessage(channel, ID, "unsubscribe")
}
