package contracts

import (
	contractutils "github.com/osmosis-labs/osmosis/v27/evm/contracts/utils"
	evmtypes "github.com/osmosis-labs/osmosis/v27/x/vm/types"
)

func LoadIcs20CallerContract() (evmtypes.CompiledContract, error) {
	return contractutils.LoadContractFromJSONFile("ICS20Caller.json")
}
