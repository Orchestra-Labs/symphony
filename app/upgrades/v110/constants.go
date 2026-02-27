package v110

import (
	store "cosmossdk.io/store/types"
	"github.com/osmosis-labs/osmosis/v27/app/upgrades"
	evmtypes "github.com/osmosis-labs/osmosis/v27/x/evm/types"
)

const UpgradeName = "v1.1.0"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			evmtypes.ModuleName, // Add EVM module store
		},
	},
}
