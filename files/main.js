new Vue({
    el: '#app',
    data() {
        return {
            player : {
                id : "",
                name : "",
                money : ""
            },
            is_online : true,
            is_loading : false,
            host : {
                name : "",
                protocol : "",
                port : ""
            }
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
        window.$('.modal').modal()
        this.initPlayer()
    },
    computed : {
        getPageName(){
            let param = new URLSearchParams(window.location.search)
            let name = param.get('page')
            return name ? name : "login-page"
        },
    },
    methods : {
        register(nm){

            this.is_loading = true

            axios.post(this.baseUrl() + "register", JSON.stringify({name : nm}))
                .then(response => {

                    this.is_loading = false
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

                })
                .catch(e => {
                    console.log(e)
                    this.is_loading = false
                })

        },
        exit(){

            this.is_loading = true

            axios.post(this.baseUrl() + "logout", JSON.stringify(this.player))
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
        switchPage(name){
            if ('URLSearchParams' in window) {
                var searchParams = new URLSearchParams(window.location.search);
                searchParams.set('page', name);
                window.location.search = searchParams.toString();
            }
        },
        switchLocalize(lang){
            if ('URLSearchParams' in window) {
                var searchParams = new URLSearchParams(window.location.search);
                searchParams.set('loc', name);
                window.location.search = searchParams.toString();
            }
        },
        localizeCheck(lang){
            let def = "ind"
            let param = new URLSearchParams(window.location.search)
            let name = param.get('loc')
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
        },
        baseUrl(){
            return this.host.protocol.concat(this.host.name + ":" + this.host.port + "/")
        },
    }
})
