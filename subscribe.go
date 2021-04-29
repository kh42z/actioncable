package actioncable

import (
	"encoding/json"
)

func formatSubscribeMessage(channel string, id int, cmd string) *message {
	data, _ := json.Marshal(identifier{
		Channel: channel,
		ID:      id,
	})
	return &message{
		Command:    cmd,
		Identifier: string(data),
	}
}

// Subscribe starts following a channel
func (ac *Client) Subscribe(channel string, id int) {
	ac.emit <- formatSubscribeMessage(channel, id, "subscribe")
}

// Unsubscribe stops following a channel
func (ac *Client) Unsubscribe(channel string, id int) {
	ac.emit <- formatSubscribeMessage(channel, id, "unsubscribe")
}
