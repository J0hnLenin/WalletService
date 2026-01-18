package walletservice

import (
	"context"
	"errors"
	"fmt"

	custom_errors "github.com/J0hnLenin/WalletService/internal/errors"
	"github.com/J0hnLenin/WalletService/internal/models"
	"github.com/google/uuid"
)

func (w *WalletService) ApplyOperation(ctx context.Context, op *models.WalletOperation) (*models.Wallet, error) {
	if op.AmountChange == 0 {
		// Если сумма операции равна нулю, то возвращаем текущий баланс
		return w.GetWalletByID(ctx, op.WalletID)
	}

	if (op.AmountChange < 0) {
		// Для операции списания проверяем существование кошелька
		// Для операции зачисления проверка не требуется 
		
		walletExist, err := w.walletExist(ctx, op.WalletID)
		if err != nil {
			return nil, fmt.Errorf("can't check wallet: %w", err)
		}

		if !walletExist {
			// Если кошелёк не существует, то всегда считаем, что у него нулевой баланс
			// т.к. по нему не было операций или все операции были неуспешными
			return nil, &custom_errors.ErrInsufficientBalance{WalletID: op.WalletID}
		}
	}

	newBalance, err := w.walletStorage.ApplyOperation(ctx, op)
	
	if err != nil {
		if errors.Is(err, &custom_errors.ErrInsufficientBalance{}) {
			return nil, err
		}

		return nil, fmt.Errorf("storage error while apply operation: %w", err)
	}

	// Если ошибок нет, то возвращаем текущий баланс 
	// после выполнения транзакции
	wallet := &models.Wallet{
		ID: op.WalletID,
		Balance: newBalance,
	}
	return wallet, nil
}

func (w *WalletService) walletExist(ctx context.Context, id uuid.UUID) (bool, error) {
	_, err := w.walletStorage.GetWalletByID(ctx, id)
	if err != nil {
		if errors.Is(err, &custom_errors.ErrWalletNotExists{}) {
			return false, nil
		}
		return false, fmt.Errorf("get wallet with id '%s' storage error: %w", id.String(), err)
	}

	return true, nil
}