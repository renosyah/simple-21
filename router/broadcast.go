package router

import (
	"time"

	"github.com/renosyah/simple-21/model"
)

func (h *RouterHub) setPlayerOnlineStatus(player model.Player, isOnline, broadcast bool) {

	timeSet := time.Now().Local()
	timeExp := timeSet.Add(time.Hour*time.Duration(0) +
		time.Minute*time.Duration(h.Config.PlayerSessionTime) +
		time.Second*time.Duration(0))

	e := model.LOBBY_EVENT_ON_DISCONNECTED
	if isOnline {
		e = model.LOBBY_EVENT_ON_JOIN
	}

	if broadcast {
		h.Lobbies.EventBroadcast <- model.EventData{Name: e, Data: player}
	}

	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()

	p, ok := h.Players[player.ID]
	if !ok {
		return
	}
	p.IsOnline = isOnline
	if isOnline {
		p.SessionExpired = timeExp
	}
}

func (h *RoomsHub) setPlayerOnlineStatus(player model.Player, isOnline bool, broadcast bool) {

	e := model.ROOM_EVENT_ON_DISCONNECTED
	if isOnline {
		e = model.ROOM_EVENT_ON_JOIN
	}

	if broadcast {
		h.EventBroadcast <- model.RoomEventData{Name: e, Data: player}
	}

	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()

	p, ok := h.RoomPlayers[player.ID]
	if !ok {
		return
	}

	p.IsOnline = isOnline

	// if p.Status == model.PLAYER_STATUS_SET_BET || p.Status == model.PLAYER_STATUS_AT_TURN {
	// 	p.Status = model.PLAYER_STATUS_OUT
	// }

}
