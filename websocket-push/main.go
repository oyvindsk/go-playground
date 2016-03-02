package main

import (
	"io"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true }, // FIXME : Remove
	}
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handleWebsocket)

	err := http.ListenAndServe(":80", mux)

	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

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

	for {
		mt, data, err := ws.ReadMessage()
		ctx := log.Fields{"mt": mt, "data": data, "err": err}
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

			ws.WriteMessage(mt, []byte("Takk!"))
		default:
			log.WithFields(ctx).Warning("Unknown Message!")
		}
	}

}
