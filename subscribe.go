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

// Subscribe starts following a channel
func (ac *Client) Subscribe(channel string, ID int) {
	ac.emit <- formatSubscribeMessage(channel, ID, "subscribe")
}

// Unsubscribe stops following a channel
func (ac *Client) Unsubscribe(channel string, ID int) {
	ac.emit <- formatSubscribeMessage(channel, ID, "unsubscribe")
}
