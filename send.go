package actioncable

import (
	"encoding/json"
)

func (ac *Client) SendMessage(channelName string, channelID int, content string) {
	data, _ := json.Marshal(command{
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

type command struct {
	Channel string `json:"channel"`
	ID      int    `json:"id"`
}

type message struct {
	Command    string `json:"command"`
	Data       string `json:"data,omitempty"`
	Identifier string `json:"identifier"`
}
