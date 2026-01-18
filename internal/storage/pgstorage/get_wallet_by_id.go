package pgstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	custom_errors "github.com/J0hnLenin/WalletService/internal/errors"
	"github.com/J0hnLenin/WalletService/internal/models"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (pg *PGStorage) GetWalletByID(ctx context.Context, id uuid.UUID) (*models.Wallet, error) {
	
	shardIndex, bucketIndex := pg.shardAndBucketByWalletID(id)
	query := squirrel.Select(balanceColumnName).
		From(tableWithBucket(bucketIndex)).
		Where(squirrel.Eq{idColumnName: id}).
		PlaceholderFormat(squirrel.Dollar)

	queryText, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("generate query error: %w", err)
	}

	wallet := &models.Wallet{
		ID: id,
	}
	
	res := pg.shards[shardIndex].db.QueryRow(ctx, queryText, args...)
	err = res.Scan(&wallet.Balance);

	if  err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &custom_errors.ErrWalletNotExists{}
		}

		return nil, fmt.Errorf("query row error: %w", err)
	}

	return wallet, nil
}