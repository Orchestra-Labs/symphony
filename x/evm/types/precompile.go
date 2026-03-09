package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

// StatefulPrecompiledContract is the interface for stateful precompiled contracts that need SDK context.
type StatefulPrecompiledContract interface {
	// Address returns the address where the precompiled contract is deployed.
	Address() common.Address

	// RequiredGas calculates the gas required to execute the precompiled contract.
	RequiredGas(input []byte) uint64

	// Run executes the precompiled contract logic with SDK context.
	Run(ctx sdk.Context, caller common.Address, input []byte, suppliedGas uint64, readOnly bool) ([]byte, uint64, error)
}
