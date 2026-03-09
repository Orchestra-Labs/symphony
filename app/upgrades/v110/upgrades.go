package v110

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/osmosis-labs/osmosis/v27/app/keepers"
	"github.com/osmosis-labs/osmosis/v27/app/upgrades"
	evmtypes "github.com/osmosis-labs/osmosis/v27/x/evm/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_bpm upgrades.BaseAppParamManager,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(context context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(context)
		ctx.Logger().Info("🚀 Starting v1.1.0 upgrade (EVM Support)")

		// Run migrations before applying any other state changes.
		// NOTE: DO NOT PUT ANY STATE CHANGES BEFORE RunMigrations().
		newVM, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			ctx.Logger().Error("❌ Migration failed:", "error", err)
			return nil, err
		}

		// Initialize EVM module parameters
		ctx.Logger().Info("⚙️  Initializing EVM module parameters")
		evmParams := evmtypes.DefaultParams()
		evmParams.ChainID = "9999" // Set to actual chain ID before mainnet
		if err := keepers.EVMKeeper.SetParams(ctx, evmParams); err != nil {
			ctx.Logger().Error("❌ Failed to set EVM params:", "error", err)
			return nil, err
		}

		ctx.Logger().Info("✅ v1.1.0 upgrade completed successfully!")
		ctx.Logger().Info("📝 EVM module enabled with the following parameters:")
		ctx.Logger().Info("   - Chain ID:", "chain_id", evmParams.ChainID)
		ctx.Logger().Info("   - Contract Creation:", "enabled", evmParams.EnableCreate)
		ctx.Logger().Info("   - Contract Calls:", "enabled", evmParams.EnableCall)
		ctx.Logger().Info("🌉 Note: JSON-RPC server must be enabled in app.toml for full EVM support")

		return newVM, nil
	}
}
