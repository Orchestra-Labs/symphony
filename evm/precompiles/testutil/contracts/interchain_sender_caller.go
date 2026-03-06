package contracts

import (
	contractutils "github.com/osmosis-labs/osmosis/v27/evm/contracts/utils"
	evmtypes "github.com/osmosis-labs/osmosis/v27/x/vm/types"
)

func LoadInterchainSenderCallerContract() (evmtypes.CompiledContract, error) {
	return contractutils.LoadContractFromJSONFile("InterchainSenderCaller.json")
}
