package router

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/renosyah/simple-21/api"
	"github.com/renosyah/simple-21/model"
)

func (h *RouterHub) HandleAddRoom(w http.ResponseWriter, r *http.Request) {
	var param model.AddRoom

	err := ParseBodyData(r.Context(), r, &param)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(h.Rooms) >= h.Config.MaxRoom {
		api.HttpResponseException(w, r, http.StatusInsufficientStorage)
		return
	}

	h.openRoom(param.HostID, param.Name, param.Players)

	h.Lobbies.EventBroadcast <- model.EventData{
		Name: model.LOBBY_EVENT_ON_ROOM_CREATED,
		Data: param,
	}

	api.HttpResponse(w, r, param, http.StatusOK)
}

func (h *RouterHub) HandleListRoom(w http.ResponseWriter, r *http.Request) {
	rooms := []model.Room{}

	pID := r.FormValue("id-player")

	for _, r := range h.Rooms {
		rooms = append(rooms, model.Room{
			ID:        r.Room.ID,
			Name:      r.Room.Name,
			Removable: r.Room.OwnerID == pID,
		})
	}

	sort.Slice(rooms, func(i, j int) bool {
		return rooms[i].Name < rooms[j].Name
	})

	api.HttpResponse(w, r, rooms, http.StatusOK)
}

func (h *RouterHub) HandleDetailRoom(w http.ResponseWriter, r *http.Request) {
	var param model.Room

	err := ParseBodyData(r.Context(), r, &param)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	room, ok := h.Rooms[param.ID]
	if !ok {
		api.HttpResponseException(w, r, http.StatusNotFound)
		return
	}

	players := []model.RoomPlayer{}

	for _, p := range room.RoomPlayers {
		pcopy := p.Copy()
		if pAcc, ok := h.Players[p.ID]; ok {
			pcopy.Money = pAcc.Money
		}
		players = append(players, pcopy)
	}

	sort.Slice(players, func(i, j int) bool {
		return players[i].TurnOrder < players[j].TurnOrder
	})

	rm := model.Room{
		ID:      room.Room.ID,
		Name:    room.Room.Name,
		Dealer:  *room.Dealer,
		Players: players,
	}

	api.HttpResponse(w, r, rm, http.StatusOK)
}

func (h *RouterHub) HandleDetailRoomPlayer(w http.ResponseWriter, r *http.Request) {

	pID := r.FormValue("id-player")
	rID := r.FormValue("id-room")

	if pID == "" || rID == "" {
		api.HttpResponseException(w, r, http.StatusBadRequest)
		return
	}

	room, ok := h.Rooms[rID]
	if !ok {
		api.HttpResponseException(w, r, http.StatusNotFound)
		return
	}

	player, ok := room.RoomPlayers[pID]
	if !ok {
		api.HttpResponseException(w, r, http.StatusNotFound)
		return
	}

	api.HttpResponse(w, r, player, http.StatusOK)
}

func (h *RouterHub) HandlePlaceBet(w http.ResponseWriter, r *http.Request) {

	var param model.RoomBet

	err := ParseBodyData(r.Context(), r, &param)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	if param.PlayerID == "" || param.RoomID == "" || param.Bet == 0 {
		api.HttpResponseException(w, r, http.StatusBadRequest)
		return
	}

	p, ok := h.Players[param.PlayerID]
	if !ok {
		api.HttpResponseException(w, r, http.StatusNotFound)
		return
	}

	h.ConnectionMx.Lock()
	p.Money -= param.Bet
	h.ConnectionMx.Unlock()

	room, ok := h.Rooms[param.RoomID]
	if !ok {
		api.HttpResponseException(w, r, http.StatusNotFound)
		return
	}

	player, ok := room.RoomPlayers[param.PlayerID]
	if !ok {
		api.HttpResponseException(w, r, http.StatusNotFound)
		return
	}

	room.ConnectionMx.Lock()
	player.Bet = param.Bet
	player.Status = model.PLAYER_STATUS_SET_BET
	room.ConnectionMx.Unlock()

	room.EventBroadcast <- model.RoomEventData{
		Name: model.ROOM_EVENT_ON_PLAYER_SET_BET,
		Data: model.Player{Name: p.Name},
	}

	if room.isPlayersStatusSame(model.PLAYER_STATUS_SET_BET) {

		go func() {

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

			// second given
			room.givePlayerOneCard(room.Dealer.ID, false)
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

			// set to play
			// set first player turn
			room.ConnectionMx.Lock()
			room.Status = model.ROOM_STATUS_ON_PLAY
			if pTurn, ok := room.RoomPlayers[room.TurnsOrder[room.TurnPost]]; ok {
				pTurn.Status = model.PLAYER_STATUS_AT_TURN
			}
			room.ConnectionMx.Unlock()

			room.EventBroadcast <- model.RoomEventData{
				Name: model.ROOM_EVENT_ON_GAME_START,
			}

		}()

		room.EventBroadcast <- model.RoomEventData{
			Name: model.ROOM_EVENT_ON_GAME_START,
		}

	}

	api.HttpResponse(w, r, true, http.StatusOK)
}

func (h *RouterHub) HandlePlayerActionTurnRoom(w http.ResponseWriter, r *http.Request) {
	var param model.RoomTurn

	err := ParseBodyData(r.Context(), r, &param)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if param.PlayerID == "" || param.RoomID == "" {
		api.HttpResponseException(w, r, http.StatusBadRequest)
		return
	}

	playerAcc, ok := h.Players[param.PlayerID]
	if !ok {
		api.HttpResponseException(w, r, http.StatusNotFound)
		return
	}

	room, ok := h.Rooms[param.RoomID]
	if !ok {
		api.HttpResponseException(w, r, http.StatusNotFound)
		return
	}

	player, ok := room.RoomPlayers[param.PlayerID]
	if !ok {
		api.HttpResponseException(w, r, http.StatusNotFound)
		return
	}

	if param.Choosed == model.ROOM_TURN_CHOOSE_HIT {

		room.givePlayerOneCard(param.PlayerID, true)

		room.EventBroadcast <- model.RoomEventData{
			Name: model.ROOM_EVENT_ON_CARD_GIVEN,
		}

	}

	room.ConnectionMx.Lock()
	defer room.ConnectionMx.Unlock()

	player.Status = model.PLAYER_STATUS_FINISH_TURN
	player.SumUpTotal()

	room.TurnPost++
	if room.TurnPost <= len(room.TurnsOrder)-1 {
		if pTurn, ok := room.RoomPlayers[room.TurnsOrder[room.TurnPost]]; ok {
			pTurn.Status = model.PLAYER_STATUS_AT_TURN
		}
	}

	if player.Total == 21 {

		player.Status = model.PLAYER_STATUS_OUT
		playerAcc.Money += player.Bet * 2
		player.Bet = 0
		room.EventBroadcast <- model.RoomEventData{
			Name: model.ROOM_EVENT_ON_PLAYER_BLACKJACK_WIN,
			Data: model.Player{Name: player.Name},
		}

	} else if player.Total > 21 {

		player.Status = model.PLAYER_STATUS_BUST
		player.Bet = 0
		room.EventBroadcast <- model.RoomEventData{
			Name: model.ROOM_EVENT_ON_PLAYER_BUST,
			Data: model.Player{Name: player.Name},
		}

	} else {

		room.EventBroadcast <- model.RoomEventData{
			Name: model.ROOM_EVENT_ON_PLAYER_END_TURN,
			Data: model.Player{Name: player.Name},
		}
	}

	// end round sum up
	if room.isPlayersStatusSame(model.PLAYER_STATUS_FINISH_TURN) {

		go func() {

			time.Sleep(2 * time.Second)
			h.EndRound(param.RoomID)

			room.EventBroadcast <- model.RoomEventData{
				Name: model.ROOM_EVENT_ON_GAME_END,
			}

			time.Sleep(5 * time.Second)
			room.resetRoom()
			room.EventBroadcast <- model.RoomEventData{
				Name: model.ROOM_EVENT_ON_GAME_START,
			}
		}()

		api.HttpResponse(w, r, true, http.StatusOK)
		return
	}

	api.HttpResponse(w, r, true, http.StatusOK)
}

func (h *RouterHub) HandleRemoveRoom(w http.ResponseWriter, r *http.Request) {
	var param model.DeleteRoom

	err := ParseBodyData(r.Context(), r, &param)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	room, ok := h.Rooms[param.ID]
	if !ok {
		api.HttpResponseException(w, r, http.StatusNotFound)
		return
	}

	cRoom := model.Room{
		ID:   room.Room.ID,
		Name: room.Room.Name,
	}

	if param.PlayerID != room.Room.OwnerID {
		api.HttpResponseException(w, r, http.StatusForbidden)
		return
	}

	if len(room.RoomSubscriber) > 0 {
		api.HttpResponseException(w, r, http.StatusForbidden)
		return
	}

	room.EventBroadcast <- model.RoomEventData{
		Status: model.ROOM_STATUS_NOT_USE,
	}

	h.Lobbies.EventBroadcast <- model.EventData{
		Name: model.LOBBY_EVENT_ON_ROOM_REMOVE,
		Data: cRoom,
	}

	api.HttpResponse(w, r, model.Player{}, http.StatusOK)
}
