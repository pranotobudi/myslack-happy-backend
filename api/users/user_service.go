package users

import (
	"log"

	"github.com/pranotobudi/myslack-happy-backend/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IUserService interface {
	GetUser(email string) (*mongodb.User, error)
	UserAuth(userAuth mongodb.UserAuth) (*mongodb.User, error)
	UpdateUserRooms(userMongo mongodb.User) (*mongodb.User, error)
}
type userService struct {
	repo mongodb.IMongoDB
}

// NewUserService will return userService object
func NewUserService() *userService {
	r := mongodb.NewMongoDB()
	return &userService{repo: r}
}

// GetUser will return User based on email
func (s *userService) GetUser(email string) (*mongodb.User, error) {
	filter := bson.M{"email": email}
	return s.repo.GetUser(filter)
}

// UserAuth will return user if exist or create new user if not exist
func (s *userService) UserAuth(userAuth mongodb.UserAuth) (*mongodb.User, error) {
	log.Println("UserService - UserAuth: ", userAuth)
	filter := bson.M{"email": userAuth.Email}
	userPtr, err := s.repo.GetUser(filter)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, err)
		// return
		log.Println("inside room_io_handler-UserAuth error: ", err)
		return nil, err
	}
	log.Println("UserService - userPtr: ", *userPtr)

	if userPtr != nil {
		// two possibility:
		// empty User{}, means: user found but failed to decode
		// non-empty User, means: user found and success to decode
		// for both result we'll return it anyway

		log.Println("UserService - userPtr is not nil")
		return userPtr, nil
	}

	// user == nil, user not found
	// register
	userDoc := bson.D{{"email", userAuth.Email}, {"username", ""}, {"user_image", userAuth.UserImage}, {"rooms", bson.A{}}}
	userID, err := s.repo.AddUser(userDoc)

	// return User data as response
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	filter = bson.M{"_id": objID}
	userPtr, err = s.repo.GetUser(filter)
	if err != nil {
		return nil, err
	}

	return userPtr, nil
}

// UpdateUserRooms will update rooms field for each user
func (s *userService) UpdateUserRooms(userMongo mongodb.User) (*mongodb.User, error) {
	filter := bson.M{"email": userMongo.Email}
	opts := options.Update().SetUpsert(true)

	// remove all user rooms first
	update := bson.D{{"$set", bson.M{"rooms": []string{}}}}
	err := s.repo.UpdateUser(filter, update, opts)
	if err != nil {
		return nil, err
	}

	// add room one by one
	for _, room := range userMongo.Rooms {
		update := bson.D{{"$push", bson.M{"rooms": room}}}
		err = s.repo.UpdateUser(filter, update, opts)
		if err != nil {
			return nil, err
		}
	}

	// get user with updated rooms element
	userPtr, err := s.repo.GetUser(filter)
	if err != nil {
		return nil, err
	}

	return userPtr, nil
}
