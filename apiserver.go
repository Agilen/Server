package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	tokens "github.com/Agilen/Server/Tokens"
	"github.com/Agilen/Server/config"
	"github.com/Agilen/Server/mycrypto"
	"github.com/Agilen/Server/store"
	"github.com/Agilen/Server/store/sqlstore"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type HttpError struct {
	Controller string
	Method     string
	Body       string
	Error      string
}

type Server struct {
	config         *config.Config
	router         *echo.Echo
	logger         *logrus.Logger
	store          store.Store
	cc             *mycrypto.CryptoContext
	LinkStore      map[string]tokens.LinkToken
	UserTokenStore map[string]tokens.UserToken
}

func NewServer(store store.Store, config *config.Config) (*Server, error) {
	var err error
	s := &Server{
		router:         echo.New(),
		logger:         logrus.New(),
		store:          store,
		config:         config,
		LinkStore:      make(map[string]tokens.LinkToken),
		UserTokenStore: make(map[string]tokens.UserToken),
	}

	s.cc, err = mycrypto.NewCryptoContext()
	if err != nil {
		return nil, err
	}

	s.configureRouter()

	return s, nil
}

func (s *Server) configureRouter() {
	s.router.POST("/reg", s.CreateUserController)
	s.router.POST("/login", s.LoginController)
	s.router.GET("/user/verify", s.VerifyUserController)
	s.router.GET("/public", s.GetPublicInfoController)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		encData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		data, err := s.cc.Decrypt(encData)
		if err != nil {
			fmt.Println(err)
			return
		}

		r.Body = ioutil.NopCloser(bytes.NewReader(data))
	}

	s.router.ServeHTTP(w, r)
}

func Start(config *config.Config) error {
	db, err := sqlstore.NewDB("DB.db")
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	defer sqlDB.Close()
	store := sqlstore.New(db)

	s, err := NewServer(store, config)
	if err != nil {
		return err
	}
	fmt.Println("start")
	return http.ListenAndServe(":10000", s)
}

func (s *Server) HttpErrorHandler(c echo.Context, err error, httpErr int) error {
	url := c.Request().URL.String()
	method := c.Request().Method
	body, _ := ioutil.ReadAll(c.Request().Body)
	k, _ := json.MarshalIndent(HttpError{Controller: url, Method: method, Body: string(body), Error: err.Error()}, "", "\t")
	s.logger.Error(string(k))

	return c.JSON(httpErr, err.Error())
}
