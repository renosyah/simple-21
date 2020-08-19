package router

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/renosyah/simple-21/api"
	"github.com/renosyah/simple-21/model"
	uuid "github.com/satori/go.uuid"
)

func (h *RouterHub) HandleAddPlayer(w http.ResponseWriter, r *http.Request) {
	var param model.Player

	err := ParseBodyData(r.Context(), r, &param)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(h.Players) >= h.Config.MaxPlayer {
		api.HttpResponseException(w, r, http.StatusInsufficientStorage)
		return
	}

	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()

	param.ID = fmt.Sprint(uuid.NewV4())
	param.Money = h.Config.StarterMoney
	param.IsOnline = true
	param.SessionExpired = h.createPlayerSessionTime()

	h.Players[param.ID] = &param

	api.HttpResponse(w, r, param, http.StatusOK)
}

func (h *RouterHub) HandleListPlayer(w http.ResponseWriter, r *http.Request) {
	players := []model.Player{}

	for _, p := range h.Players {
		players = append(players, *p)
	}

	sort.Slice(players, func(i, j int) bool {
		return players[i].Name < players[j].Name
	})

	api.HttpResponse(w, r, players, http.StatusOK)
}

func (h *RouterHub) HandleDetailPlayer(w http.ResponseWriter, r *http.Request) {
	var param model.Player

	err := ParseBodyData(r.Context(), r, &param)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p, ok := h.Players[param.ID]
	if !ok {
		api.HttpResponseException(w, r, http.StatusNotFound)
		return
	}

	api.HttpResponse(w, r, p, http.StatusOK)
}

func (h *RouterHub) HandleDetailPlayerMoney(w http.ResponseWriter, r *http.Request) {
	pID := r.FormValue("id-player")

	if pID == "" {
		api.HttpResponseException(w, r, http.StatusBadRequest)
		return
	}

	p, ok := h.Players[pID]
	if !ok {
		api.HttpResponseException(w, r, http.StatusNotFound)
		return
	}

	api.HttpResponse(w, r, p.Money, http.StatusOK)
}

func (h *RouterHub) HandleRemovePlayer(w http.ResponseWriter, r *http.Request) {
	var param model.Player

	err := ParseBodyData(r.Context(), r, &param)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()

	if _, ok := h.Players[param.ID]; ok {
		delete(h.Players, param.ID)
	}

	h.Lobbies.EventBroadcast <- model.EventData{
		Name: model.LOBBY_EVENT_ON_LOGOUT,
		Data: param,
	}

	api.HttpResponse(w, r, model.Player{}, http.StatusOK)
}
