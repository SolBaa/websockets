package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
)

var wsChan = make(chan WsPayload)
var clients = make(map[WebSocketConnection]string)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./templates"),
	jet.InDevelopmentMode(),
)

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Home(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "home.jet", nil)
	if err != nil {
		log.Println(err)
	}

}

type WebSocketConnection struct {
	*websocket.Conn
}
type WsJSONResponse struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}

type WsPayload struct {
	Action   string              `json:"action"`
	Username string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"-"`
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Client connected to endpoint")

	var response WsJSONResponse
	response.Message = "Connected to server"

	conn := WebSocketConnection{Conn: ws}
	clients[conn] = ""

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}

	go ListenForWS(&conn)

}

func ListenforWSChannel() {
	var response WsJSONResponse

	for {
		e := <-wsChan
		response.Action = "Got Here"
		response.Message = fmt.Sprintf("Some Message, and action wwasd %s", e.Action)
		braodcasToAll(response)
	}

}

func braodcasToAll(response WsJSONResponse) {

	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println(err)
			_ = client.Close()
			delete(clients, client)
		}
	}

}

func ListenForWS(conn *WebSocketConnection) {

	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	var payload WsPayload

	for {
		err := conn.ReadJSON(payload)
		if err != nil {
			log.Println(err)
		} else {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	views, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Println(err)
		return err
	}

	err = views.Execute(w, data, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
