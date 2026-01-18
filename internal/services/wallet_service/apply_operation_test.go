package walletservice

import (
	"errors"

	custom_errors "github.com/J0hnLenin/WalletService/internal/errors"
	"github.com/J0hnLenin/WalletService/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestApplyOperationDepositSuccess проверяет успешное пополнение кошелька
func (w *WalletServiceSuite) TestApplyOperationDepositSuccess() {
	id := uuid.New()
	amountChange := int64(1000)
	newBalance := int64(1000)

	op := &models.WalletOperation{
		WalletID:     id,
		AmountChange: amountChange,
	}

	w.walletStorage.On("ApplyOperation", w.ctx, op).
		Return(newBalance, nil).
		Times(1)

	wallet, err := w.walletService.ApplyOperation(w.ctx, op)

	assert.Nil(w.T(), err)
	assert.Equal(w.T(), id, wallet.ID)
	assert.Equal(w.T(), newBalance, wallet.Balance)
}

// TestApplyOperationZeroAmount проверяет операцию с нулевым изменением баланса
func (w *WalletServiceSuite) TestApplyOperationZeroAmount() {
	id := uuid.New()
	op := &models.WalletOperation{
		WalletID:     id,
		AmountChange: 0,
	}

	existingWallet := &models.Wallet{
		ID:      id,
		Balance: 500,
	}

	// Метод должен вызвать GetWalletByID для получения текущего баланса
	w.walletStorage.On("GetWalletByID", w.ctx, id).
		Return(existingWallet, nil).
		Times(1)

	wallet, err := w.walletService.ApplyOperation(w.ctx, op)

	assert.Nil(w.T(), err)
	assert.Equal(w.T(), existingWallet, wallet)
}

// TestApplyOperationWithdrawalSuccess проверяет успешное списание средств
func (w *WalletServiceSuite) TestApplyOperationWithdrawalSuccess() {
	id := uuid.New()
	amountChange := int64(-500)
	newBalance := int64(500)

	op := &models.WalletOperation{
		WalletID:     id,
		AmountChange: amountChange,
	}

	existingWallet := &models.Wallet{
		ID:      id,
		Balance: 1000,
	}

	w.walletStorage.On("GetWalletByID", w.ctx, id).
		Return(existingWallet, nil).
		Times(1)

	w.walletStorage.On("ApplyOperation", w.ctx, op).
		Return(newBalance, nil).
		Times(1)

	wallet, err := w.walletService.ApplyOperation(w.ctx, op)

	assert.Nil(w.T(), err)
	assert.Equal(w.T(), id, wallet.ID)
	assert.Equal(w.T(), newBalance, wallet.Balance)
}

// TestApplyOperationWithdrawalWalletNotExists проверяет списание с несуществующего кошелька
func (w *WalletServiceSuite) TestApplyOperationWithdrawalWalletNotExists() {
	id := uuid.New()
	amountChange := int64(-500)

	op := &models.WalletOperation{
		WalletID:     id,
		AmountChange: amountChange,
	}

	w.walletStorage.On("GetWalletByID", w.ctx, id).
		Return(nil, &custom_errors.ErrWalletNotExists{}).
		Times(1)

	_, err := w.walletService.ApplyOperation(w.ctx, op)

	assert.Error(w.T(), err)
	assert.IsType(w.T(), &custom_errors.ErrInsufficientBalance{}, err)
}

// TestApplyOperationWithdrawalCheckWalletError проверяет ошибку при проверке существования кошелька
func (w *WalletServiceSuite) TestApplyOperationWithdrawalCheckWalletError() {
	wantErrorString := "Database connection error"
	wantErr := errors.New(wantErrorString)
	id := uuid.New()
	amountChange := int64(-500)

	op := &models.WalletOperation{
		WalletID:     id,
		AmountChange: amountChange,
	}

	w.walletStorage.On("GetWalletByID", w.ctx, id).
		Return(nil, wantErr).
		Times(1)

	_, err := w.walletService.ApplyOperation(w.ctx, op)

	assert.EqualError(w.T(), err, "can't check wallet: get wallet with id '"+id.String()+"' storage error: "+wantErrorString)
}

// TestApplyOperationInsufficientBalance проверяет ошибку недостаточного баланса
func (w *WalletServiceSuite) TestApplyOperationInsufficientBalance() {
	id := uuid.New()
	amountChange := int64(-1500)

	op := &models.WalletOperation{
		WalletID:     id,
		AmountChange: amountChange,
	}

	existingWallet := &models.Wallet{
		ID:      id,
		Balance: 1000,
	}

	w.walletStorage.On("GetWalletByID", w.ctx, id).
		Return(existingWallet, nil).
		Times(1)

	w.walletStorage.On("ApplyOperation", w.ctx, op).
		Return(int64(0), &custom_errors.ErrInsufficientBalance{WalletID: id}).
		Times(1)

	_, err := w.walletService.ApplyOperation(w.ctx, op)

	assert.Error(w.T(), err)
	assert.IsType(w.T(), &custom_errors.ErrInsufficientBalance{}, err)
}

// TestApplyOperationStorageError проверяет ошибку хранилища при выполнении операции
func (w *WalletServiceSuite) TestApplyOperationStorageError() {
	wantErrorString := "Transaction failed"
	wantErr := errors.New(wantErrorString)
	id := uuid.New()
	amountChange := int64(1000)

	op := &models.WalletOperation{
		WalletID:     id,
		AmountChange: amountChange,
	}

	w.walletStorage.On("ApplyOperation", w.ctx, op).
		Return(int64(0), wantErr).
		Times(1)

	_, err := w.walletService.ApplyOperation(w.ctx, op)

	assert.EqualError(w.T(), err, "storage error while apply operation: "+wantErrorString)
}

// TestApplyOperationDepositForNewWallet проверяет пополнение для нового кошелька
func (w *WalletServiceSuite) TestApplyOperationDepositForNewWallet() {
	id := uuid.New()
	amountChange := int64(1000)
	newBalance := int64(1000)

	op := &models.WalletOperation{
		WalletID:     id,
		AmountChange: amountChange,
	}

	// Для операции пополнения не должно быть вызова проверки существования кошелька
	w.walletStorage.On("ApplyOperation", w.ctx, op).
		Return(newBalance, nil).
		Times(1)

	wallet, err := w.walletService.ApplyOperation(w.ctx, op)

	assert.Nil(w.T(), err)
	assert.Equal(w.T(), id, wallet.ID)
	assert.Equal(w.T(), newBalance, wallet.Balance)
}

// TestWalletExistSuccess проверяет успешную проверку существования кошелька
func (w *WalletServiceSuite) TestWalletExistSuccess() {
	id := uuid.New()
	existingWallet := &models.Wallet{
		ID:      id,
		Balance: 1000,
	}

	w.walletStorage.On("GetWalletByID", w.ctx, id).
		Return(existingWallet, nil).
		Times(1)

	exists, err := w.walletService.walletExist(w.ctx, id)

	assert.Nil(w.T(), err)
	assert.True(w.T(), exists)
}

// TestWalletExistNotExists проверяет проверку несуществующего кошелька
func (w *WalletServiceSuite) TestWalletExistNotExists() {
	id := uuid.New()

	w.walletStorage.On("GetWalletByID", w.ctx, id).
		Return(nil, &custom_errors.ErrWalletNotExists{}).
		Times(1)

	exists, err := w.walletService.walletExist(w.ctx, id)

	assert.Nil(w.T(), err)
	assert.False(w.T(), exists)
}

// TestWalletExistError проверяет ошибку при проверке существования кошелька
func (w *WalletServiceSuite) TestWalletExistError() {
	wantErrorString := "Database connection error"
	wantErr := errors.New(wantErrorString)
	id := uuid.New()

	w.walletStorage.On("GetWalletByID", w.ctx, id).
		Return(nil, wantErr).
		Times(1)

	exists, err := w.walletService.walletExist(w.ctx, id)

	assert.EqualError(w.T(), err, "get wallet with id '"+id.String()+"' storage error: "+wantErrorString)
	assert.False(w.T(), exists)
}

// TestApplyOperationConcurrentDeposit проверяет корректность обработки нескольких пополнений
func (w *WalletServiceSuite) TestApplyOperationConcurrentDeposit() {
	id := uuid.New()
	op1 := &models.WalletOperation{
		WalletID:     id,
		AmountChange: 1000,
	}
	op2 := &models.WalletOperation{
		WalletID:     id,
		AmountChange: 500,
	}

	// Первое пополнение
	w.walletStorage.On("ApplyOperation", w.ctx, op1).
		Return(int64(1000), nil).
		Times(1)

	// Второе пополнение
	w.walletStorage.On("ApplyOperation", w.ctx, op2).
		Return(int64(1500), nil).
		Times(1)

	wallet1, err1 := w.walletService.ApplyOperation(w.ctx, op1)
	wallet2, err2 := w.walletService.ApplyOperation(w.ctx, op2)

	assert.Nil(w.T(), err1)
	assert.Nil(w.T(), err2)
	assert.Equal(w.T(), int64(1000), wallet1.Balance)
	assert.Equal(w.T(), int64(1500), wallet2.Balance)
}

// TestApplyOperationComplexScenario проверяет комплексный сценарий
func (w *WalletServiceSuite) TestApplyOperationComplexScenario() {
	id := uuid.New()
	
	// 1. Пополнение нового кошелька
	depositOp := &models.WalletOperation{
		WalletID:     id,
		AmountChange: 2000,
	}

	// 2. Списание части средств
	withdrawalOp := &models.WalletOperation{
		WalletID:     id,
		AmountChange: -1500,
	}

	// 3. Еще одно списание (должно вызвать ошибку)
	failedWithdrawalOp := &models.WalletOperation{
		WalletID:     id,
		AmountChange: -1000,
	}

	existingWallet := &models.Wallet{
		ID:      id,
		Balance: 500, // после первого списания
	}

	// Пополнение
	w.walletStorage.On("ApplyOperation", w.ctx, depositOp).
		Return(int64(2000), nil).
		Times(1)

	// Проверка существования для списания
	w.walletStorage.On("GetWalletByID", w.ctx, id).
		Return(&models.Wallet{ID: id, Balance: 2000}, nil).
		Times(1)

	// Первое успешное списание
	w.walletStorage.On("ApplyOperation", w.ctx, withdrawalOp).
		Return(int64(500), nil).
		Times(1)

	// Проверка существования для второго списания
	w.walletStorage.On("GetWalletByID", w.ctx, id).
		Return(existingWallet, nil).
		Times(1)

	// Второе списание (неуспешное)
	w.walletStorage.On("ApplyOperation", w.ctx, failedWithdrawalOp).
		Return(int64(0), &custom_errors.ErrInsufficientBalance{WalletID: id}).
		Times(1)

	// Выполняем операции
	wallet1, err1 := w.walletService.ApplyOperation(w.ctx, depositOp)
	wallet2, err2 := w.walletService.ApplyOperation(w.ctx, withdrawalOp)
	_, err3 := w.walletService.ApplyOperation(w.ctx, failedWithdrawalOp)

	// Проверяем результаты
	assert.Nil(w.T(), err1)
	assert.Equal(w.T(), int64(2000), wallet1.Balance)

	assert.Nil(w.T(), err2)
	assert.Equal(w.T(), int64(500), wallet2.Balance)

	assert.Error(w.T(), err3)
	assert.IsType(w.T(), &custom_errors.ErrInsufficientBalance{}, err3)
}