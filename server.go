package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"monkiato/word-challenge/internal/logic"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{} // use default options
	clients = make(map[*websocket.Conn]*ClientInfo) // connected clients and score associated
	broadcast = make(chan BroadcastMessage)
	logicWords = logic.NewWords()
)

const (
	ClientMessageWord = "word"
	ClientMessageUser = "user"
	BroadcastMessageLog = "log"
	BroadcastMessagePlayerUpdates = "players"
)

type ClientInfo struct {
	username string
	score int64
}

type BroadcastMessage struct {
	Type string `json:"type"`
	Message string `json:"message"`
}

func main() {
	go idleEvaluation()

	http.HandleFunc("/connect", connect)
	http.HandleFunc("/", home)
	addr := "localhost:8080"
	log.Printf("server up: %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func idleEvaluation() {
	for {
		message := <-broadcast

		for client := range clients {
			err := client.WriteJSON(message)
			if err != nil {
				log.Println("write:", err)
				delete(clients, client)
				break
			}
		}
	}
}

func connect(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	clients[c] = &ClientInfo{}

	if len(clients) >= 2 {
		go startCountDown()
	}

	for {
		var message BroadcastMessage
		err := c.ReadJSON(&message)
		if err != nil {
			log.Println("read:", err)
			delete(clients, c)
			break
		}

		log.Printf("recv: %s: %s", message.Type, message.Message)
		processClientMessage(c, message)
	}
}

func processClientMessage(client *websocket.Conn, message BroadcastMessage) {
	switch message.Type {
	case ClientMessageUser:
		clients[client].username = message.Message
		log.Printf("registered username: %s", message.Message)

		// send score from existing users to update their players list
		//TODO: check if I can reduce amount of message for info that most of users already had from before
		for client := range clients {
			sendBroadcastScore(client)
		}

		break
	case ClientMessageWord:
		if logicWords.CurrentWord != "" {
			// challenge has started so we evaluate words coming from clients
			success, score := logicWords.EvaluateSuccess(message.Message)
			if success {
				log.Printf("%s scored %d points", clients[client].username, score)

				clients[client].score += score
				sendBroadcast(BroadcastMessage{
					Type:    BroadcastMessageLog,
					Message: fmt.Sprintf("<br />new word! <b>%s</b><br />", logicWords.CurrentWord),
				})
				sendBroadcastScore(client)
			}
		}
		break
	}
}

func sendBroadcastScore(client *websocket.Conn) {
	sendBroadcast(BroadcastMessage{
		Type:    BroadcastMessagePlayerUpdates,
		Message: fmt.Sprintf("{\"user\": \"%s\", \"score\":%d}", clients[client].username, clients[client].score),
	})
}

func startCountDown() {
	countdown := time.Duration(5)
	endTime := time.Now().Add(countdown * time.Second)
	for v := range time.Tick(time.Second) {
		log.Printf("%d seconds left", int(math.Ceil(endTime.Sub(v).Seconds())))
		sendBroadcast(BroadcastMessage{
			Type: BroadcastMessageLog,
			Message: fmt.Sprintf("%d seconds left", int(endTime.Sub(v).Seconds())),
		})
		if v.After(endTime)  {
			break
		}
	}

	logicWords.Start()
	sendBroadcast(BroadcastMessage{
		Type: BroadcastMessageLog,
		Message: fmt.Sprintf("<br />challenge started! current work: <b>%s</b><br />", logicWords.CurrentWord),
	})
}

func sendBroadcast(message BroadcastMessage) {
	broadcast <- message
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/connect")
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {

	var username = document.getElementById("username");
    var output = document.getElementById("output");
    var input = document.getElementById("input");
	var scores = document.getElementById("scores");
    var ws;
	var scoresData = {}

    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };

	var refreshScores = function() {
		scores.textContent = ""
		Object.keys(scoresData).forEach(function(key) {
			var d = document.createElement("div");
			d.innerHTML = scoresData[key] + ' - ' + key
			scores.appendChild(d);
		});
    }

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }

		if (username.value == '') {
			print("missing username");
			return false;
		}
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("Connected to server");
			// report username
			ws.send(
				JSON.stringify({
					type: 'user',
					message: username.value
				}));
        }
        ws.onclose = function(evt) {
            print("Server connection closed");
            ws = null;
        }
        ws.onmessage = function(evt) {
			msg = JSON.parse(evt.data);
			switch (msg.type) {
				case 'log':
            		print(msg.message);
					break;
				case 'players':
					data = JSON.parse(msg.message)
					scoresData[data.user] = data.score
					refreshScores()
					break;
				default:
					print("Unexpected message type: " + msg.Type)
			}
        }
        ws.onerror = function(evt) {
            print("Error: " + evt.data);
        }
        return false;
    };

	input.addEventListener("keyup", function(event) {
  		if (event.keyCode === 13) {
			// Cancel the default action, if needed
			event.preventDefault();
			if (!ws) {
				return false;
			}
            ws.send(
				JSON.stringify({
					type: 'word',
					message: input.value
				})
			);
			input.value = ""
			return false;
		}
	});

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
<tr><td valign="top" width="50%">
<p>Put your name and then click "Open" to create a connection to the server.
Once the challence starts, you have to write words as fast as you can!
</p>
<p><input id="username" type="text" placeholder="put your name"></p>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" placeholder="write words"></p>

<h2>Scores</h2>
<div id="scores"></div>

</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))