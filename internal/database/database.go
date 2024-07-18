package database

import (
	"backend/internal/sqlc"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/netip"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog/log"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	Visitor(ip string) (string, error)
}

type service struct {
	db *pgxpool.Pool
}

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	schema     = os.Getenv("DB_SCHEMA")
	dbInstance *service
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
	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.Ping(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatal().AnErr("db down: %v", err).Send() // Log the error and terminate the program
		return stats
	}

	return stats
}

func (s *service) Visitor(ip string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	q, err := s.db.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return "", err
	}
	queries := sql.New(q)

	// Parse ip address into type
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return "", err
	}

	// Check if visitor exists with ip
	visitorByIp, err := queries.GetVisitorByIp(ctx, &addr)
	if errors.Is(err, pgx.ErrNoRows) {
		// Create a new visitor
		visitorByIp, err := queries.CreateVisitor(ctx, &addr)
		if err != nil {
			return "", err
		}
		err = q.Commit(ctx)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%x", visitorByIp.ID.Bytes), nil
	} else if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", visitorByIp.ID.Bytes), nil
}

func (s *service) PushNode(nodeID string, visitorUUID string, quantity int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	q, err := s.db.Begin(ctx)
	defer func(q pgx.Tx, ctx context.Context) {
		_ = q.Rollback(ctx)
	}(q, ctx)
	if err != nil {
		return err
	}
	queries := sql.New(q)

	// Check if visitor exists
	visitorById, err := queries.GetVisitorById(ctx, pgtype.UUID{Bytes: [16]byte([]byte(visitorUUID)), Valid: true})
	if err != nil {
		return err
	}

	// Get type of node
	nodeById, err := queries.GetNodeById(ctx, nodeID)
	if err != nil {
		return err
	}

	switch nodeById.Type {
	case sql.NodeTypeENTRY:
		var entryLogType sql.EntryLogsType
		entryLog, err := queries.GetEntryLogByNodeId(ctx, pgtype.Text{String: nodeById.ID, Valid: true})
		if errors.Is(err, pgx.ErrNoRows) {
			entryLogType = sql.EntryLogsTypeENTERED
		} else if err != nil {
			return err
		}

		if entryLog.Type == sql.EntryLogsTypeENTERED {
			entryLogType = sql.EntryLogsTypeLEFT
		} else {
			entryLogType = sql.EntryLogsTypeENTERED
		}

		err = queries.CreateEntryLog(ctx, sql.CreateEntryLogParams{
			NodeID: pgtype.Text{
				String: nodeById.ID,
				Valid:  true,
			},
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

	case sql.NodeTypeFOODSTALL:
		err := queries.CreateFoodStallLog(ctx, sql.CreateFoodStallLogParams{
			NodeID: pgtype.Text{
				String: nodeById.ID,
				Valid:  true,
			},
			VisitorID: visitorById.ID,
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
			NodeID: pgtype.Text{
				String: nodeById.ID,
				Valid:  true,
			},
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

	return nil
}
