package router

import (
	"net/http"
	"sort"

	"github.com/renosyah/simple-21/api"
	"github.com/renosyah/simple-21/model"
)

func (h *RouterHub) HandleListRoomScore(w http.ResponseWriter, r *http.Request) {
	rID := r.FormValue("id-room")

	if rID == "" {
		api.HttpResponseException(w, r, http.StatusBadRequest)
		return
	}

	room, ok := h.Rooms[rID]
	if !ok {
		api.HttpResponseException(w, r, http.StatusBadRequest)
		return
	}

	scores := []model.Score{}
	for _, sc := range room.Scores {
		scores = append(scores, sc.Copy())
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Round < scores[j].Round
	})

	api.HttpResponse(w, r, scores, http.StatusOK)
}
