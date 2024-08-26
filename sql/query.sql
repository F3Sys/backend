-- name: GetVisitorByIp :one
SELECT *
FROM visitors
WHERE ip = $1
LIMIT 1;
-- name: GetVisitorById :one
SELECT *
FROM visitors
WHERE id = $1
LIMIT 1;
-- name: GetVisitorByIdAndRandom :one
SELECT *
FROM visitors
WHERE id = $1
    AND random = $2
LIMIT 1;
-- name: CreateVisitor :one
INSERT INTO visitors (ip, random)
VALUES ($1, $2)
RETURNING *;
-- name: UpdateVisitorQuantity :exec
UPDATE visitors
SET quantity = $1,
    updated_at = now()
WHERE id = $2;
-- name: CreateBattery :exec
INSERT INTO batteries (
        node_id,
        level,
        charging_time,
        discharging_time,
        charging,
        updated_at
    )
VALUES ($1, $2, $3, $4, $5, now());
-- name: UpdateBattery :exec
UPDATE batteries
SET level = coalesce(sqlc.narg('level'), level),
    charging_time = coalesce(sqlc.narg('charging_time'), charging_time),
    discharging_time = coalesce(sqlc.narg('discharging_time'), discharging_time),
    charging = coalesce(sqlc.narg('charging'), charging),
    updated_at = sqlc.arg('id')
WHERE node_id = sqlc.arg('node_id');
-- name: GetNodeById :one
SELECT *
FROM nodes
WHERE id = $1
LIMIT 1;
-- name: GetNodeByKey :one
SELECT *
FROM nodes
WHERE key = $1
LIMIT 1;
-- name: GetEntryLogByVisitorId :one
SELECT *
FROM entry_logs
WHERE visitor_id = $1
ORDER BY id DESC
LIMIT 1;
-- name: CreateEntryLog :exec
INSERT INTO entry_logs (node_id, visitor_id, type)
VALUES ($1, $2, $3);
-- name: CreateExhibitionLog :exec
INSERT INTO exhibition_logs (node_id, visitor_id)
VALUES ($1, $2);
-- name: CreateFoodStallLog :exec
INSERT INTO food_stall_logs (node_id, visitor_id, quantity)
VALUES ($1, $2, $3);
-- name: GetEntryLogByNodeId :many
SELECT *
FROM entry_logs
WHERE node_id = $1
ORDER BY id DESC
LIMIT 10;
-- name: GetFoodStallLogByNodeId :many
SELECT *
FROM food_stall_logs
WHERE node_id = $1
ORDER BY id DESC
LIMIT 10;
-- name: GetExhibitionLogByNodeId :many
SELECT *
FROM exhibition_logs
WHERE node_id = $1
ORDER BY id DESC
LIMIT 10;
-- name: UpdateFoodStallLog :exec
UPDATE food_stall_logs
SET quantity = $1,
    updated_at = now()
WHERE id = $2;