package controller

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/nayeemnishaat/go-web-app/web/lib"
)

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

var clients = make(map[*websocket.Conn]string)
var clientsMutex sync.Mutex

// var wsChan = make(chan WsPayload)

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

	conn := ws
	clients[conn] = ""

	go app.ListenForWS(conn)
}

// var ch = make(chan bool) // For unbuffered chan, both the sender and the receiver must be ready for the operation to proceed. The sender will be blocked until the receiver is ready to receive the value, and vice versa.
// var ch = make(chan bool, 1) // for buffered chan, the send operation will only block when the buffer is full, and the receive operation will only block when the buffer is empty as it will wait there forever for receving a value

func (app *Application) ListenForWS(conn *websocket.Conn) {
	defer func() {
		conn.Close()
		clientsMutex.Lock()
		delete(clients, conn)
		clientsMutex.Unlock()
	}()

	// Register a close handler
	// conn.SetCloseHandler(func(code int, text string) error {
	// 	fmt.Printf("Connection closed with code %d: %s\n", code, text)

	// 	// ch <- true // unbuffered chan, blocked if receiver is not ready
	// 	ch <- true // buffered chan, blocked if buffer is full
	// 	return nil
	// })

	// Recover gracefully (it will fire if error in for loop occurs 1000 times)
	defer func() {
		if r := recover(); r != nil {
			app.ErrorLog.Println("Error:", fmt.Sprintf("%v", r))
		}
	}()

	var payload lib.WsPayload

	for {
		err := conn.ReadJSON(&payload)

		// fmt.Println(err, payload, clients)
		// fmt.Printf("Number of running goroutines: %d\n", runtime.NumGoroutine())

		if err != nil {
			break
			// return // Alt: Also can use return and it will stll trigger defer clauses
		} else {
			payload.Conn = conn
			app.WsChan <- payload
		}

		// Alt: with channel
		// val := <-ch // This will always block until a value is received for buffered/unbuffered chan

		// select {
		// case val := <-ch:
		// 	if val {
		// 		return
		// 	}
		// default: // Important: Without this default: condition it will act as blocking (<-ch) for unbuffered/buffered chan
		// 	fmt.Println("Buffered channel is empty")
		// }

		// payload.Conn = *conn
		// wsChan <- payload
	}
}

func (app *Application) ListenToWsChannel() {
	var response WsJsonResponse

	for {
		e, ok := <-app.WsChan

		if !ok {
			fmt.Println("wsChan closed!")
			return
		}

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
			clientsMutex.Lock()
			delete(clients, client)
			clientsMutex.Unlock()
		}
	}
}
