package user

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const TIMEOUT = 10*time.Second

type repo struct {
	client *mongo.Client
}

func NewMongoRepository(client *mongo.Client) Repository {
	return &repo{
		client: client,
	}
}

func (r *repo) FindByChatID(chatID int64) (*User, error) {
	var user User
	filter := bson.M{"chat_id": chatID}
	collection := r.client.Database("sakimoto").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT)
	err := collection.FindOne(ctx, filter).Decode(&user)
	return &user, err
}

func (r *repo) FindAll() ([] *User, error) {
	var users []*User
	var user User
	collection := r.client.Database("sakimoto").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT)
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		err := cur.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *repo) InsertUser(user *User) error {
	collection := r.client.Database("sakimoto").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT)
	user.CreatedAt = time.Now()
	_, err := collection.InsertOne(ctx, &user)
	return err
}
