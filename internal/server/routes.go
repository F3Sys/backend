package server

import (
	"backend/internal/database"
	"backend/internal/sqlc"
	"errors"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"net/netip"

	"github.com/jackc/pgx/v5"
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

func TypeMiddleware(nodeTypes ...sqlc.NodeType) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			node := c.Get("node").(sqlc.Node)
			for _, nodeType := range nodeTypes {
				if node.Type == nodeType {
					return next(c)
				}
			}
			slog.Error("node type middleware", "error", "invalid node type")
			return echo.ErrBadRequest
		}
	}
}

type (
	Host struct {
		Echo *echo.Echo
	}
)

func (s *Server) ApiRoutes() *echo.Echo {
	api := echo.New()
	api.IPExtractor = func(r *http.Request) string {
		return r.Header.Get("Fly-Client-IP")
	}
	api.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// 	// AllowOrigins: []string{"http://localhost:4000", "https://aicj.io"},
		AllowMethods: []string{http.MethodGet, http.MethodPatch, http.MethodPost},
	}))
	// api.Use(middleware.CORS())

	// api.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(20))))
	// limiterStore := middleware.NewRateLimiterMemoryStore(rate.Limit(10))

	api.GET("/ip", s.PingHandler)

	api.GET("/visitor", s.VisitorHandler) // middleware.RateLimiter(limiterStore)

	api.POST("/node", s.NodeIpHandler)

	// api.POST("/vote", s.VoteHandler) // middleware.RateLimiter(limiterStore)

	protected := api.Group("/protected", middleware.KeyAuth(func(key string, c echo.Context) (bool, error) {
		node, ok, err := s.DB.Password(key)
		if err != nil {
			return false, nil
		}

		c.Set("node", node)

		return ok, nil
	}))

	protected.GET("/info", s.NodeInfoHandler)

	protected.GET("/foods", s.NodeFoodsHandler, TypeMiddleware(sqlc.NodeTypeFOODSTALL))

	protected.GET("/table", s.NodeTableHandler, TypeMiddleware(sqlc.NodeTypeENTRY, sqlc.NodeTypeFOODSTALL, sqlc.NodeTypeEXHIBITION))

	protected.GET("/count", s.NodeCountHandler, TypeMiddleware(sqlc.NodeTypeENTRY, sqlc.NodeTypeFOODSTALL, sqlc.NodeTypeEXHIBITION))

	protected.GET("/food_count", s.NodeFoodCountHandler, TypeMiddleware(sqlc.NodeTypeFOODSTALL))

	protected.GET("/entry_count", s.NodeEntryTypeCountHandler, TypeMiddleware(sqlc.NodeTypeENTRY))

	protected.GET("/visitor/:f3sid", s.NodeVisitorLookupHandler, TypeMiddleware(sqlc.NodeTypeENTRY, sqlc.NodeTypeFOODSTALL, sqlc.NodeTypeEXHIBITION))

	push := protected.Group("/push")

	push.POST("/entry", s.NodePushEntryHandler, TypeMiddleware(sqlc.NodeTypeENTRY))

	push.POST("/foodstall", s.NodePushFoodStallHandler, TypeMiddleware(sqlc.NodeTypeFOODSTALL))

	push.POST("/exhibition", s.NodePushExhibitionHandler, TypeMiddleware(sqlc.NodeTypeEXHIBITION))

	push.PATCH("/foodstall", s.NodeUpdateFoodStallHandler, TypeMiddleware(sqlc.NodeTypeFOODSTALL))

	data := protected.Group("/data")

	data.GET("/entry", s.NodeEntryPerHourCountHandler, TypeMiddleware(sqlc.NodeTypeENTRY))

	data.GET("/foodstall", s.NodeFoodStallPerHourCountHandler, TypeMiddleware(sqlc.NodeTypeFOODSTALL))

	data.GET("/exhibition", s.NodeExhibitionPerHourCountHandler, TypeMiddleware(sqlc.NodeTypeEXHIBITION))

	protected.PATCH("/status", s.NodeStatusHandler, TypeMiddleware(sqlc.NodeTypeENTRY, sqlc.NodeTypeFOODSTALL, sqlc.NodeTypeEXHIBITION))

	return api
}

func (s *Server) RegisterRoutes() *echo.Echo { // Hosts
	hosts := map[string]*Host{}

	//-----
	// API
	//-----

	api := s.ApiRoutes()

	hosts["api.aicj.io"] = &Host{api}

	//------
	// Public Website
	//------

	public := echo.New()
	public.IPExtractor = func(r *http.Request) string {
		return r.Header.Get("Fly-Client-IP")
	}

	hosts["aicj.io"] = &Host{public}

	public.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, c.Request().Header.Get("Fly-Client-IP"))
	})

	public.GET("/panic", func(c echo.Context) error {
		panic("panic")
	})

	//---------
	// Node Website
	//---------

	node := echo.New()
	node.IPExtractor = func(r *http.Request) string {
		return r.Header.Get("Fly-Client-IP")
	}

	hosts["node.aicj.io"] = &Host{node}

	node.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "node")
	})

	// Server
	e := echo.New()
	e.Use(middleware.Recover())
	e.IPExtractor = func(r *http.Request) string {
		return r.Header.Get("Fly-Client-IP")
	}
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogMethod:   true,
		LogStatus:   true,
		LogRemoteIP: true,
		LogLatency:  true,
		LogHost:     true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			slog.Info("request",
				"host", v.Host,
				"URI", v.URI,
				"method", v.Method,
				"status", v.Status,
				"remote_ip", v.RemoteIP,
				"latency", v.Latency,
			)

			return nil
		},
	}))
	e.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		host := hosts[req.Host]

		if host == nil {
			err = echo.ErrNotFound
		} else {
			host.Echo.ServeHTTP(res, req)
		}

		return
	})

	return e
}

func (s *Server) PingHandler(c echo.Context) error {
	ip := c.RealIP()
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		slog.Error("ParseAddr", "error", err)
		return echo.ErrInternalServerError
	}
	return c.String(http.StatusOK, addr.String())
}

func (s *Server) VisitorHandler(c echo.Context) error {
	ip := c.RealIP()

	sqid, err := Sqids()
	if err != nil {
		slog.Error("sqids initialization", "error", err)
		return echo.ErrInternalServerError
	}

	addr, err := netip.ParseAddr(ip)
	if err != nil {
		slog.Error("ParseAddr", "error", err)
		return echo.ErrBadRequest
	}

	visitorF3SiD, err := s.DB.GetVisitor(addr, sqid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			random := rand.Int32N(99)

			visitorF3SiD, err := s.DB.CreateVisitor(addr, random, sqid)
			if err != nil {
				slog.Error("visitor", "error", err)
				return echo.ErrBadRequest
			}

			return c.JSON(http.StatusOK, map[string]string{
				"f3sid": visitorF3SiD,
			})
		}
		slog.Error("visitor", "error", err)
		return echo.ErrBadRequest
	}

	return c.JSON(http.StatusOK, map[string]string{
		"f3sid": visitorF3SiD,
	})
}

type vote struct {
	ModelID      int    `json:"model_id"`
	VisitorF3SiD string `json:"f3sid"`
}

func (s *Server) VoteHandler(c echo.Context) error {
	var vote vote

	err := c.Bind(&vote)
	if err != nil {
		slog.Error("bind", "error", err)
		return echo.ErrBadRequest
	}

	sqid, err := Sqids()
	if err != nil {
		slog.Error("sqids initialization", "error", err)
		return echo.ErrInternalServerError
	}

	voteVisitorID := sqid.Decode(vote.VisitorF3SiD)
	if len(voteVisitorID) != 2 {
		slog.Error("sqids decode", "error", "invalid sqids")
		return echo.ErrBadRequest
	}

	err = s.DB.Vote(int64(vote.ModelID), int64(voteVisitorID[0]), int32(voteVisitorID[1]))
	if err != nil {
		slog.Error("vote", "error", err)
		return echo.ErrBadRequest
	}

	return c.NoContent(http.StatusOK)
}

type pushEntry struct {
	VisitorF3SiD string `json:"f3sid"`
}

func (s *Server) NodePushEntryHandler(c echo.Context) error {
	var push pushEntry

	err := c.Bind(&push)
	if err != nil {
		slog.Error("bind", "error", err)
		return echo.ErrBadRequest
	}

	sqid, err := Sqids()
	if err != nil {
		slog.Error("sqids initialization", "error", err)
		return echo.ErrInternalServerError
	}

	pushVisitorID := sqid.Decode(push.VisitorF3SiD)
	if len(pushVisitorID) != 2 {
		slog.Error("sqids decode", "error", "invalid sqids")
		return echo.ErrBadRequest
	}

	node := c.Get("node").(sqlc.Node)

	err = s.DB.PushEntry(node, int64(pushVisitorID[0]), int32(pushVisitorID[1]))
	if err != nil {
		slog.Error("push entry", "error", err)
		return echo.ErrBadRequest
	}

	return c.NoContent(http.StatusOK)
}

type pushFoodStall struct {
	VisitorF3SiD string  `json:"f3sid"`
	Foods        []foods `json:"foods"`
}

type foods struct {
	ID       int `json:"id"`
	Quantity int `json:"quantity"`
}

func (s *Server) NodePushFoodStallHandler(c echo.Context) error {
	var push pushFoodStall

	err := c.Bind(&push)
	if err != nil {
		slog.Error("bind", "error", err)
		return echo.ErrBadRequest
	}

	sqid, err := Sqids()
	if err != nil {
		slog.Error("sqids initialization", "error", err)
		return echo.ErrInternalServerError
	}

	pushVisitorID := sqid.Decode(push.VisitorF3SiD)
	if len(pushVisitorID) != 2 {
		slog.Error("sqids decode", "error", "invalid sqids")
		return echo.ErrBadRequest
	}

	node := c.Get("node").(sqlc.Node)

	if len(push.Foods) == 0 {
		slog.Error("foods", "error", "no foods")
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

	err = s.DB.PushFoodStall(node, int64(pushVisitorID[0]), int32(pushVisitorID[1]), foods)
	if err != nil {
		slog.Error("push foodstall", "error", err)
		return echo.ErrBadRequest
	}

	return c.NoContent(http.StatusOK)
}

type pushExhibition struct {
	VisitorF3SiD string `json:"f3sid"`
}

func (s *Server) NodePushExhibitionHandler(c echo.Context) error {
	var push pushExhibition

	err := c.Bind(&push)
	if err != nil {
		slog.Error("bind", "error", err)
		return echo.ErrBadRequest
	}

	sqid, err := Sqids()
	if err != nil {
		slog.Error("sqids initialization", "error", err)
		return echo.ErrInternalServerError
	}

	pushVisitorID := sqid.Decode(push.VisitorF3SiD)
	if len(pushVisitorID) != 2 {
		slog.Error("sqids decode", "error", "invalid sqids")
		return echo.ErrBadRequest
	}

	node := c.Get("node").(sqlc.Node)

	err = s.DB.PushExhibition(node, int64(pushVisitorID[0]), int32(pushVisitorID[1]))
	if err != nil {
		slog.Error("push exhibition", "error", err)
		return echo.ErrBadRequest
	}

	return c.NoContent(http.StatusOK)
}

type status struct {
	Charging        bool `json:"charging"`
	ChargingTime    int  `json:"charging_time"`
	DischargingTime int  `json:"discharging_time"`
	Level           int  `json:"level"`
}

func (s *Server) NodeStatusHandler(c echo.Context) error {
	var status status
	err := c.Bind(&status)
	if err != nil {
		slog.Error("bind", "error", err)
		return echo.ErrBadRequest
	}

	node := c.Get("node").(sqlc.Node)

	err = s.DB.StatusNode(node.ID, int32(status.Level), int32(status.ChargingTime), int32(status.DischargingTime), status.Charging)
	if err != nil {
		slog.Error("status", "error", err)
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
		slog.Error("foods", "error", err)
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

type visitorLookup struct {
	visitorF3SiD string `param:"f3sid"`
}

func (s *Server) NodeVisitorLookupHandler(c echo.Context) error {
	var visitorLookup visitorLookup

	err := c.Bind(&visitorLookup)
	if err != nil {
		slog.Error("bind", "error", err)
		return echo.ErrBadRequest
	}

	sqid, err := Sqids()
	if err != nil {
		slog.Error("sqids initialization", "error", err)
		return echo.ErrInternalServerError
	}

	visitorLookupVisitorID := sqid.Decode(visitorLookup.visitorF3SiD)
	if len(visitorLookupVisitorID) != 2 {
		slog.Error("sqids decode", "error", "invalid sqids")
		return echo.ErrBadRequest
	}

	isFirst, err := s.DB.IsVisitorFirst(int64(visitorLookupVisitorID[0]))
	if err != nil {
		slog.Error("is visitor first", "error", err)
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
		slog.Error("sqids initialization", "error", err)
		return echo.ErrInternalServerError
	}

	switch node.Type {
	case sqlc.NodeTypeENTRY:
		entryRow, err := s.DB.EntryRows(node, sqid)
		if err != nil {
			slog.Error("entry row", "error", err)
			return echo.ErrBadRequest
		}

		return c.JSON(http.StatusOK, entryRow)
	case sqlc.NodeTypeFOODSTALL:
		foodstallRawLog, err := s.DB.FoodstallRows(node, sqid)
		if err != nil {
			slog.Error("foodstall row", "error", err)
			return echo.ErrBadRequest
		}

		return c.JSON(http.StatusOK, foodstallRawLog)

	case sqlc.NodeTypeEXHIBITION:
		exhibitionRowLog, err := s.DB.ExhibitionRows(node, sqid)
		if err != nil {
			slog.Error("exhibition row", "error", err)
			return echo.ErrBadRequest
		}

		return c.JSON(http.StatusOK, exhibitionRowLog)
	}

	return echo.ErrBadRequest
}

type updatePush struct {
	ID       int `json:"id"`
	FoodID   int `json:"food_id"`
	Quantity int `json:"quantity"`
}

func (s *Server) NodeUpdateFoodStallHandler(c echo.Context) error {
	var push updatePush

	err := c.Bind(&push)
	if err != nil {
		slog.Error("bind", "error", err)
		return echo.ErrBadRequest
	}

	node := c.Get("node").(sqlc.Node)

	err = s.DB.UpdateFoodLog(node, int64(push.ID), int64(push.FoodID), int32(push.Quantity))
	if err != nil {
		slog.Error("update push", "error", err)
		return echo.ErrBadRequest
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) NodeIpHandler(c echo.Context) error {
	ip := c.RealIP()

	slog.Info("turi")
	slog.Info("ip", "ip", ip)

	addr, err := netip.ParseAddr(ip)
	if err != nil {
		slog.Error("ParseAddr", "error", err)
		return echo.ErrInternalServerError
	}

	node, err := s.DB.IpNode(addr)
	if err != nil {
		slog.Error("ip node", "error", err)
		return echo.ErrBadRequest
	}

	return c.JSON(http.StatusOK, map[string]string{
		"key": node.Key.String,
	})
}

func (s *Server) NodeCountHandler(c echo.Context) error {
	node := c.Get("node").(sqlc.Node)

	switch node.Type {
	case sqlc.NodeTypeENTRY:
		count, err := s.DB.CountEntry(node)
		if err != nil {
			slog.Error("entry row", "error", err)
			return echo.ErrBadRequest
		}

		return c.JSON(http.StatusOK, map[string]int64{"count": count})
	case sqlc.NodeTypeFOODSTALL:
		count, err := s.DB.CountFoodStall(node)
		if err != nil {
			slog.Error("foodstall row", "error", err)
			return echo.ErrBadRequest
		}

		return c.JSON(http.StatusOK, map[string]int64{"count": count})
	case sqlc.NodeTypeEXHIBITION:
		count, err := s.DB.CountExhibition(node)
		if err != nil {
			slog.Error("exhibition row", "error", err)
			return echo.ErrBadRequest
		}

		return c.JSON(http.StatusOK, map[string]int64{"count": count})
	}

	return echo.ErrBadRequest
}

func (s *Server) NodeFoodCountHandler(c echo.Context) error {
	node := c.Get("node").(sqlc.Node)

	foodCount, err := s.DB.CountFood(node)
	if err != nil {
		slog.Error("foodstall row", "error", err)
		return echo.ErrBadRequest
	}

	foodsArray := make([]map[string]interface{}, len(foodCount))
	for i, food := range foodCount {
		foodsArray[i] = map[string]interface{}{
			"id":    food.ID,
			"name":  food.Name,
			"count": food.Count,
		}
	}

	return c.JSON(http.StatusOK, foodsArray)
}

func (s *Server) NodeEntryTypeCountHandler(c echo.Context) error {
	node := c.Get("node").(sqlc.Node)

	exhibitionCount, err := s.DB.CountEntryType(node)
	if err != nil {
		slog.Error("exhibition row", "error", err)
		return echo.ErrBadRequest
	}

	exhibitions := make([]map[string]interface{}, len(exhibitionCount))
	for i, exhibition := range exhibitionCount {
		exhibitions[i] = map[string]interface{}{
			"type":  exhibition.Type,
			"count": exhibition.Count,
		}
	}

	return c.JSON(http.StatusOK, exhibitions)
}

func (s *Server) NodeEntryPerHourCountHandler(c echo.Context) error {
	node := c.Get("node").(sqlc.Node)

	exhibitionCount, err := s.DB.CountEntryPerHour(node)
	if err != nil {
		slog.Error("exhibition row", "error", err)
		return echo.ErrBadRequest
	}

	return c.JSON(http.StatusOK, exhibitionCount)
}

func (s *Server) NodeFoodStallPerHourCountHandler(c echo.Context) error {
	node := c.Get("node").(sqlc.Node)

	exhibitionCount, err := s.DB.CountFoodStallPerHour(node)
	if err != nil {
		slog.Error("exhibition row", "error", err)
		return echo.ErrBadRequest
	}

	return c.JSON(http.StatusOK, exhibitionCount)
}

func (s *Server) NodeExhibitionPerHourCountHandler(c echo.Context) error {
	node := c.Get("node").(sqlc.Node)

	exhibitionCount, err := s.DB.CountExhibitionPerHour(node)
	if err != nil {
		slog.Error("exhibition row", "error", err)
		return echo.ErrBadRequest
	}

	return c.JSON(http.StatusOK, exhibitionCount)
}
