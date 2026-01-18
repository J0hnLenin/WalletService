package models

import (
	"github.com/google/uuid"
)

type Wallet struct {
	ID uuid.UUID 
	Balance int64
}

type WalletOperation struct {
	WalletID uuid.UUID
	AmountChange int64
}