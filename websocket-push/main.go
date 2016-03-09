package main

import (
	"io"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

// messae represents a single message
type message struct {
	From string
	Msg  []byte
	Time time.Time
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true }, // FIXME : Remove
	}
)

// One global room, for now
var room1 *room

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		log.WithField("PORT", port).Fatal("$PORT must be set")
	}

	room1 = startRoom()

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handleWebsocket)
	mux.HandleFunc("/api/hello", handleAPIHello)

	err := http.ListenAndServe(":"+port, mux)

	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

}

func handleAPIHello(w http.ResponseWriter, r *http.Request) {
	room1.toAll <- message{Msg: []byte("Hello From the API =)")}
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithField("err", err).Println("Upgrading to websockets")
		http.Error(w, "Error Upgrading to websockets", 400)
		return
	}

	// OK websokset connection. Add to room
	room1.join <- ws

	for {
		mt, data, err := ws.ReadMessage()
		ctx := log.Fields{"mt": mt, "data": string(data), "err": err}
		if err != nil {
			if err == io.EOF {
				log.WithFields(ctx).Info("Websocket closed!")
			} else {
				log.WithFields(ctx).Error("Error reading websocket message")
			}
			break
		}
		switch mt {
		case websocket.TextMessage:
			//msg, err := validateMessage(data)
			//if err != nil {
			ctx["msg"] = data
			//ctx["err"] = err
			log.WithFields(ctx).Println("Read from ws")
			//	break
			//}
			//rw.publish(data)

			msg := message{
				Msg: data,
			}
			room1.toAll <- msg

			//ws.WriteMessage(mt, []byte("Takk!"))
		default:
			log.WithFields(ctx).Warning("Unknown Message!")
		}
	}

}
