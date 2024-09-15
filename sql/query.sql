-- name: GetVisitorByIp :one
SELECT *
FROM visitors
WHERE ip = ?
LIMIT 1;
-- name: GetVisitorById :one
SELECT *
FROM visitors
WHERE id = ?
LIMIT 1;
-- name: GetVisitorByIdAndRandom :one
SELECT *
FROM visitors
WHERE id = ?
    AND random = ?
LIMIT 1;
-- name: CreateVisitor :one
INSERT INTO visitors (ip, random)
VALUES (?, ?)
RETURNING *;
-- name: CreateBattery :exec
INSERT INTO batteries (
        node_id,
        level,
        charging_time,
        discharging_time,
        charging,
        updated_at
    )
VALUES (?, ?, ?, ?, ?, date('now'));
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
WHERE id = ?
LIMIT 1;
-- name: GetNodeByKey :one
SELECT *
FROM nodes
WHERE key = ?
LIMIT 1;
-- name: GetNodeByIp :one
SELECT *
FROM nodes
WHERE ip = ?
LIMIT 1;
-- name: DeleteNodeIp :exec
UPDATE nodes
SET ip = NULL
WHERE ip = ?;
-- name: GetFoodById :one
SELECT *
FROM foods
WHERE id = ?
LIMIT 1;
-- name: GetFoodsByNodeId :many
SELECT *
FROM foods
WHERE node_id = ?;
-- name: GetEntryLogByVisitorId :one
SELECT *
FROM entry_logs
WHERE visitor_id = ?
ORDER BY id DESC
LIMIT 1;
-- name: CreateEntryLog :exec
INSERT INTO entry_logs (node_id, visitor_id, type)
VALUES (?, ?, ?);
-- name: CreateExhibitionLog :exec
INSERT INTO exhibition_logs (node_id, visitor_id)
VALUES (?, ?);
-- name: CreateFoodStallLog :exec
INSERT INTO food_stall_logs (node_id, visitor_id, food_id, quantity)
VALUES (?, ?, ?, ?);
-- name: GetEntryLogByNodeId :many
SELECT *
FROM entry_logs
WHERE node_id = ?
ORDER BY id DESC
LIMIT 10;
-- name: GetFoodStallLogByNodeId :many
SELECT *
FROM food_stall_logs
WHERE node_id = ?
ORDER BY id DESC
LIMIT 10;
-- name: GetExhibitionLogByNodeId :many
SELECT *
FROM exhibition_logs
WHERE node_id = ?
ORDER BY id DESC
LIMIT 10;
-- name: UpdateFoodStallLog :exec
UPDATE food_stall_logs
SET quantity = ?,
    updated_at = date('now')
WHERE id = ?;
-- name: CountEntryLogByNodeId :one
SELECT COUNT(*)
FROM entry_logs
WHERE node_id = ?;
-- name: CountFoodStallLogByNodeId :one
SELECT SUM(quantity)
FROM food_stall_logs
WHERE node_id = ?;
-- name: CountExhibitionLogByNodeId :one
SELECT COUNT(*)
FROM exhibition_logs
WHERE node_id = ?;