package v10

import (
	"context"
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

		// Run migrations before applying any other state changes.
		// NOTE: DO NOT PUT ANY STATE CHANGES BEFORE RunMigrations().
		migrations, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return nil, err
		}

		return migrations, nil
	}
}
