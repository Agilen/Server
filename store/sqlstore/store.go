package sqlstore

import (
	"github.com/Agilen/Server/model"
	"github.com/Agilen/Server/store"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Store struct {
	db             *gorm.DB
	UserRepository *UserRepository
}

func New(db *gorm.DB) *Store {
	db.AutoMigrate(&model.User{})

	return &Store{
		db: db,
	}
}

func NewDB(databaseURL string) (*gorm.DB, error) {

	// db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	// if err != nil {
	// 	return nil, err
	// }
	db, err := gorm.Open(sqlite.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (s *Store) User() store.UserRepository {
	if s.UserRepository != nil {
		return s.UserRepository
	}

	s.UserRepository = &UserRepository{
		store: s,
	}

	return s.UserRepository
}
