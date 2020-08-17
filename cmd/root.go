package cmd

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/renosyah/simple-21/model"
	"github.com/renosyah/simple-21/router"
)

var (
	routerHub *router.RouterHub
	cfgFile   string
)

var rootCmd = &cobra.Command{
	Use: "app",
	PreRun: func(cmd *cobra.Command, args []string) {

	},
	Run: func(cmd *cobra.Command, args []string) {

		rand.Seed(time.Now().UTC().UnixNano())
		r := mux.NewRouter()

		r.HandleFunc("/addplayer", routerHub.HandleAddPlayer)
		r.HandleFunc("/player", routerHub.HandleDetailPlayer)
		r.HandleFunc("/player/money", routerHub.HandleDetailPlayerMoney)
		r.HandleFunc("/players", routerHub.HandleListPlayer)
		r.HandleFunc("/delplayer", routerHub.HandleRemovePlayer)

		r.HandleFunc("/addroom", routerHub.HandleAddRoom)
		r.HandleFunc("/room", routerHub.HandleDetailRoom)
		r.HandleFunc("/room/player", routerHub.HandleDetailRoomPlayer)
		r.HandleFunc("/room/setbet", routerHub.HandlePlaceBet)
		r.HandleFunc("/room/action", routerHub.HandlePlayerActionTurnRoom)
		r.HandleFunc("/room/scores", routerHub.HandleListRoomScore)
		r.HandleFunc("/rooms", routerHub.HandleListRoom)
		r.HandleFunc("/delroom", routerHub.HandleRemoveRoom)

		r.HandleFunc("/moneys", routerHub.HandleListMoney)
		r.HandleFunc("/money/buy", routerHub.HandlePurchaseMoney)

		r.HandleFunc("/random-name", router.HandleGetRandomName)
		r.HandleFunc("/card-group", router.HandleGetCardsGroup)

		r.HandleFunc("/ws-lobby", routerHub.HandleStreamLobby)
		r.HandleFunc("/ws-room", routerHub.HandleStreamRoom)

		r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("files"))))

		r.NotFoundHandler = r.NewRoute().HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}).GetHandler()

		port := viper.GetInt("app.port")
		p := os.Getenv("PORT")
		if p != "" {
			port, _ = strconv.Atoi(p)
		}

		server := &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			Handler:      r,
			ReadTimeout:  time.Duration(viper.GetInt("read_timeout")) * time.Second,
			WriteTimeout: time.Duration(viper.GetInt("write_timeout")) * time.Second,
			IdleTimeout:  time.Duration(viper.GetInt("idle_timeout")) * time.Second,
		}

		done := make(chan bool, 1)
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, os.Interrupt)

		go func() {
			<-quit
			log.Println("Server is shutting down...")

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			server.SetKeepAlivesEnabled(false)
			if err := server.Shutdown(ctx); err != nil {
				log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
			}
			close(done)
		}()

		log.Println("Server is ready to handle requests at", fmt.Sprintf(":%d", port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", fmt.Sprintf(":%d", port), err)
		}

		<-done
		log.Println("Server stopped")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is github.com/renosyah/simple-21/.server.toml)")
	cobra.OnInitialize(initConfig, initRouterHub)
}

func initRouterHub() {
	routerHub = router.NewRouterHub(model.GameConfig{
		MaxPlayer:         viper.GetInt("game.player"),
		MaxRoom:           viper.GetInt("game.room"),
		StarterMoney:      viper.GetInt("game.money"),
		PlayerSessionTime: viper.GetInt("game.player_session"),
		RoomSessionTime:   viper.GetInt("game.room_session"),
	})
}

func initConfig() {
	viper.SetConfigType("toml")
	if cfgFile != "" {

		viper.SetConfigFile(cfgFile)
	} else {

		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.AddConfigPath("/etc/simple-21")
		viper.SetConfigName(".server")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
