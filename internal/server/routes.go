package server

import (
	sql "backend/internal/sqlc"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/netip"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.IPExtractor = echo.ExtractIPFromXFFHeader()
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

	e.GET("/health", s.HealthHandler)

	e.GET("/visitor", s.VisitorHandler)

	protected := e.Group("/protected", middleware.KeyAuth(func(key string, c echo.Context) (bool, error) {
		node, ok, err := s.DB.Password(key)
		if err != nil {
			return false, nil
		}

		c.Set("node", node)

		return ok, nil
	}))

	protected.POST("/push", s.NodePushHandler)

	protected.POST("/status", s.NodeStatusHandler)

	return e
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) HealthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.DB.Health())
}

func (s *Server) VisitorHandler(c echo.Context) error {
	ip := c.RealIP()
	//ip := "127.0.2.1"
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		log.Info().Err(err).Send()
		return echo.ErrBadRequest
	}
	uuid, err := s.DB.Visitor(addr)
	if err != nil {
		log.Info().Err(err).Send()
		return echo.ErrBadRequest
	}

	return c.JSON(http.StatusOK, map[string]string{
		"uuid": uuid,
	})
}

type Push struct {
	visitorUUID string `query:"visitor_uuid"`
	quantity    int    `query:"quantity"`
}

func (s *Server) NodePushHandler(c echo.Context) error {
	var push Push
	err := c.Bind(&push)
	if err != nil {
		return echo.ErrBadRequest
	}

	node := c.Get("node").(*sql.Node)

	err = s.DB.PushNode(int(node.ID), push.visitorUUID, push.quantity)
	if err != nil {
		return echo.ErrBadRequest
	}

	return c.NoContent(http.StatusOK)
}

type Status struct {
	charging        bool `query:"charging"`
	chargingTime    int  `query:"charging_time"`
	dischargingTime int  `query:"discharging_time"`
	level           int  `query:"level"`
}

func (s *Server) NodeStatusHandler(c echo.Context) error {
	var status Status
	err := c.Bind(&status)
	if err != nil {
		return echo.ErrBadRequest
	}

	node := c.Get("node").(*sql.Node)

	err = s.DB.StatusNode(int(node.ID), status.charging, status.chargingTime, status.dischargingTime, status.level)
	if err != nil {
		return echo.ErrBadRequest
	}

	return c.NoContent(http.StatusOK)
}
