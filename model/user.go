package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID        string
	NickName  string
	Login     string `gorm:"not null;unique"`
	Mail      string `gorm:"not null; unique"`
	Password  string `gorm:"not null;<-"`
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime `gorm:"index"`
}
