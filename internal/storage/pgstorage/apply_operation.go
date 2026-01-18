package pgstorage

import (
	"context"
	"errors"
	"fmt"

	custom_errors "github.com/J0hnLenin/WalletService/internal/errors"
	"github.com/J0hnLenin/WalletService/internal/models"
	"github.com/jackc/pgx/v5"
)

func (pg *PGStorage) ApplyOperation(ctx context.Context, op *models.WalletOperation) (newBalance int64, err error) {
    shardIndex, bucketIndex := pg.shardAndBucketByWalletID(op.WalletID)
    tableName := tableWithBucket(bucketIndex)

    // Сначала пытаемся создать запись с нулевым балансом, если ее нет
    insertQuery := fmt.Sprintf(`
        INSERT INTO %s (%s, %s, %s)
        VALUES ($1, 0, CURRENT_TIMESTAMP)
        ON CONFLICT (%s) DO NOTHING
    `,
        tableName,
        idColumnName, balanceColumnName, createdAtColumnName,
        idColumnName,
    )

    // Затем обновляем баланс с проверкой
    updateQuery := fmt.Sprintf(`
        UPDATE %s 
        SET %s = %s + $2,
            %s = CURRENT_TIMESTAMP
        WHERE %s = $1
        AND %s + $2 >= 0
        RETURNING %s
    `,
        tableName,
        balanceColumnName, balanceColumnName,
        updatedAtColumnName,
        idColumnName,
        balanceColumnName,
        balanceColumnName,
    )

    args := []interface{}{op.WalletID, op.AmountChange}

    transaction, err := pg.shards[shardIndex].db.Begin(ctx)
    if err != nil {
        return 0, fmt.Errorf("begin transaction error: %w", err)
    }
    
    defer func() {
        if err != nil {
            transaction.Rollback(ctx)
        }
    }()

    _, err = transaction.Exec(ctx, insertQuery, op.WalletID)
    if err != nil {
        return 0, fmt.Errorf("insert if new error: %w", err)
    }

    var balance int64
    err = transaction.QueryRow(ctx, updateQuery, args...).Scan(&balance)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            // Недостаточно средств
            return 0, &custom_errors.ErrInsufficientBalance{WalletID: op.WalletID}
        }
        return 0, fmt.Errorf("update balance error: %w", err)
    }

    // Фиксируем транзакцию
    if err = transaction.Commit(ctx); err != nil {
        return 0, fmt.Errorf("commit transaction error: %w", err)
    }

    return balance, nil
}