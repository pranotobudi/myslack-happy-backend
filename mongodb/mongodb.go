package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/pranotobudi/myslack-happy-backend/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	Username  string   `json:"username"`
	UserImage string   `json:"user_image"`
	Rooms     []string `json:"rooms"`
}

type UserAuth struct {
	Email     string `json:"email"`
	UserImage string `json:"user_image"`
}

type Message struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	RoomID    string    `json:"room_id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	UserImage string    `json:"user_image"`
	Timestamp time.Time `json:"timestamp"`
}
type ClientMessage struct {
	Message   string    `json:"message"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	UserImage string    `json:"user_image"`
	RoomID    string    `json:"room_id"`
	Timestamp time.Time `json:"timestamp"`
}
type RoomMongo struct {
	_ID  primitive.ObjectID
	Name string
}

// Room is neutral without ObjectID
type Room struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type IMongoDB interface {
	createCollection(name string)
	getCollection(name string) *mongo.Collection
	InsertDoc(name string, doc bson.D)
	DataSeeder()
	GetRooms() ([]Room, error)
	GetRoom(filter interface{}) (*Room, error)
	GetAnyRoom() (*Room, error)
	AddRoom(roomName string) (string, error)
	AddRooms(rooms []interface{}) ([]string, error)
	GetMessages(filter interface{}) ([]Message, error)
	GetMessage(filter interface{}) (Message, error)
	AddMessage(message interface{}) (string, error)
	AddMessages(messages []interface{}) ([]string, error)
	GetUsers(filter interface{}) ([]User, error)
	GetUser(filter interface{}) (*User, error)
	AddUser(user interface{}) (string, error)
	UpdateUser(filter interface{}, update interface{}, options *options.UpdateOptions) error
	AddUsers(users []interface{}) ([]string, error)
}
type MongoDB struct {
	client *mongo.Client
	config config.MongoDb
}

var MongoDBInstance *MongoDB
var once sync.Once

func NewMongoDB() *MongoDB {
	once.Do(func() {
		dbConfig := config.MongoDbConfig()
		clientOptions := options.Client().ApplyURI("mongodb+srv://pranotobudi:myslack-db-password@myslack-db.bovrx.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Database Connection success ...")
		mongodb := &MongoDB{
			client: client,
			config: dbConfig,
		}
		MongoDBInstance = mongodb

		// only if needed
		// MongoDBInstance.DataSeeder()
	})

	return MongoDBInstance
}
func (m *MongoDB) createCollection(name string) {
	coll := m.client.Database("myslack-db").Collection(name)
	log.Println("create collection, name:", coll.Name())
}
func (m *MongoDB) getCollection(name string) *mongo.Collection {
	return m.client.Database(m.config.Name).Collection(name)
}

// func (m *MongoDB) insertDoc(coll *mongo.Collection, name string, doc bson.D) {
// 	// doc := bson.D{{"title", "Invisible Cities"}, {"author", "Italo Calvino"}, {"year_published", 1974}}
// 	result, err := coll.InsertOne(context.TODO(), doc)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Println("Inserted document with _id: \n", result.InsertedID)
// }

func (m *MongoDB) InsertDoc(name string, doc bson.D) {
	// doc := bson.D{{"title", "Invisible Cities"}, {"author", "Italo Calvino"}, {"year_published", 1974}}
	coll := m.getCollection(name)
	result, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted document with _id: \n", result.InsertedID)
}

func (m *MongoDB) DataSeeder() {
	// create collection
	m.createCollection("rooms")
	m.createCollection("users")
	m.createCollection("messages")
	rooms := []interface{}{
		bson.D{{"name", "room1"}},
		bson.D{{"name", "room2"}},
		bson.D{{"name", "room3"}},
	}
	roomIds, err := m.AddRooms(rooms)
	if err != nil {
		fmt.Println("error AddRooms: ", err)
	}
	for _, id := range roomIds {
		fmt.Println("doc room id: ", id)
	}
	// roomsColl := m.getCollection("rooms")
	// results, _ := roomsColl.InsertMany(context.TODO(), rooms)
	// for _, id := range results.InsertedIDs {
	// 	fmt.Println("doc room id: ", id)
	// }

	users := []interface{}{
		bson.D{{"email", "ocean.king.digital@gmail.com"}, {"username", "ocean.king.digital"}, {"user_image", "localhost"}, {"rooms", bson.A{roomIds[0], roomIds[1]}}},
		bson.D{{"email", "lumion.design.studio@gmail.com"}, {"username", "lumion.design.studio"}, {"user_image", "localhost"}, {"rooms", bson.A{roomIds[0], roomIds[1]}}},
	}
	userIds, err := m.AddUsers(users)
	if err != nil {
		fmt.Println("error AddUsers: ", err)
	}
	for _, id := range userIds {
		fmt.Println("doc user id: ", id)
	}

	// usersColl := m.getCollection("users")
	// results, _ = usersColl.InsertMany(context.TODO(), users)

	// for _, result := range results.InsertedIDs {
	// 	fmt.Println("doc users id: ", result)
	// }
	var messages []interface{}
	for _, userId := range userIds {
		for _, roomId := range roomIds {
			message := bson.D{{"message", "a" + userId}, {"user_id", userId}, {"room_id", roomId}, {"username", userId + roomId}, {"user_image", "http://localhost"}, {"timestamp", time.Now()}}
			messages = append(messages, message)
		}
	}

	messageIds, err := m.AddMessages(messages)
	if err != nil {
		fmt.Println("error AddMessages: ", err)
	}
	for _, id := range messageIds {
		fmt.Println("doc message id: ", id)
	}

	// messagesColl := m.getCollection("messages")
	// results, _ = messagesColl.InsertMany(context.TODO(), messages)

	// for _, result := range results.InsertedIDs {
	// 	fmt.Println("doc messages id: ", result)
	// }

}

func (m *MongoDB) GetRooms() ([]Room, error) {
	coll := m.getCollection("rooms")
	log.Println("getRooms coll: ", coll)
	filter := bson.D{}

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, nil
		// panic(err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, nil
		// panic(err)
	}
	var finalResult []Room
	for _, result := range results {
		fmt.Println(result)
		var room Room
		room.ID = fmt.Sprintf("%v", result["_id"].(primitive.ObjectID).Hex())
		room.Name = fmt.Sprintf("%v", result["name"])
		finalResult = append(finalResult, room)
	}
	return finalResult, nil
}

func (m *MongoDB) GetRoom(filter interface{}) (*Room, error) {
	coll := m.getCollection("rooms")
	// filter := bson.D{}
	log.Println("getRoom coll: ", coll)

	var roomMongo bson.M
	coll.FindOne(context.TODO(), filter).Decode(&roomMongo)
	log.Println("inside GetRoom, roomMongo: ", roomMongo)

	var room Room
	if roomMongo == nil {
		return &room, errors.New("room not found")
	}
	room.ID = roomMongo["_id"].(primitive.ObjectID).Hex()
	room.Name = roomMongo["name"].(string)
	log.Println("inside GetRoom, room: ", room)
	return &room, nil
}

func (m *MongoDB) GetAnyRoom() (*Room, error) {
	coll := m.getCollection("rooms")
	log.Println("GetAnyRoom coll: ", coll)

	// filter := bson.D{}
	var roomMongo bson.M
	log.Println("inside GetAnyRoom, roomMongo before: ", roomMongo)
	coll.FindOne(context.TODO(), bson.M{}).Decode(&roomMongo)
	log.Println("inside GetAnyRoom, roomMongo after: ", roomMongo)
	var room Room
	if roomMongo == nil {
		return &room, errors.New("room not found")
	}
	room.ID = roomMongo["_id"].(primitive.ObjectID).Hex()
	room.Name = roomMongo["name"].(string)
	log.Println("inside GetAnyRoom, room: ", room)

	return &room, nil
}

func (m *MongoDB) AddRoom(roomName string) (string, error) {

	coll := m.getCollection("rooms")
	doc := bson.D{{"name", roomName}}
	result, err := coll.InsertOne(context.TODO(), doc)

	if err != nil {
		log.Println("failed to insert room: ", err)
		return "", err
	}

	return fmt.Sprintf("%v", result.InsertedID), nil
}

func (m *MongoDB) AddRooms(rooms []interface{}) ([]string, error) {

	coll := m.getCollection("rooms")
	// doc := bson.D{{"name", roomName}}
	results, err := coll.InsertMany(context.TODO(), rooms)
	if err != nil {
		log.Println("failed to insert rooms: ", err)
		return nil, err
	}

	var retValues []string
	for _, result := range results.InsertedIDs {
		fmt.Println("doc room id: ", result)
		retValues = append(retValues, result.(primitive.ObjectID).Hex())
	}

	return retValues, nil
}

func (m *MongoDB) GetMessages(filter interface{}) ([]Message, error) {
	log.Println("INSIDE REPO GetMessages")
	coll := m.getCollection("messages")
	// oid, err := primitive.ObjectIDFromHex(roomId)
	// if err != nil {
	// 	return nil, nil
	// }
	// log.Println("mongoDB-GetMesssages, roomId: ", roomId)
	// filter := bson.M{"room_id": roomId}
	// filter := bson.M{}

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, nil
		// panic(err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, nil
		// panic(err)
	}
	var finalResult []Message
	for _, result := range results {
		log.Println("mongoDB-GetMessages message: ", result)
		var message Message
		message.ID = result["_id"].(primitive.ObjectID).Hex()
		message.Message = result["message"].(string)
		message.RoomID = result["room_id"].(string)
		message.Timestamp = result["timestamp"].(primitive.DateTime).Time()
		message.Username = result["username"].(string)
		message.UserID = result["user_id"].(string)
		message.UserImage = result["user_image"].(string)
		finalResult = append(finalResult, message)
	}
	return finalResult, nil
}

func (m *MongoDB) GetMessage(filter interface{}) (Message, error) {
	coll := m.getCollection("messages")
	// filter := bson.D{}

	var messageMongo bson.M
	coll.FindOne(context.TODO(), filter).Decode(&messageMongo)
	log.Println("inside GetMessage, messageMongo: ", messageMongo)

	var message Message
	message.ID = messageMongo["_id"].(primitive.ObjectID).Hex()
	message.Message = messageMongo["message"].(string)
	message.RoomID = messageMongo["room_id"].(string)
	message.Timestamp = messageMongo["timestamp"].(primitive.DateTime).Time()
	message.Username = messageMongo["username"].(string)
	message.UserID = messageMongo["user_id"].(string)
	message.UserImage = messageMongo["user_image"].(string)

	log.Println("inside GetMessage, message: ", message)
	return message, nil
}

func (m *MongoDB) AddMessage(message interface{}) (string, error) {

	coll := m.getCollection("messages")
	// doc := bson.D{{"name", roomName}}
	result, err := coll.InsertOne(context.TODO(), message)

	if err != nil {
		log.Println("failed to insert message: ", err)
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
	// return fmt.Sprintf("%v", result.InsertedID), nil
}

func (m *MongoDB) AddMessages(messages []interface{}) ([]string, error) {

	coll := m.getCollection("messages")
	// doc := bson.D{{"name", roomName}}
	results, err := coll.InsertMany(context.TODO(), messages)
	if err != nil {
		log.Println("failed to insert messages: ", err)
		return nil, err
	}

	var retValues []string
	for _, result := range results.InsertedIDs {
		fmt.Println("doc message id: ", result)
		retValues = append(retValues, result.(primitive.ObjectID).Hex())
	}

	return retValues, nil
}

func (m *MongoDB) GetUsers(filter interface{}) ([]User, error) {
	coll := m.getCollection("users")
	// oid, err := primitive.ObjectIDFromHex(roomId)
	// if err != nil {
	// 	return nil, nil
	// }
	// log.Println("mongoDB-GetMesssages, roomId: ", roomId)
	// filter := bson.M{"room_id": roomId}
	// filter := bson.M{}

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, nil
		// panic(err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, nil
		// panic(err)
	}
	var finalResult []User
	for _, result := range results {
		log.Println("mongoDB- user: ", result)
		var user User
		user.ID = result["_id"].(primitive.ObjectID).Hex()
		user.Email = result["email"].(string)
		roomsPrimitiveArray := result["rooms"].(primitive.A)
		roomsInterface := []interface{}(roomsPrimitiveArray)
		for _, v := range roomsInterface {
			user.Rooms = append(user.Rooms, v.(string))
		}
		user.Username = result["username"].(string)
		user.UserImage = result["user_image"].(string)
		finalResult = append(finalResult, user)
	}
	return finalResult, nil
}

func (m *MongoDB) GetUser(filter interface{}) (*User, error) {
	coll := m.getCollection("users")
	// filter := bson.D{}

	result := coll.FindOne(context.TODO(), filter)
	if result.Err() != nil {
		log.Println("inside GetUser, user not found: ", result.Err())
		return nil, result.Err()
	}
	var userMongo bson.M
	err := result.Decode(&userMongo)
	if err != nil {
		log.Println("inside GetUser, fail to decode user: ", err)
		return &User{}, err
	}
	log.Println("inside GetUser, userMongo: ", userMongo)

	var user User
	user.ID = userMongo["_id"].(primitive.ObjectID).Hex()
	user.Email = userMongo["email"].(string)
	roomsPrimitiveArray := userMongo["rooms"].(primitive.A)
	roomsInterface := []interface{}(roomsPrimitiveArray)
	for _, v := range roomsInterface {
		user.Rooms = append(user.Rooms, v.(string))
	}
	user.Username = userMongo["username"].(string)
	user.UserImage = userMongo["user_image"].(string)

	log.Println("inside GetUser, user: ", user)
	return &user, nil
}

func (m *MongoDB) AddUser(user interface{}) (string, error) {

	coll := m.getCollection("users")
	// doc := bson.D{{"name", roomName}}
	result, err := coll.InsertOne(context.TODO(), user)

	if err != nil {
		log.Println("failed to insert user: ", err)
		return "", err
	}

	return fmt.Sprintf("%v", result.InsertedID), nil
	// return result.InsertedID.(string), nil
}
func (m *MongoDB) UpdateUser(filter interface{}, update interface{}, options *options.UpdateOptions) error {

	coll := m.getCollection("users")
	// doc := bson.D{{"name", roomName}}
	result, err := coll.UpdateOne(context.TODO(), filter, update, options)
	if err != nil {
		// if result.MatchedCount == 0 {
		// 	log.Println("failed to update user: ")
		// }
		// log.Println("failed to insert user, result: ", result.UpsertedID)
		log.Println("failed to insert user, error: ", err)
		return err
	}

	log.Println("UpdateUser MatchedCount: ", result.MatchedCount, " UpsertedCount: ", result.UpsertedCount)
	// userFilter := bson.M{"_id": result.UpsertedID.(primitive.ObjectID)}

	// return fmt.Sprintf("%v", result.UpsertedID), nil
	return nil
}

func (m *MongoDB) AddUsers(users []interface{}) ([]string, error) {

	coll := m.getCollection("users")
	// doc := bson.D{{"name", roomName}}
	results, err := coll.InsertMany(context.TODO(), users)
	if err != nil {
		log.Println("failed to insert users: ", err)
		return nil, err
	}

	var retValues []string
	for _, result := range results.InsertedIDs {
		fmt.Println("doc user id: ", result)
		retValues = append(retValues, result.(primitive.ObjectID).Hex())
	}

	return retValues, nil
}
