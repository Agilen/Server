package sqlstore

import (
	"errors"
	"fmt"

	"github.com/Agilen/Server/model"
	"github.com/google/uuid"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) CreateUser(u *model.User) (string, error) {

	var count int64
	if err := r.store.db.Where("login = ?", u.Login).Find(&[]model.User{}).Count(&count).Error; err != nil {
		return "", err
	}
	if count > 0 {
		return "", errors.New("name is already in use")
	}
	u.ID = uuid.NewString()
	u.IsActive = false
	u.NickName = u.Login

	if err := r.store.db.Create(&u).Scan(&u).Error; err != nil {
		return "", err
	}

	return u.ID, nil
}

func (r *UserRepository) CheckUser(u *model.User) error {

	if err := r.store.db.Where(u).Find(&u).Error; err != nil {
		return err
	}

	if u.ID == "" {
		return fmt.Errorf("login or password is wrong")
	}
	fmt.Println(u)

	if !u.IsActive {
		return errors.New("pls active your account")
	}

	return nil
}

func (r *UserRepository) FindUser(u *model.User) (*[]model.User, error) {
	var uu *[]model.User

	if err := r.store.db.Where(u).Find(uu).Error; err != nil {
		return nil, err
	}

	return uu, nil
}

func (r *UserRepository) ChangeStatus(id string) error {
	if err := r.store.db.Table("users").Where("id = ?", id).Update("is_active", 1).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) DeleteUser(id string) error {
	if err := r.store.db.Delete(model.User{ID: id}).Error; err != nil {
		return err
	}
	return nil
}
