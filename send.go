package actioncable

import (
	"encoding/json"
)

// SendMessage sends a message on the specified channel
func (ac *Client) SendMessage(channelName string, channelID int, content string) {
	data, _ := json.Marshal(identifier{
		Channel: channelName,
		ID:      channelID,
	})
	ac.emit <- &message{
		Command:    "message",
		Identifier: string(data),
		Data:       content,
	}
}

func (ac *Client) send() {
	for {
		select {
		case <-ac.stop:
			return
		case m := <-ac.emit:
			if err := ac.ws.WriteJSON(m); err != nil {
				ac.logger.Println("Unable to send msg:", err)
			}
		}
	}
}

type message struct {
	Command    string `json:"command"`
	Data       string `json:"data,omitempty"`
	Identifier string `json:"identifier"`
}
