package server

import (
	"net/http"
	"strconv"
	"strings"

	tokens "github.com/Agilen/Server/Tokens"
	"github.com/google/uuid"
)

func (s *Server) DeletePort(port string) {
	if val, ok := s.PortsStore[port]; ok && val {
		delete(s.PortsStore, port)
	}
}

func InitPortsMap() map[string]bool {
	m := make(map[string]bool)
	for i := 10000; i < 20000; i++ {

		err := http.ListenAndServe(":"+strconv.Itoa(i), nil)

		if err != nil {
			m[strconv.Itoa(i)] = true
			continue
		}
		m[strconv.Itoa(i)] = false
	}

	return m
}

func (s *Server) CreateVerifyLink(userId string) string {
	id := uuid.New().String()

	s.LinkStore[id] = *tokens.NewToken(userId, s.config.TokenTTL)

	return strings.Join([]string{s.config.BaseUrl, s.config.Port, "/user/verify?id=", id}, "")
}
