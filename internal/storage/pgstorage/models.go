package pgstorage

type bucketNum uint16

const (
	bucketPrefix = "bucket_" // "bucket_n"
	tableName    = "wallets"

	idColumnName        = "id"
	balanceColumnName   = "balance"
	createdAtColumnName = "created_at"
	updatedAtColumnName = "updated_at"
)
