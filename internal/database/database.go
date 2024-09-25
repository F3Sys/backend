package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/netip"
	"os"
	"time"

	"backend/internal/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sqids/sqids-go"

	_ "github.com/joho/godotenv/autoload"
)

// Service represents a DbService that interacts with a database.
type Service interface {
	Password(key string) (sqlc.Node, bool, error)

	Vote(modelID int64, visitorID int64, visitorRandom int32) error

	GetVisitor(ip netip.Addr, sqid *sqids.Sqids) (string, error)

	CreateVisitor(ip netip.Addr, rand int32, sqid *sqids.Sqids) (string, error)

	IpNode(ip netip.Addr) (sqlc.Node, error)

	PushEntry(node sqlc.Node, visitorID int64, visitorRandom int32) error

	PushFoodStall(node sqlc.Node, visitorID int64, visitorRandom int32, foods []Foods) error

	PushExhibition(node sqlc.Node, visitorID int64, visitorRandom int32) error

	UpdatePushNode(node sqlc.Node, id int64, foodID int64, quantity int32) error

	StatusNode(nodeID int64, level int32, chargingTime int32, dischargingTime int32, charging bool) error

	IsVisitorFirst(visitorID int64) (bool, error)

	EntryRow(node sqlc.Node, sqid *sqids.Sqids) ([]EntryRowLog, error)

	FoodstallRow(node sqlc.Node, sqid *sqids.Sqids) ([]FoodstallRowLog, error)

	ExhibitionRow(node sqlc.Node, sqid *sqids.Sqids) ([]ExhibitionRowLog, error)

	Foods(node sqlc.Node) ([]NodeFood, error)

	CountEntry(node sqlc.Node) (int64, error)

	CountFoodStall(node sqlc.Node) (int64, error)

	CountExhibition(node sqlc.Node) (int64, error)

	CountFood(node sqlc.Node) ([]NodeFoodCount, error)

	CountEntryType(node sqlc.Node) ([]NodeEntryCount, error)
}

type DbService struct {
	DB *pgxpool.Pool
}

var (
	databaseURL = os.Getenv("DATABASE_URL")
	dbInstance  *DbService
)

func New() Service {
	if dbInstance != nil {
		return dbInstance
	}

	dbpool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		slog.Default().Error("failed to create connection pool", "error", err)
		os.Exit(1)
	}
	// defer dbpool.Close()

	if dbpool.Ping(context.Background()) != nil {
		slog.Default().Error("failed to ping db", "error", err)
		os.Exit(1)
	}

	dbInstance = &DbService{
		DB: dbpool,
	}
	return dbInstance
}

func (s *DbService) Password(key string) (sqlc.Node, bool, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return sqlc.Node{}, false, err
	}
	queries := sqlc.New(q)

	nodeByKey, err := queries.GetNodeByKey(ctx, pgtype.Text{
		String: key,
		Valid:  true,
	})
	if err != nil {
		return sqlc.Node{}, false, err
	}

	return nodeByKey, true, nil
}

func (s *DbService) Vote(modelID int64, visitorID int64, visitorRandom int32) error {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return err
	}
	queries := sqlc.New(q)

	err = queries.UpdateVisitorModel(ctx, sqlc.UpdateVisitorModelParams{
		ModelID: pgtype.Int8{Int64: modelID, Valid: true},
		ID:      visitorID,
		Random:  visitorRandom,
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

func (s *DbService) GetVisitor(ip netip.Addr, sqid *sqids.Sqids) (string, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
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

func (s *DbService) CreateVisitor(ip netip.Addr, rand int32, sqid *sqids.Sqids) (string, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
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

	err = q.Commit(ctx)
	if err != nil {
		return "", err
	}

	visitorF3SiD, err := sqid.Encode([]uint64{uint64(visitorByIp.ID), uint64(visitorByIp.Random)})
	if err != nil {
		return "", err
	}

	return visitorF3SiD, nil
}

func (s *DbService) PushEntry(node sqlc.Node, visitorID int64, visitorRandom int32) error {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
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

	entryLog, err := queries.GetEntryLogByVisitorId(ctx, pgtype.Int8{
		Int64: visitorById.ID,
		Valid: true,
	})

	var entryLogType sqlc.EntryLogsType

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// First time entering
			entryLogType = sqlc.EntryLogsTypeENTERED
		} else {
			return err
		}
	} else {
		if entryLog.Type == sqlc.EntryLogsTypeENTERED {
			entryLogType = sqlc.EntryLogsTypeLEFT
		} else {
			entryLogType = sqlc.EntryLogsTypeENTERED
		}
	}

	err = queries.CreateEntryLog(ctx, sqlc.CreateEntryLogParams{
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

func (s *DbService) PushFoodStall(node sqlc.Node, visitorID int64, visitorRandom int32, foods []Foods) error {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
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

func (s *DbService) PushExhibition(node sqlc.Node, visitorID int64, visitorRandom int32) error {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
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

func (s *DbService) UpdatePushNode(node sqlc.Node, id int64, foodID int64, quantity int32) error {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return err
	}
	queries := sqlc.New(q)

	if node.Type == sqlc.NodeTypeFOODSTALL {
		err = queries.UpdateFoodStallLog(ctx, sqlc.UpdateFoodStallLogParams{
			ID:       id,
			FoodID:   pgtype.Int8{Int64: foodID, Valid: true},
			Quantity: quantity,
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
	queries := sqlc.New(q)

	err = queries.UpdateBattery(ctx, sqlc.UpdateBatteryParams{
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
	queries := sqlc.New(q)

	// Check if visitor exists
	_, err = queries.GetEntryLogByVisitorId(ctx,
		pgtype.Int8{
			Int64: visitorID,
			Valid: true,
		})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
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
	FoodID    int64     `json:"food_id"`
	FoodName  string    `json:"food_name"`
	Quantity  int32     `json:"quantity"`
	Price     int32     `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

type ExhibitionRowLog struct {
	Id        int64     `json:"id"`
	F3SiD     string    `json:"f3sid"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *DbService) EntryRow(node sqlc.Node, sqid *sqids.Sqids) ([]EntryRowLog, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return nil, err
	}
	queries := sqlc.New(q)

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

		rowLog = append(rowLog, EntryRowLog{row.ID, visitorF3SiD, string(row.Type), row.CreatedAt.Time})
	}

	return rowLog, nil
}

func (s *DbService) FoodstallRow(node sqlc.Node, sqid *sqids.Sqids) ([]FoodstallRowLog, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return nil, err
	}
	queries := sqlc.New(q)

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

		rowLog = append(rowLog, FoodstallRowLog{row.ID, visitorF3SiD, foodByID.ID, foodByID.Name, row.Quantity, foodByID.Price, row.CreatedAt.Time})
	}

	return rowLog, nil
}

func (s *DbService) ExhibitionRow(node sqlc.Node, sqid *sqids.Sqids) ([]ExhibitionRowLog, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return nil, err
	}
	queries := sqlc.New(q)

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

func (s *DbService) IpNode(ip netip.Addr) (sqlc.Node, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return sqlc.Node{}, err
	}
	queries := sqlc.New(q)

	nodeByIp, err := queries.GetNodeByIp(ctx, &ip)
	if err != nil {
		return sqlc.Node{}, err
	}

	err = queries.DeleteNodeIp(ctx, &ip)
	if err != nil {
		return sqlc.Node{}, err
	}

	err = q.Commit(ctx)
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

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return nil, err
	}
	queries := sqlc.New(q)

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

func (s *DbService) CountEntry(node sqlc.Node) (int64, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return 0, err
	}
	queries := sqlc.New(q)

	count, err := queries.CountEntryLogByNodeId(ctx, pgtype.Int8{
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

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return 0, err
	}
	queries := sqlc.New(q)

	count, err := queries.CountFoodStallLogByNodeId(ctx, pgtype.Int8{
		Int64: node.ID,
		Valid: true,
	})
	if err != nil {
		return 0, err
	}

	return count, nil
}
func (s *DbService) CountExhibition(node sqlc.Node) (int64, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return 0, err
	}
	queries := sqlc.New(q)

	count, err := queries.CountExhibitionLogByNodeId(ctx, pgtype.Int8{
		Int64: node.ID,
		Valid: true,
	})
	if err != nil {
		return 0, err
	}

	return count, nil
}

type NodeFoodCount struct {
	ID    int
	Name  string
	Count int
}

func (s *DbService) CountFood(node sqlc.Node) ([]NodeFoodCount, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return []NodeFoodCount{}, err
	}
	queries := sqlc.New(q)

	foods, err := queries.GetFoodsByNodeId(ctx, pgtype.Int8{
		Int64: node.ID,
		Valid: true,
	})
	if err != nil {
		return nil, err
	}

	var foodsList []NodeFoodCount

	for _, food := range foods {
		foodCountById, err := queries.CountFood(ctx, sqlc.CountFoodParams{
			NodeID: pgtype.Int8{Int64: node.ID, Valid: true},
			FoodID: pgtype.Int8{Int64: food.ID, Valid: true},
		})
		if err != nil {
			return []NodeFoodCount{}, err
		}
		foodsList = append(foodsList, NodeFoodCount{int(food.ID), food.Name, int(foodCountById)})
	}

	return foodsList, nil
}

type NodeEntryCount struct {
	Type  sqlc.EntryLogsType
	Count int
}

func (s *DbService) CountEntryType(node sqlc.Node) ([]NodeEntryCount, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return []NodeEntryCount{}, err
	}
	queries := sqlc.New(q)

	var nodeEntryCounts []NodeEntryCount
	for _, entryType := range []sqlc.EntryLogsType{sqlc.EntryLogsTypeENTERED, sqlc.EntryLogsTypeLEFT} {
		count, err := queries.CountEntryLogTypeByNodeId(ctx, sqlc.CountEntryLogTypeByNodeIdParams{
			NodeID: pgtype.Int8{Int64: node.ID, Valid: true},
			Type:   entryType,
		})
		if err != nil {
			return []NodeEntryCount{}, err
		}
		nodeEntryCounts = append(nodeEntryCounts, NodeEntryCount{entryType, int(count)})
	}

	return nodeEntryCounts, nil
}
