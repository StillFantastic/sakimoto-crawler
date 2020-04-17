package user

import "time"

type User struct {
	Username  string    `json:"username" bson:"username"`
	ChatID    int64    `json:"chat_id" bson:"chat_id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
