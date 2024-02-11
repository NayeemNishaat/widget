package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocketConnection struct {
	*websocket.Conn
}

type WsPayload struct {
	Action      string              `json:"action"`
	Message     string              `json:"message"`
	Username    string              `json:"username"`
	MessageType string              `json:"message_type"`
	UsedID      int                 `json:"user_id"`
	Conn        WebSocketConnection `json:"-"`
}

type WsJsonResponse struct {
	Action  string `json:"action"`
	Message string `json:"message"`
	UsedID  int    `json:"user_id"`
}

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // This is used to secure the WS connection
}

var clients = make(map[WebSocketConnection]string)
var wsChan = make(chan WsPayload)

func (app *Application) WsEndpoint(w http.ResponseWriter, r *http.Request) {
	wUpgrade := w
	if u, ok := w.(interface{ Unwrap() http.ResponseWriter }); ok {
		wUpgrade = u.Unwrap()
	}

	ws, err := upgradeConnection.Upgrade(wUpgrade, r, nil)
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	app.InfoLog.Printf(fmt.Sprintf("Client connected from %s", r.RemoteAddr))

	var response WsJsonResponse
	response.Message = "Connected to server"

	err = ws.WriteJSON(response)
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	conn := WebSocketConnection{Conn: ws}
	clients[conn] = ""

	go app.ListenForWS(&conn)
}

// var ch = make(chan bool) // For unbuffered chan, both the sender and the receiver must be ready for the operation to proceed. The sender will be blocked until the receiver is ready to receive the value, and vice versa.
// var ch = make(chan bool, 1) // for buffered chan, the send operation will only block when the buffer is full, and the receive operation will only block when the buffer is empty

func (app *Application) ListenForWS(conn *WebSocketConnection) {
	defer delete(clients, *conn)
	defer conn.Close()

	// Register a close handler
	// conn.SetCloseHandler(func(code int, text string) error {
	// 	fmt.Printf("Connection closed with code %d: %s\n", code, text)

	// 	// ch <- true // unbuffered chan, blocked because receiver is not ready
	// 	// ch <- true
	// 	return nil
	// })

	// Recover gracefully (it will fire if error in for loop occurs 1000 times)
	defer func() {
		if r := recover(); r != nil {
			app.ErrorLog.Println("Error:", fmt.Sprintf("%v", r))
		}
	}()

	var payload WsPayload

	for {
		err := conn.ReadJSON(&payload)

		// fmt.Println(err, payload, clients)
		// fmt.Printf("Number of running goroutines: %d\n", runtime.NumGoroutine())

		if err != nil {
			// return
			break
		} else {
			payload.Conn = *conn
			wsChan <- payload
		}

		// Alt: with channel
		// v := <-ch // unbuffered, will block until a value is sent to the channel

		// select {
		// case val := <-ch:
		// 	if val {
		// 		return
		// 	}
		// default:
		// 	fmt.Println("Buffered channel is empty")
		// }

		// payload.Conn = *conn
		// wsChan <- payload
	}
}

func (app *Application) ListenToWsChannel() {
	var response WsJsonResponse

	for {
		e := <-wsChan

		switch e.Action {
		case "deleteUser":
			response.Action = "logout"
			response.Message = "Your account has been deleted!"
			response.UsedID = e.UsedID
			app.broadcastToAll(response)
		default:
		}
	}
}

func (app *Application) broadcastToAll(response WsJsonResponse) {
	for client := range clients {
		// broadcast to every connected client
		err := client.WriteJSON(response)
		if err != nil {
			app.ErrorLog.Printf("Wesocket error on %s:%s", response.Action, err)
			_ = client.Close()
			delete(clients, client)
		}
	}
}
