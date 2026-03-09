package keeper

import (
	"github.com/osmosis-labs/osmosis/v27/x/evm/keeper/precompiles/bank"
	"github.com/osmosis-labs/osmosis/v27/x/evm/keeper/precompiles/oracle"
	// "github.com/osmosis-labs/osmosis/v27/x/evm/keeper/precompiles/staking"
)

// InitializePrecompiles registers all custom precompiled contracts with the EVM keeper.
// This should be called after the keeper is fully initialized with all dependencies.
func (k *Keeper) InitializePrecompiles(oracleKeeper oracle.OracleKeeper) {
	// Register Oracle precompile (0x0000000000000000000000000000000000000800)
	// Provides price feed access for smart contracts
	if oracleKeeper != nil {
		oraclePrecompile := oracle.NewOraclePrecompile(oracleKeeper)
		k.RegisterPrecompile(oracle.PrecompileAddress, oraclePrecompile)
	}

	// Create adapter wrappers for keepers to match precompile interfaces
	bankAdapter := &bankKeeperAdapter{keeper: k.bankKeeper}
	accountAdapter := &accountKeeperAdapter{keeper: k.accountKeeper}

	// Register Bank precompile (0x0000000000000000000000000000000000000801)
	// Provides ERC-20 compatible interface over x/bank module
	bankPrecompile := bank.NewBankPrecompile(bankAdapter, accountAdapter)
	k.RegisterPrecompile(bank.PrecompileAddress, bankPrecompile)

	// Register Staking precompile (0x0000000000000000000000000000000000000802)
	// Provides staking operations (delegate, undelegate) from smart contracts
	// Note: Staking precompile currently not fully compatible with types.StakingKeeper
	// Commenting out until interface is aligned
	// stakingPrecompile := staking.NewStakingPrecompile(k.stakingKeeper, bankAdapter)
	// k.RegisterPrecompile(staking.PrecompileAddress, stakingPrecompile)
}
