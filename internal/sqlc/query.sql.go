// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package sqlc

import (
	"context"
	"net/netip"

	"github.com/jackc/pgx/v5/pgtype"
)

const countEntryLogByNodeId = `-- name: CountEntryLogByNodeId :one
SELECT COUNT(*)
FROM entry_logs
WHERE node_id = $1
`

func (q *Queries) CountEntryLogByNodeId(ctx context.Context, nodeID pgtype.Int8) (int64, error) {
	row := q.db.QueryRow(ctx, countEntryLogByNodeId, nodeID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countExhibitionLogByNodeId = `-- name: CountExhibitionLogByNodeId :one
SELECT COUNT(*)
FROM exhibition_logs
WHERE node_id = $1
`

func (q *Queries) CountExhibitionLogByNodeId(ctx context.Context, nodeID pgtype.Int8) (int64, error) {
	row := q.db.QueryRow(ctx, countExhibitionLogByNodeId, nodeID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countFood = `-- name: CountFood :one
SELECT SUM(quantity)
FROM food_stall_logs
WHERE node_id = $1
    AND food_id = $2
`

type CountFoodParams struct {
	NodeID pgtype.Int8
	FoodID pgtype.Int8
}

func (q *Queries) CountFood(ctx context.Context, arg CountFoodParams) (int64, error) {
	row := q.db.QueryRow(ctx, countFood, arg.NodeID, arg.FoodID)
	var sum int64
	err := row.Scan(&sum)
	return sum, err
}

const countFoodStallLogByNodeId = `-- name: CountFoodStallLogByNodeId :one
SELECT SUM(quantity)
FROM food_stall_logs
WHERE node_id = $1
`

func (q *Queries) CountFoodStallLogByNodeId(ctx context.Context, nodeID pgtype.Int8) (int64, error) {
	row := q.db.QueryRow(ctx, countFoodStallLogByNodeId, nodeID)
	var sum int64
	err := row.Scan(&sum)
	return sum, err
}

const createBattery = `-- name: CreateBattery :exec
INSERT INTO batteries (
        node_id,
        level,
        charging_time,
        discharging_time,
        charging,
        updated_at
    )
VALUES ($1, $2, $3, $4, $5, now())
`

type CreateBatteryParams struct {
	NodeID          pgtype.Int8
	Level           pgtype.Int4
	ChargingTime    pgtype.Int4
	DischargingTime pgtype.Int4
	Charging        pgtype.Bool
}

func (q *Queries) CreateBattery(ctx context.Context, arg CreateBatteryParams) error {
	_, err := q.db.Exec(ctx, createBattery,
		arg.NodeID,
		arg.Level,
		arg.ChargingTime,
		arg.DischargingTime,
		arg.Charging,
	)
	return err
}

const createEntryLog = `-- name: CreateEntryLog :exec
INSERT INTO entry_logs (node_id, visitor_id, type)
VALUES ($1, $2, $3)
`

type CreateEntryLogParams struct {
	NodeID    pgtype.Int8
	VisitorID pgtype.Int8
	Type      EntryLogsType
}

func (q *Queries) CreateEntryLog(ctx context.Context, arg CreateEntryLogParams) error {
	_, err := q.db.Exec(ctx, createEntryLog, arg.NodeID, arg.VisitorID, arg.Type)
	return err
}

const createExhibitionLog = `-- name: CreateExhibitionLog :exec
INSERT INTO exhibition_logs (node_id, visitor_id)
VALUES ($1, $2)
`

type CreateExhibitionLogParams struct {
	NodeID    pgtype.Int8
	VisitorID pgtype.Int8
}

func (q *Queries) CreateExhibitionLog(ctx context.Context, arg CreateExhibitionLogParams) error {
	_, err := q.db.Exec(ctx, createExhibitionLog, arg.NodeID, arg.VisitorID)
	return err
}

const createFoodStallLog = `-- name: CreateFoodStallLog :exec
INSERT INTO food_stall_logs (node_id, visitor_id, food_id, quantity)
VALUES ($1, $2, $3, $4)
`

type CreateFoodStallLogParams struct {
	NodeID    pgtype.Int8
	VisitorID pgtype.Int8
	FoodID    pgtype.Int8
	Quantity  int32
}

func (q *Queries) CreateFoodStallLog(ctx context.Context, arg CreateFoodStallLogParams) error {
	_, err := q.db.Exec(ctx, createFoodStallLog,
		arg.NodeID,
		arg.VisitorID,
		arg.FoodID,
		arg.Quantity,
	)
	return err
}

const createVisitor = `-- name: CreateVisitor :one
INSERT INTO visitors (ip, random)
VALUES ($1, $2)
RETURNING id, model_id, ip, random, created_at, updated_at
`

type CreateVisitorParams struct {
	Ip     netip.Addr
	Random int32
}

func (q *Queries) CreateVisitor(ctx context.Context, arg CreateVisitorParams) (Visitor, error) {
	row := q.db.QueryRow(ctx, createVisitor, arg.Ip, arg.Random)
	var i Visitor
	err := row.Scan(
		&i.ID,
		&i.ModelID,
		&i.Ip,
		&i.Random,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteNodeIp = `-- name: DeleteNodeIp :exec
UPDATE nodes
SET ip = NULL
WHERE ip = $1
`

func (q *Queries) DeleteNodeIp(ctx context.Context, ip *netip.Addr) error {
	_, err := q.db.Exec(ctx, deleteNodeIp, ip)
	return err
}

const getEntryLogByNodeId = `-- name: GetEntryLogByNodeId :many
SELECT id, node_id, visitor_id, type, created_at, updated_at
FROM entry_logs
WHERE node_id = $1
ORDER BY id DESC
LIMIT 10
`

func (q *Queries) GetEntryLogByNodeId(ctx context.Context, nodeID pgtype.Int8) ([]EntryLog, error) {
	rows, err := q.db.Query(ctx, getEntryLogByNodeId, nodeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []EntryLog
	for rows.Next() {
		var i EntryLog
		if err := rows.Scan(
			&i.ID,
			&i.NodeID,
			&i.VisitorID,
			&i.Type,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEntryLogByVisitorId = `-- name: GetEntryLogByVisitorId :one
SELECT id, node_id, visitor_id, type, created_at, updated_at
FROM entry_logs
WHERE visitor_id = $1
ORDER BY id DESC
LIMIT 1
`

func (q *Queries) GetEntryLogByVisitorId(ctx context.Context, visitorID pgtype.Int8) (EntryLog, error) {
	row := q.db.QueryRow(ctx, getEntryLogByVisitorId, visitorID)
	var i EntryLog
	err := row.Scan(
		&i.ID,
		&i.NodeID,
		&i.VisitorID,
		&i.Type,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getExhibitionLogByNodeId = `-- name: GetExhibitionLogByNodeId :many
SELECT id, node_id, visitor_id, created_at, updated_at
FROM exhibition_logs
WHERE node_id = $1
ORDER BY id DESC
LIMIT 10
`

func (q *Queries) GetExhibitionLogByNodeId(ctx context.Context, nodeID pgtype.Int8) ([]ExhibitionLog, error) {
	rows, err := q.db.Query(ctx, getExhibitionLogByNodeId, nodeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ExhibitionLog
	for rows.Next() {
		var i ExhibitionLog
		if err := rows.Scan(
			&i.ID,
			&i.NodeID,
			&i.VisitorID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFoodById = `-- name: GetFoodById :one
SELECT id, node_id, name, price, created_at, updated_at
FROM foods
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetFoodById(ctx context.Context, id int64) (Food, error) {
	row := q.db.QueryRow(ctx, getFoodById, id)
	var i Food
	err := row.Scan(
		&i.ID,
		&i.NodeID,
		&i.Name,
		&i.Price,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getFoodStallLogByNodeId = `-- name: GetFoodStallLogByNodeId :many
SELECT id, node_id, visitor_id, food_id, quantity, created_at, updated_at
FROM food_stall_logs
WHERE node_id = $1
ORDER BY id DESC
LIMIT 10
`

func (q *Queries) GetFoodStallLogByNodeId(ctx context.Context, nodeID pgtype.Int8) ([]FoodStallLog, error) {
	rows, err := q.db.Query(ctx, getFoodStallLogByNodeId, nodeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FoodStallLog
	for rows.Next() {
		var i FoodStallLog
		if err := rows.Scan(
			&i.ID,
			&i.NodeID,
			&i.VisitorID,
			&i.FoodID,
			&i.Quantity,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFoodsByNodeId = `-- name: GetFoodsByNodeId :many
SELECT id, node_id, name, price, created_at, updated_at
FROM foods
WHERE node_id = $1
`

func (q *Queries) GetFoodsByNodeId(ctx context.Context, nodeID pgtype.Int8) ([]Food, error) {
	rows, err := q.db.Query(ctx, getFoodsByNodeId, nodeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Food
	for rows.Next() {
		var i Food
		if err := rows.Scan(
			&i.ID,
			&i.NodeID,
			&i.Name,
			&i.Price,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNodeById = `-- name: GetNodeById :one
SELECT id, key, name, ip, type, created_at, updated_at
FROM nodes
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetNodeById(ctx context.Context, id int64) (Node, error) {
	row := q.db.QueryRow(ctx, getNodeById, id)
	var i Node
	err := row.Scan(
		&i.ID,
		&i.Key,
		&i.Name,
		&i.Ip,
		&i.Type,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getNodeByIp = `-- name: GetNodeByIp :one
SELECT id, key, name, ip, type, created_at, updated_at
FROM nodes
WHERE ip = $1
LIMIT 1
`

func (q *Queries) GetNodeByIp(ctx context.Context, ip *netip.Addr) (Node, error) {
	row := q.db.QueryRow(ctx, getNodeByIp, ip)
	var i Node
	err := row.Scan(
		&i.ID,
		&i.Key,
		&i.Name,
		&i.Ip,
		&i.Type,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getNodeByKey = `-- name: GetNodeByKey :one
SELECT id, key, name, ip, type, created_at, updated_at
FROM nodes
WHERE key = $1
LIMIT 1
`

func (q *Queries) GetNodeByKey(ctx context.Context, key pgtype.Text) (Node, error) {
	row := q.db.QueryRow(ctx, getNodeByKey, key)
	var i Node
	err := row.Scan(
		&i.ID,
		&i.Key,
		&i.Name,
		&i.Ip,
		&i.Type,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getVisitorById = `-- name: GetVisitorById :one
SELECT id, model_id, ip, random, created_at, updated_at
FROM visitors
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetVisitorById(ctx context.Context, id int64) (Visitor, error) {
	row := q.db.QueryRow(ctx, getVisitorById, id)
	var i Visitor
	err := row.Scan(
		&i.ID,
		&i.ModelID,
		&i.Ip,
		&i.Random,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getVisitorByIdAndRandom = `-- name: GetVisitorByIdAndRandom :one
SELECT id, model_id, ip, random, created_at, updated_at
FROM visitors
WHERE id = $1
    AND random = $2
LIMIT 1
`

type GetVisitorByIdAndRandomParams struct {
	ID     int64
	Random int32
}

func (q *Queries) GetVisitorByIdAndRandom(ctx context.Context, arg GetVisitorByIdAndRandomParams) (Visitor, error) {
	row := q.db.QueryRow(ctx, getVisitorByIdAndRandom, arg.ID, arg.Random)
	var i Visitor
	err := row.Scan(
		&i.ID,
		&i.ModelID,
		&i.Ip,
		&i.Random,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getVisitorByIp = `-- name: GetVisitorByIp :one
SELECT id, model_id, ip, random, created_at, updated_at
FROM visitors
WHERE ip = $1
LIMIT 1
`

func (q *Queries) GetVisitorByIp(ctx context.Context, ip netip.Addr) (Visitor, error) {
	row := q.db.QueryRow(ctx, getVisitorByIp, ip)
	var i Visitor
	err := row.Scan(
		&i.ID,
		&i.ModelID,
		&i.Ip,
		&i.Random,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateBattery = `-- name: UpdateBattery :exec
UPDATE batteries
SET level = coalesce($1, level),
    charging_time = coalesce($2, charging_time),
    discharging_time = coalesce($3, discharging_time),
    charging = coalesce($4, charging),
    updated_at = $5
WHERE node_id = $6
`

type UpdateBatteryParams struct {
	Level           pgtype.Int4
	ChargingTime    pgtype.Int4
	DischargingTime pgtype.Int4
	Charging        pgtype.Bool
	ID              pgtype.Timestamp
	NodeID          pgtype.Int8
}

func (q *Queries) UpdateBattery(ctx context.Context, arg UpdateBatteryParams) error {
	_, err := q.db.Exec(ctx, updateBattery,
		arg.Level,
		arg.ChargingTime,
		arg.DischargingTime,
		arg.Charging,
		arg.ID,
		arg.NodeID,
	)
	return err
}

const updateFoodStallLog = `-- name: UpdateFoodStallLog :exec
UPDATE food_stall_logs
SET quantity = $1,
    updated_at = now()
WHERE id = $2
`

type UpdateFoodStallLogParams struct {
	Quantity int32
	ID       int64
}

func (q *Queries) UpdateFoodStallLog(ctx context.Context, arg UpdateFoodStallLogParams) error {
	_, err := q.db.Exec(ctx, updateFoodStallLog, arg.Quantity, arg.ID)
	return err
}

const updateVisitorModel = `-- name: UpdateVisitorModel :exec
UPDATE visitors
SET model_id = $1,
    updated_at = now()
WHERE id = $2
    AND random = $3
`

type UpdateVisitorModelParams struct {
	ModelID pgtype.Int8
	ID      int64
	Random  int32
}

func (q *Queries) UpdateVisitorModel(ctx context.Context, arg UpdateVisitorModelParams) error {
	_, err := q.db.Exec(ctx, updateVisitorModel, arg.ModelID, arg.ID, arg.Random)
	return err
}
