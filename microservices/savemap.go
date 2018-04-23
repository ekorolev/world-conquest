package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	Map := [][]int{}
	part := float64(100. / 200.)
	fmt.Println(part)

	content, err := ioutil.ReadFile("map.data")
	if err != nil {
		fmt.Println("Ничего страшного")
	}

	err = json.Unmarshal(content, &Map)
	if err != nil {
		fmt.Println("Ничего страшного")
	}

	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:9001/savemap-ws", nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for _ = range ticker.C {
			msg, err := json.Marshal(Map)
			if err != nil {
				fmt.Println(err.Error())
			}
			if err = ioutil.WriteFile("map.data", msg, 0777); err != nil {
				fmt.Println(err.Error())
			}
		}
	}()

	for {
		var msg []byte
		if c != nil {
			_, msg, err = c.ReadMessage()
		}
		if err != nil {
			time.Sleep(time.Second)
			c, _, err = websocket.DefaultDialer.Dial("ws://localhost:9001/savemap-ws", nil)
			if err != nil {
				fmt.Printf("%s, but try later\n", err.Error())
			}
			continue
		}
		err = json.Unmarshal(msg, &Map)
		if err != nil {
			log.Fatal(err)
		}

	}
}
