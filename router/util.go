package router

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
	"github.com/renosyah/simple-21/api"
	"github.com/renosyah/simple-21/model"
	"github.com/renosyah/simple-21/util"
)

func HandleGetRandomName(w http.ResponseWriter, r *http.Request) {
	wttl := r.FormValue("title")
	api.HttpResponse(w, r, util.GenerateRandomName(wttl != ""), http.StatusOK)
}

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
				h.Rooms[id].EventBroadcast <- model.RoomEventData{Status: model.ROOM_STATUS_NOT_USE}
				break
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func (r *RoomsHub) givePlayerOneCard(id string, show bool) {
	r.ConnectionMx.Lock()
	defer r.ConnectionMx.Unlock()

	if len(r.Cards) <= 0 {
		return
	}

	card := model.Card{}
	for _, c := range r.Cards {
		card = c.Copy(show)
		break
	}

	fmt.Println(string(model.ToJson(card)))

	if _, ok := r.Cards[card.ID]; ok {
		delete(r.Cards, card.ID)
	}

	p, ok := r.RoomPlayers[id]
	if !ok {

		// give it to dealer
		r.Dealer.Cards = append(r.Dealer.Cards, card)
		return
	}

	p.Cards = append(p.Cards, card)
}

func (r *RoomsHub) isPlayersStatusSame(status int) bool {
	for _, i := range r.RoomPlayers {
		if i.Status != status {
			return false
		}
	}
	return true
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
