package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	Map := [][]int{}
	Stats := []float64{100, 0, 0, 0, 0, 0}
	part := float64(100. / 200.)
	fmt.Println(part)

	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:9001/statistic-ws", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	for {
		var msg []byte
		if c != nil {
			_, msg, err = c.ReadMessage()
		}
		if err != nil {
			time.Sleep(time.Second)
			c, _, err = websocket.DefaultDialer.Dial("ws://localhost:9001/statistic-ws", nil)
			if err != nil {
				fmt.Printf("%s, but try later\n", err.Error())
			}
			continue
		}
		err = json.Unmarshal(msg, &Map)
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < 6; i++ {
			Stats[i] = 0
		}
		for x := 0; x < 20; x++ {
			for y := 0; y < 10; y++ {
				Stats[Map[x][y]] += float64(part)
			}
		}
		m, _ := json.Marshal(Stats)
		c.WriteMessage(websocket.TextMessage, m)
	}
}
