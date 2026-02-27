package staking

import (
	"errors"
	"math/big"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
)

// StakingPrecompile provides staking operations from EVM smart contracts.
// Address: 0x0000000000000000000000000000000000000802
type StakingPrecompile struct {
	stakingKeeper StakingKeeper
	bankKeeper    BankKeeper
}

// StakingKeeper defines the expected interface for the staking keeper.
type StakingKeeper interface {
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (stakingtypes.Validator, error)
	Delegate(ctx sdk.Context, delAddr sdk.AccAddress, bondAmt math.Int, tokenSrc stakingtypes.BondStatus, validator stakingtypes.Validator, subtractAccount bool) (newShares math.LegacyDec, err error)
	GetDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.Delegation, error)
	Undelegate(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, sharesAmount math.LegacyDec) (completionTime time.Time, unbondedAmount math.Int, err error)
}

// BankKeeper defines the expected interface for the bank keeper.
type BankKeeper interface {
	DelegateCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

var (
	// PrecompileAddress is the fixed address for the staking precompile
	PrecompileAddress = common.HexToAddress("0x0000000000000000000000000000000000000802")

	// WeiPerNote is the scaling factor: 1 note = 10^12 wei
	WeiPerNote = new(big.Int).Exp(big.NewInt(10), big.NewInt(12), nil)

	// ErrExecutionReverted is returned when precompile execution fails
	ErrExecutionReverted = errors.New("execution reverted")
)

// NewStakingPrecompile creates a new staking precompile instance.
func NewStakingPrecompile(stakingKeeper StakingKeeper, bankKeeper BankKeeper) *StakingPrecompile {
	return &StakingPrecompile{
		stakingKeeper: stakingKeeper,
		bankKeeper:    bankKeeper,
	}
}

// Address returns the precompile address.
func (p *StakingPrecompile) Address() common.Address {
	return PrecompileAddress
}

// RequiredGas returns the gas required to execute the precompile.
func (p *StakingPrecompile) RequiredGas(input []byte) uint64 {
	if len(input) < 4 {
		return 0
	}
	methodID := input[:4]

	// delegate, undelegate - expensive operations
	if isDelegateMethod(methodID) || isUndelegateMethod(methodID) {
		return 50000
	}
	// getDelegation - read operation
	if isGetDelegationMethod(methodID) {
		return 5000
	}

	return 10000
}

// Run executes the precompile logic with SDK context.
// Supported methods:
//   - delegate(string validatorAddr, uint256 amount) returns (bool)
//   - getDelegation(address delegator, string validatorAddr) returns (uint256)
//   - undelegate(string validatorAddr, uint256 amount) returns (bool)
func (p *StakingPrecompile) Run(ctx sdk.Context, caller common.Address, input []byte, suppliedGas uint64, readOnly bool) ([]byte, uint64, error) {
	if len(input) < 4 {
		return nil, suppliedGas, ErrExecutionReverted
	}

	methodID := input[:4]
	gasRequired := p.RequiredGas(input)

	if gasRequired > suppliedGas {
		return nil, suppliedGas, errors.New("out of gas")
	}

	// delegate(string,uint256) - 0x026e402b
	if isDelegateMethod(methodID) {
		if readOnly {
			return nil, suppliedGas, errors.New("cannot delegate in read-only mode")
		}
		result, err := p.delegate(ctx, caller, input[4:])
		if err != nil {
			return nil, suppliedGas, err
		}
		return result, suppliedGas - gasRequired, nil
	}

	// getDelegation(address,string) - 0x4d99dd16
	if isGetDelegationMethod(methodID) {
		result, err := p.getDelegation(ctx, input[4:])
		if err != nil {
			return nil, suppliedGas, err
		}
		return result, suppliedGas - gasRequired, nil
	}

	// undelegate(string,uint256) - 0x4f498b8a
	if isUndelegateMethod(methodID) {
		if readOnly {
			return nil, suppliedGas, errors.New("cannot undelegate in read-only mode")
		}
		result, err := p.undelegate(ctx, caller, input[4:])
		if err != nil {
			return nil, suppliedGas, err
		}
		return result, suppliedGas - gasRequired, nil
	}

	return nil, suppliedGas, ErrExecutionReverted
}

// delegate delegates tokens to a validator.
// Input: validatorAddr (string) + amount (uint256)
// Output: success (bool)
func (p *StakingPrecompile) delegate(ctx sdk.Context, caller common.Address, input []byte) ([]byte, error) {
	// Parse validator address string and amount
	validatorAddr, amount, err := parseStringAndUint256(input)
	if err != nil {
		return nil, ErrExecutionReverted
	}

	// Convert caller to SDK address
	delAddr := sdk.AccAddress(caller.Bytes())

	// Parse validator address
	valAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return encodeBool(false), nil
	}

	// Get validator
	validator, err := p.stakingKeeper.GetValidator(ctx, valAddr)
	if err != nil {
		return encodeBool(false), nil
	}

	// Convert wei to note
	noteAmount := new(big.Int).Div(amount, WeiPerNote)
	bondAmt := math.NewIntFromBigInt(noteAmount)

	// Perform delegation
	_, err = p.stakingKeeper.Delegate(ctx, delAddr, bondAmt, stakingtypes.Unbonded, validator, true)
	if err != nil {
		return encodeBool(false), nil
	}

	return encodeBool(true), nil
}

// getDelegation returns the delegation amount for a delegator-validator pair.
// Input: delegator (address) + validatorAddr (string)
// Output: amount (uint256)
func (p *StakingPrecompile) getDelegation(ctx sdk.Context, input []byte) ([]byte, error) {
	if len(input) < 32 {
		return nil, ErrExecutionReverted
	}

	// Parse delegator address (first 32 bytes)
	delAddr := sdk.AccAddress(common.BytesToAddress(input[12:32]).Bytes())

	// Parse validator address string (remaining bytes)
	validatorAddr, err := parseStringArg(input[32:])
	if err != nil {
		return nil, ErrExecutionReverted
	}

	valAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return encodeUint256(big.NewInt(0)), nil
	}

	// Get delegation
	delegation, err := p.stakingKeeper.GetDelegation(ctx, delAddr, valAddr)
	if err != nil {
		return encodeUint256(big.NewInt(0)), nil
	}

	// Convert shares to tokens (simplified - assumes 1:1 ratio)
	noteAmount := delegation.Shares.TruncateInt().BigInt()
	weiAmount := new(big.Int).Mul(noteAmount, WeiPerNote)

	return encodeUint256(weiAmount), nil
}

// undelegate undelegates tokens from a validator.
// Input: validatorAddr (string) + amount (uint256)
// Output: success (bool)
func (p *StakingPrecompile) undelegate(ctx sdk.Context, caller common.Address, input []byte) ([]byte, error) {
	// Parse validator address string and amount
	validatorAddr, amount, err := parseStringAndUint256(input)
	if err != nil {
		return nil, ErrExecutionReverted
	}

	delAddr := sdk.AccAddress(caller.Bytes())
	valAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return encodeBool(false), nil
	}

	// Convert wei to note shares
	noteAmount := new(big.Int).Div(amount, WeiPerNote)
	shares := math.LegacyNewDecFromBigInt(noteAmount)

	// Perform undelegation
	_, _, err = p.stakingKeeper.Undelegate(ctx, delAddr, valAddr, shares)
	if err != nil {
		return encodeBool(false), nil
	}

	return encodeBool(true), nil
}

// Method ID helpers
func isDelegateMethod(methodID []byte) bool {
	// delegate(string,uint256) = keccak256("delegate(string,uint256)")[:4] = 0x026e402b
	expected := []byte{0x02, 0x6e, 0x40, 0x2b}
	return compareBytes(methodID, expected)
}

func isGetDelegationMethod(methodID []byte) bool {
	// getDelegation(address,string) = keccak256("getDelegation(address,string)")[:4] = 0x4d99dd16
	expected := []byte{0x4d, 0x99, 0xdd, 0x16}
	return compareBytes(methodID, expected)
}

func isUndelegateMethod(methodID []byte) bool {
	// undelegate(string,uint256) = keccak256("undelegate(string,uint256)")[:4] = 0x4f498b8a
	expected := []byte{0x4f, 0x49, 0x8b, 0x8a}
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

// parseStringAndUint256 parses a string followed by uint256 from ABI-encoded data.
func parseStringAndUint256(data []byte) (string, *big.Int, error) {
	if len(data) < 64 {
		return "", nil, ErrExecutionReverted
	}

	// Read string offset (first 32 bytes)
	strOffset := new(big.Int).SetBytes(data[0:32]).Uint64()

	// Read uint256 amount (second 32 bytes)
	amount := new(big.Int).SetBytes(data[32:64])

	// Parse string at offset
	if uint64(len(data)) < strOffset+32 {
		return "", nil, ErrExecutionReverted
	}

	strLength := new(big.Int).SetBytes(data[strOffset : strOffset+32]).Uint64()
	if uint64(len(data)) < strOffset+32+strLength {
		return "", nil, ErrExecutionReverted
	}

	strData := data[strOffset+32 : strOffset+32+strLength]

	return string(strData), amount, nil
}

// parseStringArg parses an ABI-encoded string argument.
func parseStringArg(data []byte) (string, error) {
	if len(data) < 64 {
		return "", ErrExecutionReverted
	}

	// Read offset (first 32 bytes)
	offset := new(big.Int).SetBytes(data[0:32]).Uint64()

	// If offset is 0 or 32, string data starts at byte 32
	if offset <= 32 {
		offset = 0
		// Read length from bytes 32-64
		length := new(big.Int).SetBytes(data[32:64]).Uint64()
		if uint64(len(data)) < 64+length {
			return "", ErrExecutionReverted
		}
		return string(data[64 : 64+length]), nil
	}

	// Otherwise read at the offset
	if uint64(len(data)) < offset+32 {
		return "", ErrExecutionReverted
	}
	length := new(big.Int).SetBytes(data[offset : offset+32]).Uint64()
	if uint64(len(data)) < offset+32+length {
		return "", ErrExecutionReverted
	}
	return string(data[offset+32 : offset+32+length]), nil
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
