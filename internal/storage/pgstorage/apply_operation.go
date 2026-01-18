package pgstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	custom_errors "github.com/J0hnLenin/WalletService/internal/errors"
	"github.com/J0hnLenin/WalletService/internal/models"
	"github.com/jackc/pgx/v5"
)

func (pg *PGStorage) ApplyOperation(ctx context.Context, op *models.WalletOperation) (newBalance int64, err error) {
    shardIndex, bucketIndex := pg.shardAndBucketByWalletID(op.WalletID)
    tableName := tableWithBucket(bucketIndex)

    upsert := fmt.Sprintf(
        "WITH upsert_balance AS ( "+
            "INSERT INTO %s (%s, %s, %s) "+
            "VALUES ($1, GREATEST($2, 0), CURRENT_TIMESTAMP) "+
            "ON CONFLICT (%s) DO UPDATE "+
            "SET %s = %s.%s + EXCLUDED.%s, "+
            "%s = CURRENT_TIMESTAMP "+
            "WHERE %s.%s + EXCLUDED.%s >= 0 "+
            "RETURNING %s "+
        ") "+
        "SELECT %s FROM upsert_balance",
        tableName,
        idColumnName, balanceColumnName, createdAtColumnName,
        idColumnName,
        balanceColumnName, tableName, balanceColumnName, balanceColumnName,
        updatedAtColumnName,
        tableName, balanceColumnName, balanceColumnName,
        balanceColumnName,
        balanceColumnName,
    )

    args := []interface{}{op.WalletID, op.AmountChange}

    transactionOptions :=  pgx.TxOptions{
        IsoLevel:   pgx.Serializable,
        AccessMode: pgx.ReadWrite,
    }

    transaction, err := pg.shards[shardIndex].db.BeginTx(ctx, transactionOptions)
    if err != nil {
        return 0, fmt.Errorf("begin transaction error: %w", err)
    }
    
    defer func() {
        if err != nil {
            // Если во время выполнения запроса возникла ошибка, 
            // то откатываем транзакцию 
            transaction.Rollback(ctx)
        }
    }()

    var balance int64
    err = transaction.QueryRow(ctx, upsert, args...).Scan(&balance)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return 0, &custom_errors.ErrInsufficientBalance{WalletID: op.WalletID}
        }
        return 0, fmt.Errorf("execute operation error: %w", err)
    }

    // Фиксируем транзакцию
    if err = transaction.Commit(ctx); err != nil {
        return 0, fmt.Errorf("commit transaction error: %w", err)
    }

    return balance, nil
}