package testutil

import (
	anteinterfaces "github.com/osmosis-labs/osmosis/v27/evm/ante/interfaces"
	"github.com/osmosis-labs/osmosis/v27/x/vm/statedb"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewStateDB returns a new StateDB for testing purposes.
func NewStateDB(ctx sdk.Context, evmKeeper anteinterfaces.EVMKeeper) *statedb.StateDB {
	return statedb.New(ctx, evmKeeper, statedb.NewEmptyTxConfig())
}
