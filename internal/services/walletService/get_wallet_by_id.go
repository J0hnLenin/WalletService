package walletservice

import (
	"context"
	"errors"
	"fmt"

	custom_errors "github.com/J0hnLenin/WalletService/internal/errors"
	"github.com/J0hnLenin/WalletService/internal/models"
	"github.com/google/uuid"
)

func (w *WalletService) GetWalletByID(ctx context.Context, id uuid.UUID) (*models.Wallet, error) {
	wallet, err := w.walletStorage.GetWalletByID(ctx, id)

	if err != nil {
		if errors.Is(err, &custom_errors.ErrWalletNotExists{}) {
			// Если в базе данных ничего не найдено, значит по кошельку не было операций, 
			// либо все все операции были не успешны.
			// Считаем что у такого кошелька всегда нулевой баланс.

			wallet = &models.Wallet{
				ID: id,
				Balance: 0,
			}
			return wallet, nil
		}
		
		return nil, fmt.Errorf("get wallet with id '%s' storage error: %w", id.String(), err)
	}

	return wallet, nil
}