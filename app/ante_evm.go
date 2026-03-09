package app

import (
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"

	evmkeeper "github.com/osmosis-labs/osmosis/v27/x/evm/keeper"
	evmtypes "github.com/osmosis-labs/osmosis/v27/x/evm/types"
)

// EVMAnteHandler returns an AnteHandler for processing Ethereum transactions.
// This is a simplified version - full implementation would include:
// - EthSigVerificationDecorator: Verify Ethereum ECDSA signatures
// - EthAccountVerificationDecorator: Verify sender account exists
// - EthNonceVerificationDecorator: Verify nonce is correct
// - EthGasConsumeDecorator: Consume gas for the transaction
// - CanTransferDecorator: Verify sender has sufficient balance
// - EthIncrementSenderSequenceDecorator: Increment sender nonce
func NewEVMAnteHandler(
	evmKeeper *evmkeeper.Keeper,
	accountKeeper ante.AccountKeeper,
) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(), // Set up context
		NewEthSetupContextDecorator(evmKeeper),
		NewEthValidateBasicDecorator(evmKeeper),
		NewEthSigVerificationDecorator(evmKeeper),
		NewEthAccountVerificationDecorator(evmKeeper, accountKeeper),
		NewEthNonceVerificationDecorator(evmKeeper),
		NewEthGasConsumeDecorator(evmKeeper),
		NewCanTransferDecorator(evmKeeper),
		NewEthIncrementSenderSequenceDecorator(evmKeeper, accountKeeper),
	)
}

// IsEVMTx checks if a transaction contains an Ethereum transaction message.
func IsEVMTx(tx sdk.Tx) bool {
	for _, msg := range tx.GetMsgs() {
		if _, ok := msg.(*evmtypes.MsgEthereumTx); ok {
			return true
		}
	}
	return false
}

// EthSetupContextDecorator sets up the context for EVM transactions.
type EthSetupContextDecorator struct {
	evmKeeper *evmkeeper.Keeper
}

func NewEthSetupContextDecorator(evmKeeper *evmkeeper.Keeper) EthSetupContextDecorator {
	return EthSetupContextDecorator{evmKeeper: evmKeeper}
}

func (escd EthSetupContextDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// Set up EVM-specific context (block context, chain ID, etc.)
	// In a full implementation, this would configure the EVM environment
	return next(ctx, tx, simulate)
}

// EthValidateBasicDecorator validates basic Ethereum transaction fields.
type EthValidateBasicDecorator struct {
	evmKeeper *evmkeeper.Keeper
}

func NewEthValidateBasicDecorator(evmKeeper *evmkeeper.Keeper) EthValidateBasicDecorator {
	return EthValidateBasicDecorator{evmKeeper: evmKeeper}
}

func (vbd EthValidateBasicDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// Validate basic fields of Ethereum transactions
	for _, msg := range tx.GetMsgs() {
		if ethMsg, ok := msg.(*evmtypes.MsgEthereumTx); ok {
			if err := ethMsg.ValidateBasic(); err != nil {
				return ctx, err
			}
		}
	}
	return next(ctx, tx, simulate)
}

// EthSigVerificationDecorator verifies Ethereum ECDSA signatures.
type EthSigVerificationDecorator struct {
	evmKeeper *evmkeeper.Keeper
}

func NewEthSigVerificationDecorator(evmKeeper *evmkeeper.Keeper) EthSigVerificationDecorator {
	return EthSigVerificationDecorator{evmKeeper: evmKeeper}
}

func (svd EthSigVerificationDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// Verify ECDSA signature and recover sender address
	// In a full implementation, this would:
	// 1. Recover the sender address from the signature
	// 2. Verify the signature is valid
	// 3. Set the recovered address in the message

	// For now, we assume the signature has been verified
	// and the 'From' field is already set in MsgEthereumTx
	return next(ctx, tx, simulate)
}

// EthAccountVerificationDecorator verifies the sender account exists.
type EthAccountVerificationDecorator struct {
	evmKeeper     *evmkeeper.Keeper
	accountKeeper ante.AccountKeeper
}

func NewEthAccountVerificationDecorator(evmKeeper *evmkeeper.Keeper, accountKeeper ante.AccountKeeper) EthAccountVerificationDecorator {
	return EthAccountVerificationDecorator{
		evmKeeper:     evmKeeper,
		accountKeeper: accountKeeper,
	}
}

func (avd EthAccountVerificationDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	for _, msg := range tx.GetMsgs() {
		if ethMsg, ok := msg.(*evmtypes.MsgEthereumTx); ok {
			from := ethMsg.GetFrom()
			if len(from) == 0 {
				return ctx, sdkerrors.ErrInvalidAddress.Wrap("sender address not set")
			}

			// Verify account exists (will be created if it doesn't)
			_ = avd.evmKeeper.GetOrCreateAccount(ctx, sdk.AccAddress(from))
		}
	}
	return next(ctx, tx, simulate)
}

// EthNonceVerificationDecorator verifies the transaction nonce.
type EthNonceVerificationDecorator struct {
	evmKeeper *evmkeeper.Keeper
}

func NewEthNonceVerificationDecorator(evmKeeper *evmkeeper.Keeper) EthNonceVerificationDecorator {
	return EthNonceVerificationDecorator{evmKeeper: evmKeeper}
}

func (nvd EthNonceVerificationDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	for _, msg := range tx.GetMsgs() {
		if ethMsg, ok := msg.(*evmtypes.MsgEthereumTx); ok {
			from := sdk.AccAddress(ethMsg.GetFrom())
			currentNonce := nvd.evmKeeper.GetNonce(ctx, from)

			if ethMsg.Data.Nonce != currentNonce {
				if ethMsg.Data.Nonce < currentNonce {
					return ctx, evmtypes.ErrNonceTooLow
				}
				return ctx, evmtypes.ErrNonceTooHigh
			}
		}
	}
	return next(ctx, tx, simulate)
}

// EthGasConsumeDecorator consumes gas for Ethereum transactions.
type EthGasConsumeDecorator struct {
	evmKeeper *evmkeeper.Keeper
}

func NewEthGasConsumeDecorator(evmKeeper *evmkeeper.Keeper) EthGasConsumeDecorator {
	return EthGasConsumeDecorator{evmKeeper: evmKeeper}
}

func (gcd EthGasConsumeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// Set gas limit for the transaction
	// In a full implementation, this would calculate intrinsic gas
	for _, msg := range tx.GetMsgs() {
		if ethMsg, ok := msg.(*evmtypes.MsgEthereumTx); ok {
			// Set the gas limit from the Ethereum transaction
			ctx = ctx.WithGasMeter(storetypes.NewGasMeter(ethMsg.Data.Gas))
		}
	}
	return next(ctx, tx, simulate)
}

// CanTransferDecorator checks if the sender has sufficient balance.
type CanTransferDecorator struct {
	evmKeeper *evmkeeper.Keeper
}

func NewCanTransferDecorator(evmKeeper *evmkeeper.Keeper) CanTransferDecorator {
	return CanTransferDecorator{evmKeeper: evmKeeper}
}

func (ctd CanTransferDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	for _, msg := range tx.GetMsgs() {
		if ethMsg, ok := msg.(*evmtypes.MsgEthereumTx); ok {
			from := sdk.AccAddress(ethMsg.GetFrom())
			balance := ctd.evmKeeper.GetBalance(ctx, from)

			// Calculate total cost: value + (gas * gasPrice)
			totalCost := ethMsg.GetValue()
			gasCost := ethMsg.GetGasPrice()
			gasCost.Mul(gasCost, math.NewInt(int64(ethMsg.Data.Gas)).BigInt())
			totalCost.Add(totalCost, gasCost)

			if balance.Cmp(totalCost) < 0 {
				return ctx, evmtypes.ErrInsufficientFunds
			}
		}
	}
	return next(ctx, tx, simulate)
}

// EthIncrementSenderSequenceDecorator increments the sender's nonce.
type EthIncrementSenderSequenceDecorator struct {
	evmKeeper     *evmkeeper.Keeper
	accountKeeper ante.AccountKeeper
}

func NewEthIncrementSenderSequenceDecorator(evmKeeper *evmkeeper.Keeper, accountKeeper ante.AccountKeeper) EthIncrementSenderSequenceDecorator {
	return EthIncrementSenderSequenceDecorator{
		evmKeeper:     evmKeeper,
		accountKeeper: accountKeeper,
	}
}

func (issd EthIncrementSenderSequenceDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// Nonce will be incremented in the message handler after successful execution
	// This decorator is a placeholder for consistency with the decorator chain
	return next(ctx, tx, simulate)
}
