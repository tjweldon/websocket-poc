// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

// LPFloat serialises as a number represented to a fixed number
// decimal places eg. 1.00
type LPFloat struct {
	Value  float64 // the actual value
	Digits int     // the number of digits used in json
}

// MarshalJSON serialises the LPFloat Type
func (l LPFloat) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("%.*f", l.Digits, l.Value)
	return []byte(s), nil
}

type Coords struct {
	X LPFloat `json:"x"`
	Y LPFloat `json:"y"`
}

func tick(w http.ResponseWriter, r *http.Request) {
	// Upgrade the web request to a socket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	// Socket close on function return
	defer conn.Close()

	// Initiate Clock
	ticker := time.NewTicker(time.Second / 70)
	defer ticker.Stop()

	// Signal to kill the socket
	done := make(chan bool)
	go func() {
		time.Sleep(10 * time.Second)
		done <- true
	}()

	// Stream loop
	for {
		select {
		case <-done:
			// Send a farewell
			err = conn.WriteMessage(websocket.TextMessage, []byte("Done!"))
			return
		case t := <-ticker.C:
			// Calculate coordinates
			timeFloat := float64(t.UnixMilli()%1000) / 1e3
			coords := Coords{
				X: LPFloat{Value: 100*math.Cos(2*math.Pi*timeFloat) + 200, Digits: 2},
				Y: LPFloat{Value: 100*math.Sin(2*math.Pi*timeFloat) + 200, Digits: 2},
			}

			// JSON serialise
			rawJson, err := json.Marshal(coords)
			if err != nil {
				panic(err)
			}

			// SEND IT LADS, YEET
			err = conn.WriteMessage(websocket.TextMessage, rawJson)
		}
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	otherTemplate.Execute(w, "ws://"+r.Host+"/tick")
}

func debug(w http.ResponseWriter, r *http.Request) {
	debugTemplate.Execute(w, "ws://"+r.Host+"/tick")
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	//// Static assets (doesn't seem to be necessary)
	http.Handle("/src/", http.FileServer(http.Dir(".")))
	http.Handle("/resources/", http.FileServer(http.Dir(".")))

	// Sockets
	http.HandleFunc("/tick", tick)

	// Pages
	http.HandleFunc("/", index)
	http.HandleFunc("/debug", debug)

	// Start server
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func getFileText(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(content)
}

var otherTemplate = template.Must(
	template.New("index").Parse(getFileText("./home.html")),
)

var debugTemplate = template.Must(
	template.New("index").Parse(getFileText("./debug.html")),
)