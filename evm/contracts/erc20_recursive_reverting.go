package contracts

import (
	contractutils "github.com/osmosis-labs/osmosis/v27/evm/contracts/utils"
	evmtypes "github.com/osmosis-labs/osmosis/v27/x/vm/types"
)

func LoadERC20RecursiveReverting() (evmtypes.CompiledContract, error) {
	return contractutils.LoadContractFromJSONFile("solidity/ERC20RecursiveRevertingPrecompileCall.json")
}
