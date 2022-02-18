package msgserver

import (
	"log"

	"github.com/pranotobudi/myslack-happy-backend/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Hub struct {
	participants map[string]map[string]*wsClient
	// string type is for roomId,
	// *client because we want to hold the reference (memory address) only

	broadcastMsg chan mongodb.Message
	register     chan *wsClient
	unregister   chan *wsClient

	// broadcastMsg     chan []byte
	// broadcastMsg chan ClientMsg
}

type ClientMsg struct {
	// deprecated, only for initial app
	ClientId string `json:"client_id"`
	Text     string `json:"text"`
	RoomId   string `json:"room_id"`
}

// NewHub creates newHub object
func NewHub() *Hub {
	log.Println("newHub")
	return &Hub{
		participants: make(map[string]map[string]*wsClient),
		register:     make(chan *wsClient),
		unregister:   make(chan *wsClient),
		broadcastMsg: make(chan mongodb.Message),
		// broadcastMsg: make(chan ClientMsg),
	}
}

// addClient will add client c to the room roomName
func (h *Hub) addClient(roomName string, c *wsClient) {
	h.participants[roomName][c.clientId] = c
	// h.participants[roomName] = append(h.participants[roomName], c)
}

// Run will run hub in infinite looping
func (h *Hub) Run() {
	log.Println("inside Run")
	for {
		select {
		case msg := <-h.broadcastMsg:
			log.Println("inside Run: new message, send to room participants, clientMsg:", msg)
			log.Println("<- h.broadcastMsg total member: ", len(h.participants[msg.RoomID]))
			// broadcast to its room participants
			for _, client := range h.participants[msg.RoomID] {
				log.Println("inside Run - h.broadcasting, room: ", msg.RoomID, " clientID: ", client.clientId)
				select {
				case client.send <- msg:
				default:
					log.Println("inside Run- h.broadcastMsg default")
					close(client.send)
					delete(h.participants[msg.RoomID], client.clientId)
				}
			}
		case client := <-h.register:
			log.Println("inside Run: register new client:", client.clientId)
			err := h.registerClient(client)
			if err != nil {
				log.Println("client registration to hub failed..: ", err)
			}
			log.Println("inside Run: client registration Success..")

		case client := <-h.unregister:
			log.Println("inside Run: unregister client:", client)
			err := h.unregisterClient(client)
			if err != nil {
				log.Println("client unregistration from hub failed..: ", err)
			}
			log.Println("inside Run: unregister client Success..")
		}
	}
}

// registerClient will register the client to the hub
func (h *Hub) registerClient(c *wsClient) error {
	userId := c.clientId
	objID, err := primitive.ObjectIDFromHex(userId)
	log.Println("func registerClient - objID: ", objID, " userID: ", userId)
	if err != nil {
		log.Println("convert string to objectID failed")
		return err
		// panic(err)
	}

	filter := bson.M{"_id": objID}
	user, err := c.mongodbConn.GetUser(filter)
	if err != nil {
		log.Println("loadRooms - failed to get user: ", err)
		return err
	}

	for _, room := range user.Rooms {
		if h.participants[room] == nil {
			// because each room is a map which has not been initialized, don't forget make(map[*client]bool)
			h.participants[room] = make(map[string]*wsClient)
		}
		log.Println("--- before total member in: ", room, ": ", len(h.participants[room]), "roomID: ")
		// add client to map of map
		h.participants[room][c.clientId] = c
		log.Println("--- after total member in: ", room, ": ", len(h.participants[room]))
	}
	return nil
}

// unregisterClient will unregister the client from the hub
func (h *Hub) unregisterClient(c *wsClient) error {
	userId := c.clientId
	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Println("convert string to objectID failed")
		return err
		// panic(err)
	}

	filter := bson.M{"_id": objID}
	user, err := c.mongodbConn.GetUser(filter)
	if err != nil {
		log.Println("loadRooms - failed to get user: ", err)
		return err
	}

	for _, room := range user.Rooms {
		// delete client from map
		// close(client.send) // remove send channel from memory, actually no need for this line, it will be garbage collected automatically later, but for channel it is better do it manually. it should be done first before remove the client
		// when connection cut, it is closed automatically, so no need to close it, otherside panic
		delete(h.participants[room], c.clientId)

		if len(h.participants) == 0 {
			//delete the room from map
			delete(h.participants, room)
		}
	}
	return nil
}

// func (h *Hub) loadRoomsById(clientId string) []string {
// 	// load rooms from database, if any room not load to hub, add to it
// 	// because each room is a map which has not initialized, don't forget make(map[*client]bool)
// 	return []string{"room1", "room2"}
// }
