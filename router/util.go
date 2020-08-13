package router

import (
	"context"
	"encoding/json"
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

func HandleGetCardsGroup(w http.ResponseWriter, r *http.Request) {
	api.HttpResponse(w, r, model.CARD_GROUPS, http.StatusOK)
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
			if len(room.RoomSubscriber) == 0 && time.Now().Local().After(room.SessionExpired) {
				h.Lobbies.EventBroadcast <- model.EventData{
					Name: model.LOBBY_EVENT_ON_ROOM_REMOVE,
					Data: room.Room,
				}
				h.Rooms[id].EventBroadcast <- model.RoomEventData{Status: model.ROOM_STATUS_NOT_USE}
				break
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func (room *RoomsHub) startGame() {

	// first given
	room.givePlayerOneCard(room.Dealer.ID, true)
	room.EventBroadcast <- model.RoomEventData{
		Name: model.ROOM_EVENT_ON_CARD_GIVEN,
	}
	time.Sleep(2 * time.Second)

	for _, id := range room.TurnsOrder {
		room.givePlayerOneCard(id, true)
		room.EventBroadcast <- model.RoomEventData{
			Name: model.ROOM_EVENT_ON_CARD_GIVEN,
		}
		time.Sleep(2 * time.Second)
	}
	time.Sleep(3 * time.Second)

	// second given and check for blackjack
	room.givePlayerOneCard(room.Dealer.ID, false)
	evtDlr := room.blackjackForEvt(room.Dealer.ID, model.ROOM_EVENT_ON_CARD_GIVEN)
	room.EventBroadcast <- model.RoomEventData{
		Name: evtDlr,
	}
	time.Sleep(2 * time.Second)

	for _, id := range room.TurnsOrder {
		room.givePlayerOneCard(id, true)
		evt := room.blackjackForEvt(id, model.ROOM_EVENT_ON_CARD_GIVEN)
		room.EventBroadcast <- model.RoomEventData{
			Name: evt,
		}
		time.Sleep(2 * time.Second)
	}

	room.ConnectionMx.Lock()
	defer room.ConnectionMx.Unlock()

	// set to play
	// set first player turn
	room.Status = model.ROOM_STATUS_ON_PLAY
	if pTurn, ok := room.RoomPlayers[room.TurnsOrder[room.TurnPost]]; ok {
		pTurn.Status = model.PLAYER_STATUS_AT_TURN
	}

	room.EventBroadcast <- model.RoomEventData{
		Name: model.ROOM_EVENT_ON_GAME_START,
	}
}

func (r *RoomsHub) givePlayerOneCard(id string, show bool) {
	r.ConnectionMx.Lock()
	defer r.ConnectionMx.Unlock()

	if len(r.Cards) <= 0 {
		r.Room.CanDrawCard = false
		return
	}

	card := model.Card{}
	for _, c := range r.Cards {
		card = c.Copy(show)
		break
	}

	if _, ok := r.Cards[card.ID]; ok {
		delete(r.Cards, card.ID)
	}

	p, ok := r.RoomPlayers[id]
	if !ok {

		// give it to dealer
		r.Dealer.Cards = append(r.Dealer.Cards, card)
		r.Dealer.SumUpTotal()
		return
	}

	p.Cards = append(p.Cards, card)
	p.SumUpTotal()
}

func (r *RoomsHub) blackjackForEvt(id string, dfEvt string) string {
	r.ConnectionMx.Lock()
	defer r.ConnectionMx.Unlock()

	evt := dfEvt

	if player, ok := r.RoomPlayers[id]; ok {
		if player.Total == 21 {
			player.Status = model.PLAYER_STATUS_REWARDED
			evt = model.ROOM_EVENT_ON_PLAYER_BLACKJACK_WIN
		} else if player.Total > 21 {
			player.Status = model.PLAYER_STATUS_BUST
			evt = model.ROOM_EVENT_ON_PLAYER_BUST
		}

	} else if r.Dealer.ID == id {
		if r.Dealer.Total == 21 {
			evt = model.ROOM_EVENT_ON_PLAYER_BLACKJACK_WIN
		} else if r.Dealer.Total > 21 {
			evt = model.ROOM_EVENT_ON_PLAYER_BUST
		}
	}

	return evt
}

func (h *RouterHub) EndRound(id string) {

	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()

	r, ok := h.Rooms[id]
	if !ok {
		return
	}

	r.ConnectionMx.Lock()
	defer r.ConnectionMx.Unlock()

	for _, p := range r.RoomPlayers {

		// dealer bust
		// all player win
		// except who is buts
		if r.Dealer.Total > 21 && p.Total <= 21 {

			if pAcc, okAcc := h.Players[p.ID]; okAcc {
				pAcc.Money += p.Bet * 2
				p.Status = model.PLAYER_STATUS_REWARDED
			}

		} else {

			// if player is 21
			// win sweet
			if p.Total == 21 {

				if pAcc, okAcc := h.Players[p.ID]; okAcc {
					pAcc.Money += (p.Bet * 2) + (p.Bet / 2)
					p.Status = model.PLAYER_STATUS_REWARDED
				}

				// if player is score is higher
				// win
			} else if r.isMineHighThanOther(p) {

				if pAcc, okAcc := h.Players[p.ID]; okAcc {
					pAcc.Money += (p.Bet * 2)
					p.Status = model.PLAYER_STATUS_REWARDED
				}

				// lose bet
			} else {

				p.Status = model.PLAYER_STATUS_LOSE

			}
		}
	}

	scoreRound := r.Round
	r.Scores[scoreRound] = model.RecordScore(scoreRound, r.Dealer, r.RoomPlayers)

	r.Round++
}

func (r *RoomsHub) isMineHighThanOther(p *model.RoomPlayer) bool {
	isHigher := true

	if p.Total > 21 {
		return false
	}

	if p.Total <= r.Dealer.Total {
		return false
	}

	for _, rp := range r.RoomPlayers {
		if rp.ID != p.ID && p.Total < rp.Total {
			isHigher = false
			break
		}
	}

	return isHigher
}

func (r *RoomsHub) isPlayersStatusSame(status int) bool {
	for _, i := range r.RoomPlayers {
		if i.Status != status && i.Status != model.PLAYER_STATUS_BUST && i.Status != model.PLAYER_STATUS_REWARDED {
			return false
		}
	}
	return true
}

func (r *RoomsHub) removeFromTurnOrder(id string) {
	r.ConnectionMx.Lock()
	defer r.ConnectionMx.Unlock()

	r.TurnsOrder = []string{}
	for _, p := range r.Room.Players {
		if p.ID != id {
			r.TurnsOrder = append(r.TurnsOrder, p.ID)
		}
	}
}

func (r *RoomsHub) resetRoom() {
	r.ConnectionMx.Lock()
	defer r.ConnectionMx.Unlock()

	cards := model.NewCards(r.Room.CardGroups)
	r.Cards = make(map[string]*model.Card)
	for _, c := range cards {
		r.Cards[c.ID] = c.CopyPointer()
	}

	r.TurnPost = 0
	r.TurnsOrder = []string{}
	for _, p := range r.Room.Players {
		r.TurnsOrder = append(r.TurnsOrder, p.ID)
	}

	r.Dealer.Status = model.PLAYER_STATUS_IDLE
	r.Dealer.Bet = 0
	r.Dealer.Cards = []model.Card{}
	r.Dealer.Total = 0
	r.Dealer.TotalShow = 0

	for _, i := range r.RoomPlayers {
		i.Status = model.PLAYER_STATUS_IDLE
		i.Bet = 0
		i.Cards = []model.Card{}
		i.Total = 0
		i.TotalShow = 0
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
