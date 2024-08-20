package server

import (
	sql "backend/internal/sqlc"
	"net/http"
	"net/netip"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/sqids/sqids-go"
)

func Sqids() (*sqids.Sqids, error) {
	minlength, _ := strconv.Atoi(os.Getenv("SQIDS_MINLENGTH"))
	alphabet := os.Getenv("SQIDS_ALPHABET")

	sqid, err := sqids.New(sqids.Options{
		MinLength: uint8(minlength),
		Alphabet:  alphabet,
	})
	if err != nil {
		return nil, err
	}

	return sqid, nil
}

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.IPExtractor = echo.ExtractIPFromXFFHeader()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogRemoteIP: true,
		LogLatency:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Info().
				Str("URI", v.URI).
				Int("status", v.Status).
				Str("remote_ip", v.RemoteIP).
				Dur("latency", v.Latency).
				Msg("request")

			return nil
		},
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	e.GET("/ping", s.PingHandler)

	e.GET("/visitor", s.VisitorHandler)

	protected := e.Group("/protected", middleware.KeyAuth(func(key string, c echo.Context) (bool, error) {
		node, ok, err := s.DB.Password(key)
		if err != nil {
			return false, nil
		}

		c.Set("node", node)

		return ok, nil
	}))

	protected.GET("/ping", s.PingHandler)

	protected.GET("/info", s.NodeInfoHandler)

	protected.GET("/visitor/:sqid", s.NodeVisitorLookupHandler)

	protected.POST("/push", s.NodePushHandler)

	protected.PATCH("/status", s.NodeStatusHandler)

	return e
}

func (s *Server) PingHandler(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func (s *Server) VisitorHandler(c echo.Context) error {
	ip := c.RealIP()
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return echo.ErrInternalServerError
	}

	sqid, err := Sqids()
	if err != nil {
		return echo.ErrInternalServerError
	}

	visitorSQID, err := s.DB.Visitor(addr, sqid)
	if err != nil {
		return echo.ErrBadRequest
	}

	return c.JSON(http.StatusOK, map[string]string{
		"sqid": visitorSQID,
	})
}

type Push struct {
	VisitorSQID string `json:"sqid"`
	Quantity    int    `json:"quantity"`
}

func (s *Server) NodePushHandler(c echo.Context) error {
	var push Push

	err := c.Bind(&push)
	if err != nil {
		return echo.ErrBadRequest
	}

	sqid, err := Sqids()
	if err != nil {
		return echo.ErrInternalServerError
	}

	pushVisitorID := sqid.Decode(push.VisitorSQID)

	node := c.Get("node").(*sql.Node)

	err = s.DB.PushNode(node, int64(pushVisitorID[0]), int32(pushVisitorID[1]), push.Quantity)
	if err != nil {
		return echo.ErrBadRequest
	}

	return c.NoContent(http.StatusOK)
}

type Status struct {
	Charging        bool `json:"charging"`
	ChargingTime    int  `json:"charging_time"`
	DischargingTime int  `json:"discharging_time"`
	Level           int  `json:"level"`
}

func (s *Server) NodeStatusHandler(c echo.Context) error {
	var status Status
	err := c.Bind(&status)
	if err != nil {
		return echo.ErrBadRequest
	}

	node := c.Get("node").(*sql.Node)

	err = s.DB.StatusNode(node.ID, int32(status.Level), int32(status.ChargingTime), int32(status.DischargingTime), status.Charging)
	if err != nil {
		return echo.ErrBadRequest
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) NodeInfoHandler(c echo.Context) error {
	node := c.Get("node").(*sql.Node)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"name":  node.Name,
		"type":  node.Type,
		"price": node.Price,
	})
}

type VisitorLookup struct {
	VisitorSQID string `param:"sqid"`
}

func (s *Server) NodeVisitorLookupHandler(c echo.Context) error {
	var visitorLookup VisitorLookup

	err := c.Bind(&visitorLookup)
	if err != nil {
		return echo.ErrBadRequest
	}

	sqid, err := Sqids()
	if err != nil {
		return echo.ErrInternalServerError
	}

	visitorLookupVisitorID := sqid.Decode(visitorLookup.VisitorSQID)

	isFirst, err := s.DB.IsVisitorFirst(int64(visitorLookupVisitorID[0]))
	if err != nil {
		return echo.ErrBadRequest
	}

	return c.JSON(http.StatusOK, map[string]bool{
		"is_first": isFirst,
	})
}
