package main

import (
	"github.com/alehano/wsgame/utils"
	"github.com/gorilla/websocket"
)

type Player struct {
	Id      string
	Name    string
	Connect *websocket.Conn
}

// Send message to single player
func (p *Player) SendMessage(message string) error {
	return p.Connect.WriteMessage(websocket.TextMessage, []byte(message))
}

type Game struct {
	Players map[string]*Player
}

// Initialize new game and return pointer to game object
func NewGame() *Game {
	game := new(Game)
	game.Players = make(map[string]*Player)
	return game
}

// Add new player to game and return player object
func (g *Game) NewPlayer(n string, c *websocket.Conn) *Player {
	p := new(Player)
	p.Id = utils.RandString(16)
	p.Name = n
	p.Connect = c
	g.Players[p.Id] = p
	return p
}

// Delete player from game object
func (g *Game) DeletePlayer(id string) {
	delete(g.Players, id)
}

// Send message to all players
func (g *Game) SendAll(message string) {
	for key, player := range g.Players {
		err := player.SendMessage(message)
		if err != nil {
			g.DeletePlayer(key)
		}
	}
}
