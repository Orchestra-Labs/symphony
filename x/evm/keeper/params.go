package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/osmosis-labs/osmosis/v27/x/evm/types"
)

// GetParams returns the total set of EVM parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyPrefixParams)
	if bz == nil {
		return types.DefaultParams()
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams sets the EVM parameters to the store.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.KeyPrefixParams, bz)

	return nil
}

// GetEnableCreate returns whether contract creation is enabled.
func (k Keeper) GetEnableCreate(ctx sdk.Context) bool {
	params := k.GetParams(ctx)
	return params.EnableCreate
}

// GetEnableCall returns whether contract calls are enabled.
func (k Keeper) GetEnableCall(ctx sdk.Context) bool {
	params := k.GetParams(ctx)
	return params.EnableCall
}

// GetChainID returns the EVM chain ID.
func (k Keeper) GetChainID(ctx sdk.Context) string {
	params := k.GetParams(ctx)
	return params.ChainID
}

// GetExtraEIPs returns the list of extra EIPs to enable.
func (k Keeper) GetExtraEIPs(ctx sdk.Context) []int64 {
	params := k.GetParams(ctx)
	return params.ExtraEIPs
}
