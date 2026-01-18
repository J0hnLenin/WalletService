package walletservice

import (
	"context"
	"testing"

	mocks "github.com/J0hnLenin/WalletService/internal/services/wallet_service/mocks"
	"github.com/stretchr/testify/suite"
)

type WalletServiceSuite struct {
	suite.Suite
	ctx           context.Context
	walletStorage *mocks.WalletStorage
	walletService *WalletService
}

func (w *WalletServiceSuite) SetupTest() {
	w.ctx = context.Background()
	w.walletStorage = mocks.NewWalletStorage(w.T())
	w.walletService = NewWalletService(w.ctx, w.walletStorage)
}

func TestWalletServiceSuite(t *testing.T) {
	suite.Run(t, new(WalletServiceSuite))
}