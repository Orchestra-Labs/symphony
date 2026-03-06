package testutil

import (
	contractutils "github.com/osmosis-labs/osmosis/v27/evm/contracts/utils"
	evmtypes "github.com/osmosis-labs/osmosis/v27/x/vm/types"
)

func LoadCounterWithCallbacksContract() (evmtypes.CompiledContract, error) {
	return contractutils.LoadContractFromJSONFile("CounterWithCallbacks.json")
}
