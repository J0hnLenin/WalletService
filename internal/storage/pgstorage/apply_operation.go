package pgstorage

import (
	"context"
	"fmt"

	"github.com/J0hnLenin/WalletService/internal/models"
	"github.com/Masterminds/squirrel"
)

func (pg *PGStorage) ApplyOperation(ctx context.Context, op *models.WalletOperation) error {

    shardIndex, bucketIndex := pg.shardAndBucketByWalletID(op.WalletID)    
    tableName := tableWithBucket(bucketIndex)
    
    upsertCTE := fmt.Sprintf(
        "WITH upsert_balance AS ( "+
            "INSERT INTO %s (%s, %s, %s) "+
            "VALUES ($1, GREATEST($2, 0), CURRENT_TIMESTAMP) "+
            "ON CONFLICT (%s) DO UPDATE "+
            "SET %s = %s.%s + EXCLUDED.%s, "+
            "%s = CURRENT_TIMESTAMP "+
            "WHERE %s.%s + EXCLUDED.%s >= 0 "+
            "RETURNING %s "+
        ") ",
        tableName, 
        idColumnName, amountColumnName, createdAtColumnName,
        idColumnName,
        amountColumnName, tableName, amountColumnName, amountColumnName,
        updatedAtColumnName,
        tableName, amountColumnName, amountColumnName,
        amountColumnName,
    )
    
    query := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
        Select(amountColumnName).
        Prefix(upsertCTE).
        From("upsert_balance").
        Suffix(fmt.Sprintf(
            "UNION ALL "+
            "SELECT NULL "+
            "WHERE NOT EXISTS (SELECT 1 FROM %s WHERE %s = $1) "+
            "AND $2 < 0",
            tableName, idColumnName,
        )).
        Prefix("BEGIN;").
        Suffix("; COMMIT;")
    
    queryText, args, err := query.ToSql()
    if err != nil {
        return fmt.Errorf("generate query error: %w", err)
    }
    
    args = []interface{}{op.WalletID, op.AmountChange}
    
    var balance *int64
    err = pg.shards[shardIndex].db.QueryRow(ctx, queryText, args...).Scan(&balance)
    if err != nil {
        return fmt.Errorf("execute operation error: %w", err)
    }
    
    // Если balance == nil, это означает попытку списания с несуществующего счёта
    if balance == nil && op.AmountChange < 0 {
        return fmt.Errorf("can't withdraw from non-existent wallet")
    }
    
    return nil
}