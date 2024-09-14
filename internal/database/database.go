package database

import (
	"context"
	gosql "database/sql"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"os"
	"time"

	sql "backend/internal/sqlc"

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
	// Health returns a map of health status information.
	// The keys and values in the map are DbService-specific.
	Health() map[string]string

	Password(key string) (sql.Node, bool, error)

	Visitor(ip string, sqid *sqids.Sqids) (string, error)

	IpNode(ip string) (sql.Node, error)

	PushEntry(node sql.Node, visitorID int64, visitorRandom int64) error

	PushFoodStall(node sql.Node, visitorID int64, visitorRandom int64, foods []Foods) error

	PushExhibition(node sql.Node, visitorID int64, visitorRandom int64) error

	UpdatePushNode(node sql.Node, id int64, quantity int64) error

	StatusNode(nodeID int64, level int64, chargingTime int64, dischargingTime int64, charging bool) error

	IsVisitorFirst(visitorID int64) (bool, error)

	EntryRow(node sql.Node, sqid *sqids.Sqids) ([]EntryRowLog, error)

	FoodstallRow(node sql.Node, sqid *sqids.Sqids) ([]FoodstallRowLog, error)

	ExhibitionRow(node sql.Node, sqid *sqids.Sqids) ([]ExhibitionRowLog, error)

	Foods(node sql.Node) ([]NodeFood, error)
}

type DbService struct {
	DB *gosql.DB
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
	db, err := gosql.Open("sqlite3", dsn)
	if err != nil {
		slog.Default().Error("database config parse error", "error", err)
	}
	dbInstance = &DbService{
		DB: db,
	}
	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *DbService) Health() map[string]string {
	_, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.DB.Ping()
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

	q, err := s.DB.BeginTx(ctx, nil)
	defer q.Rollback()
	if err != nil {
		return sql.Node{}, false, err
	}
	queries := sql.New(q)

	nodeByKey, err := queries.GetNodeByKey(ctx, gosql.NullString{
		String: key,
		Valid:  true,
	})
	if err != nil {
		return sql.Node{}, false, err
	}

	return nodeByKey, true, nil
}

func (s *DbService) Visitor(ip string, sqid *sqids.Sqids) (string, error) {
	ctx := context.Background()

	q, err := s.DB.BeginTx(ctx, nil)
	defer q.Rollback()
	if err != nil {
		return "", err
	}
	queries := sql.New(q)

	// Check if visitor exists with ip
	visitorByIp, err := queries.GetVisitorByIp(ctx, ip)
	if err != nil {
		if errors.Is(err, gosql.ErrNoRows) {
			random := rand.Int32()

			// Create a new visitor
			visitorByIp, err := queries.CreateVisitor(ctx, sql.CreateVisitorParams{
				Ip:     ip,
				Random: int64(random),
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

func (s *DbService) PushEntry(node sql.Node, visitorID int64, visitorRandom int64) error {
	ctx := context.Background()

	q, err := s.DB.BeginTx(ctx, nil)
	defer q.Rollback()
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

	entryLog, err := queries.GetEntryLogByVisitorId(ctx, gosql.NullInt64{
		Int64: visitorById.ID,
		Valid: true,
	})

	var entryLogType string

	if err != nil {
		if errors.Is(err, gosql.ErrNoRows) {
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

	err = queries.CreateEntryLog(ctx, sql.CreateEntryLogParams{
		NodeID: gosql.NullInt64{
			Int64: node.ID,
			Valid: true,
		},
		VisitorID: gosql.NullInt64{
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

func (s *DbService) PushFoodStall(node sql.Node, visitorID int64, visitorRandom int64, foods []Foods) error {
	ctx := context.Background()

	q, err := s.DB.BeginTx(ctx, nil)
	defer q.Rollback()
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

		err = queries.CreateFoodStallLog(ctx, sql.CreateFoodStallLogParams{
			NodeID:    gosql.NullInt64{Int64: node.ID, Valid: true},
			VisitorID: gosql.NullInt64{Int64: visitorById.ID, Valid: true},
			FoodID:    gosql.NullInt64{Int64: foodById, Valid: true},
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

func (s *DbService) PushExhibition(node sql.Node, visitorID int64, visitorRandom int64) error {
	ctx := context.Background()

	q, err := s.DB.BeginTx(ctx, nil)
	defer q.Rollback()
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
		NodeID:    gosql.NullInt64{Int64: node.ID, Valid: true},
		VisitorID: gosql.NullInt64{Int64: visitorById.ID, Valid: true},
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

func (s *DbService) UpdatePushNode(node sql.Node, id int64, quantity int64) error {
	ctx := context.Background()

	q, err := s.DB.BeginTx(ctx, nil)
	defer q.Rollback()
	if err != nil {
		return err
	}
	queries := sql.New(q)

	if node.Type == FOODSTALL {
		err = queries.UpdateFoodStallLog(ctx, sql.UpdateFoodStallLogParams{
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
	queries := sql.New(q)

	err = queries.UpdateBattery(ctx, sql.UpdateBatteryParams{
		NodeID:          gosql.NullInt64{Int64: nodeID, Valid: true},
		Level:           gosql.NullInt64{Int64: level, Valid: true},
		ChargingTime:    gosql.NullInt64{Int64: chargingTime, Valid: true},
		DischargingTime: gosql.NullInt64{Int64: dischargingTime, Valid: true},
		Charging:        gosql.NullBool{Bool: charging, Valid: true},
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
	queries := sql.New(q)

	// Check if visitor exists
	_, err = queries.GetEntryLogByVisitorId(ctx,
		gosql.NullInt64{
			Int64: visitorID,
			Valid: true,
		})
	if err != nil {
		if errors.Is(err, gosql.ErrNoRows) {
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

func (s *DbService) EntryRow(node sql.Node, sqid *sqids.Sqids) ([]EntryRowLog, error) {
	ctx := context.Background()

	q, err := s.DB.BeginTx(ctx, nil)
	defer q.Rollback()
	if err != nil {
		return nil, err
	}
	queries := sql.New(q)

	rows, err := queries.GetEntryLogByNodeId(ctx, gosql.NullInt64{Int64: node.ID, Valid: true})
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

	q, err := s.DB.BeginTx(ctx, nil)
	defer q.Rollback()
	if err != nil {
		return nil, err
	}
	queries := sql.New(q)

	rows, err := queries.GetFoodStallLogByNodeId(ctx, gosql.NullInt64{Int64: node.ID, Valid: true})
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

	q, err := s.DB.BeginTx(ctx, nil)
	defer q.Rollback()
	if err != nil {
		return nil, err
	}
	queries := sql.New(q)

	rows, err := queries.GetExhibitionLogByNodeId(ctx, gosql.NullInt64{Int64: node.ID, Valid: true})
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

func (s *DbService) IpNode(ip string) (sql.Node, error) {
	ctx := context.Background()

	q, err := s.DB.BeginTx(ctx, nil)
	defer q.Rollback()
	if err != nil {
		return sql.Node{}, err
	}
	queries := sql.New(q)

	nodeByIp, err := queries.GetNodeByIp(ctx, gosql.NullString{String: ip, Valid: true})
	if err != nil {
		return sql.Node{}, err
	}

	err = queries.DeleteNodeIp(ctx, gosql.NullString{String: ip, Valid: true})
	if err != nil {
		return sql.Node{}, err
	}

	err = q.Commit()
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

	q, err := s.DB.BeginTx(ctx, nil)
	defer q.Rollback()
	if err != nil {
		return nil, err
	}
	queries := sql.New(q)

	foods, err := queries.GetFoodsByNodeId(ctx, gosql.NullInt64{
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
