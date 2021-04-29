package actioncable

import (
	"encoding/json"
	"errors"
	"log"
	"testing"
	"time"
)

type wsReadMock struct {
	Welcomed    bool
	ReadLimit   int
	ReadPayload []byte
}

func (ws *wsReadMock) ReadJSON(v interface{}) error {
	if ws.Welcomed == false {
		json.Unmarshal([]byte("{\"type\":\"welcome\"}"), v)
		ws.Welcomed = true
		return nil
	}
	err := json.Unmarshal(ws.ReadPayload, v)
	if err != nil {
		return err
	}
	ws.ReadLimit--
	if ws.ReadLimit < 0 {
		return errors.New("Done")
	}
	return nil
}

func (ws *wsReadMock) WriteJSON(_ interface{}) error {
	return nil
}

type wsWriteMock struct {
	WriteLimit   int
	WritePayload []byte
	Welcomed     bool
	Over         chan struct{}
}

func (ws *wsWriteMock) ReadJSON(v interface{}) error {
	if !ws.Welcomed {
		json.Unmarshal([]byte("{\"type\":\"welcome\"}"), v)
		ws.Welcomed = true
		return nil
	}
	<-ws.Over
	return nil
}

func (ws *wsWriteMock) WriteJSON(v interface{}) error {
	var err error
	ws.WritePayload, err = json.Marshal(v)
	if err != nil {
		log.Fatal("WriteJSON: unable to json")
	}
	ws.WriteLimit--
	if ws.WriteLimit < 0 {
		return errors.New("Done")
	}
	return nil
}

func (ws *wsWriteMock) WaitForWrite() {
	for {
		if len(ws.WritePayload) > 0 {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (ws *wsWriteMock) CancelRead() {
	close(ws.Over)
}

func TestNewClient(t *testing.T) {
	NewClient(&wsReadMock{}, WithLogger(&log.Logger{}))
}
