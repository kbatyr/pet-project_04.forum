package models

import (
	"errors"
)

// Если в БД отсутсвует строка/и по запросу
var ErrNoRecord = errors.New("models: no matching entry found ")
var ErrUserAlreadyExists = errors.New("such user already exists")
var ErrNoSuchUserOrPSW = errors.New("incorrect username or password")

type User struct {
	ID               int
	UserName         string
	Email            string
	Password         string
	Role             string
	RegistrationDate string
}

type CreateUserDTO struct {
	UserName    string
	Email       string
	Password    string
	PasswordRPT string
}

type Post struct {
	ID           int
	Title        string
	Author       string
	Content      string
	CreationDate string
	Categories   []string
	Comments     []*Comment
	Rating       *Rating
}

type Session struct {
	UUID    string
	Expires string
	User_ID int
}

type Category struct {
	ID   int
	Name string
}

type Comment struct {
	ID      int
	Message string
	Date    string
	Author  string
	Rating  *Rating
}

type Rating struct {
	Likes    int64
	Dislikes int64
}
