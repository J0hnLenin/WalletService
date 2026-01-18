package pgstorage

import "github.com/google/uuid"

type bucketNum uint16

type wallet struct {
	id        uuid.UUID `db:"id"`
	amount    int64     `db:"amount"`
}

const (
	bucketPrefix = "bucket_" // "bucket_n"
	tableName    = "wallets"

	idColumnName        = "id"
	amountColumnName = "game_proto"
	createdAtColumnName = "created_at"
	updatedAtColumnName = "updated_at"
)
