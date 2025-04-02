package v10

import (
	store "cosmossdk.io/store/types"
	"github.com/osmosis-labs/osmosis/v26/app/upgrades"
	epochtypes "github.com/osmosis-labs/osmosis/v26/x/epochs/types"
	txfeestypes "github.com/osmosis-labs/osmosis/v26/x/txfees/types"
)

const UpgradeName = "v10"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added:   []string{txfeestypes.StoreKey, epochtypes.StoreKey},
		Deleted: []string{},
	},
}
