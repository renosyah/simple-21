new Vue({
    el: '#app',
    data() {
        return {          
            page_name : "main-page",
            is_online : true,
            is_loading : false,
            host : {
                name : "",
                protocol : "",
                port : ""
            },
            localize : {
                choosed : "ind",
                list : [{
                    label :"Indonesia" ,
                    value :"ind" 
                },{
                    label :"English" ,
                    value :"en" 
                }]
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
    },
    computed : {
        baseUrl(){
            return this.host.protocol.concat(this.host.name + ":" + this.host.port + "/")
        }
    },
    methods : {
        switchPage(name){
            this.page_name = name
        },
        backPress(){
            if (event.state && event.state.noBackExitsApp) {
                window.history.pushState({ noBackExitsApp: true }, '')
            }
        },
        localizeCheck(lang){ 
            return this.localize.choosed == lang 
        },
        setCurrentHost(){
            this.host.name = window.location.hostname
            this.host.port = location.port
            this.host.protocol = location.protocol.concat("//")
        }
    }
})
