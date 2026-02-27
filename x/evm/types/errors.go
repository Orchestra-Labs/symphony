package types

import (
	errorsmod "cosmossdk.io/errors"
)

// EVM module sentinel errors
var (
	// ErrInvalidChainID is returned when the chain ID is invalid.
	ErrInvalidChainID = errorsmod.Register(ModuleName, 2, "invalid chain ID")

	// ErrInvalidState is returned when the EVM state is invalid.
	ErrInvalidState = errorsmod.Register(ModuleName, 3, "invalid EVM state")

	// ErrExecutionReverted is returned when EVM execution is reverted.
	ErrExecutionReverted = errorsmod.Register(ModuleName, 4, "execution reverted")

	// ErrInvalidTransaction is returned when the transaction is invalid.
	ErrInvalidTransaction = errorsmod.Register(ModuleName, 5, "invalid transaction")

	// ErrInvalidSignature is returned when the signature is invalid.
	ErrInvalidSignature = errorsmod.Register(ModuleName, 6, "invalid signature")

	// ErrGasOverflow is returned when gas computation overflows.
	ErrGasOverflow = errorsmod.Register(ModuleName, 7, "gas overflow")

	// ErrInvalidGasLimit is returned when the gas limit is invalid.
	ErrInvalidGasLimit = errorsmod.Register(ModuleName, 8, "invalid gas limit")

	// ErrInvalidGasPrice is returned when the gas price is invalid.
	ErrInvalidGasPrice = errorsmod.Register(ModuleName, 9, "invalid gas price")

	// ErrInvalidValue is returned when the transaction value is invalid.
	ErrInvalidValue = errorsmod.Register(ModuleName, 10, "invalid transaction value")

	// ErrInvalidNonce is returned when the nonce is invalid.
	ErrInvalidNonce = errorsmod.Register(ModuleName, 11, "invalid nonce")

	// ErrNonceTooLow is returned when the nonce is too low.
	ErrNonceTooLow = errorsmod.Register(ModuleName, 12, "nonce too low")

	// ErrNonceTooHigh is returned when the nonce is too high.
	ErrNonceTooHigh = errorsmod.Register(ModuleName, 13, "nonce too high")

	// ErrInsufficientFunds is returned when the account has insufficient funds.
	ErrInsufficientFunds = errorsmod.Register(ModuleName, 14, "insufficient funds")

	// ErrIntrinsicGas is returned when the intrinsic gas is too low.
	ErrIntrinsicGas = errorsmod.Register(ModuleName, 15, "intrinsic gas too low")

	// ErrGasLimitReached is returned when the gas limit is reached.
	ErrGasLimitReached = errorsmod.Register(ModuleName, 16, "gas limit reached")

	// ErrInvalidCode is returned when the contract code is invalid.
	ErrInvalidCode = errorsmod.Register(ModuleName, 17, "invalid code")

	// ErrCreateDisabled is returned when contract creation is disabled.
	ErrCreateDisabled = errorsmod.Register(ModuleName, 18, "contract creation disabled")

	// ErrCallDisabled is returned when contract calls are disabled.
	ErrCallDisabled = errorsmod.Register(ModuleName, 19, "contract calls disabled")

	// ErrInvalidAddress is returned when the address is invalid.
	ErrInvalidAddress = errorsmod.Register(ModuleName, 20, "invalid address")

	// ErrAccountNotFound is returned when the account is not found.
	ErrAccountNotFound = errorsmod.Register(ModuleName, 21, "account not found")

	// ErrCodeHashNotFound is returned when the code hash is not found.
	ErrCodeHashNotFound = errorsmod.Register(ModuleName, 22, "code hash not found")

	// ErrInvalidBlockContext is returned when the block context is invalid.
	ErrInvalidBlockContext = errorsmod.Register(ModuleName, 23, "invalid block context")
)
