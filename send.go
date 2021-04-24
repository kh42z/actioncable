package actioncable

import (
	"encoding/json"
)

type Command struct {
	Channel string `json:"channel"`
	ID      int    `json:"id"`
}

type Message struct {
	Command    string `json:"command"`
	Data       string `json:"data,omitempty"`
	Identifier string `json:"identifier"`
}

func (ac *Client) SendMessage(channelName string, channelID int, content string) {
	data, _ := json.Marshal(Command{
		Channel: channelName,
		ID:      channelID,
	})
	ac.emit <- &Message{
		Command:    "message",
		Identifier: string(data),
		Data:       content,
	}
}

func (ac *Client) send() {
	for {
		select {
		case <-ac.quit:
			return
		case m := <-ac.emit:
			if err := ac.ws.WriteJSON(m); err != nil {
				ac.logger.Println("Unable to send msg:", err)
			}
		}
	}
}
