package pgstorage

import (
	"context"
	"fmt"

	"github.com/J0hnLenin/WalletService/internal/models"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (pg *PGStorage) GetWalletByID(ctx context.Context, id uuid.UUID) (*models.Wallet, error) {
	
	shardIndex, bucketIndex := pg.shardAndBucketByWalletID(id)
	query := squirrel.Select(amountColumnName).
		From(tableWithBacket(bucketIndex)).
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
	err = res.Scan(&wallet.Amount);
	if  err != nil {
		return nil, fmt.Errorf("query row error: %w", err)
	}
	
	return wallet, nil
}