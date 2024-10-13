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
-- name: UpdateVisitorModel :exec
UPDATE visitors
SET model_id = $1,
    updated_at = now()
WHERE id = $2
    AND random = $3;
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
-- name: GetBatteries :many
SELECT *
FROM batteries;
-- name: UpdateBattery :exec
UPDATE batteries
SET level = coalesce(sqlc.narg('level'), level),
    charging_time = coalesce(sqlc.narg('charging_time'), charging_time),
    discharging_time = coalesce(sqlc.narg('discharging_time'), discharging_time),
    charging = coalesce(sqlc.narg('charging'), charging),
    updated_at = now()
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
-- name: GetNodeByIp :one
SELECT *
FROM nodes
WHERE ip = $1
LIMIT 1;
-- name: GetNodeByOTP :one
SELECT *
FROM nodes
WHERE otp = $1 AND updated_at >= now() - INTERVAL '5 minute'
LIMIT 1;
-- name: DeleteNodeIp :exec
UPDATE nodes
SET ip = NULL
WHERE ip = $1;
-- name: DeleteNodeOTP :exec
UPDATE nodes
SET otp = NULL, updated_at = now()
WHERE otp = $1;
-- name: GetFoodById :one
SELECT *
FROM foods
WHERE id = $1
LIMIT 1;
-- name: GetFoodsByNodeId :many
SELECT f.*
FROM foods f
JOIN node_foods nf ON f.id = nf.food_id
WHERE nf.node_id = $1;
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
INSERT INTO food_stall_logs (node_food_id, visitor_id, quantity)
VALUES ($1, $2, $3);
-- name: CreateFoodStalllogByNodeFoodId :exec
INSERT INTO food_stall_logs (node_food_id, visitor_id, quantity)
VALUES (
    (SELECT id
     FROM node_foods
     WHERE food_id = $1
       AND node_id = $2),
    $3,
    $4
);
-- name: GetFoodByNodeFoodId :one
SELECT f.*
FROM foods f
JOIN node_foods nf ON f.id = nf.food_id
WHERE nf.id = $1
LIMIT 1;
-- name: GetFoodNodeByFoodAndNodeId :one
SELECT nf.*
FROM node_foods nf
JOIN foods f ON nf.food_id = f.id
WHERE f.id = $1
  AND nf.node_id = $2
LIMIT 1;
-- name: GetEntryLogByNodeId :many
SELECT *
FROM entry_logs
WHERE node_id = $1
ORDER BY id DESC
LIMIT 10;
-- name: GetFoodStallLogByNodeId :many
SELECT *
FROM food_stall_logs
WHERE node_food_id IN (
    SELECT id
    FROM node_foods
    WHERE node_id = $1
)
ORDER BY id DESC
LIMIT 10;
-- name: GetExhibitionLogByNodeId :many
SELECT *
FROM exhibition_logs
WHERE node_id = $1
ORDER BY id DESC
LIMIT 10;
-- name: UpdateFoodStallLog :exec
UPDATE food_stall_logs fsl
SET quantity = $2,
    node_food_id = (
        SELECT nf.id
        FROM node_foods nf
        WHERE nf.food_id = $3
          AND nf.node_id = $4  -- Check if the node_id owns the food
        LIMIT 1
    ),
    updated_at = NOW()
WHERE fsl.id = $1;
-- name: CountFoodStallLogByNodeId :one
SELECT SUM(quantity)
FROM food_stall_logs
WHERE node_food_id IN (
    SELECT id
    FROM node_foods
    WHERE node_id = $1
);
-- name: CountExhibitionLogByNodeId :one
SELECT COUNT(*)
FROM exhibition_logs
WHERE node_id = $1;
-- name: CountFood :one
SELECT SUM(fsl.quantity)
FROM food_stall_logs fsl
JOIN node_foods nf ON fsl.node_food_id = nf.id
WHERE nf.food_id = $1;
-- name: CountEntryLog :one
SELECT COUNT(*)
FROM entry_logs;
-- name: CountFoodStallLogByNodeIdOwned :one
SELECT SUM(fsl.quantity)
FROM food_stall_logs fsl
JOIN node_foods nf ON fsl.node_food_id = nf.id
WHERE nf.node_id = $1;
-- name: QuantityFoodStallLogByNodeIdOwned :one
SELECT SUM(fsl.quantity * f.quantity)
FROM food_stall_logs fsl
JOIN node_foods nf ON fsl.node_food_id = nf.id
JOIN foods f ON nf.food_id = f.id
WHERE nf.node_id = $1;
-- name: CountEntryLogTypeByType :one
SELECT COUNT(*)
FROM entry_logs
WHERE type = $1;
-- name: CountEntryPerHalfHourByEntryType :many
SELECT COUNT(*) AS count,
  DATE_PART('hour', el.created_at AT TIME ZONE '+09:00') AS hour,
  FLOOR(DATE_PART('minute', el.created_at AT TIME ZONE '+09:00') / 30) * 30 AS minute
FROM entry_logs el
WHERE el.type = $1
  AND DATE(el.created_at AT TIME ZONE '+09:00') = CURRENT_DATE
  AND DATE_PART('hour', el.created_at AT TIME ZONE '+09:00') BETWEEN 8 AND 18
GROUP BY hour, minute
ORDER BY hour DESC, minute DESC;
-- name: CountFoodStallPerHalfHourByFoodId :many
SELECT SUM(fsl.quantity) AS count,
  DATE_PART('hour', fsl.created_at AT TIME ZONE '+09:00') AS hour,
  FLOOR(DATE_PART('minute', fsl.created_at AT TIME ZONE '+09:00') / 30) * 30 AS minute
FROM food_stall_logs fsl
JOIN node_foods nf ON fsl.node_food_id = nf.id
WHERE nf.food_id = $1
  AND DATE(fsl.created_at AT TIME ZONE '+09:00') = CURRENT_DATE
  AND DATE_PART('hour', fsl.created_at AT TIME ZONE '+09:00') BETWEEN 8 AND 18
GROUP BY hour, minute
ORDER BY hour DESC, minute DESC;
-- name: QuantityFoodStallPerHourByFoodId :many
SELECT SUM(fsl.quantity * f.quantity) AS quantity,
  DATE_PART('hour', fsl.created_at AT TIME ZONE '+09:00') AS hour,
  FLOOR(DATE_PART('minute', fsl.created_at AT TIME ZONE '+09:00') / 30) * 30 AS minute
FROM food_stall_logs fsl
JOIN node_foods nf ON fsl.node_food_id = nf.id
JOIN foods f ON nf.food_id = f.id
WHERE nf.food_id = $1
  AND DATE(fsl.created_at AT TIME ZONE '+09:00') = CURRENT_DATE
  AND DATE_PART('hour', fsl.created_at AT TIME ZONE '+09:00') BETWEEN 8 AND 18
GROUP BY hour, minute
ORDER BY hour DESC, minute DESC;
-- name: CountFoodStallPerHalfHourByNodeId :many
SELECT SUM(fsl.quantity) AS count,
  DATE_PART('hour', fsl.created_at AT TIME ZONE '+09:00') AS hour,
  FLOOR(DATE_PART('minute', fsl.created_at AT TIME ZONE '+09:00') / 30) * 30 AS minute
FROM food_stall_logs fsl
JOIN node_foods nf ON fsl.node_food_id = nf.id
WHERE nf.node_id = $1
  AND DATE(fsl.created_at AT TIME ZONE '+09:00') = CURRENT_DATE
  AND DATE_PART('hour', fsl.created_at AT TIME ZONE '+09:00') BETWEEN 8 AND 18
GROUP BY hour, minute
ORDER BY hour DESC, minute DESC;
-- name: QuantityFoodStallPerHalfHourByNodeId :many
SELECT SUM(fsl.quantity * f.quantity) AS quantity,
  DATE_PART('hour', fsl.created_at AT TIME ZONE '+09:00') AS hour,
  FLOOR(DATE_PART('minute', fsl.created_at AT TIME ZONE '+09:00') / 30) * 30 AS minute
FROM food_stall_logs fsl
JOIN node_foods nf ON fsl.node_food_id = nf.id
JOIN foods f ON nf.food_id = f.id
WHERE nf.node_id = $1
  AND DATE(fsl.created_at AT TIME ZONE '+09:00') = CURRENT_DATE
  AND DATE_PART('hour', fsl.created_at AT TIME ZONE '+09:00') BETWEEN 8 AND 18
GROUP BY hour, minute
ORDER BY hour DESC, minute DESC;
-- name: CountExhibitionPerHalfHourByNodeId :many
SELECT COUNT(*) AS count,
  DATE_PART('hour', el.created_at AT TIME ZONE '+09:00') AS hour,
  FLOOR(DATE_PART('minute', el.created_at AT TIME ZONE '+09:00') / 30) * 30 AS minute
FROM exhibition_logs el
WHERE el.node_id = $1
  AND DATE(el.created_at AT TIME ZONE '+09:00') = CURRENT_DATE
  AND DATE_PART('hour', el.created_at AT TIME ZONE '+09:00') BETWEEN 8 AND 18
GROUP BY hour, minute
ORDER BY hour DESC, minute DESC;
