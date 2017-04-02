package main

import (
	"fmt"
	"log"

	"vallon.me/place"

	"github.com/gorilla/websocket"
)

const wsURL = "WS_URL_HERE"

func main() {
	c, _, err = websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	for tile := range place.ReadChanges(c) {
		fmt.Println(tile)
	}
}
