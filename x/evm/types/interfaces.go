package types

import (
	"context"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// AccountKeeper defines the expected account keeper interface for the EVM module.
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	SetAccount(ctx context.Context, account sdk.AccountI)
	NewAccountWithAddress(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	GetModuleAddress(moduleName string) sdk.AccAddress
	GetSequence(context.Context, sdk.AccAddress) (uint64, error)
}

// BankKeeper defines the expected bank keeper interface for the EVM module.
type BankKeeper interface {
	SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
}

// StakingKeeper defines the expected staking keeper interface for the EVM module.
// Used for accessing validator information in precompiles.
type StakingKeeper interface {
	GetValidator(ctx context.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, err error)
	GetAllValidators(ctx context.Context) (validators []stakingtypes.Validator, err error)
	GetDelegation(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (delegation stakingtypes.Delegation, err error)
}

// EvmHooks defines the interface for EVM hooks.
// Hooks allow other modules to react to EVM events.
type EvmHooks interface {
	// PostTxProcessing is called after an EVM transaction is successfully processed.
	// It receives the transaction sender, recipient, and whether the transaction succeeded.
	PostTxProcessing(ctx sdk.Context, msg *MsgEthereumTx, receipt *TxReceipt) error
}

// MsgEthereumTx is defined here as a placeholder for the interfaces file.
// The actual implementation will be in msg.go.
type MsgEthereumTx struct {
	// Data is the Ethereum transaction payload (RLP-encoded).
	Data []byte
	// Hash is the Ethereum transaction hash (keccak256).
	Hash string
	// From is the Ethereum address of the sender (hex string without 0x).
	From string
}

// TxReceipt represents an Ethereum transaction receipt.
type TxReceipt struct {
	// TxHash is the transaction hash.
	TxHash string
	// ContractAddress is the address of the created contract (if any).
	ContractAddress string
	// GasUsed is the amount of gas used by the transaction.
	GasUsed uint64
	// Status indicates whether the transaction succeeded (1) or failed (0).
	Status uint64
	// Logs contains the event logs emitted by the transaction.
	Logs []*Log
	// Bloom is the Bloom filter for the logs.
	Bloom []byte
}

// Log represents an Ethereum event log.
type Log struct {
	// Address is the contract address that emitted the log.
	Address string
	// Topics are the indexed event parameters.
	Topics []string
	// Data is the non-indexed event data.
	Data []byte
	// BlockNumber is the block number.
	BlockNumber uint64
	// TxHash is the transaction hash.
	TxHash string
	// TxIndex is the transaction index in the block.
	TxIndex uint64
	// BlockHash is the block hash.
	BlockHash string
	// Index is the log index in the transaction.
	Index uint64
}

// StateDB is a minimal interface for EVM state database operations.
// The full implementation will be in keeper/statedb.go.
type StateDB interface {
	CreateAccount(addr sdk.AccAddress)
	SubBalance(addr sdk.AccAddress, amount *big.Int)
	AddBalance(addr sdk.AccAddress, amount *big.Int)
	GetBalance(addr sdk.AccAddress) *big.Int
	GetNonce(addr sdk.AccAddress) uint64
	SetNonce(addr sdk.AccAddress, nonce uint64)
	GetCodeHash(addr sdk.AccAddress) []byte
	GetCode(addr sdk.AccAddress) []byte
	SetCode(addr sdk.AccAddress, code []byte)
	GetState(addr sdk.AccAddress, key []byte) []byte
	SetState(addr sdk.AccAddress, key, value []byte)
	Commit() error
}
