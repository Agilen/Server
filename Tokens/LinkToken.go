package tokens

import (
	"errors"
	"fmt"
	"time"
)

type LinkToken struct {
	ID        string
	TimeStamp int
}

func NewToken(ID string, lifetime int) *LinkToken {
	fmt.Println(time.Now().Add(time.Minute * time.Duration(lifetime)).Clock())
	return &LinkToken{
		TimeStamp: int(time.Now().Add(time.Minute * time.Duration(lifetime)).Unix()),
		ID:        ID,
	}
}

func (t *LinkToken) CheckToken() error {
	fmt.Println(time.Now().Clock())

	fmt.Println(time.Now().Unix())
	if t.TimeStamp < int(time.Now().Unix()) {
		return errors.New("link is no longer active")
	}

	return nil
}
