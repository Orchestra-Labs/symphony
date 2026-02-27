package oracle

import (
	"errors"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// OraclePrecompile provides access to Symphony's oracle price feeds from EVM smart contracts.
// Address: 0x0000000000000000000000000000000000000800
type OraclePrecompile struct {
	oracleKeeper OracleKeeper
}

// OracleKeeper defines the expected interface for the oracle keeper.
type OracleKeeper interface {
	GetMelodyExchangeRate(ctx sdk.Context, denom string) (osmomath.Dec, error)
}

var (
	// PrecompileAddress is the fixed address for the oracle precompile
	PrecompileAddress = common.HexToAddress("0x0000000000000000000000000000000000000800")

	// ErrExecutionReverted is returned when precompile execution fails
	ErrExecutionReverted = errors.New("execution reverted")
)

// NewOraclePrecompile creates a new oracle precompile instance.
func NewOraclePrecompile(oracleKeeper OracleKeeper) *OraclePrecompile {
	return &OraclePrecompile{
		oracleKeeper: oracleKeeper,
	}
}

// Address returns the precompile address.
func (p *OraclePrecompile) Address() common.Address {
	return PrecompileAddress
}

// RequiredGas returns the gas required to execute the precompile.
func (p *OraclePrecompile) RequiredGas(input []byte) uint64 {
	// Base cost for oracle read operation
	return 5000
}

// Run executes the precompile logic with SDK context.
// Supported methods:
//   - getPrice(string denom) returns (uint256 price)
//     Method signature: 0x41976e09
func (p *OraclePrecompile) Run(ctx sdk.Context, caller common.Address, input []byte, suppliedGas uint64, readOnly bool) ([]byte, uint64, error) {
	// Check minimum input length (4 bytes for method selector)
	if len(input) < 4 {
		return nil, suppliedGas, ErrExecutionReverted
	}

	// Parse method selector (first 4 bytes)
	methodID := input[:4]

	// getPrice(string) - method ID: 0x41976e09
	if isGetPriceMethod(methodID) {
		result, err := p.getPrice(ctx, input[4:])
		if err != nil {
			return nil, suppliedGas, err
		}

		gasUsed := p.RequiredGas(input)
		if gasUsed > suppliedGas {
			return nil, suppliedGas, errors.New("out of gas")
		}

		return result, suppliedGas - gasUsed, nil
	}

	return nil, suppliedGas, ErrExecutionReverted
}

// getPrice retrieves the price for a given denomination.
// Input: ABI-encoded string (denom)
// Output: ABI-encoded uint256 (price in wei units, scaled by 10^18)
func (p *OraclePrecompile) getPrice(ctx sdk.Context, input []byte) ([]byte, error) {
	// Parse ABI-encoded string parameter
	denom, err := parseStringArg(input)
	if err != nil {
		return nil, ErrExecutionReverted
	}

	// Get price from oracle keeper (exchange rate)
	exchangeRate, err := p.oracleKeeper.GetMelodyExchangeRate(ctx, denom)
	if err != nil {
		// Return 0 if price not found
		return encodeUint256(big.NewInt(0)), nil
	}

	// Convert osmomath.Dec to big.Int (scale by 10^18 for EVM wei units)
	// osmomath.Dec has 18 decimal places internally
	priceInt := exchangeRate.BigInt()

	return encodeUint256(priceInt), nil
}

// isGetPriceMethod checks if the method ID matches getPrice(string)
func isGetPriceMethod(methodID []byte) bool {
	// getPrice(string) = keccak256("getPrice(string)")[:4] = 0x41976e09
	expected := []byte{0x41, 0x97, 0x6e, 0x09}
	if len(methodID) != 4 {
		return false
	}
	for i := 0; i < 4; i++ {
		if methodID[i] != expected[i] {
			return false
		}
	}
	return true
}

// parseStringArg parses an ABI-encoded string argument.
// ABI encoding: offset(32) + length(32) + data(padded to 32)
func parseStringArg(data []byte) (string, error) {
	if len(data) < 64 {
		return "", ErrExecutionReverted
	}

	// Read length from bytes 32-64
	length := new(big.Int).SetBytes(data[32:64]).Uint64()

	if uint64(len(data)) < 64+length {
		return "", ErrExecutionReverted
	}

	// Extract string data
	strData := data[64 : 64+length]
	return string(strData), nil
}

// encodeUint256 encodes a big.Int as a 32-byte array for ABI return.
func encodeUint256(value *big.Int) []byte {
	result := make([]byte, 32)
	if value != nil {
		value.FillBytes(result)
	}
	return result
}
