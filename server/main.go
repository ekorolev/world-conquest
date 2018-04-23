package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

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

func wsStatisticHandler(g *Game) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
		if _, ok := err.(websocket.HandshakeError); ok {
			http.Error(w, "Not a websocket handshake", 400)
			return
		} else if err != nil {
			return
		}

		// время от времени посылаем данные о карте
		// для рассчета статистики
		statisticTicker := time.NewTicker(500 * time.Millisecond)
		go func() {
			for _ = range statisticTicker.C {
				msg, _ := json.Marshal(g.Map)
				ws.WriteMessage(websocket.TextMessage, msg)
			}
		}()

		// ну и ждем новой инфы от сервиса статистики
		// обновляем по приходу
		for {
			messageType, message, err := ws.ReadMessage()
			if err != nil {
				fmt.Println(err.Error())
				break
			}
			if messageType == websocket.TextMessage {

				data := make([]float64, 6)
				err := json.Unmarshal(message, &data)
				if err != nil {
					fmt.Println("error tut: ")
					fmt.Println(err.Error())
				}
				g.SetStats(data)
				g.SendStats()
			}
		}
	}
}

func wsSavemapHandler(g *Game) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
		if _, ok := err.(websocket.HandshakeError); ok {
			http.Error(w, "Not a websocket handshake", 400)
			return
		} else if err != nil {
			return
		}

		// время от времени посылаем данные о карте
		// для сохранения
		savemapTicker := time.NewTicker(500 * time.Millisecond)
		go func() {
			for _ = range savemapTicker.C {
				msg, _ := json.Marshal(g.Map)
				ws.WriteMessage(websocket.TextMessage, msg)
			}
		}()
	}
}

func main() {
	game := NewGame()

	mapFile, _ := filepath.Abs("./map.data")
	content, err := ioutil.ReadFile(mapFile)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err.Error())
		panic(err)
	} else {
		err = json.Unmarshal(content, &game.Map)
		if err != nil {
			fmt.Println(err.Error())
		}

	}

	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fmt.Sprintf("client/%s", r.URL.Path[8:]))
	})
	http.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fmt.Sprintf("../client/assets/%s", r.URL.Path[8:]))
	})
	http.HandleFunc("/ws", wsHandler(game))
	http.HandleFunc("/statistic-ws", wsStatisticHandler(game))
	http.HandleFunc("/savemap-ws", wsSavemapHandler(game))

	if err := http.ListenAndServe(":9001", nil); err != nil {
		log.Fatal("ListenAndServe:", err.Error())
	}
}
