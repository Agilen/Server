package tokens

import (
	"errors"
	"time"
)

type LinkToken struct {
	ID        string
	TimeStamp int
}

func NewToken(ID string, lifetime int) *LinkToken {
	return &LinkToken{
		TimeStamp: int(time.Now().Add(time.Minute * time.Duration(lifetime)).Unix()),
		ID:        ID,
	}
}

func (t *LinkToken) CheckToken() error {
	if t.TimeStamp < int(time.Now().Unix()) {
		return errors.New("link is no longer active")
	}

	return nil
}
