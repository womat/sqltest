
SELECT * FROM solar.cost
WHERE CostId = $1

-- name: ListAuthors :many
SELECT * FROM solar.Telesol
WHERE TeleSolId = $1
ORDER BY Date;
