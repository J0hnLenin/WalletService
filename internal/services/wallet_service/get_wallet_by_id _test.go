package walletservice

import (
	"errors"

	custom_errors "github.com/J0hnLenin/WalletService/internal/errors"
	"github.com/J0hnLenin/WalletService/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestGetWalletByIDSuccess проверяет успешное получение кошелька из хранилища
func (w *WalletServiceSuite) TestGetWalletByIDSuccess() {
	id := uuid.New()
	wantWallet := &models.Wallet{
		ID:      id,
		Balance: 1000,
	}

	w.walletStorage.On("GetWalletByID", w.ctx, id).
		Return(wantWallet, nil).
		Times(1)

	wallet, err := w.walletService.GetWalletByID(w.ctx, id)

	assert.Nil(w.T(), err)
	assert.Equal(w.T(), wantWallet, wallet)
}

// TestGetWalletByIDWalletNotExists проверяет получение несуществующего кошелька
func (w *WalletServiceSuite) TestGetWalletByIDWalletNotExists() {
	id := uuid.New()

	w.walletStorage.On("GetWalletByID", w.ctx, id).
		Return(nil, &custom_errors.ErrWalletNotExists{}).
		Times(1)

	wallet, err := w.walletService.GetWalletByID(w.ctx, id)

	assert.Nil(w.T(), err)
	assert.Equal(w.T(), id, wallet.ID)
	assert.Equal(w.T(), int64(0), wallet.Balance)
}

// TestGetWalletByIDStorageError проверяет ошибку хранилища при получении кошелька
func (w *WalletServiceSuite) TestGetWalletByIDStorageError() {
	wantErrorString := "Database connection error"
	wantErr := errors.New(wantErrorString)
	id := uuid.New()

	w.walletStorage.On("GetWalletByID", w.ctx, id).
		Return(nil, wantErr).
		Times(1)

	_, err := w.walletService.GetWalletByID(w.ctx, id)

	assert.EqualError(w.T(), err, "get wallet with id '"+id.String()+"' storage error: "+wantErrorString)
}
