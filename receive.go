package actioncable

import (
	"encoding/json"
	"errors"
)

func (ac *Client) receive() error {
	for {
		var event Event
		if err := ac.ws.ReadJSON(&event); err != nil {
			return err
		}

		if len(event.Type) != 0 {
			if err := ac.handleInternalEvent(&event); err != nil {
				return err
			}
		}
		ac.handler.Handle(ac, &event)
	}
}

func (ac *Client) handleInternalEvent(e *Event) error {
	switch e.Type {
	case "ping":
	case "confirm_subscription":
	case "welcome":
	case "disconnect":
		ac.exit()
		return errors.New("disconnect")
	default:
		ac.logger.Println("handleActionCable: unknown internal type ", e.Type)
	}
	return nil
}

type Event struct {
	Message    json.RawMessage `json:"message"`
	Type       string          `json:"type"`
	Identifier json.RawMessage `json:"identifier"`
}

type Identifier struct {
	Channel string `json:"channel"`
	ID      int    `json:"id"`
}
