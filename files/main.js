
// FOR BROADCASTER RECEIVER
const LOBBY_EVENT_ON_JOIN         = "LOBBY_EVENT_ON_PLAYER_JOIN"
const LOBBY_EVENT_ON_DISCONNECTED = "LOBBY_EVENT_ON_PLAYER_DISCONNECTED"
const LOBBY_EVENT_ON_LOGOUT       = "LOBBY_EVENT_ON_PLAYER_LOGOUT"

const ROOM_EVENT_ON_JOIN          = "ROOM_EVENT_ON_PLAYER_JOIN"
const ROOM_EVENT_ON_DISCONNECTED  = "ROOM_EVENT_ON_PLAYER_DISCONNECTED"
const LOBBY_EVENT_ON_ROOM_CREATED = "LOBBY_EVENT_ON_ROOM_CREATED"
const LOBBY_EVENT_ON_ROOM_REMOVE  = "LOBBY_EVENT_ON_ROOM_REMOVED"
const ROOM_EVENT_ON_PLAYER_SET_BET  = "ROOM_EVENT_ON_PLAYER_SET_BET"
const ROOM_EVENT_ON_PLAYER_END_TURN = "ROOM_EVENT_ON_PLAYER_END_TURN"
const ROOM_EVENT_ON_GAME_END = "ROOM_EVENT_ON_GAME_END"
const ROOM_EVENT_ON_PLAYER_BLACKJACK_WIN = "ROOM_EVENT_ON_PLAYER_BLACKJACK_WIN"
const ROOM_EVENT_ON_PLAYER_BUST = "ROOM_EVENT_ON_PLAYER_BUST"
const ROOM_EVENT_ON_GAME_START   = "ROOM_EVENT_ON_GAME_START"
const ROOM_EVENT_ON_CARD_GIVEN   = "ROOM_EVENT_ON_CARD_GIVEN"

// FOR ROOM STATUS
const ROOM_STATUS_USE     = 0
const ROOM_STATUS_ON_PLAY = 1
const ROOM_STATUS_NOT_USE = 2

// FOR ROOM PLAYER AS VALUE
const AS_VALUE_ELEVEN = 0
const AS_VALUE_ONE    = 1

// FOR ROOM PLAYER STATUS

const PLAYER_STATUS_SPECTATE = -1
const PLAYER_STATUS_INVITED = 0
const PLAYER_STATUS_SET_BET = 1
const PLAYER_STATUS_IDLE    = 2
const PLAYER_STATUS_AT_TURN = 3
const PLAYER_STATUS_FINISH_TURN = 4
const PLAYER_STATUS_OUT     = 5
const PLAYER_STATUS_BUST    = 6
const PLAYER_STATUS_REWARDED= 7
const PLAYER_STATUS_LOSE    = 8

// TOAST
const TOAST_OUT_DURRATION = 60

new Vue({
    el: '#app',
    data() {
        return {
            players : [],
            player : {
                id : "",
                name : "",
                money : 0
            },
            rooms : [],
            room : {
                id:"",
                name : "",
                dealer : {
                    id:"",
                    name:"",
                    bet:0,
                    cards : [],
                    total_show : 0,
                    total : 0,
                    status :0
                },
                players : [],
                status : 0,
            },
            player_in_room : {
                id:"",
                name:"",
                bet:0,
                cards : [],
                total_show : 0,
                total : 0,
                status :0                    
            },
            bet_holder : 50,
            add_room : {
                host_id:"",
                name:"",
                players:[],
                card_groups : []
            },
            moneys : [],
            lobby_ws : null,
            room_ws : null,
            is_online : true,
            is_loading : false,
            host : {
                name : "",
                protocol : "",
                port : "",
                ws_protocol : "",
            },
            random_name :"",
            card_groups : []
        }
    },
    created(){
        window.addEventListener('offline', () => { this.is_online = false })
        window.addEventListener('online', () => { this.is_online = true })
        window.history.pushState({ noBackExitsApp: true }, '')
        window.addEventListener('popstate', this.backPress )
        this.setCurrentHost()
    },
    mounted () {
        window.$('.dropdown-trigger').dropdown()
        window.$('.modal').modal({opacity:0.1,dismissible: false,preventScrolling:false})
        window.$('.sidenav').sidenav();
        this.initPlayer()  
        this.randomName("yes")      
    },
    watch: {
        player_in_room : (val) => {
            window.$("#modal-room-action").modal(val.status == PLAYER_STATUS_AT_TURN ? 'open' : 'close');
            window.$("#modal-room-bet").modal(val.status == PLAYER_STATUS_INVITED || val.status == PLAYER_STATUS_IDLE ? 'open' : 'close');    
            window.$("#modal-room-win").modal(val.status == PLAYER_STATUS_REWARDED ? 'open' : 'close');
            window.$("#modal-room-lose").modal(val.status == PLAYER_STATUS_LOSE ? 'open' : 'close');
        }
    },
    computed : {
        getPageName(){
            let param = new URLSearchParams(window.location.search)
            let name = param.get('page')
            return name ? name : "login-page"
        }
    },
    methods : {
        addPlayer(){
            this.is_loading = true
            axios.post(this.baseUrl() + "addplayer", JSON.stringify(this.player))
                .then(response => {
                    this.is_loading = false
                    if (response.data.status == 507) {
                        window.M.toast({outDuration : TOAST_OUT_DURRATION,html: '<b>Server is Full!</b>', classes: 'white green-text'})
                        return
                    }
                    this.player = response.data.result
                    if ('URLSearchParams' in window) {
                        var searchParams = new URLSearchParams(window.location.search);
                        searchParams.set('page','main-page');
                        searchParams.set('player-id',this.player.id);
                        window.location.search = searchParams.toString();
                    }
                })
                .catch(e => {
                    console.log(e)
                    this.is_loading = false
                })
        },
        initPlayer(){
            let param = new URLSearchParams(window.location.search)
            if (!param.get('player-id')){
                return
            }
            this.is_loading = true
            axios.post(this.baseUrl() + "player", JSON.stringify({id : param.get('player-id')}))
                .then(response => {
                    this.is_loading = false
                    if (response.data.status == 404) {
                        window.location = this.baseUrl()
                        return
                    }
                    this.player = response.data.result
                    if (param.get('page') == 'main-page'){
                        this.initLobbyWs()
                    } else if (param.get('page') == 'room-page'){
                        this.initRoomWs(param.get('id-room'))
                    }                        
                })
                .catch(e => {
                    console.log(e)
                    this.is_loading = false
                })
        },
        getPlayers(){
            axios.get(this.baseUrl() + "players")
                .then(response => {
                    this.players = response.data.result
                })
                .catch(e => {
                    console.log(e)
                })
        },
        getPlayerMoney(){
            axios.get(this.baseUrl() + "player/money" + "?id-player=" + this.player.id)
                .then(response => {
                    this.player.money = response.data.result
                })
                .catch(e => {
                    console.log(e)
                })
        },
        exit(){
            this.lobby_ws.close()
            this.is_loading = true
            axios.post(this.baseUrl() + "delplayer", JSON.stringify(this.player))
                .then(response => {
                    this.is_loading = false
                    this.player = response.data.result
                    window.location = this.baseUrl()
                })
                .catch(e => {
                    console.log(e)
                    this.is_loading = false
                })
        },
        addRoom(){
            this.add_room.host_id = this.player.id
            this.add_room.players.unshift(this.player)
            this.is_loading = true
            axios.post(this.baseUrl() + "addroom",JSON.stringify(this.add_room))
                .then(response => {
                    this.is_loading = false
                    this.add_room = {
                        host_id:"",
                        name:"",
                        players:[]
                    }

                    if (response.data.status != 200) {
                        window.M.toast({outDuration : TOAST_OUT_DURRATION,html: '<b>Failed Create Room!</b>', classes: 'white green-text'})
                        return
                    }
                })
                .catch(e => {
                    this.is_loading = false
                    console.log(e)
                })
        },
        getRoom(roomId){
            axios.post(this.baseUrl() + "room",JSON.stringify({id:roomId}))
                .then(response => {
                    if (response.data.status == 404) {
                        window.M.toast({outDuration : TOAST_OUT_DURRATION,html: '<b>Room Not Found!</b>', classes: 'white green-text'})
                        return
                    }
                    this.room = response.data.result
                    this.getPlayerInRoom(roomId)
                })
                .catch(e => {
                    console.log(e)
                })
        },
        getPlayerInRoom(roomId){
            axios.get(this.baseUrl() + "room/player" + "?id-player="+this.player.id + "&id-room="+roomId)
                .then(response => {
                    if (response.data.status == 404) {
                        return
                    }
                    this.player_in_room = response.data.result
                    this.getPlayerMoney()
                })
                .catch(e => {
                    console.log(e)
                })
        },
        getRooms(){
            axios.get(this.baseUrl() + "rooms" +"?id-player=" + this.player.id)
                .then(response => {
                    this.rooms = response.data.result
                })
                .catch(e => {
                    console.log(e)
                })
        },
        setBet(){
            this.is_loading = true
            axios.post(this.baseUrl() + "room/setbet",JSON.stringify({player_id : this.player.id, room_id : this.room.id, bet : parseInt(this.bet_holder)}))
                .then(response => {
                    this.is_loading = false
                    if (response.data.status == 404) {
                        window.M.toast({outDuration : TOAST_OUT_DURRATION,html: '<b>Room Not Found!</b>', classes: 'white green-text'})
                        return
                    }
                    this.bet_holder = this.player.money > 50 ? 50 : 0 
                })
                .catch(e => {
                    this.is_loading = false
                    console.log(e)
                })                
        },
        setAction(action){
            this.is_loading = true
            axios.post(this.baseUrl() + "room/action",JSON.stringify({player_id : this.player.id, room_id : this.room.id, choosed : action}))
                .then(response => {
                    this.is_loading = false
                    if (response.data.status == 404) {
                        window.M.toast({outDuration : TOAST_OUT_DURRATION,html: '<b>Room Not Found!</b>', classes: 'white green-text'})
                        return
                    }
                })
                .catch(e => {
                    this.is_loading = false
                    console.log(e)
                })                
        },
        deleteRoom(roomId){
            this.is_loading = true
            axios.post(this.baseUrl() + "delroom",JSON.stringify({id:roomId, player_id:this.player.id}))
                .then(response => {
                    this.is_loading = false
                    if (response.data.status != 200) {
                        window.M.toast({outDuration : TOAST_OUT_DURRATION,html: '<b>Failed Delete Room!</b>', classes: 'white green-text'})
                        return
                    }
                })
                .catch(e => {
                    this.is_loading = false
                    console.log(e)
                })
        },
        getMoneys(){
            axios.get(this.baseUrl() + "moneys" +"?id-player=" + this.player.id)
                .then(response => {
                    this.moneys = response.data.result
                })
                .catch(e => {
                    console.log(e)
                })
        },
        buyMoney(idMoney){
            this.is_loading = true
            axios.post(this.baseUrl() + "money/buy",JSON.stringify({id:idMoney,player_id:this.player.id}))
                .then(response => {
                    this.is_loading = false
                    if (response.data.status != 200) {
                        window.M.toast({outDuration : TOAST_OUT_DURRATION,html: '<b>Failed Purchase!</b>', classes: 'white green-text'})
                        return
                    }
                    this.player.money = response.data.result.money
                })
                .catch(e => {
                    this.is_loading = false
                    console.log(e)
                })
        },
        initLobbyWs(){
            this.lobby_ws = new WebSocket(this.baseWsUrl() + "ws-lobby" + "?id-player=" + this.player.id)
            this.lobby_ws.onmessage = (evt) => {
                let event = JSON.parse(evt.data)
                switch (event.name) {
                    case LOBBY_EVENT_ON_JOIN: 
                        window.M.toast({outDuration : TOAST_OUT_DURRATION,html: '<b>' + event.data.name + ' is Join!</b>', classes: 'white green-text'})
                        this.getPlayers()
                        break;
                    case LOBBY_EVENT_ON_LOGOUT: 
                        window.M.toast({outDuration : TOAST_OUT_DURRATION,html: '<b>' + event.data.name + ' is Logout!</b>', classes: 'white green-text'})
                        this.getPlayers()
                        break;
                    case LOBBY_EVENT_ON_DISCONNECTED:
                        window.M.toast({outDuration : TOAST_OUT_DURRATION,html: '<b>' + event.data.name + ' is Offline!</b>', classes: 'white green-text'})
                        this.getPlayers()
                        break;
                    case LOBBY_EVENT_ON_ROOM_CREATED:
                        window.M.toast({outDuration : TOAST_OUT_DURRATION,html: '<b>' + event.data.name + ' is Created!</b>', classes: 'white green-text'})
                        this.getRooms()
                        break;
                        case LOBBY_EVENT_ON_ROOM_REMOVE:
                        window.M.toast({outDuration : TOAST_OUT_DURRATION,html: '<b>' + event.data.name + ' is Removed!</b>', classes: 'white green-text'})
                        this.getRooms()
                        break;
                    default: break;
                }
            }
            this.lobby_ws.onopen = () => {
                this.getPlayers()
                this.getRooms()
                this.getMoneys()
                this.cardsGroup()                     
            }
            this.lobby_ws.onerror = (e) => {
                console.log(e)
                this.is_online = false
            }
        },
        initRoomWs(idR){
            let idRoom = idR
            this.room_ws = new WebSocket(this.baseWsUrl() + "ws-room" + "?id-player=" + this.player.id + "&id-room=" + idRoom)
            this.room_ws.onmessage = (evt) => {
                let event = JSON.parse(evt.data)
                switch (event.name) {
                    case ROOM_EVENT_ON_JOIN:
                        window.M.toast({outDuration : TOAST_OUT_DURRATION,html: '<b>' + event.data.name + ' is Join the room!</b>', classes: 'white green-text'})
                        this.getRoom(idRoom)
                        break;
                    case ROOM_EVENT_ON_PLAYER_SET_BET:
                        //window.M.toast({outDuration : TOAST_OUT_DURRATION,html: '<b>' + event.data.name + ' has Set his bet!</b>', classes: 'white green-text'})
                        this.getRoom(idRoom)
                        break;
                    case ROOM_EVENT_ON_GAME_START:
                        if (this.player_in_room.status == PLAYER_STATUS_SPECTATE){
                            window.$("#modal-room-score").modal('close');
                        }  
                        window.M.toast({outDuration : TOAST_OUT_DURRATION,html: '<b>game is started!</b>', classes: 'white green-text'})
                        this.getRoom(idRoom)
                        break;
                    case ROOM_EVENT_ON_CARD_GIVEN:
                        //window.M.toast({outDuration : TOAST_OUT_DURRATION,html: '<b>Dealer given card</b>', classes: 'white green-text'})
                        this.getRoom(idRoom)
                        break;
                    case ROOM_EVENT_ON_PLAYER_END_TURN:
                        //window.M.toast({outDuration : TOAST_OUT_DURRATION,html: '<b>' + event.data.name + ' has End Turn!</b>', classes: 'white green-text'})
                        this.getRoom(idRoom)
                        break
                    case ROOM_EVENT_ON_PLAYER_BLACKJACK_WIN:
                        window.M.toast({outDuration : TOAST_OUT_DURRATION,html: '<b>' + event.data.name + ' Total Card is 21!</b>', classes: 'white green-text'})
                        this.getRoom(idRoom)
                        break;
                    case ROOM_EVENT_ON_PLAYER_BUST:
                        window.M.toast({outDuration : TOAST_OUT_DURRATION,html: '<b>' + event.data.name + ' has Bust!</b>', classes: 'white green-text'})
                        this.getRoom(idRoom)
                        break
                    case ROOM_EVENT_ON_GAME_END:
                        window.M.toast({outDuration : TOAST_OUT_DURRATION,html: '<b>Round is End!</b>', classes: 'white green-text'})
                        if (this.player_in_room.status == PLAYER_STATUS_SPECTATE){
                            window.$("#modal-room-score").modal('open');
                        }
                        this.getRoom(idRoom)
                        break;
                    case ROOM_EVENT_ON_DISCONNECTED:
                        window.M.toast({outDuration : TOAST_OUT_DURRATION,html: '<b>' + event.data.name + ' has Exit room!</b>', classes: 'white green-text'})
                        this.getRoom(idRoom)
                        break;
                    default: break;
                }
            }
            this.room_ws.onopen = () => {
                this.getRoom(idRoom)
            };
            this.room_ws.onclose = () => {
                    
            };
            this.room_ws.onerror = (e) => {
                console.log(e)
                this.is_online = false
            };
        },
        randomName(title){
            axios.get(this.baseUrl() + "random-name" + "?title="+title)
                .then(response => {
                    this.random_name = response.data.result
                })
                .catch(e => {
                    console.log(e)
                })
        },
        cardsGroup(){
            axios.get(this.baseUrl() + "card-group")
                .then(response => {
                    this.card_groups = response.data.result
                    this.card_groups.forEach(element => {
                        this.add_room.card_groups.push(element)
                    });
                })
                .catch(e => {
                    console.log(e)
                })
        },
        toRoom(idRoom){
            if ('URLSearchParams' in window) {
                var searchParams = new URLSearchParams(window.location.search);
                searchParams.set('id-room', idRoom);
                searchParams.set('page', 'room-page');
                window.location.search = searchParams.toString();
            }
        },
        getPlayerStatus(p){
            return p.status == PLAYER_STATUS_INVITED ?
            'Ready' :  p.status == PLAYER_STATUS_SET_BET ? 
            'Set Bet' : p.status == PLAYER_STATUS_IDLE ? 
            'Idle' : p.status == PLAYER_STATUS_AT_TURN ? 
            'In Turn' : p.status == PLAYER_STATUS_FINISH_TURN ? 
            'End Turn' :  p.status == PLAYER_STATUS_OUT ? 
            'Out' :  p.status == PLAYER_STATUS_BUST ? 
            'Bust' : p.status == PLAYER_STATUS_REWARDED ? 
            'Win' : 'Lose' 
        },
        isPlayerInRoom(roomId){
            
            let exist = false
            let roomPos = 0
            let i;
            for (i = 0; i < this.rooms.length; i++) {
                if (this.rooms[i].id == roomId){
                    roomPos = i
                    break;
                }
            }
            i = 0;
            let players = this.rooms[roomPos].players
            for (i = 0; i < players.length; i++) {
                if (players[i].id == this.player.id){
                    exist = true
                    break;
                }
            }

            return exist
        },
        switchPage(name){
            if ('URLSearchParams' in window) {
                var searchParams = new URLSearchParams(window.location.search);
                searchParams.set('page', name);
                window.location.search = searchParams.toString();
            }
        },
        switchLang(lang){
            if ('URLSearchParams' in window) {
                var searchParams = new URLSearchParams(window.location.search);
                searchParams.set('lang', lang);
                window.location.search = searchParams.toString();
            }
        },
        langCheck(lang){
            let def = "ind"
            let param = new URLSearchParams(window.location.search)
            let name = param.get('lang')
            return name ? (name == lang) : (def == lang)
        },
        backPress(){
            if (event.state && event.state.noBackExitsApp) {
                window.history.pushState({ noBackExitsApp: true }, '')
            }
        },
        setCurrentHost(){
            this.host.name = window.location.hostname
            this.host.port = location.port
            this.host.protocol = location.protocol.concat("//")
            this.host.ws_protocol = this.host.protocol == "https://" ? "wss://"  : "ws://" 
        },
        baseUrl(){
            return this.host.protocol.concat(this.host.name + ":" + this.host.port + "/")
        },
        baseWsUrl(){
            return this.host.ws_protocol.concat(this.host.name + ":" + this.host.port + "/")
        }
    }
})

