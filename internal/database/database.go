package database

import (
	sql "backend/internal/sqlc"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net/netip"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sqids/sqids-go"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

// Service represents a DbService that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are DbService-specific.
	Health() map[string]string

	Password(key string) (sql.Node, bool, error)

	Visitor(ip netip.Addr, sqid *sqids.Sqids) (string, error)

	IpNode(ip netip.Addr) (sql.Node, error)

	PushEntry(node sql.Node, visitorID int64, visitorRandom int32) error

	PushFoodStall(node sql.Node, visitorID int64, visitorRandom int32, foods []Foods) error

	PushExhibition(node sql.Node, visitorID int64, visitorRandom int32) error

	UpdatePushNode(node sql.Node, id int64, quantity int32) error

	StatusNode(nodeID int64, level int32, chargingTime int32, dischargingTime int32, charging bool) error

	IsVisitorFirst(visitorID int64) (bool, error)

	EntryRow(node sql.Node, sqid *sqids.Sqids) ([]EntryRowLog, error)

	FoodstallRow(node sql.Node, sqid *sqids.Sqids) ([]FoodstallRowLog, error)

	ExhibitionRow(node sql.Node, sqid *sqids.Sqids) ([]ExhibitionRowLog, error)

	Foods(node sql.Node) ([]NodeFood, error)
}

type DbService struct {
	DB *pgxpool.Pool
}

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	schema     = os.Getenv("DB_SCHEMA")
	dbInstance *DbService
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	config, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema))
	if err != nil {
		slog.Default().Error("database config parse error", "error", err)
	}
	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		slog.Default().Error("database connection failed", "error", err)
	}
	//defer db.Close()
	dbInstance = &DbService{
		DB: db,
	}
	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *DbService) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.DB.Ping(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		slog.Default().Error("database down", "error", err)
		return stats
	}

	return stats
}

func (s *DbService) Password(key string) (sql.Node, bool, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		if q != nil {
			_ = q.Rollback(ctx)
		}
	}(q, ctx)
	if err != nil {
		return sql.Node{}, false, err
	}
	queries := sql.New(q)

	nodeByKey, err := queries.GetNodeByKey(ctx, pgtype.Text{
		String: key,
		Valid:  true,
	})
	if err != nil {
		return sql.Node{}, false, err
	}

	return nodeByKey, true, nil
}

func (s *DbService) Visitor(ip netip.Addr, sqid *sqids.Sqids) (string, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		if q != nil {
			_ = q.Rollback(ctx)
		}
	}(q, ctx)
	if err != nil {
		return "", err
	}
	queries := sql.New(q)

	// Check if visitor exists with ip
	visitorByIp, err := queries.GetVisitorByIp(ctx, ip)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			rand := rand.Int32()

			// Create a new visitor
			visitorByIp, err := queries.CreateVisitor(ctx, sql.CreateVisitorParams{
				Ip:     ip,
				Random: rand,
			})
			if err != nil {
				return "", err
			}
			err = q.Commit(ctx)
			if err != nil {
				return "", err
			}

			visitorF3SiD, err := sqid.Encode([]uint64{uint64(visitorByIp.ID), uint64(visitorByIp.Random)})
			if err != nil {
				return "", err
			}

			return visitorF3SiD, nil
		} else {
			return "", err
		}
	}

	visitorF3SiD, err := sqid.Encode([]uint64{uint64(visitorByIp.ID), uint64(visitorByIp.Random)})
	if err != nil {
		return "", err
	}

	return visitorF3SiD, nil
}

func (s *DbService) PushEntry(node sql.Node, visitorID int64, visitorRandom int32) error {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return err
	}
	queries := sql.New(q)

	// Check if visitor exists
	visitorById, err := queries.GetVisitorByIdAndRandom(ctx, sql.GetVisitorByIdAndRandomParams{
		ID:     visitorID,
		Random: visitorRandom,
	})
	if err != nil {
		return err
	}

	entryLog, err := queries.GetEntryLogByVisitorId(ctx, pgtype.Int8{Int64: visitorById.ID, Valid: true})

	var entryLogType sql.EntryLogsType

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// First time entering
			entryLogType = sql.EntryLogsTypeENTERED
		} else {
			return err
		}
	} else {
		if entryLog.Type == sql.EntryLogsTypeENTERED {
			entryLogType = sql.EntryLogsTypeLEFT
		} else {
			entryLogType = sql.EntryLogsTypeENTERED
		}
	}

	err = queries.CreateEntryLog(ctx, sql.CreateEntryLogParams{
		NodeID: pgtype.Int8{
			Int64: node.ID,
			Valid: true,
		},
		VisitorID: pgtype.Int8{
			Int64: visitorById.ID,
			Valid: true,
		},
		Type: entryLogType,
	})
	if err != nil {
		return err
	}

	err = q.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

type Foods struct {
	ID       int
	Quantity int
}

func (s *DbService) PushFoodStall(node sql.Node, visitorID int64, visitorRandom int32, foods []Foods) error {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return err
	}
	queries := sql.New(q)

	// Check if visitor exists
	visitorById, err := queries.GetVisitorByIdAndRandom(ctx, sql.GetVisitorByIdAndRandomParams{
		ID:     visitorID,
		Random: visitorRandom,
	})
	if err != nil {
		return err
	}

	foodMap := make(map[int]int64)
	// Check each food in foods
	for _, food := range foods {
		foodById, err := queries.GetFoodIdById(ctx, int64(food.ID))
		if err != nil {
			return err
		}
		foodMap[food.ID] = foodById
	}

	for _, food := range foods {
		foodById, exists := foodMap[food.ID]
		if !exists {
			return fmt.Errorf("food with id %d not found in map", food.ID)
		}

		err = queries.CreateFoodStallLog(ctx, sql.CreateFoodStallLogParams{
			NodeID:    pgtype.Int8{Int64: node.ID, Valid: true},
			VisitorID: pgtype.Int8{Int64: visitorById.ID, Valid: true},
			FoodID:    pgtype.Int8{Int64: foodById, Valid: true},
			Quantity:  int32(food.Quantity),
		})
		if err != nil {
			return err
		}
	}

	err = q.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *DbService) PushExhibition(node sql.Node, visitorID int64, visitorRandom int32) error {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return err
	}
	queries := sql.New(q)

	// Check if visitor exists
	visitorById, err := queries.GetVisitorByIdAndRandom(ctx, sql.GetVisitorByIdAndRandomParams{
		ID:     visitorID,
		Random: visitorRandom,
	})
	if err != nil {
		return err
	}

	err = queries.CreateExhibitionLog(ctx, sql.CreateExhibitionLogParams{
		NodeID:    pgtype.Int8{Int64: node.ID, Valid: true},
		VisitorID: pgtype.Int8{Int64: visitorById.ID, Valid: true},
	})
	if err != nil {
		return err
	}

	err = q.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *DbService) UpdatePushNode(node sql.Node, id int64, quantity int32) error {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return err
	}
	queries := sql.New(q)

	if node.Type == sql.NodeTypeFOODSTALL {
		err = queries.UpdateFoodStallLog(ctx, sql.UpdateFoodStallLogParams{
			Quantity: quantity,
			ID:       id,
		})
		if err != nil {
			return err
		}

		err = q.Commit(ctx)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

func (s *DbService) StatusNode(nodeID int64, level int32, chargingTime int32, dischargingTime int32, charging bool) error {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return err
	}
	queries := sql.New(q)

	err = queries.UpdateBattery(ctx, sql.UpdateBatteryParams{
		NodeID:          pgtype.Int8{Int64: nodeID, Valid: true},
		Level:           pgtype.Int4{Int32: level, Valid: true},
		ChargingTime:    pgtype.Int4{Int32: chargingTime, Valid: true},
		DischargingTime: pgtype.Int4{Int32: dischargingTime, Valid: true},
		Charging:        pgtype.Bool{Bool: charging, Valid: true},
	})
	if err != nil {
		return err
	}

	err = q.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *DbService) IsVisitorFirst(visitorID int64) (bool, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return false, err
	}
	queries := sql.New(q)

	// Check if visitor exists
	_, err = queries.GetEntryLogByVisitorId(ctx, pgtype.Int8{Int64: visitorID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return true, nil
		}
		return false, err
	}

	return false, err
}

type EntryRowLog struct {
	Id        int64             `json:"id"`
	F3SiD     string            `json:"f3sid"`
	Type      sql.EntryLogsType `json:"type"`
	CreatedAt time.Time         `json:"created_at"`
}

type FoodstallRowLog struct {
	Id        int64     `json:"id"`
	F3SiD     string    `json:"f3sid"`
	FoodName  string    `json:"name"`
	Quantity  int32     `json:"quantity"`
	Price     int32     `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

type ExhibitionRowLog struct {
	Id        int64     `json:"id"`
	F3SiD     string    `json:"f3sid"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *DbService) EntryRow(node sql.Node, sqid *sqids.Sqids) ([]EntryRowLog, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return nil, err
	}
	queries := sql.New(q)

	rows, err := queries.GetEntryLogByNodeId(ctx, pgtype.Int8{Int64: node.ID, Valid: true})
	if err != nil {
		return nil, err
	}

	var rowLog []EntryRowLog

	for _, row := range rows {
		visitorByID, err := queries.GetVisitorById(ctx, row.VisitorID.Int64)
		if err != nil {
			return nil, err
		}

		visitorF3SiD, err := sqid.Encode([]uint64{uint64(visitorByID.ID), uint64(visitorByID.Random)})
		if err != nil {
			return nil, err
		}

		rowLog = append(rowLog, EntryRowLog{row.ID, visitorF3SiD, row.Type, row.CreatedAt.Time})
	}

	return rowLog, nil
}

func (s *DbService) FoodstallRow(node sql.Node, sqid *sqids.Sqids) ([]FoodstallRowLog, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return nil, err
	}
	queries := sql.New(q)

	rows, err := queries.GetFoodStallLogByNodeId(ctx, pgtype.Int8{Int64: node.ID, Valid: true})
	if err != nil {
		return nil, err
	}

	var rowLog []FoodstallRowLog

	for _, row := range rows {
		visitorByID, err := queries.GetVisitorById(ctx, row.VisitorID.Int64)
		if err != nil {
			return nil, err
		}

		visitorF3SiD, err := sqid.Encode([]uint64{uint64(visitorByID.ID), uint64(visitorByID.Random)})
		if err != nil {
			return nil, err
		}

		foodByID, err := queries.GetFoodById(ctx, row.FoodID.Int64)
		if err != nil {
			return nil, err
		}

		rowLog = append(rowLog, FoodstallRowLog{row.ID, visitorF3SiD, foodByID.Name, row.Quantity, foodByID.Price, row.CreatedAt.Time})
	}

	return rowLog, nil
}

func (s *DbService) ExhibitionRow(node sql.Node, sqid *sqids.Sqids) ([]ExhibitionRowLog, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return nil, err
	}
	queries := sql.New(q)

	rows, err := queries.GetExhibitionLogByNodeId(ctx, pgtype.Int8{Int64: node.ID, Valid: true})
	if err != nil {
		return nil, err
	}

	var rowLog []ExhibitionRowLog

	for _, row := range rows {
		visitorByID, err := queries.GetVisitorById(ctx, row.VisitorID.Int64)
		if err != nil {
			return nil, err
		}

		visitorF3SiD, err := sqid.Encode([]uint64{uint64(visitorByID.ID), uint64(visitorByID.Random)})
		if err != nil {
			return nil, err
		}

		rowLog = append(rowLog, ExhibitionRowLog{row.ID, visitorF3SiD, row.CreatedAt.Time})
	}

	return rowLog, nil
}

func (s *DbService) IpNode(ip netip.Addr) (sql.Node, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return sql.Node{}, err
	}
	queries := sql.New(q)

	nodeByIp, err := queries.GetNodeByIp(ctx, &ip)
	if err != nil {
		return sql.Node{}, err
	}

	err = queries.DeleteNodeIp(ctx, &ip)
	if err != nil {
		return sql.Node{}, err
	}

	err = q.Commit(ctx)
	if err != nil {
		return sql.Node{}, err
	}

	return nodeByIp, nil
}

type NodeFood struct {
	ID    int
	Name  string
	Price int
}

func (s *DbService) Foods(node sql.Node) ([]NodeFood, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return nil, err
	}
	queries := sql.New(q)

	foods, err := queries.GetFoodsByNodeId(ctx, pgtype.Int8{
		Int64: node.ID,
		Valid: true,
	})
	if err != nil {
		return nil, err
	}

	var foodsList []NodeFood

	for _, food := range foods {
		foodsList = append(foodsList, NodeFood{int(food.ID), food.Name, int(food.Price)})
	}

	return foodsList, nil
}
