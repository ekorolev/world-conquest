package main

import (
	"encoding/json"
	"errors"
	"fmt"
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
	Id           string          `json:"id"`
	Name         string          `json:"name"`
	Connect      *websocket.Conn `json:"-"`
	X            int             `json:"x"`
	Y            int             `json:"y"`
	Rad          int             `json:"radius"`
	Angle        int             `json:"angle"`
	Actions      []Action        `json:"-"`
	Team         int             `json:"team"`
	LastClick    int64           `json:"-"`
	CurrentPause int64           `json:"-"`
}

func (p *Player) SendMap(m [][]int) {
	msg, err := json.Marshal(m)
	fmt.Println(msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	p.Connect.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Map:%s", msg)))
}

func (p *Player) AllowClick() bool {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	fmt.Println(now, p.LastClick, p.CurrentPause)
	if p.LastClick == 0 {
		return true
	}
	if p.LastClick+p.CurrentPause > now {
		return false
	} else {
		return true
	}
}
func (p *Player) DoClick() {
	p.LastClick = time.Now().UnixNano() / int64(time.Millisecond)
	p.CurrentPause = 5000
}

func (p *Player) ChangeTeam(team int) {
	p.Team = team
}

func (p *Player) SendPlayers(pl map[string]*Player) {
	msg, err := json.Marshal(pl)
	if err != nil {
		fmt.Println(err)
		return
	}
	p.Connect.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Players:%s", msg)))
}

func (p *Player) SendAboutInfo() {
	msg, err := json.Marshal(p)
	if err != nil {
		return
	}
	p.Connect.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("AboutInfo:%s", msg)))
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
func (p *Player) Receiver(g *Game, onbroken func(id string)) {
	for {
		_, message, err := p.Connect.ReadMessage()
		if err != nil {
			fmt.Println("player disconnected")
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
		case "Coords":
			x, _ := strconv.Atoi(parts[1])
			y, _ := strconv.Atoi(parts[2])
			r, _ := strconv.Atoi(parts[3])
			p.X = x
			p.Y = y
			p.Angle = r
		case "MarkCell":
			x, _ := strconv.Atoi(parts[1])
			y, _ := strconv.Atoi(parts[2])
			if !p.AllowClick() {
				p.SendMessage("Error:Нельзя делать клик, рано")
				break
			}
			p.DoClick()
			g.MarkCell(x, y, p)
		case "ChangeTeam":
			team, _ := strconv.Atoi(parts[1])
			p.ChangeTeam(team)
			g.SendPlayerToAll(p)
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
	Map     [][]int
	Stats   []float64
}

// Initialize new game and return pointer to game object
func NewGame() *Game {
	game := new(Game)
	game.Players = make(map[string]*Player)
	game.Map = make([][]int, 20)
	for i := 0; i < 20; i++ {
		game.Map[i] = make([]int, 10)
		for j := 0; j < 10; j++ {
			game.Map[i][j] = 0
		}
	}
	return game
}

// Add new player to game and return player object
func (g *Game) NewPlayer(n string, c *websocket.Conn) *Player {
	p := new(Player)
	p.Id = utils.RandString(16)
	p.Name = n
	p.Connect = c
	g.Players[p.Id] = p

	p.X = rand.Intn(800)
	p.Y = rand.Intn(400)
	p.Angle = rand.Intn(360)
	p.Rad = 10
	p.Team = rand.Intn(5) + 1

	go p.Receiver(g, g.BrokenPlayer)

	return p
}

func (g *Game) MarkCell(x int, y int, p *Player) {
	g.Map[x][y] = p.Team
	g.SendAll(fmt.Sprintf("UpdateCell:%d:%d:%d", x, y, p.Team))
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
	g.SendAll(fmt.Sprintf("Logout:%s", player.Id))
	g.DeletePlayer(id)
}

// Send message to all players
func (g *Game) SendAll(message string) {
	for key, player := range g.Players {
		err := player.SendMessage(message)
		if err != nil {
			g.DeletePlayer(key)
			g.BrokenPlayer(key)
		}
	}
}
func (g *Game) SetStats(m []float64) {
	g.Stats = m
}
func (g *Game) SendStats() {
	statsmsg, _ := json.Marshal(g.Stats)
	msg := fmt.Sprintf("Stats:%s", statsmsg)
	g.SendAll(msg)
}

func (g *Game) SendPlayerToAll(p *Player) {
	msg, _ := json.Marshal(p)
	query := fmt.Sprintf("UpdatePlayer:%s:%s", p.Id, msg)
	g.SendAll(query)
}
