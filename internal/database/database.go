package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"backend/internal/sqlc"

	"github.com/sqids/sqids-go"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

const (
	ENTRY      = "ENTRY"
	FOODSTALL  = "FOODSTALL"
	EXHIBITION = "EXHIBITION"

	ENTERED = "ENTERED"
	LEFT    = "LEFT"
)

// Service represents a DbService that interacts with a database.
type Service interface {
	Password(key string) (sqlc.Node, bool, error)

	GetVisitor(ip string, sqid *sqids.Sqids) (string, error)

	CreateVisitor(ip string, rand int64, sqid *sqids.Sqids) (string, error)

	IpNode(ip string) (sqlc.Node, error)

	PushEntry(node sqlc.Node, visitorID int64, visitorRandom int64) error

	PushFoodStall(node sqlc.Node, visitorID int64, visitorRandom int64, foods []Foods) error

	PushExhibition(node sqlc.Node, visitorID int64, visitorRandom int64) error

	UpdatePushNode(node sqlc.Node, id int64, quantity int64) error

	StatusNode(nodeID int64, level int64, chargingTime int64, dischargingTime int64, charging bool) error

	IsVisitorFirst(visitorID int64) (bool, error)

	EntryRow(node sqlc.Node, sqid *sqids.Sqids) ([]EntryRowLog, error)

	FoodstallRow(node sqlc.Node, sqid *sqids.Sqids) ([]FoodstallRowLog, error)

	ExhibitionRow(node sqlc.Node, sqid *sqids.Sqids) ([]ExhibitionRowLog, error)

	Foods(node sqlc.Node) ([]NodeFood, error)

	CountEntry(node sqlc.Node) (int64, error)

	CountFoodStall(node sqlc.Node) (int64, error)

	CountExhibition(node sqlc.Node) (int64, error)
}

type DbService struct {
	DB *sql.DB
}

var (
	dsn        = os.Getenv("DSN")
	dbInstance *DbService
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_fk=true&_journal=WAL", dsn))
	if err != nil {
		slog.Default().Error("database config parse error", "error", err)
	}
	if db.Ping() != nil {
		slog.Default().Error("database connection error", "error", err)
	}
	dbInstance = &DbService{
		DB: db,
	}
	return dbInstance
}

func (s *DbService) Password(key string) (sqlc.Node, bool, error) {
	ctx := context.Background()

	q, err := s.DB.BeginTx(ctx, nil)
	defer q.Rollback()
	if err != nil {
		return sqlc.Node{}, false, err
	}
	queries := sqlc.New(q)

	nodeByKey, err := queries.GetNodeByKey(ctx, sql.NullString{
		String: key,
		Valid:  true,
	})
	if err != nil {
		return sqlc.Node{}, false, err
	}

	return nodeByKey, true, nil
}

func (s *DbService) GetVisitor(ip string, sqid *sqids.Sqids) (string, error) {
	ctx := context.Background()

	q, err := s.DB.BeginTx(ctx, nil)
	defer q.Rollback()
	if err != nil {
		return "", err
	}
	queries := sqlc.New(q)

	visitorByIp, err := queries.GetVisitorByIp(ctx, ip)
	if err != nil {
		return "", err
	}

	visitorF3SiD, err := sqid.Encode([]uint64{uint64(visitorByIp.ID), uint64(visitorByIp.Random)})
	if err != nil {
		return "", err
	}

	return visitorF3SiD, nil
}

func (s *DbService) CreateVisitor(ip string, rand int64, sqid *sqids.Sqids) (string, error) {
	ctx := context.Background()

	q, err := s.DB.BeginTx(ctx, nil)
	defer q.Rollback()
	if err != nil {
		return "", err
	}
	queries := sqlc.New(q)

	// Create a new visitor
	visitorByIp, err := queries.CreateVisitor(ctx, sqlc.CreateVisitorParams{
		Ip:     ip,
		Random: rand,
	})
	if err != nil {
		return "", err
	}

	err = q.Commit()
	if err != nil {
		return "", err
	}

	visitorF3SiD, err := sqid.Encode([]uint64{uint64(visitorByIp.ID), uint64(visitorByIp.Random)})
	if err != nil {
		return "", err
	}

	return visitorF3SiD, nil
}

func (s *DbService) PushEntry(node sqlc.Node, visitorID int64, visitorRandom int64) error {
	ctx := context.Background()

	q, err := s.DB.BeginTx(ctx, nil)
	defer q.Rollback()
	if err != nil {
		return err
	}
	queries := sqlc.New(q)

	// Check if visitor exists
	visitorById, err := queries.GetVisitorByIdAndRandom(ctx, sqlc.GetVisitorByIdAndRandomParams{
		ID:     visitorID,
		Random: visitorRandom,
	})
	if err != nil {
		return err
	}

	entryLog, err := queries.GetEntryLogByVisitorId(ctx, sql.NullInt64{
		Int64: visitorById.ID,
		Valid: true,
	})

	var entryLogType string

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// First time entering
			entryLogType = ENTERED
		} else {
			return err
		}
	} else {
		if entryLog.Type == ENTERED {
			entryLogType = LEFT
		} else {
			entryLogType = ENTERED
		}
	}

	err = queries.CreateEntryLog(ctx, sqlc.CreateEntryLogParams{
		NodeID: sql.NullInt64{
			Int64: node.ID,
			Valid: true,
		},
		VisitorID: sql.NullInt64{
			Int64: visitorById.ID,
			Valid: true,
		},
		Type: entryLogType,
	})
	if err != nil {
		return err
	}

	err = q.Commit()
	if err != nil {
		return err
	}

	return nil
}

type Foods struct {
	ID       int
	Quantity int
}

func (s *DbService) PushFoodStall(node sqlc.Node, visitorID int64, visitorRandom int64, foods []Foods) error {
	ctx := context.Background()

	q, err := s.DB.BeginTx(ctx, nil)
	defer q.Rollback()
	if err != nil {
		return err
	}
	queries := sqlc.New(q)

	// Check if visitor exists
	visitorById, err := queries.GetVisitorByIdAndRandom(ctx, sqlc.GetVisitorByIdAndRandomParams{
		ID:     visitorID,
		Random: visitorRandom,
	})
	if err != nil {
		return err
	}

	foodMap := make(map[int]int64)
	// Check each food in foods
	for _, food := range foods {
		foodById, err := queries.GetFoodById(ctx, int64(food.ID))
		if err != nil {
			return err
		}
		foodMap[food.ID] = foodById.ID
	}

	for _, food := range foods {
		foodById, exists := foodMap[food.ID]
		if !exists {
			return fmt.Errorf("food with id %d not found in map", food.ID)
		}

		err = queries.CreateFoodStallLog(ctx, sqlc.CreateFoodStallLogParams{
			NodeID:    sql.NullInt64{Int64: node.ID, Valid: true},
			VisitorID: sql.NullInt64{Int64: visitorById.ID, Valid: true},
			FoodID:    sql.NullInt64{Int64: foodById, Valid: true},
			Quantity:  int64(food.Quantity),
		})
		if err != nil {
			return err
		}
	}

	err = q.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *DbService) PushExhibition(node sqlc.Node, visitorID int64, visitorRandom int64) error {
	ctx := context.Background()

	q, err := s.DB.BeginTx(ctx, nil)
	defer q.Rollback()
	if err != nil {
		return err
	}
	queries := sqlc.New(q)

	// Check if visitor exists
	visitorById, err := queries.GetVisitorByIdAndRandom(ctx, sqlc.GetVisitorByIdAndRandomParams{
		ID:     visitorID,
		Random: visitorRandom,
	})
	if err != nil {
		return err
	}

	err = queries.CreateExhibitionLog(ctx, sqlc.CreateExhibitionLogParams{
		NodeID:    sql.NullInt64{Int64: node.ID, Valid: true},
		VisitorID: sql.NullInt64{Int64: visitorById.ID, Valid: true},
	})
	if err != nil {
		return err
	}

	err = q.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *DbService) UpdatePushNode(node sqlc.Node, id int64, quantity int64) error {
	ctx := context.Background()

	q, err := s.DB.BeginTx(ctx, nil)
	defer q.Rollback()
	if err != nil {
		return err
	}
	queries := sqlc.New(q)

	if node.Type == FOODSTALL {
		err = queries.UpdateFoodStallLog(ctx, sqlc.UpdateFoodStallLogParams{
			Quantity: quantity,
			ID:       id,
		})
		if err != nil {
			return err
		}

		err = q.Commit()
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

func (s *DbService) StatusNode(nodeID int64, level int64, chargingTime int64, dischargingTime int64, charging bool) error {
	ctx := context.Background()

	q, err := s.DB.BeginTx(ctx, nil)
	defer q.Rollback()
	if err != nil {
		return err
	}
	queries := sqlc.New(q)

	err = queries.UpdateBattery(ctx, sqlc.UpdateBatteryParams{
		NodeID:          sql.NullInt64{Int64: nodeID, Valid: true},
		Level:           sql.NullInt64{Int64: level, Valid: true},
		ChargingTime:    sql.NullInt64{Int64: chargingTime, Valid: true},
		DischargingTime: sql.NullInt64{Int64: dischargingTime, Valid: true},
		Charging:        sql.NullBool{Bool: charging, Valid: true},
	})
	if err != nil {
		return err
	}

	err = q.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *DbService) IsVisitorFirst(visitorID int64) (bool, error) {
	ctx := context.Background()

	q, err := s.DB.BeginTx(ctx, nil)
	defer q.Rollback()
	if err != nil {
		return false, err
	}
	queries := sqlc.New(q)

	// Check if visitor exists
	_, err = queries.GetEntryLogByVisitorId(ctx,
		sql.NullInt64{
			Int64: visitorID,
			Valid: true,
		})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}
		return false, err
	}

	return false, err
}

type EntryRowLog struct {
	Id        int64     `json:"id"`
	F3SiD     string    `json:"f3sid"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type FoodstallRowLog struct {
	Id        int64     `json:"id"`
	F3SiD     string    `json:"f3sid"`
	FoodName  string    `json:"name"`
	Quantity  int64     `json:"quantity"`
	Price     int64     `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

type ExhibitionRowLog struct {
	Id        int64     `json:"id"`
	F3SiD     string    `json:"f3sid"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *DbService) EntryRow(node sqlc.Node, sqid *sqids.Sqids) ([]EntryRowLog, error) {
	ctx := context.Background()

	q, err := s.DB.BeginTx(ctx, nil)
	defer q.Rollback()
	if err != nil {
		return nil, err
	}
	queries := sqlc.New(q)

	rows, err := queries.GetEntryLogByNodeId(ctx, sql.NullInt64{Int64: node.ID, Valid: true})
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

func (s *DbService) FoodstallRow(node sqlc.Node, sqid *sqids.Sqids) ([]FoodstallRowLog, error) {
	ctx := context.Background()

	q, err := s.DB.BeginTx(ctx, nil)
	defer q.Rollback()
	if err != nil {
		return nil, err
	}
	queries := sqlc.New(q)

	rows, err := queries.GetFoodStallLogByNodeId(ctx, sql.NullInt64{Int64: node.ID, Valid: true})
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

func (s *DbService) ExhibitionRow(node sqlc.Node, sqid *sqids.Sqids) ([]ExhibitionRowLog, error) {
	ctx := context.Background()

	q, err := s.DB.BeginTx(ctx, nil)
	defer q.Rollback()
	if err != nil {
		return nil, err
	}
	queries := sqlc.New(q)

	rows, err := queries.GetExhibitionLogByNodeId(ctx, sql.NullInt64{Int64: node.ID, Valid: true})
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

func (s *DbService) IpNode(ip string) (sqlc.Node, error) {
	ctx := context.Background()

	q, err := s.DB.BeginTx(ctx, nil)
	defer q.Rollback()
	if err != nil {
		return sqlc.Node{}, err
	}
	queries := sqlc.New(q)

	nodeByIp, err := queries.GetNodeByIp(ctx, sql.NullString{String: ip, Valid: true})
	if err != nil {
		return sqlc.Node{}, err
	}

	err = queries.DeleteNodeIp(ctx, sql.NullString{String: ip, Valid: true})
	if err != nil {
		return sqlc.Node{}, err
	}

	err = q.Commit()
	if err != nil {
		return sqlc.Node{}, err
	}

	return nodeByIp, nil
}

type NodeFood struct {
	ID    int
	Name  string
	Price int
}

func (s *DbService) Foods(node sqlc.Node) ([]NodeFood, error) {
	ctx := context.Background()

	q, err := s.DB.BeginTx(ctx, nil)
	defer q.Rollback()
	if err != nil {
		return nil, err
	}
	queries := sqlc.New(q)

	foods, err := queries.GetFoodsByNodeId(ctx, sql.NullInt64{
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

func (s *DbService) CountEntry(node sqlc.Node) (int64, error) {
	ctx := context.Background()

	q, err := s.DB.BeginTx(ctx, nil)
	defer q.Rollback()
	if err != nil {
		return 0, err
	}
	queries := sqlc.New(q)

	count, err := queries.CountEntryLogByNodeId(ctx, sql.NullInt64{
		Int64: node.ID,
		Valid: true,
	})
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *DbService) CountFoodStall(node sqlc.Node) (int64, error) {
	ctx := context.Background()

	q, err := s.DB.BeginTx(ctx, nil)
	defer q.Rollback()
	if err != nil {
		return 0, err
	}
	queries := sqlc.New(q)

	count, err := queries.CountFoodStallLogByNodeId(ctx, sql.NullInt64{
		Int64: node.ID,
		Valid: true,
	})
	if err != nil {
		return 0, err
	}
	if !count.Valid {
		return 0, nil
	}

	return int64(count.Float64), nil
}

func (s *DbService) CountExhibition(node sqlc.Node) (int64, error) {
	ctx := context.Background()

	q, err := s.DB.BeginTx(ctx, nil)
	defer q.Rollback()
	if err != nil {
		return 0, err
	}
	queries := sqlc.New(q)

	count, err := queries.CountExhibitionLogByNodeId(ctx, sql.NullInt64{
		Int64: node.ID,
		Valid: true,
	})
	if err != nil {
		return 0, err
	}

	return count, nil
}
