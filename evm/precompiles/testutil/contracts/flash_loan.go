package contracts

import (
	contractutils "github.com/osmosis-labs/osmosis/v27/evm/contracts/utils"
	evmtypes "github.com/osmosis-labs/osmosis/v27/x/vm/types"
)

func LoadFlashLoanContract() (evmtypes.CompiledContract, error) {
	return contractutils.LoadContractFromJSONFile("FlashLoan.json")
}
