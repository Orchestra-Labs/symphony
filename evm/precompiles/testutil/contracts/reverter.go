package contracts

import (
	contractutils "github.com/osmosis-labs/osmosis/v27/evm/contracts/utils"
	evmtypes "github.com/osmosis-labs/osmosis/v27/x/vm/types"
)

func LoadReverterContract() (evmtypes.CompiledContract, error) {
	return contractutils.LoadContractFromJSONFile("Reverter.json")
}
