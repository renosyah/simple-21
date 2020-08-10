package router

import (
	"net/http"
	"sort"

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
		players = append(players, *p)
	}

	sort.Slice(players, func(i, j int) bool {
		return players[i].Name < players[j].Name
	})

	rm := model.Room{
		ID:      room.Room.ID,
		Name:    room.Room.Name,
		Dealer:  *room.Dealer,
		Players: players,
	}

	api.HttpResponse(w, r, rm, http.StatusOK)
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

	room.EventBroadcast <- model.RoomEventData{
		Status: model.ROOM_STATUS_NOT_USE,
	}

	h.Lobbies.EventBroadcast <- model.EventData{
		Name: model.LOBBY_EVENT_ON_ROOM_REMOVE,
		Data: cRoom,
	}

	api.HttpResponse(w, r, model.Player{}, http.StatusOK)
}
