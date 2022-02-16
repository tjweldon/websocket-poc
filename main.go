// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func tick(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()
	done := make(chan bool)
	go func() {
		time.Sleep(10 * time.Second)
		done <- true
	}()
	for {
		select {
		case <-done:
			err = c.WriteMessage(websocket.TextMessage, []byte("Done!"))
			return
		case t := <-ticker.C:
			err = c.WriteMessage(websocket.TextMessage, []byte(t.Format("15:04:05.999999999")))

		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/tick")
}

func other(w http.ResponseWriter, r *http.Request) {
	otherTemplate.Execute(w, "ws://"+r.Host+"/tick")
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	// Static assets (doesn't seem to be necessary)
	//http.Handle("/assets/", http.FileServer(http.Dir("./assets")))
	//http.Handle("/resources/", http.FileServer(http.Dir("./resources")))

	// Sockets
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/tick", tick)

	// Pages
	http.HandleFunc("/other", other)
	http.HandleFunc("/", home)
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
	template.New("").Parse(getFileText("./home.html")),
)

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.innerHTML = message;
        output.scroll(0, output.scrollHeight);
    };

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("T: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

});
</script>
</head>
<body>
<table>
<tr>
	<td width="20%">
	<form>
	<button id="open">Open</button>
	<button id="close">Close</button>
	</form>
	</td>
	<td valign="top" width="80%">
	<p>Click "Open" to create a connection to the server, 
	"Send" to send a message to the server and "Close" to close the connection. 
	You can change the message and send multiple times.
	<p>
	</td>
	</tr>
<tr>
	<td width="20%">
	<div id="output" style="max-height: 70vh;overflow-y: scroll; font-family: 'Courier New', monospace; font-size: 10pt"></div>
	</td>
	<td width="80%">
    <main></main>
	</td>
</tr>
</table>

</body>
</html>
`))
