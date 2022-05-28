package server

import (
	"fmt"
	"net/http"

	"github.com/Agilen/Server/mailing"
	"github.com/Agilen/Server/model"
	"github.com/labstack/echo"
)

type (
	RegistrationRequest struct {
		Login    string
		Mail     string
		Password string
	}

	RegistrationResponce struct {
		Nickname string
		Login    string
	}

	LoginRequest struct {
		Login    string
		Password string
	}

	LoginResponce struct {
		Nickname string
		Login    string
	}

	FindUserRequest struct {
		Parm string
	}

	FindUserResponce struct {
		Users []model.User
	}

	VerifyRequrst struct {
		id string
	}
)

func (s *Server) CreateUserController(c echo.Context) error {
	req := new(RegistrationRequest)
	if err := c.Bind(req); err != nil {
		return s.HttpErrorHandler(c, err, http.StatusBadRequest)
	}

	u := &model.User{
		Login:    req.Login,
		Mail:     req.Mail,
		Password: req.Password,
	}

	id, err := s.store.User().CreateUser(u)
	if err != nil {
		return s.HttpErrorHandler(c, err, http.StatusInternalServerError)
	}

	err = mailing.SendMail(req.Mail, s.CreateVerifyLink(id))
	if err != nil {

		e := s.store.User().DeleteUser(id)
		if err != nil {
			fmt.Println("error", e)
		}
		return s.HttpErrorHandler(c, err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, RegistrationResponce{Nickname: u.NickName, Login: u.Login})
}

func (s *Server) VerifyUserController(c echo.Context) error {

	fmt.Println("hello")
	id := c.QueryParam("id")
	if id == "" {
		return s.HttpErrorHandler(c, fmt.Errorf("id is nil"), http.StatusBadRequest)
	}
	fmt.Println(id)
	if token, ok := s.LinkStore[id]; ok {

		delete(s.LinkStore, id)

		err := token.CheckToken()
		if err != nil {
			return s.HttpErrorHandler(c, err, http.StatusInternalServerError)
		}

		err = s.store.User().ChangeStatus(token.ID)
		if err != nil {
			return s.HttpErrorHandler(c, err, http.StatusInternalServerError)
		}
	} else {
		return s.HttpErrorHandler(c, fmt.Errorf("this link is not avtive"), http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, nil)
}

func (s *Server) LoginController(c echo.Context) error {
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return s.HttpErrorHandler(c, err, http.StatusBadRequest)
	}

	u := &model.User{
		Login:    req.Login,
		Password: req.Password,
	}
	if err := s.store.User().CheckUser(u); err != nil {
		return s.HttpErrorHandler(c, err, http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, LoginResponce{Nickname: u.NickName, Login: u.Login})
}

func (s *Server) FindUserController(c echo.Context) error {
	req := new(FindUserRequest)
	if err := c.Bind(req); err != nil {
		return s.HttpErrorHandler(c, err, http.StatusBadRequest)
	}
	var u model.User

	if req.Parm[0] == byte(64) {
		u.Login = req.Parm[1:]
	}
	u.NickName = req.Parm

	uu, err := s.store.User().FindUser(&u)
	if err != nil {
		return s.HttpErrorHandler(c, err, http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, FindUserResponce{Users: *uu})
}
