package actioncable

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"testing"
	"time"
)

type wsMock struct {
	NoRead       bool
	ReadLimit    int
	ReadPayload  []byte
	WriteLimit   int
	WritePayload []byte
	sync.Mutex
}

func (ws *wsMock) ReadJSON(v interface{}) error {
	ws.Lock()
	defer ws.Unlock()
	if ws.NoRead {
		time.Sleep(1000 * time.Millisecond)
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

func (ws *wsMock) WriteJSON(v interface{}) error {
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

func (ws *wsMock) CancelRead() {
	ws.Lock()
	defer ws.Unlock()
	ws.NoRead = false
	ws.ReadLimit = 0
}

func TestNewClient(t *testing.T) {
	NewClient(&wsMock{}, WithLogger(&log.Logger{}))
}
