package keeper

import (
	icacallbackstypes "github.com/osmosis-labs/osmosis/v27/x/icacallbacks/types"
)

const (
	ICACallbackID_InstantiateOracle = "instantiate_oracle"
	ICACallbackID_UpdateOracle      = "update_oracle"
)

func (k Keeper) Callbacks() icacallbackstypes.ModuleCallbacks {
	return []icacallbackstypes.ICACallback{
		{
			CallbackId:   ICACallbackID_InstantiateOracle,
			CallbackFunc: icacallbackstypes.ICACallbackFunction(k.InstantiateOracleCallback),
		},
		{
			CallbackId:   ICACallbackID_UpdateOracle,
			CallbackFunc: icacallbackstypes.ICACallbackFunction(k.UpdateOracleCallback),
		},
	}
}
