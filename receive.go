package actioncable

import (
	"encoding/json"
	"errors"
)

func (ac *Client) receive() error {
	for {
		var event event
		if err := ac.ws.ReadJSON(&event); err != nil {
			return err
		}

		if len(event.Type) != 0 {
			if err := ac.handleActionCableEvent(&event); err != nil {
				return err
			}
		} else {
			ac.handleEvent(&event)
		}
	}
}

func (ac *Client) handleActionCableEvent(e *event) error {
	switch e.Type {
	case "welcome":
		ac.onWelcome()
	case "confirm_subscription":
		var i identifier
		err := json.Unmarshal([]byte(e.Identifier), &i)
		if err != nil {
			ac.logger.Println("handleActionCable: unable to unmarshal Identifier", i)
			return err
		}
		for name, e := range ac.channels {
			if name == i.Channel {
				e.OnSubscription(i.ID)
			}
		}
	case "disconnect":
		return errors.New("actioncable: disconnect")
	case "ping":
	default:
		ac.logger.Println("handleActionCable: unknown internal type ", e.Type)
	}
	return nil
}

type event struct {
	Message    json.RawMessage `json:"message"`
	Type       string          `json:"type"`
	Identifier string          `json:"identifier"`
}

type identifier struct {
	Channel string `json:"channel"`
	ID      int    `json:"id"`
}
