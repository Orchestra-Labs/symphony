package keeper

import (
	storetypes "cosmossdk.io/store/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/osmosis-labs/osmosis/v27/x/stablestaking/types"
)

type Keeper struct {
	storeKey   storetypes.StoreKey
	cdc        codec.Codec
	paramSpace paramstypes.Subspace

	epochKeeper   types.EpochKeeper
	AccountKeeper types.AccountKeeper
	BankKeeper    types.BankKeeper
	OracleKeeper  types.OracleKeeper
}

func NewKeeper(
	cdc codec.Codec,
	storeKey storetypes.StoreKey,
	paramstore paramstypes.Subspace,
	epochKeeper types.EpochKeeper,
	accKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	oracleKeeper types.OracleKeeper,
) Keeper {
	// ensure market module account is set
	if addr := accKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	if !paramstore.HasKeyTable() {
		paramstore = paramstore.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		paramSpace:    paramstore,
		epochKeeper:   epochKeeper,
		BankKeeper:    bankKeeper,
		AccountKeeper: accKeeper,
		OracleKeeper:  oracleKeeper,
	}
}

// GetParams return module params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSetIfExists(ctx, &params)
	return params
}

// SetParams set up module params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
