package database

import (
	sql "backend/internal/sqlc"
	"context"
	"errors"
	"fmt"
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
	"github.com/rs/zerolog/log"
)

// Service represents a DbService that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are DbService-specific.
	Health() map[string]string

	Password(key string) (*sql.Node, bool, error)

	Visitor(ip netip.Addr, sqid *sqids.Sqids) (string, error)

	PushNode(node *sql.Node, visitorID int64, visitorRand int32, quantity int) error

	StatusNode(nodeID int64, level int32, chargingTime int32, dischargingTime int32, charging bool) error

	IsVisitorFirst(visitorID int64) (bool, error)
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
		log.Fatal().AnErr("error", err).Send()
	}
	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal().AnErr("error", err).Send()
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
		log.Fatal().AnErr("db down: %v", err).Send() // Log the error and terminate the program
		return stats
	}

	return stats
}

func (s *DbService) Password(key string) (*sql.Node, bool, error) {
	ctx := context.Background()

	q, err := s.DB.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		if q != nil {
			_ = q.Rollback(ctx)
		}
	}(q, ctx)
	if err != nil {
		return nil, false, err
	}
	queries := sql.New(q)

	nodeByKey, err := queries.GetNodeByKey(ctx, pgtype.Text{
		String: key,
		Valid:  true,
	})
	if err != nil {
		return nil, false, err
	}

	return &nodeByKey, true, nil
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

			visitorSQID, err := sqid.Encode([]uint64{uint64(visitorByIp.ID), uint64(visitorByIp.Random)})
			if err != nil {
				return "", err
			}

			return visitorSQID, nil
		} else {
			return "", err
		}
	}

	visitorSQID, err := sqid.Encode([]uint64{uint64(visitorByIp.ID), uint64(visitorByIp.Random)})
	if err != nil {
		return "", err
	}

	return visitorSQID, nil
}

func (s *DbService) PushNode(node *sql.Node, visitorID int64, visitorRandom int32, quantity int) error {
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

	switch node.Type {
	case sql.NodeTypeENTRY:
		var entryLogType sql.EntryLogsType
		entryLog, err := queries.GetEntryLogByVisitorId(ctx, pgtype.Int8{Int64: visitorById.ID, Valid: true})

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				// First time entering
				entryLogType = sql.EntryLogsTypeENTERED

				err := queries.UpdateVisitorQuantity(ctx, sql.UpdateVisitorQuantityParams{
					Quantity: int32(quantity),
					ID:       visitorById.ID,
				})
				if err != nil {
					return err
				}
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

	case sql.NodeTypeFOODSTALL:
		err := queries.CreateFoodStallLog(ctx, sql.CreateFoodStallLogParams{
			NodeID:    pgtype.Int8{Int64: node.ID, Valid: true},
			VisitorID: pgtype.Int8{Int64: visitorById.ID, Valid: true},
			Quantity:  int32(quantity),
		})
		if err != nil {
			return err
		}

		err = q.Commit(ctx)
		if err != nil {
			return err
		}

		return nil

	case sql.NodeTypeEXHIBITION:
		err := queries.CreateExhibitionLog(ctx, sql.CreateExhibitionLogParams{
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
