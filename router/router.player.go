package router

import (
	"fmt"
	"net/http"

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

	param.ID = fmt.Sprint(uuid.NewV4())
	param.Money = 500

	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()

	h.PlayersConn[param.ID] = &PlayerConn{
		Player: param,
	}

	api.HttpResponse(w, r, param, http.StatusOK)
}

func (h *RouterHub) HandleRemovePlayer(w http.ResponseWriter, r *http.Request) {
	var param model.Player

	err := ParseBodyData(r.Context(), r, &param)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h.removePlayerConnection(param.ID)
	api.HttpResponse(w, r, param, http.StatusOK)
}
