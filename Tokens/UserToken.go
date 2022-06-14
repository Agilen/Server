package tokens

import (
	"fmt"
	"time"
)

type UserToken struct {
	PublicInfo []byte
	TimeStamp  int
}

func NewUserToken(publicInfo []byte, lifetime int) *UserToken {
	return &UserToken{
		TimeStamp:  int(time.Now().Add(time.Hour * time.Duration(lifetime)).Unix()),
		PublicInfo: publicInfo,
	}
}

func (t *UserToken) CheckToken() error {
	if t.TimeStamp < int(time.Now().Unix()) {
		return fmt.Errorf("link is no longer active")
	}

	return nil
}
