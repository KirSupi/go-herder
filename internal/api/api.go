package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-herder/internal/herder"
	"net/http"
)

type Config struct {
	Addr string
}

type API struct {
	c *Config
	h *herder.Herder
	e *gin.Engine
}

func New(c Config, h *herder.Herder) *API {
	var api = &API{
		c: &c,
		h: h,
		e: gin.Default(),
	}
	api.initRoutes()
	return api
}

func (api *API) Run() error {
	if api.c != nil {
		return errors.New("can't start API without *Config")
	}
	return api.e.Run(api.c.Addr)
}

func (api *API) initRoutes() {
	api.e.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
}
