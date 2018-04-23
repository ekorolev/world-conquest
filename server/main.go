package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

func wsHandler(g *Game) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
		if _, ok := err.(websocket.HandshakeError); ok {
			http.Error(w, "Not a websocket handshake", 400)
			return
		} else if err != nil {
			return
		}

		// После подключения к серверу мы ждем от клиента авторизации
		// Как только пользователь присылает имя, мы регистрируем
		// Его в нашей игре и выходим из игры
		for {
			messageType, message, err := ws.ReadMessage()
			parts := strings.Split(string(message), ":")
			if messageType == websocket.TextMessage && parts[0] == "AuthName" {
				p := g.NewPlayer(parts[1], ws)
				ws.WriteMessage(websocket.TextMessage, []byte("Auth:OK"))
				g.SendPlayerToAll(p)
				p.SendMap(g.Map)
				p.SendPlayers(g.Players)
				p.SendAboutInfo()
				break
			} else {
				err = ws.WriteMessage(websocket.TextMessage, []byte("Error:You must authenticate in game"))
				if err != nil {
					break
				}
			}
		}
	}
}

func main() {
	game := NewGame()
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fmt.Sprintf("../client/dist/%s", r.URL.Path[8:]))
	})
	http.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fmt.Sprintf("../client/assets/%s", r.URL.Path[8:]))
	})
	http.HandleFunc("/ws", wsHandler(game))

	if err := http.ListenAndServe(":9001", nil); err != nil {
		log.Fatal("ListenAndServe:", err.Error())
	}
}
