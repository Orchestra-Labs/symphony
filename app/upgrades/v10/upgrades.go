package v10

import (
	"context"
	"cosmossdk.io/log"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/osmosis-labs/osmosis/v26/app/keepers"
	"github.com/osmosis-labs/osmosis/v26/app/upgrades"
	txfeestypes "github.com/osmosis-labs/osmosis/v26/x/txfees/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	bpm upgrades.BaseAppParamManager,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx1 := sdk.UnwrapSDKContext(ctx)
		fromVM1 := make(map[string]uint64)
		for moduleName := range mm.Modules {
			fromVM1[moduleName] = 1
		}

		logger := log.NewNopLogger()
		logger.Warn("Run migrations")
		ctx1.Logger().Warn("Run migrations")
		// Run migrations before applying any other state changes.
		// NOTE: DO NOT PUT ANY STATE CHANGES BEFORE RunMigrations().
		newVM, err := mm.RunMigrations(ctx1, configurator, fromVM1)
		if err != nil {
			ctx1.Logger().Error("❌ Migration failed:", "error", err)
			return nil, err
		}

		ctx1.Logger().Warn("Setting txfees module genesis with actual v5 desired genesis")
		initGenParams := txfeestypes.DefaultGenesis()
		keepers.TxFeesKeeper.InitGenesis(ctx1, *initGenParams)

		txfeesParams := txfeestypes.DefaultParams()
		keepers.TxFeesKeeper.SetParams(ctx1, txfeesParams)

		// N.B.: this is done to avoid initializing genesis for ibcratelimit module.
		// Otherwise, it would overwrite migrations with InitGenesis().
		// See RunMigrations() for details.
		// fromVM[markettypes.ModuleName] = 0
		//fromVM[stablestakingincentvicestypes.ModuleName] = 0
		//fromVM[txfeestypes.ModuleName] = 0
		//fromVM[epochtypes.ModuleName] = 0

		return newVM, nil
	}
}
