package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-herder/internal/herder"
	"go-herder/internal/repository"
	"net/http"
	"strconv"
)

type Config struct {
	Addr string `yaml:"addr"`
}

type API struct {
	c *Config
	h *herder.Herder
	r *repository.Repository
	e *gin.Engine
}

func New(c Config, h *herder.Herder, r *repository.Repository) *API {
	var api = &API{
		c: &c,
		h: h,
		r: r,
		e: gin.Default(),
	}
	api.initRoutes()
	return api
}

func (api *API) Run() error {
	if api.c == nil {
		return errors.New("can't run API without *Config")
	}
	if api.h == nil {
		return errors.New("can't run API without *Herder")
	}
	if api.r == nil {
		return errors.New("can't run API without *Repository")
	}
	return api.e.Run(api.c.Addr)
}

func (api *API) initRoutes() {
	apiRoutes := api.e.Group("/api", api.middlewareJsonHeaders)
	{
		apiRoutes.GET("/ping", api.ping)

		authGroup := apiRoutes.Group("/auth")
		{
			authGroup.GET("/", api.authGet)
			authGroup.POST("/", api.authPost)
		}

		herderGroup := apiRoutes.Group("/herder", api.middlewareAuth)
		{
			herderGroup.GET("/state", api.herderState)
			herderGroup.GET("/restart", api.herderRestart)
			herderGroup.GET("/kill", api.herderKill)
			herderGroup.GET("/run", api.herderRun)

			processesGroup := herderGroup.Group("/processes", api.checkProcessExistsMiddleware)
			{
				processesGroup.GET("/:id/state", api.processState)
				processesGroup.GET("/:id/restart", api.processRestart)
				processesGroup.GET("/:id/kill", api.processKill)
				processesGroup.GET("/:id/run", api.processRun)
			}
		}
	}
	//api.e.Static("/", "/gui")
}
func (api *API) ping(c *gin.Context) {
	c.Header("Content-Type", "text/plain")
	c.String(http.StatusOK, "pong")
}

func (api *API) middlewareJsonHeaders(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.Next()
}
func (api *API) checkProcessExistsMiddleware(c *gin.Context) {
	processID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Error(&Err{Code: 400, Message: "bad process id"}))
		return
	}
	err = api.h.CheckProcessExists(processID)
	if err != nil {
		c.JSON(http.StatusBadRequest, Error(&Err{Code: 400, Message: err.Error()}))
		return
	}
	c.Set("id", processID)
	c.Next()
}

func (api *API) authGet(c *gin.Context) {
	sessionID, err := c.Cookie("session")
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorUnauthorized)
		return
	}
	session, err := api.r.GetSession(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Error(&Err{Code: http.StatusInternalServerError, Message: err.Error()}))
		return
	}
	if session.UserAgent != c.Request.UserAgent() {
		err = api.r.DeleteSession(sessionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Error(&Err{Code: http.StatusInternalServerError, Message: err.Error()}))
			return
		}
		c.JSON(http.StatusInternalServerError, Error(&Err{Code: http.StatusUnauthorized, Message: "bad user agent, please auth again"}))
		return
	}
	c.JSON(http.StatusOK, Ok)
}
func (api *API) authPost(c *gin.Context) {
	sessionID, err := c.Cookie("session")
	if err == nil {
		err = api.r.DeleteSession(sessionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Error(&Err{Code: http.StatusInternalServerError, Message: err.Error()}))
			return
		}
	}

	session, err := api.r.CreateSession(c.RemoteIP(), c.Request.UserAgent())
	if err != nil {
		c.JSON(http.StatusInternalServerError, Error(&Err{Code: http.StatusInternalServerError, Message: err.Error()}))
		return
	}
	if session.UserAgent != c.Request.UserAgent() {
		err = api.r.DeleteSession(sessionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Error(&Err{Code: http.StatusInternalServerError, Message: err.Error()}))
			return
		}
		c.JSON(http.StatusInternalServerError, Error(&Err{Code: http.StatusUnauthorized, Message: "bad user agent, please auth again"}))
		return
	}
	c.JSON(http.StatusOK, Ok)
}

func (api *API) middlewareAuth(c *gin.Context) {
	c.Next()
}
func (api *API) herderState(c *gin.Context) {
	if c.Query("type") != "json" {
		c.Header("Content-Type", "text/plain")
	}
	states, err := api.h.GetStates()

	if c.Query("type") == "json" {
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorInternalServer)
		} else {
			c.JSON(http.StatusOK, Ok.WithData(states))
		}
	} else {
		if err != nil {
			c.String(http.StatusInternalServerError, "internal server error")
		} else {
			c.String(http.StatusOK, herder.StringifyStates(states))
		}
	}
}
func (api *API) herderRun(c *gin.Context) {
	err := api.h.RunAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorInternalServer)
		return
	}
	c.JSON(http.StatusOK, Ok)
}
func (api *API) herderRestart(c *gin.Context) {
	api.h.RestartAll()
	c.JSON(http.StatusOK, Ok)
}
func (api *API) herderKill(c *gin.Context) {
	api.h.KillAll()
	c.JSON(http.StatusOK, Ok)
}

func (api *API) processState(c *gin.Context) {
	if c.Query("type") == "json" {
		c.Header("Content-Type", "application/json")
	} else {
		c.Header("Content-Type", "text/plain")
	}
	state, err := api.h.GetState(c.GetInt("id"))

	if c.Query("type") == "json" {
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorInternalServer)
		} else {
			c.JSON(http.StatusOK, Ok.WithData(state))
		}
	} else {
		if err != nil {
			c.String(http.StatusInternalServerError, "internal server error")
		} else {
			c.String(http.StatusOK, herder.StringifyStates([]herder.ProcessState{state}))
		}
	}
}
func (api *API) processRun(c *gin.Context) {
	err := api.h.Run(c.GetInt("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Error(&Err{Code: http.StatusInternalServerError, Message: err.Error()}))
		return
	}
	c.JSON(http.StatusOK, Ok)
}
func (api *API) processRestart(c *gin.Context) {
	err := api.h.Restart(c.GetInt("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Error(&Err{Code: http.StatusInternalServerError, Message: err.Error()}))
		return
	}
	c.JSON(http.StatusOK, Ok)
}
func (api *API) processKill(c *gin.Context) {
	err := api.h.Kill(c.GetInt("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Error(&Err{Code: http.StatusInternalServerError, Message: err.Error()}))
		return
	}
	c.JSON(http.StatusOK, Ok)
}
