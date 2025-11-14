package keeper

import (
	"github.com/osmosis-labs/osmosis/v27/x/icacallbacks/types"
)

var _ types.QueryServer = Keeper{}
