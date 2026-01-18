package errors

import (
	"fmt"

	"github.com/google/uuid"
)

type ErrInsufficientBalance struct {
    WalletID uuid.UUID
}

func (e *ErrInsufficientBalance) Error() string {
    return fmt.Sprintf("not enough balance to withdraw from wallet %s", 
        e.WalletID.String())
}

type ErrWalletNotExists struct {
    WalletID uuid.UUID
}

func (e *ErrWalletNotExists) Error() string {
    return fmt.Sprintf("wallet %s not exists", 
        e.WalletID.String())
}

func (e *ErrInsufficientBalance) Is(tgt error) bool {
    _, ok := tgt.(*ErrInsufficientBalance)
    if !ok{
        return false
    }
    return true
}

func (e *ErrWalletNotExists) Is(tgt error) bool {
    _, ok := tgt.(*ErrWalletNotExists)
    if !ok{
        return false
    }
    return true
}