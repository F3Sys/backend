package server

import (
	"backend/internal/database"
	"backend/internal/sqlc"
	"database/sql"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"net/netip"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sqids/sqids-go"
)

func Sqids() (*sqids.Sqids, error) {
	sqid, err := sqids.New(sqids.Options{
		MinLength: uint8(7),
		Alphabet:  "23456789CFGHJMPQRVWX",
		Blocklist: []string{},
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
			slog.Default().Info("request",
				"URI", v.URI,
				"status", v.Status,
				"remote_ip", v.RemoteIP,
				"latency", v.Latency,
			)

			return nil
		},
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowMethods: []string{http.MethodGet, http.MethodPatch, http.MethodPost},
	}))

	e.GET("/ip", s.PingHandler)

	e.GET("/visitor", s.GetVisitorHandler)

	e.POST("/visitor", s.PostVisitorHandler)

	e.POST("/node", s.NodeIpHandler)

	protected := e.Group("/protected", middleware.KeyAuth(func(key string, c echo.Context) (bool, error) {
		node, ok, err := s.DB.Password(key)
		if err != nil {
			return false, nil
		}

		c.Set("node", node)

		return ok, nil
	}))

	protected.GET("/info", s.NodeInfoHandler)

	protected.GET("/foods", s.NodeFoodsHandler)

	protected.GET("/table", s.NodeTableHandler)

	protected.GET("/count", s.NodeCountHandler)

	protected.GET("/visitor/:f3sid", s.NodeVisitorLookupHandler)

	protected.POST("/push/entry", s.NodePushEntryHandler)

	protected.POST("/push/foodstall", s.NodePushFoodStallHandler)

	protected.POST("/push/exhibition", s.NodePushExhibitionHandler)

	protected.PATCH("/push", s.NodeUpdatePushHandler)

	protected.PATCH("/status", s.NodeStatusHandler)

	return e
}

func (s *Server) PingHandler(c echo.Context) error {
	ip := c.RealIP()
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		slog.Default().Error("ParseAddr", "error", err)
		return echo.ErrInternalServerError
	}
	return c.String(http.StatusOK, addr.String())
}

func (s *Server) GetVisitorHandler(c echo.Context) error {
	ip := c.RealIP()

	sqid, err := Sqids()
	if err != nil {
		slog.Default().Error("sqids initialization", "error", err)
		return echo.ErrInternalServerError
	}

	visitorF3SiD, err := s.DB.GetVisitor(ip, sqid)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.NoContent(http.StatusNotFound)
		}
		slog.Default().Error("visitor", "error", err)
		return echo.ErrBadRequest
	}

	return c.JSON(http.StatusOK, map[string]string{
		"f3sid": visitorF3SiD,
	})
}

func (s *Server) PostVisitorHandler(c echo.Context) error {
	ip := c.RealIP()

	sqid, err := Sqids()
	if err != nil {
		slog.Default().Error("sqids initialization", "error", err)
		return echo.ErrInternalServerError
	}

	random := rand.Int32()

	visitorF3SiD, err := s.DB.CreateVisitor(ip, int64(random), sqid)
	if err != nil {
		slog.Default().Error("visitor", "error", err)
		return echo.ErrBadRequest
	}

	return c.JSON(http.StatusOK, map[string]string{
		"f3sid": visitorF3SiD,
	})
}

type PushEntry struct {
	VisitorF3SiD string `json:"f3sid"`
}

func (s *Server) NodePushEntryHandler(c echo.Context) error {
	var push PushEntry

	err := c.Bind(&push)
	if err != nil {
		slog.Default().Error("bind", "error", err)
		return echo.ErrBadRequest
	}

	sqid, err := Sqids()
	if err != nil {
		slog.Default().Error("sqids initialization", "error", err)
		return echo.ErrInternalServerError
	}

	pushVisitorID := sqid.Decode(push.VisitorF3SiD)
	if len(pushVisitorID) != 2 {
		slog.Default().Error("sqids decode", "error", "invalid sqids")
		return echo.ErrBadRequest
	}

	node := c.Get("node").(sqlc.Node)

	if node.Type != database.ENTRY {
		slog.Default().Error("node type", "error", "invalid node type")
		return echo.ErrBadRequest
	}

	err = s.DB.PushEntry(node, int64(pushVisitorID[0]), int64(pushVisitorID[1]))
	if err != nil {
		slog.Default().Error("push entry", "error", err)
		return echo.ErrBadRequest
	}

	return c.NoContent(http.StatusOK)
}

type PushFoodStall struct {
	VisitorF3SiD string  `json:"f3sid"`
	Foods        []Foods `json:"foods"`
}

type Foods struct {
	ID       int `json:"id"`
	Quantity int `json:"quantity"`
}

func (s *Server) NodePushFoodStallHandler(c echo.Context) error {
	var push PushFoodStall

	err := c.Bind(&push)
	if err != nil {
		slog.Default().Error("bind", "error", err)
		return echo.ErrBadRequest
	}

	sqid, err := Sqids()
	if err != nil {
		slog.Default().Error("sqids initialization", "error", err)
		return echo.ErrInternalServerError
	}

	pushVisitorID := sqid.Decode(push.VisitorF3SiD)
	if len(pushVisitorID) != 2 {
		slog.Default().Error("sqids decode", "error", "invalid sqids")
		return echo.ErrBadRequest
	}

	node := c.Get("node").(sqlc.Node)

	if node.Type != database.FOODSTALL {
		slog.Default().Error("node type", "error", "invalid node type")
		return echo.ErrBadRequest
	}

	if len(push.Foods) == 0 {
		slog.Default().Error("foods", "error", "no foods")
		return echo.ErrBadRequest
	}

	// Convert push.Foods to database.Foods
	foods := make([]database.Foods, len(push.Foods))
	for i, food := range push.Foods {
		foods[i] = database.Foods{
			ID:       food.ID,
			Quantity: food.Quantity,
		}
	}

	err = s.DB.PushFoodStall(node, int64(pushVisitorID[0]), int64(pushVisitorID[1]), foods)
	if err != nil {
		slog.Default().Error("push foodstall", "error", err)
		return echo.ErrBadRequest
	}

	return c.NoContent(http.StatusOK)
}

type PushExhibition struct {
	VisitorF3SiD string `json:"f3sid"`
}

func (s *Server) NodePushExhibitionHandler(c echo.Context) error {
	var push PushExhibition

	err := c.Bind(&push)
	if err != nil {
		slog.Default().Error("bind", "error", err)
		return echo.ErrBadRequest
	}

	sqid, err := Sqids()
	if err != nil {
		slog.Default().Error("sqids initialization", "error", err)
		return echo.ErrInternalServerError
	}

	pushVisitorID := sqid.Decode(push.VisitorF3SiD)
	if len(pushVisitorID) != 2 {
		slog.Default().Error("sqids decode", "error", "invalid sqids")
		return echo.ErrBadRequest
	}

	node := c.Get("node").(sqlc.Node)

	if node.Type != database.EXHIBITION {
		slog.Default().Error("node type", "error", "invalid node type")
		return echo.ErrBadRequest
	}

	err = s.DB.PushExhibition(node, int64(pushVisitorID[0]), int64(pushVisitorID[1]))
	if err != nil {
		slog.Default().Error("push exhibition", "error", err)
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
		slog.Default().Error("bind", "error", err)
		return echo.ErrBadRequest
	}

	node := c.Get("node").(sqlc.Node)

	err = s.DB.StatusNode(node.ID, int64(status.Level), int64(status.ChargingTime), int64(status.DischargingTime), status.Charging)
	if err != nil {
		slog.Default().Error("status", "error", err)
		return echo.ErrBadRequest
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) NodeInfoHandler(c echo.Context) error {
	node := c.Get("node").(sqlc.Node)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"name": node.Name,
		"type": node.Type,
	})
}

func (s *Server) NodeFoodsHandler(c echo.Context) error {
	node := c.Get("node").(sqlc.Node)

	foods, err := s.DB.Foods(node)
	if err != nil {
		slog.Default().Error("foods", "error", err)
		return echo.ErrBadRequest
	}

	foodsArray := make([]map[string]interface{}, len(foods))
	for i, food := range foods {
		foodsArray[i] = map[string]interface{}{
			"id":    food.ID,
			"name":  food.Name,
			"price": food.Price,
		}
	}

	return c.JSON(http.StatusOK, foodsArray)
}

type VisitorLookup struct {
	visitorF3SiD string `param:"f3sid"`
}

func (s *Server) NodeVisitorLookupHandler(c echo.Context) error {
	var visitorLookup VisitorLookup

	err := c.Bind(&visitorLookup)
	if err != nil {
		slog.Default().Error("bind", "error", err)
		return echo.ErrBadRequest
	}

	sqid, err := Sqids()
	if err != nil {
		slog.Default().Error("sqids initialization", "error", err)
		return echo.ErrInternalServerError
	}

	visitorLookupVisitorID := sqid.Decode(visitorLookup.visitorF3SiD)
	if len(visitorLookupVisitorID) != 2 {
		slog.Default().Error("sqids decode", "error", "invalid sqids")
		return echo.ErrBadRequest
	}

	isFirst, err := s.DB.IsVisitorFirst(int64(visitorLookupVisitorID[0]))
	if err != nil {
		slog.Default().Error("is visitor first", "error", err)
		return echo.ErrBadRequest
	}

	return c.JSON(http.StatusOK, map[string]bool{
		"is_first": isFirst,
	})
}

func (s *Server) NodeTableHandler(c echo.Context) error {
	node := c.Get("node").(sqlc.Node)

	sqid, err := Sqids()
	if err != nil {
		slog.Default().Error("sqids initialization", "error", err)
		return echo.ErrInternalServerError
	}

	switch node.Type {
	case database.ENTRY:
		entryRow, err := s.DB.EntryRow(node, sqid)
		if err != nil {
			slog.Default().Error("entry row", "error", err)
			return echo.ErrBadRequest
		}

		return c.JSON(http.StatusOK, entryRow)
	case database.FOODSTALL:
		foodstallRawLog, err := s.DB.FoodstallRow(node, sqid)
		if err != nil {
			slog.Default().Error("foodstall row", "error", err)
			return echo.ErrBadRequest
		}

		return c.JSON(http.StatusOK, foodstallRawLog)

	case database.EXHIBITION:
		exhibitionRowLog, err := s.DB.ExhibitionRow(node, sqid)
		if err != nil {
			slog.Default().Error("exhibition row", "error", err)
			return echo.ErrBadRequest
		}

		return c.JSON(http.StatusOK, exhibitionRowLog)
	}

	return echo.ErrBadRequest
}

type UpdatePush struct {
	Id       int `json:"id"`
	Quantity int `json:"quantity"`
}

func (s *Server) NodeUpdatePushHandler(c echo.Context) error {
	var push UpdatePush

	err := c.Bind(&push)
	if err != nil {
		slog.Default().Error("bind", "error", err)
		return echo.ErrBadRequest
	}

	node := c.Get("node").(sqlc.Node)

	err = s.DB.UpdatePushNode(node, int64(push.Id), int64(push.Quantity))
	if err != nil {
		slog.Default().Error("update push", "error", err)
		return echo.ErrBadRequest
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) NodeIpHandler(c echo.Context) error {
	ip := c.RealIP()

	node, err := s.DB.IpNode(ip)
	if err != nil {
		slog.Default().Error("ip node", "error", err)
		return echo.ErrBadRequest
	}

	return c.JSON(http.StatusOK, map[string]string{
		"key": node.Key.String,
	})
}

func (s *Server) NodeCountHandler(c echo.Context) error {
	node := c.Get("node").(sqlc.Node)

	switch node.Type {
	case database.ENTRY:
		count, err := s.DB.CountEntry(node)
		if err != nil {
			slog.Default().Error("entry row", "error", err)
			return echo.ErrBadRequest
		}

		return c.JSON(http.StatusOK, map[string]int64{"count": count})
	case database.FOODSTALL:
		count, err := s.DB.CountFoodStall(node)
		if err != nil {
			slog.Default().Error("foodstall row", "error", err)
			return echo.ErrBadRequest
		}

		return c.JSON(http.StatusOK, map[string]int64{"count": count})
	case database.EXHIBITION:
		count, err := s.DB.CountExhibition(node)
		if err != nil {
			slog.Default().Error("exhibition row", "error", err)
			return echo.ErrBadRequest
		}

		return c.JSON(http.StatusOK, map[string]int64{"count": count})
	}

	return echo.ErrBadRequest
}
