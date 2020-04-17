package user

type Reader interface {
	FindByChatID(charID int64) (*User ,error)
	FindAll() ([]*User, error)
}

type Writer interface {
	InsertUser(user *User) error
}

type Repository interface {
	Reader
	Writer
}
