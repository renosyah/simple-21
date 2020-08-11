package router

import (
	"net/http"

	"github.com/renosyah/simple-21/api"
	"github.com/renosyah/simple-21/model"
)

func (h *RouterHub) HandleListMoney(w http.ResponseWriter, r *http.Request) {
	pID := r.FormValue("id-player")

	if pID == "" {
		api.HttpResponseException(w, r, http.StatusBadRequest)
		return
	}

	_, ok := h.Players[pID]
	if !ok {
		api.HttpResponseException(w, r, http.StatusNotFound)
		return
	}

	api.HttpResponse(w, r, h.ListMoneyShops, http.StatusOK)
}

func (h *RouterHub) HandlePurchaseMoney(w http.ResponseWriter, r *http.Request) {
	var param model.PuchaseMoney

	err := ParseBodyData(r.Context(), r, &param)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p, ok := h.Players[param.PlayerID]
	if !ok {
		api.HttpResponseException(w, r, http.StatusNotFound)
		return
	}

	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()

	if m, ok := h.ListMoneyShops[param.ID]; ok {
		p.Money += m.Amount
	}

	api.HttpResponse(w, r, p, http.StatusOK)
}
