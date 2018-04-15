package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/alehano/wsgame/utils"
	"github.com/gorilla/websocket"
)

type Action struct {
	TickN int
	Type  string
	Value int
}

type Player struct {
	Id      string          `json:"-"`
	Name    string          `json:"name"`
	Connect *websocket.Conn `json:"-"`
	X       int             `json:"x"`
	Y       int             `json:"y"`
	Rad     int             `json:"radius"`
	Angle   int             `json:"angle"`
	Actions []Action        `json:"-"`
}

// Send message to single player
func (p *Player) SendMessage(message string) error {
	return p.Connect.WriteMessage(websocket.TextMessage, []byte(message))
}

// Add action to player actions buffer (auto-convert value)
func (p *Player) AddAction(a Action) {
	p.Actions = append(p.Actions, a)
}

// Goroutine for receive commands from client
func (p *Player) Receiver(onbroken func(id string)) {
	for {
		_, message, err := p.Connect.ReadMessage()
		if err != nil {
			onbroken(p.Id)
			return
		}
		parts := strings.Split(string(message), ":")

		switch parts[0] {
		case "Action":
			v, err := strconv.Atoi(parts[2])
			if err != nil {
				break
			}
			p.AddAction(Action{Type: parts[1], Value: v})
		default:
			err = p.Connect.WriteMessage(websocket.TextMessage, []byte("Error:Invalid type of command"))
			if err != nil {
				onbroken(p.Id)
			}
		}
	}
}

type Game struct {
	TickN   int64              `json:"tickN"`
	Players map[string]*Player `json:"players"`
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

	p.X = int((rand.Int() / math.MaxInt32) * 800)
	p.Y = int((rand.Int() / math.MaxInt32) * 400)
	p.Angle = int((rand.Int() / math.MaxInt32) * 360)
	p.Rad = 10

	go p.Receiver(g.BrokenPlayer)

	return p
}

// Delete player from game object
func (g *Game) DeletePlayer(id string) {
	delete(g.Players, id)
}

// Get player by id
func (g *Game) GetPlayer(id string) (*Player, error) {
	value, ok := g.Players[id]
	if !ok {
		return nil, errors.New("Player not found")
	}
	return value, nil
}

// Call when player loses connection with server
func (g *Game) BrokenPlayer(id string) {
	player, err := g.GetPlayer(id)
	if err != nil {
		return
	}
	g.SendAll(fmt.Sprintf("Logout:%s", player.Name))
	g.DeletePlayer(id)
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

// Send state of game to all players
func (g *Game) SendState() {
	gameJSON, err := json.Marshal(g)
	if err != nil {
		log.Fatal(err.Error())
	}

	g.SendAll(fmt.Sprintf("State:%s", string(gameJSON)))
}

// Main game loop
func (g *Game) TickTack() {
	g.TickN = 0
	for {
		timeTickBegin := time.Now().UnixNano() / int64(time.Millisecond)

		// calculate
		g.TickN += 1

		// sendState
		g.SendState()

		timeTickEnd := time.Now().UnixNano() / int64(time.Millisecond)
		timeTickDiff := 34 - (timeTickEnd - timeTickBegin)
		if timeTickDiff > 0 {
			time.Sleep(time.Duration(timeTickDiff) * time.Millisecond)
		}
	}
}
