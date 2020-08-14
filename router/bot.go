package router

import (
	"github.com/renosyah/simple-21/model"
)

func (h *RoomsHub) runBotFunction(hub *RouterHub, bot *model.RoomPlayer) {
	go h.botSubBroadcastsEvent(hub, bot)
}

func (r *RoomsHub) botSubBroadcastsEvent(hub *RouterHub, bot *model.RoomPlayer) {
	subReceiver := r.subscribeRoom(bot.ID)
	defer r.unSubscribeRoom(bot.ID)

	for {
		select {
		case msg := <-subReceiver:

			switch msg.Status {
			case model.ROOM_STATUS_CLEAR_BOT:
				return
			default:
			}

			switch msg.Name {
			case model.ROOM_EVENT_ON_JOIN:
				/* this event is for client */
			case model.ROOM_EVENT_ON_PLAYER_SET_BET:
				/* this event is for client */
			case model.ROOM_EVENT_ON_GAME_START:
				/* this event is for client */
			case model.ROOM_EVENT_ON_CARD_GIVEN:
				/* this event is for client */
			case model.ROOM_EVENT_ON_PLAYER_END_TURN:
				/* this event is for client */
			case model.ROOM_EVENT_ON_PLAYER_BLACKJACK_WIN:
				/* this event is for client */
			case model.ROOM_EVENT_ON_PLAYER_BUST:
				/* this event is for client */
			case model.ROOM_EVENT_ON_GAME_END:
				/* this event is for client */
			case model.ROOM_EVENT_ON_DISCONNECTED:
				/* this event is for client */
			default:
			}
		}
	}
}
