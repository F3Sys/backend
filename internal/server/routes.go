package server

import (
	"github.com/rs/zerolog/log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Info().
				Str("URI", v.URI).
				Int("status", v.Status).
				Msg("request")

			return nil
		},
	}))
	e.Use(middleware.Recover())

	e.GET("/", s.HelloWorldHandler)

	e.GET("/health", s.healthHandler)

	e.GET("/visitor", s.visitorHandler)

	return e
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}

func (s *Server) visitorHandler(c echo.Context) error {
	ip := c.RealIP()
	uuid, err := s.db.Visitor(ip)
	if err != nil {
		log.Info().Err(err).Send()
		return echo.ErrBadRequest
	}

	return c.String(http.StatusOK, uuid)
}

func (s *Server) pushNodeHandler(c echo.Context) error {
	return c.String(http.StatusOK, "update")
}
