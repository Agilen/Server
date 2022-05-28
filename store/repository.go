package store

import "github.com/Agilen/Server/model"

type UserRepository interface {
	CreateUser(u *model.User) (string, error)
	CheckUser(u *model.User) error
	FindUser(u *model.User) (*[]model.User, error)
	ChangeStatus(id string) error
	DeleteUser(id string) error
}
