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

	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()

	h.openRoom(param.HostID, param.Name, param.Players)

	h.Lobbies.EventBroadcast <- model.EventData{
		Name: model.LOBBY_EVENT_ROOM_CREATED,
		Data: param,
	}

	api.HttpResponse(w, r, param, http.StatusOK)
}

func (h *RouterHub) HandleListRoom(w http.ResponseWriter, r *http.Request) {
	rooms := []model.Room{}

	for _, r := range h.Rooms {
		rooms = append(rooms, model.Room{
			ID:   r.Room.ID,
			Name: r.Room.Name,
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

	api.HttpResponse(w, r, room.Room, http.StatusOK)
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

	if param.PlayerID != room.Room.OwnerID {
		api.HttpResponseException(w, r, http.StatusForbidden)
		return
	}

	room.EventBroadcast <- model.RoomEventData{
		Name:   model.LOBBY_EVENT_ROOM_REMOVE,
		Status: ROOM_STATUS_NOT_USE,
		Data:   nil,
	}

	api.HttpResponse(w, r, model.Player{}, http.StatusOK)
}
