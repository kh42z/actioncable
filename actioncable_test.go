package actioncable

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
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
	Over         bool
	WriteLimit   int
	WritePayload []byte
	sync.Mutex
}

func (ws *wsWriteMock) ReadJSON(v interface{}) error {
	ws.Lock()
	defer ws.Unlock()
	if ws.Over {
		return errors.New("Done")
	}
	json.Unmarshal([]byte("{\"type\":\"welcome\"}"), v)
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
	ws.Lock()
	defer ws.Unlock()
	ws.Over = true
}

func TestNewClient(t *testing.T) {
	NewClient(&wsReadMock{}, WithLogger(&log.Logger{}))
}
