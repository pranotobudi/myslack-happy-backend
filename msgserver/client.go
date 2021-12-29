package msgserver

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pranotobudi/myslack-monolith-backend/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

type client struct {
	conn              *websocket.Conn
	clientId          string
	clientToServerMsg chan string
	serverToClientMsg chan string
	// client hold hub's memory address, because we need to notify incoming message to hub's room
	// hub (see hub structure) also hold client's memory address in memory
	// so, it is just circular references (the memory address), not the physically embedded circular recursive
	hub *Hub

	// Time allowed to write a message to the peer.
	writeWait time.Duration

	// Time allowed to read the next pong message from the peer.
	pongWait time.Duration

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod time.Duration

	// Maximum message size allowed from peer.
	maxMessageSize int64

	// Buffered channel of outbound messages.
	// send chan []byte
	// send chan ClientMsg
	send        chan mongodb.ClientMessage
	mongodbConn *mongodb.MongoDB
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func NewClient(conn *websocket.Conn, hub *Hub, mongodbConn *mongodb.MongoDB) *client {
	return &client{
		conn:              conn,
		clientId:          "",
		clientToServerMsg: make(chan string),
		serverToClientMsg: make(chan string),
		hub:               hub,
		// Time allowed to write a message to the peer.
		writeWait: 10 * time.Second,

		// Time allowed to read the next pong message from the peer.
		pongWait: 60 * time.Second,

		// Send pings to peer with this period. Must be less than pongWait.
		pingPeriod: ((60 * time.Second) * 9) / 10,

		// Maximum message size allowed from peer.
		maxMessageSize: 512,

		// Buffered channel of outbound messages.
		// send: make(chan []byte, 256),
		send:        make(chan mongodb.ClientMessage),
		mongodbConn: mongodbConn,
	}
}

func InitWebsocket(hub *Hub, mongodbConn *mongodb.MongoDB) func(c *gin.Context) {

	// init websocket
	log.Println("initWebsocket")
	return func(c *gin.Context) {
		var upgrader websocket.Upgrader
		upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("websocket connection failed: ", err)
			return
		}
		log.Println("inside InitWebsocket! connection success")

		// notify hub for client initiation event
		client := NewClient(conn, hub, mongodbConn)
		// hub.addClient("room1", client)
		// log.Println("register client to hub (will load client snapshot to hub)...", client)
		// client.hub.register <- client

		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go client.writePump()
		go client.readPump()
	}
}

func (c *client) clientSendToServer() {
	log.Println("clientSendToServer")
}

func (c *client) serverSendToClient() {
	log.Println("serverSendToClient")
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *client) readPump() {
	log.Println("ReadPump run...")
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(c.maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(c.pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(c.pongWait)); return nil })
	for {
		// var clientMsg ClientMsg
		var clientMsg mongodb.ClientMessage
		err := c.conn.ReadJSON(&clientMsg)
		log.Println(fmt.Sprintf("clientMsg: %+v", clientMsg))
		// log.Println("client: ", clientMsg, "id: ", clientMsg.ClientId, "text: ", clientMsg.Text, "roomID: ", clientMsg.RoomId)
		log.Println("text: ", clientMsg.Message, "roomID: ", clientMsg.RoomID, " userID: ", clientMsg.UserID, "timestamp: ", clientMsg.Timestamp)
		if err != nil {
			log.Println("inside readPump - ERROR READJSON")
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// special message - initial message for user info
		if clientMsg.Message == "[USERINFO]" {
			c.clientId = clientMsg.UserID
			log.Println("inside readPump - special Message: ", c.clientId)
			log.Println("register client to hub (will load client snapshot to hub)...", c)
			c.hub.register <- c
		} else {
			// normal message

			// save to mongoDB
			message := bson.D{{"message", clientMsg.Message}, {"user_id", clientMsg.UserID}, {"room_id", clientMsg.RoomID}, {"username", clientMsg.Username}, {"user_image", clientMsg.UserImage}, {"timestamp", clientMsg.Timestamp}}
			docId, err := c.mongodbConn.AddMessage(message)
			if err != nil {
				log.Println("inside readPump - normal Message, add message to MongoDB FAILED")
			}
			log.Println("inside readPump - normal Message, add message to MongoDB success, id: ", docId)
			// TODO

			// broadcast to other clients
			c.hub.broadcastMsg <- clientMsg
			log.Println("inside readPump - normal Message: ", clientMsg.Message)
		}
		// _, message, err := c.conn.ReadMessage()
		// message := clientMsg.Text
		// message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		// c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *client) writePump() {
	log.Println("writePump run...")
	ticker := time.NewTicker(c.pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case clientMsg, ok := <-c.send:
			log.Println("inside writePump.. <-c.send event, ok:", ok)
			c.conn.SetWriteDeadline(time.Now().Add(c.writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.conn.WriteJSON(clientMsg)
			if err != nil {
				log.Println("inside writePump.. err: ", err.Error())
				c.conn.Close()
				return
			}
			// w, err := c.conn.NextWriter(websocket.TextMessage)
			// if err != nil {
			// 	return
			// }
			// w.Write(message)

			// // Add queued chat messages to the current websocket message.
			// n := len(c.send)
			// for i := 0; i < n; i++ {
			// 	w.Write(newline)
			// 	w.Write(<-c.send)
			// }

			// if err := w.Close(); err != nil {
			// 	return
			// }
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(c.writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}