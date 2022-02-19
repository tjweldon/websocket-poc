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
	"strconv"
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

type Frame []Coords

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

// streamFrames handles the websocket that will stream the animation frames.
// It sets up:
//  - The listen goroutine to handle buffering requests from the socket client.
//  - The frameGenerator goroutine to generate the requested number of frames.
//  - A loop to serialise and return contiguous chunks of frame data to the client.
func streamFrames(w http.ResponseWriter, r *http.Request) {
	// Upgrade the web request to a socket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	// Socket close on function return
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)

	// Channel setup
	frameStream := make(chan []Frame)
	frameRequest := make(chan int)

	// Frame data calculation goroutine
	go frameGenerator(frameStream, frameRequest)

	// Listens for buffering requests
	go listen(conn, frameRequest)

	// Instructs frameGenerator to begin with a request for 60 frames
	frameRequest <- 60

	for {
		select {
		case frames := <-frameStream:
			rawJson, err := json.Marshal(frames)
			if err != nil {
				return
			}
			err = conn.WriteMessage(websocket.TextMessage, rawJson)
			if err != nil {
				return
			}
		}
	}
}

// frameGenerator is intended to be run asynchronously and will await a
// message on the frameRequest channel in the form of an integer number
// of frames. On receiving such a message it will calculate the next
// sequence of frames until it has the number requested. They are then
// sent into the frameStream channel.
func frameGenerator(frameStream chan []Frame, frameRequest chan int) {
	defer close(frameStream)
	frameCount := 0
	for seqLen := range frameRequest {
		switch seqLen {
		case -1:
			return
		default:
			var frames []Frame
			for i := frameCount; i < frameCount+seqLen; i++ {
				timeFloat := float64(i) / 60
				frames = append(frames,
					Frame{
						Coords{
							X: LPFloat{Value: 100*math.Cos(2*math.Pi*timeFloat) + 200, Digits: 2},
							Y: LPFloat{Value: 100*math.Sin(2*math.Pi*timeFloat) + 200, Digits: 2},
						},
					},
				)
			}
			frameStream <- frames
		}
		frameCount += seqLen
	}
}

// listen is a thin adaptor layer goroutine that naively interprets
// the utf8 text read from the socket as an int and then supplies
// that value to the frameRequest Channel
func listen(conn *websocket.Conn, frameRequest chan int) {
	defer func() {
		frameRequest <- -1
		close(frameRequest)
	}()
	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		var frameCount int
		if mt == websocket.TextMessage {
			frameCount, err = strconv.Atoi(string(msg))
			if err != nil {
				return
			}
		} else {
			frameCount = int(msg[0])
		}

		frameRequest <- frameCount
	}
}

// index renders the root page to the response
func index(w http.ResponseWriter, r *http.Request) {
	otherTemplate.Execute(w, "ws://"+r.Host+"/tick")
}

// debug renders the root page with extra dev info
func debug(w http.ResponseWriter, r *http.Request) {
	debugTemplate.Execute(w, "ws://"+r.Host+"/tick")
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	// Pages
	http.HandleFunc("/", index)
	http.HandleFunc("/debug", debug)

	// Sockets
	http.HandleFunc("/tick", streamFrames)

	// Static assets (doesn't seem to be necessary)
	http.Handle("/src/", http.FileServer(http.Dir(".")))
	http.Handle("/resources/", http.FileServer(http.Dir(".")))

	// Start server
	log.Fatal(http.ListenAndServe(*addr, nil))
}

// A hack I'm using bc I couldn't get the template library to just read
// directly from the files.
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