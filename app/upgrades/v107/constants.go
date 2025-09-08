package v29

import (
	store "cosmossdk.io/store/types"
	"github.com/osmosis-labs/osmosis/v27/app/upgrades"
)

const UpgradeName = "v1.0.7"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades:        store.StoreUpgrades{},
}
