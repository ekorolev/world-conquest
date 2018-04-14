package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alehano/wsgame/utils"
	"github.com/gorilla/websocket"
)

type Player struct {
	Id      string
	Name    string
	Connect *websocket.Conn
}

type Game struct {
	Players map[string]*Player
}

func NewGame() *Game {
	game := new(Game)
	game.Players = make(map[string]*Player)
	return game
}

func (g *Game) NewPlayer(ws *websocket.Conn) *Player {
	pl := new(Player)
	pl.Connect = ws
	pl.Name = ""
	pl.Id = utils.RandString(16)
	return pl
}

func (g *Game) SendAll(message string) {
	for _, player := range g.Players {
		err := player.Connect.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Printf("Player %s is pidor: %s", player.Name, err.Error())
		}
	}
}

func wsHandler(g *Game) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "false")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
		if _, ok := err.(websocket.HandshakeError); ok {
			http.Error(w, "Not a websocket handshake", 400)
			return
		} else if err != nil {
			return
		}

		log.Println("New connection!")
		g.NewPlayer(ws)
	}
}

func main() {
	game := NewGame()
	fmt.Println("Server started, ok?")
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path[1:])
		fmt.Println(fmt.Sprintf("../client/dist/%s", r.URL.Path[8:]))
		http.ServeFile(w, r, fmt.Sprintf("../client/dist/%s", r.URL.Path[8:]))
	})
	http.HandleFunc("/ws", wsHandler(game))

	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for _ = range ticker.C {
			log.Println("Send message to all")
			game.SendAll("Hello!")
		}
	}()

	if err := http.ListenAndServe(":9001", nil); err != nil {
		log.Fatal("ListenAndServe:", err.Error())
	}
}
