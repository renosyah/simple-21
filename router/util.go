package router

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
	"github.com/renosyah/simple-21/model"
)

func (h *RouterHub) dropOffPlayer() {
	for {
		h.ConnectionMx.Lock()
		for k, p := range h.Players {
			_, ok := h.Lobbies.Subscriber[k]
			if !ok && !p.IsOnline && time.Now().Local().After(p.SessionExpired) {
				delete(h.Players, p.ID)
				h.Lobbies.EventBroadcast <- model.EventData{Name: model.LOBBY_EVENT_ON_LOGOUT}
				break
			}
		}
		h.ConnectionMx.Unlock()
		time.Sleep(5 * time.Second)
	}
}

func (h *RouterHub) dropEmptyRoom() {
	for {
		for id, room := range h.Rooms {
			if len(room.RoomPlayers) == 0 && time.Now().Local().After(room.SessionExpired) {
				h.Rooms[id].EventBroadcast <- model.RoomEventData{Status: ROOM_STATUS_NOT_USE}
				break
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func ParseBodyData(ctx context.Context, r *http.Request, data interface{}) error {
	bBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(err, "read")
	}

	err = json.Unmarshal(bBody, data)
	if err != nil {
		return errors.Wrap(err, "json")
	}

	valid, err := govalidator.ValidateStruct(data)
	if err != nil {
		return errors.Wrap(err, "validate")
	}

	if !valid {
		return errors.New("invalid data")
	}

	return nil
}
