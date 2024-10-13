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

	Batteries() ([]sqlc.Battery, error)

	GetVisitor(ip netip.Addr, sqid *sqids.Sqids) (string, error)

	CreateVisitor(ip netip.Addr, rand int32, sqid *sqids.Sqids) (string, error)

	IpNode(ip netip.Addr) (sqlc.Node, error)

	NodeByID(id int64) (sqlc.Node, error)

	OTPNode(otp string) (sqlc.Node, error)

	PushEntry(node sqlc.Node, visitorID int64, visitorRandom int32) error

	PushFoodStall(node sqlc.Node, visitorID int64, visitorRandom int32, foods []Foods) error

	PushExhibition(node sqlc.Node, visitorID int64, visitorRandom int32) error

	UpdateFoodLog(node sqlc.Node, id int64, foodID int64, quantity int32) error

	StatusNode(nodeID int64, level int32, chargingTime int32, dischargingTime int32, charging bool) error

	IsVisitorFirst(visitorID int64) (bool, error)

	EntryRows(node sqlc.Node, sqid *sqids.Sqids) ([]EntryRowLog, error)

	FoodstallRows(node sqlc.Node, sqid *sqids.Sqids) ([]FoodstallRowLog, error)

	ExhibitionRows(node sqlc.Node, sqid *sqids.Sqids) ([]ExhibitionRowLog, error)

	Foods(node sqlc.Node) ([]NodeFood, error)

	CountEntry(node sqlc.Node) (int64, error)

	CountFoodStall(node sqlc.Node) (int64, error)
	QuantityFoodStall(node sqlc.Node) (int64, error)

	CountExhibition(node sqlc.Node) (int64, error)

	CountFood(node sqlc.Node) ([]NodeFoodCount, error)

	CountEntryType(node sqlc.Node) ([]NodeEntryCount, error)

	CountEntryPerHalfHour(node sqlc.Node) ([]EntryPerDay, error)

	CountFoodStallPerHalfHour(node sqlc.Node) ([]FoodStallCountPerDay, error)
	QuantityFoodStallPerHalfHour(node sqlc.Node) ([]FoodStallQuantityPerDay, error)
	TotalCountFoodStallPerHalfHour(node sqlc.Node) ([]FoodCountPerHalfHour, error)
	TotalQuantityFoodStallPerHalfHour(node sqlc.Node) ([]FoodQuantityPerHalfHour, error)

	CountExhibitionPerHalfHour(node sqlc.Node) ([]ExhibitionPerHour, error)
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

	dbconfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		slog.Error("failed to parse db config", "error", err)
		os.Exit(1)
	}
	dbpool, err := pgxpool.NewWithConfig(context.Background(), dbconfig)
	if err != nil {
		slog.Error("failed to create connection pool", "error", err)
		os.Exit(1)
	}
	// defer dbpool.Close()

	if dbpool.Ping(context.Background()) != nil {
		slog.Error("failed to ping db", "error", err)
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

func (s *DbService) NodeByID(id int64) (sqlc.Node, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return sqlc.Node{}, err
	}
	queries := sqlc.New(q)

	nodeByID, err := queries.GetNodeById(ctx, id)
	if err != nil {
		return sqlc.Node{}, err
	}

	return nodeByID, nil
}

func (s *DbService) Batteries() ([]sqlc.Battery, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return []sqlc.Battery{}, err
	}
	queries := sqlc.New(q)

	batteries, err := queries.GetBatteries(ctx)
	if err != nil {
		return nil, err
	}

	return batteries, nil
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

	visitorByIp, err := queries.GetVisitorByIp(ctx, &ip)
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
		Ip:     &ip,
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

	entryLog, err := queries.GetEntryLogByVisitorId(ctx, visitorById.ID)

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
		NodeID:    node.ID,
		VisitorID: visitorById.ID,
		Type:      entryLogType,
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

		err = queries.CreateFoodStalllogByNodeFoodId(ctx, sqlc.CreateFoodStalllogByNodeFoodIdParams{
			NodeID:    node.ID,
			FoodID:    foodById,
			VisitorID: visitorById.ID,
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
		NodeID:    node.ID,
		VisitorID: visitorById.ID,
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

func (s *DbService) UpdateFoodLog(node sqlc.Node, id int64, foodID int64, quantity int32) error {
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
			FoodID:   foodID,
			NodeID:   node.ID,
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
	_, err = queries.GetEntryLogByVisitorId(ctx, visitorID)
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

func (s *DbService) EntryRows(node sqlc.Node, sqid *sqids.Sqids) ([]EntryRowLog, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return nil, err
	}
	queries := sqlc.New(q)

	rows, err := queries.GetEntryLogByNodeId(ctx, node.ID)
	if err != nil {
		return nil, err
	}

	rowLog := make([]EntryRowLog, len(rows))

	for i, row := range rows {
		visitorByID, err := queries.GetVisitorById(ctx, row.VisitorID)
		if err != nil {
			return nil, err
		}

		visitorF3SiD, err := sqid.Encode([]uint64{uint64(visitorByID.ID), uint64(visitorByID.Random)})
		if err != nil {
			return nil, err
		}

		rowLog[i] = EntryRowLog{row.ID, visitorF3SiD, string(row.Type), row.CreatedAt.Time}
	}

	return rowLog, nil
}

func (s *DbService) FoodstallRows(node sqlc.Node, sqid *sqids.Sqids) ([]FoodstallRowLog, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return nil, err
	}
	queries := sqlc.New(q)

	rows, err := queries.GetFoodStallLogByNodeId(ctx, node.ID)
	if err != nil {
		return nil, err
	}

	rowLog := make([]FoodstallRowLog, len(rows))

	for i, row := range rows {
		visitorByID, err := queries.GetVisitorById(ctx, row.VisitorID)
		if err != nil {
			return nil, err
		}

		visitorF3SiD, err := sqid.Encode([]uint64{uint64(visitorByID.ID), uint64(visitorByID.Random)})
		if err != nil {
			return nil, err
		}

		foodByNodeFoodId, err := queries.GetFoodByNodeFoodId(ctx, row.NodeFoodID)
		if err != nil {
			return nil, err
		}

		foodByID, err := queries.GetFoodById(ctx, foodByNodeFoodId.ID)
		if err != nil {
			return nil, err
		}

		rowLog[i] = FoodstallRowLog{row.ID, visitorF3SiD, foodByID.ID, foodByID.Name, row.Quantity, foodByID.Price, row.CreatedAt.Time}
	}

	return rowLog, nil
}

func (s *DbService) ExhibitionRows(node sqlc.Node, sqid *sqids.Sqids) ([]ExhibitionRowLog, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return nil, err
	}
	queries := sqlc.New(q)

	rows, err := queries.GetExhibitionLogByNodeId(ctx, node.ID)
	if err != nil {
		return nil, err
	}

	rowLog := make([]ExhibitionRowLog, len(rows))

	for i, row := range rows {
		visitorByID, err := queries.GetVisitorById(ctx, row.VisitorID)
		if err != nil {
			return nil, err
		}

		visitorF3SiD, err := sqid.Encode([]uint64{uint64(visitorByID.ID), uint64(visitorByID.Random)})
		if err != nil {
			return nil, err
		}

		rowLog[i] = ExhibitionRowLog{row.ID, visitorF3SiD, row.CreatedAt.Time}
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

func (s *DbService) OTPNode(otp string) (sqlc.Node, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return sqlc.Node{}, err
	}
	queries := sqlc.New(q)

	nodeByOTP, err := queries.GetNodeByOTP(ctx, pgtype.Text{
		String: otp,
		Valid:  true,
	})
	if err != nil {
		return sqlc.Node{}, err
	}

	err = queries.DeleteNodeOTP(ctx, pgtype.Text{
		String: otp,
		Valid:  true,
	})
	if err != nil {
		return sqlc.Node{}, err
	}

	err = q.Commit(ctx)
	if err != nil {
		return sqlc.Node{}, err
	}

	return nodeByOTP, nil
}

type NodeFood struct {
	ID       int
	Name     string
	Price    int
	Quantity int
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

	foods, err := queries.GetFoodsByNodeId(ctx, node.ID)
	if err != nil {
		return nil, err
	}

	foodsList := make([]NodeFood, len(foods))

	for i, food := range foods {
		foodsList[i] = NodeFood{int(food.ID), food.Name, int(food.Price), int(food.Quantity)}
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

	count, err := queries.CountEntryLog(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// func (s *DbService) QuantityEntry(node sqlc.Node) (int64, error) {
// 	ctx := context.Background()

// 	q, err := s.DB.Begin(ctx)
// 	if err != nil {
// 		return 0, err
// 	}
// 	queries := sqlc.New(q)

// 	count, err := queries.CountEntryLog(ctx)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return count, nil
// }

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

	count, err := queries.CountFoodStallLogByNodeIdOwned(ctx, node.ID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *DbService) QuantityFoodStall(node sqlc.Node) (int64, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return 0, err
	}
	queries := sqlc.New(q)

	count, err := queries.QuantityFoodStallLogByNodeIdOwned(ctx, node.ID)
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

	count, err := queries.CountExhibitionLogByNodeId(ctx, node.ID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

type NodeFoodCount struct {
	ID       int
	Name     string
	Count    int
	Quantity int
	Price    int
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

	foods, err := queries.GetFoodsByNodeId(ctx, node.ID)
	if err != nil {
		return nil, err
	}

	foodsList := make([]NodeFoodCount, len(foods))

	for i, food := range foods {
		foodCountById, err := queries.CountFood(ctx, food.ID)
		if err != nil {
			return []NodeFoodCount{}, err
		}
		foodsList[i] = NodeFoodCount{int(food.ID), food.Name, int(foodCountById), int(food.Quantity) * int(foodCountById), int(food.Price)}
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

	nodeEntryCounts := make([]NodeEntryCount, 2)
	for i, entryType := range []sqlc.EntryLogsType{sqlc.EntryLogsTypeENTERED, sqlc.EntryLogsTypeLEFT} {
		count, err := queries.CountEntryLogTypeByType(ctx, entryType)
		if err != nil {
			return []NodeEntryCount{}, err
		}
		nodeEntryCounts[i] = NodeEntryCount{entryType, int(count)}
	}

	return nodeEntryCounts, nil
}

type EntryPerDay struct {
	Name    string             `json:"name"`
	Entries []EntryPerHalfHour `json:"entries"`
}

type EntryPerHalfHour struct {
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
	Count  int `json:"count"`
}

func (s *DbService) CountEntryPerHalfHour(node sqlc.Node) ([]EntryPerDay, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return []EntryPerDay{}, err
	}
	queries := sqlc.New(q)

	var countPerDay = make([]EntryPerDay, 2)
	for i, entryType := range []sqlc.EntryLogsType{sqlc.EntryLogsTypeENTERED, sqlc.EntryLogsTypeLEFT} {
		rows, err := queries.CountEntryPerHalfHourByEntryType(ctx, entryType)
		if err != nil {
			return []EntryPerDay{}, err
		}

		var countPerHour = make([]EntryPerHalfHour, len(rows))
		for i, row := range rows {
			countPerHour[i] = EntryPerHalfHour{int(row.Hour), int(row.Minute), int(row.Count)}
		}
		countPerDay[i] = EntryPerDay{string(entryType), countPerHour}
	}

	return countPerDay, nil
}

type FoodStallCountPerDay struct {
	Name  string                 `json:"name"`
	Foods []FoodCountPerHalfHour `json:"foods"`
}

type FoodCountPerHalfHour struct {
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
	Count  int `json:"count"`
}

func (s *DbService) CountFoodStallPerHalfHour(node sqlc.Node) ([]FoodStallCountPerDay, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return []FoodStallCountPerDay{}, err
	}
	queries := sqlc.New(q)

	foods, err := queries.GetFoodsByNodeId(ctx, node.ID)
	if err != nil {
		return []FoodStallCountPerDay{}, err
	}

	var countPerDay = make([]FoodStallCountPerDay, len(foods))

	for i, food := range foods {
		countPerHourByFood, err := queries.CountFoodStallPerHalfHourByFoodId(ctx, food.ID)
		if err != nil {
			return []FoodStallCountPerDay{}, err
		}

		countPerHour := make([]FoodCountPerHalfHour, len(countPerHourByFood))
		for i, row := range countPerHourByFood {
			countPerHour[i] = FoodCountPerHalfHour{int(row.Hour), int(row.Minute), int(row.Count)}
		}
		countPerDay[i] = FoodStallCountPerDay{food.Name, countPerHour}
	}

	return countPerDay, nil
}

type FoodStallQuantityPerDay struct {
	Name  string                    `json:"name"`
	Foods []FoodQuantityPerHalfHour `json:"foods"`
}

type FoodQuantityPerHalfHour struct {
	Hour     int `json:"hour"`
	Minute   int `json:"minute"`
	Quantity int `json:"quantity"`
}

func (s *DbService) QuantityFoodStallPerHalfHour(node sqlc.Node) ([]FoodStallQuantityPerDay, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return []FoodStallQuantityPerDay{}, err
	}
	queries := sqlc.New(q)

	foods, err := queries.GetFoodsByNodeId(ctx, node.ID)
	if err != nil {
		return []FoodStallQuantityPerDay{}, err
	}

	var countPerDay = make([]FoodStallQuantityPerDay, len(foods))

	for i, food := range foods {
		countPerHourByFood, err := queries.QuantityFoodStallPerHourByFoodId(ctx, food.ID)
		if err != nil {
			return []FoodStallQuantityPerDay{}, err
		}

		countPerHour := make([]FoodQuantityPerHalfHour, len(countPerHourByFood))
		for i, row := range countPerHourByFood {
			countPerHour[i] = FoodQuantityPerHalfHour{int(row.Hour), int(row.Minute), int(row.Quantity)}
		}
		countPerDay[i] = FoodStallQuantityPerDay{food.Name, countPerHour}
	}

	return countPerDay, nil
}

func (s *DbService) TotalCountFoodStallPerHalfHour(node sqlc.Node) ([]FoodCountPerHalfHour, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return []FoodCountPerHalfHour{}, err
	}
	queries := sqlc.New(q)

	rows, err := queries.CountFoodStallPerHalfHourByNodeId(ctx, node.ID)
	if err != nil {
		return []FoodCountPerHalfHour{}, err
	}

	var countPerHour = make([]FoodCountPerHalfHour, len(rows))
	for i, row := range rows {
		countPerHour[i] = FoodCountPerHalfHour{int(row.Hour), int(row.Minute), int(row.Count)}
	}

	return countPerHour, nil
}

func (s *DbService) TotalQuantityFoodStallPerHalfHour(node sqlc.Node) ([]FoodQuantityPerHalfHour, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return []FoodQuantityPerHalfHour{}, err
	}
	queries := sqlc.New(q)

	rows, err := queries.QuantityFoodStallPerHalfHourByNodeId(ctx, node.ID)
	if err != nil {
		return []FoodQuantityPerHalfHour{}, err
	}

	var countPerHour = make([]FoodQuantityPerHalfHour, len(rows))
	for i, row := range rows {
		countPerHour[i] = FoodQuantityPerHalfHour{int(row.Hour), int(row.Minute), int(row.Quantity)}
	}

	return countPerHour, nil
}

// type ExhibitionPerDay struct {
// 	Name   string              `json:"name"`
// 	Counts []ExhibitionPerHour `json:"counts"`
// }

type ExhibitionPerHour struct {
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
	Count  int `json:"count"`
}

func (s *DbService) CountExhibitionPerHalfHour(node sqlc.Node) ([]ExhibitionPerHour, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return []ExhibitionPerHour{}, err
	}
	queries := sqlc.New(q)

	rows, err := queries.CountExhibitionPerHalfHourByNodeId(ctx, node.ID)
	if err != nil {
		return []ExhibitionPerHour{}, err
	}

	var countPerHour = make([]ExhibitionPerHour, len(rows))
	for i, row := range rows {
		countPerHour[i] = ExhibitionPerHour{int(row.Hour), int(row.Minute), int(row.Count)}
	}

	// var countPerDay = make([]ExhibitionPerHour, 1)
	// countPerDay[0] = ExhibitionPerHour{countPerHour}

	return countPerHour, nil
}
