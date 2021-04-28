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
			if err := ac.handleInternalEvent(&event); err != nil {
				return err
			}
		} else {
			ac.handleEvent(&event)
		}
	}
}

func (ac *Client) handleInternalEvent(e *event) error {
	switch e.Type {
	case "ping":
	case "confirm_subscription":
		if err := ac.handleSubscription(e); err != nil {
			return err
		}
	case "disconnect":
		ac.exit()
		return errors.New("disconnect")
	default:
		ac.logger.Println("handleActionCable: unknown internal type ", e.Type)
	}
	return nil
}

func (ac *Client) handleSubscription(e *event) error {
	var i identifier
	err := json.Unmarshal([]byte(e.Identifier), &i)
	if err != nil {
		return err
	}
	for name, e := range ac.channels {
		if name == i.Channel {
			e.SubscriptionHandler(ac, i.ID)
		}
	}
	return nil
}

func (ac *Client) waitWelcome() error {
	var event event
	if err := ac.ws.ReadJSON(&event); err != nil {
		return err
	}
	if len(event.Type) == 0 || event.Type != "welcome" {
		return errors.New("expecting welcome type message")
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
