package router

import (
	"math/rand"
	"time"

	"github.com/renosyah/simple-21/model"
	"github.com/renosyah/simple-21/util"
)

func (h *RoomsHub) runBotFunction(hub *RouterHub, bot *model.RoomPlayer) {
	go h.botSubBroadcastsEvent(hub, bot)
}

func (r *RoomsHub) botSubBroadcastsEvent(hub *RouterHub, bot *model.RoomPlayer) {
	subReceiver := r.subscribeRoom(bot.ID)
	defer r.unSubscribeRoom(bot.ID)

	go func() {

		botDecision := util.NewBoolgen()

		for {
			switch bot.Status {
			case model.PLAYER_STATUS_OUT:
				return
			case model.PLAYER_STATUS_AT_TURN:

				time.Sleep(2 * time.Second)

				if botDecision.Bool() {
					bot.Status = model.PLAYER_STATUS_SET_BET
					r.givePlayerOneCard(bot.ID, true)
					evt := r.blackjackForEvt(bot.ID, model.ROOM_EVENT_ON_CARD_GIVEN)

					if evt == model.ROOM_EVENT_ON_PLAYER_BUST || evt == model.ROOM_EVENT_ON_PLAYER_BLACKJACK_WIN {
						r.removeFromTurnOrder(bot.ID)
						r.Turn.TurnPost--
					}
					r.EventBroadcast <- model.RoomEventData{
						Name: evt,
						Data: model.Player{ID: bot.ID, Name: bot.Name},
					}

				} else {

					bot.Status = model.PLAYER_STATUS_FINISH_TURN
					r.removeFromTurnOrder(bot.ID)
					r.Turn.TurnPost--

				}

				r.nextTurnOrder()
				r.EventBroadcast <- model.RoomEventData{
					Name: model.ROOM_EVENT_ON_PLAYER_END_TURN,
					Data: model.Player{Name: bot.Name},
				}

				if r.isPlayersStatusSame(model.PLAYER_STATUS_FINISH_TURN) {
					hub.allPlayerTurnFinish(r)
				}

			default:
			}
		}
	}()

	for {
		select {
		case msg := <-subReceiver:

			switch msg.Name {
			case model.ROOM_EVENT_ON_BOT_REMOVE:
				bot.Status = model.PLAYER_STATUS_OUT
				return

			case model.ROOM_EVENT_ON_JOIN:
				if bot.Status == model.PLAYER_STATUS_INVITED || bot.Status == model.PLAYER_STATUS_IDLE {
					bet := rand.Intn(500-50) + 50

					r.ConnectionMx.Lock()
					if bot.Money < 50 {
						bot.Money += 500
					}
					bot.Money = (bot.Money - bet)
					bot.Bet = bet
					bot.Status = model.PLAYER_STATUS_SET_BET
					r.ConnectionMx.Unlock()
				}

			case model.ROOM_EVENT_ON_PLAYER_SET_BET:
				/* this event is for client */
			case model.ROOM_EVENT_ON_GAME_START:

				if bot.Status == model.PLAYER_STATUS_INVITED || bot.Status == model.PLAYER_STATUS_IDLE {
					bet := rand.Intn(500-50) + 50

					r.ConnectionMx.Lock()
					if bot.Money < 50 {
						bot.Money += 500
					}
					bot.Money = (bot.Money - bet)
					bot.Bet = bet
					bot.Status = model.PLAYER_STATUS_SET_BET
					r.ConnectionMx.Unlock()
				}

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
