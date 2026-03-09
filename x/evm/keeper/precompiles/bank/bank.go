package bank

import (
	"errors"
	"math/big"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

// BankPrecompile provides ERC-20 compatible interface over Cosmos SDK x/bank module.
// Address: 0x0000000000000000000000000000000000000801
type BankPrecompile struct {
	bankKeeper    BankKeeper
	accountKeeper AccountKeeper
}

// BankKeeper defines the expected interface for the bank keeper.
type BankKeeper interface {
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}

// AccountKeeper defines the expected interface for the account keeper.
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) sdk.AccountI
}

var (
	// PrecompileAddress is the fixed address for the bank precompile
	PrecompileAddress = common.HexToAddress("0x0000000000000000000000000000000000000801")

	// WeiPerNote is the scaling factor: 1 note = 10^12 wei
	WeiPerNote = new(big.Int).Exp(big.NewInt(10), big.NewInt(12), nil)

	// ErrExecutionReverted is returned when precompile execution fails
	ErrExecutionReverted = errors.New("execution reverted")
)

// NewBankPrecompile creates a new bank precompile instance.
func NewBankPrecompile(bankKeeper BankKeeper, accountKeeper AccountKeeper) *BankPrecompile {
	return &BankPrecompile{
		bankKeeper:    bankKeeper,
		accountKeeper: accountKeeper,
	}
}

// Address returns the precompile address.
func (p *BankPrecompile) Address() common.Address {
	return PrecompileAddress
}

// RequiredGas returns the gas required to execute the precompile.
func (p *BankPrecompile) RequiredGas(input []byte) uint64 {
	if len(input) < 4 {
		return 0
	}
	methodID := input[:4]

	// balanceOf - read operation
	if isBalanceOfMethod(methodID) {
		return 2000
	}
	// transfer - write operation
	if isTransferMethod(methodID) {
		return 10000
	}

	return 5000
}

// Run executes the precompile logic with SDK context.
// Supported ERC-20 methods:
//   - balanceOf(address) returns (uint256)
//   - transfer(address to, uint256 amount) returns (bool)
func (p *BankPrecompile) Run(ctx sdk.Context, caller common.Address, input []byte, suppliedGas uint64, readOnly bool) ([]byte, uint64, error) {
	if len(input) < 4 {
		return nil, suppliedGas, ErrExecutionReverted
	}

	methodID := input[:4]
	gasRequired := p.RequiredGas(input)

	if gasRequired > suppliedGas {
		return nil, suppliedGas, errors.New("out of gas")
	}

	// balanceOf(address) - 0x70a08231
	if isBalanceOfMethod(methodID) {
		result, err := p.balanceOf(ctx, input[4:])
		if err != nil {
			return nil, suppliedGas, err
		}
		return result, suppliedGas - gasRequired, nil
	}

	// transfer(address,uint256) - 0xa9059cbb
	if isTransferMethod(methodID) {
		if readOnly {
			return nil, suppliedGas, errors.New("cannot transfer in read-only mode")
		}
		result, err := p.transfer(ctx, caller, input[4:])
		if err != nil {
			return nil, suppliedGas, err
		}
		return result, suppliedGas - gasRequired, nil
	}

	return nil, suppliedGas, ErrExecutionReverted
}

// balanceOf returns the note balance for a given address (in wei units).
func (p *BankPrecompile) balanceOf(ctx sdk.Context, input []byte) ([]byte, error) {
	if len(input) < 32 {
		return nil, ErrExecutionReverted
	}

	// Parse address from first 32 bytes (address is in last 20 bytes)
	addr := common.BytesToAddress(input[12:32])
	accAddr := sdk.AccAddress(addr.Bytes())

	// Get balance from bank keeper
	balance := p.bankKeeper.GetBalance(ctx, accAddr, "note")

	// Convert note to wei (multiply by 10^12)
	weiBalance := new(big.Int).Mul(balance.Amount.BigInt(), WeiPerNote)

	return encodeUint256(weiBalance), nil
}

// transfer sends tokens from caller to recipient.
// Input: to (address, 32 bytes) + amount (uint256, 32 bytes)
// Output: success (bool, 32 bytes)
func (p *BankPrecompile) transfer(ctx sdk.Context, caller common.Address, input []byte) ([]byte, error) {
	if len(input) < 64 {
		return nil, ErrExecutionReverted
	}

	// Parse recipient address (bytes 12-32 of first 32 bytes)
	toAddr := common.BytesToAddress(input[12:32])

	// Parse amount (second 32 bytes)
	weiAmount := new(big.Int).SetBytes(input[32:64])

	// Convert wei to note (divide by 10^12)
	noteAmount := new(big.Int).Div(weiAmount, WeiPerNote)

	// Create Cosmos addresses
	fromAddr := sdk.AccAddress(caller.Bytes())
	recipientAddr := sdk.AccAddress(toAddr.Bytes())

	// Transfer using bank keeper
	coins := sdk.NewCoins(sdk.NewCoin("note", math.NewIntFromBigInt(noteAmount)))
	err := p.bankKeeper.SendCoins(ctx, fromAddr, recipientAddr, coins)
	if err != nil {
		// Return false on failure
		return encodeBool(false), nil
	}

	// Return true on success
	return encodeBool(true), nil
}

// Method ID helpers
func isBalanceOfMethod(methodID []byte) bool {
	// balanceOf(address) = keccak256("balanceOf(address)")[:4] = 0x70a08231
	expected := []byte{0x70, 0xa0, 0x82, 0x31}
	return compareBytes(methodID, expected)
}

func isTransferMethod(methodID []byte) bool {
	// transfer(address,uint256) = keccak256("transfer(address,uint256)")[:4] = 0xa9059cbb
	expected := []byte{0xa9, 0x05, 0x9c, 0xbb}
	return compareBytes(methodID, expected)
}

func compareBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// ABI encoding helpers
func encodeUint256(value *big.Int) []byte {
	result := make([]byte, 32)
	if value != nil {
		value.FillBytes(result)
	}
	return result
}

func encodeBool(value bool) []byte {
	result := make([]byte, 32)
	if value {
		result[31] = 1
	}
	return result
}
