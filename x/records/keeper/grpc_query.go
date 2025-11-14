package keeper

import (
	"github.com/osmosis-labs/osmosis/v27/x/records/types"
)

var _ types.QueryServer = Keeper{}
