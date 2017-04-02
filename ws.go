package place

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func ReadChanges(c *websocket.Conn) <-chan TileUpdate {
	ch := make(chan TileUpdate)
	go func(c *websocket.Conn, ch chan<- TileUpdate) {
		defer close(ch)
		for {
			var evt Event
			if err := c.ReadJSON(&evt); err != nil {
				log.Println(err)
				return
			}

			switch evt.Type {
			case "place":
				var tile TileUpdate
				if err := json.Unmarshal([]byte(evt.Payload), &tile); err != nil {
					log.Println(err)
					return
				}
				ch <- tile
			case "batch-place":
				var tiles []TileUpdate
				if err := json.Unmarshal([]byte(evt.Payload), &tiles); err != nil {
					log.Println(err)
					return
				}
				for _, t := range tiles {
					ch <- t
				}
			case "activity":
				log.Printf("%s\n", evt.Payload)
			default:
				log.Println("unknown event with type:", evt.Type)
			}
		}
	}(c, ch)

	return ch
}

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type TileUpdate struct {
	Tile
	Author string `json:"author"`
}
