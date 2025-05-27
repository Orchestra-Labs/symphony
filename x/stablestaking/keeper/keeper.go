package keeper

import (
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
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

func (k Keeper) SetEpochSnapshot(ctx sdk.Context, snapshot types.EpochSnapshot, denom string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.SnapshotKey))
	bz := k.cdc.MustMarshal(&snapshot)
	store.Set([]byte(fmt.Sprintf("latest:%s", denom)), bz)
}

func (k Keeper) GetEpochSnapshot(ctx sdk.Context, denom string) types.EpochSnapshot {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.SnapshotKey))
	bz := store.Get([]byte(fmt.Sprintf("latest:%s", denom)))

	if bz == nil {
		return types.EpochSnapshot{}
	}

	var snapshot types.EpochSnapshot
	k.cdc.MustUnmarshal(bz, &snapshot)
	return snapshot
}

func (k Keeper) SnapshotCurrentEpoch(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.UserStakeKey))

	var totalShares math.LegacyDec
	var totalStaked math.LegacyDec
	var stakers []*types.UserStake

	//k.IterateActiveStakers(ctx, func(addr sdk.AccAddress, stake types.UserStake) {
	//	stakers = append(stakers, &stake)
	//	totalShares = totalShares.Add(stake.Shares)
	//	totalStaked = totalStaked.Add(stake.Shares)
	//})

	snapshot := types.EpochSnapshot{
		TotalShares: totalShares,
		TotalStaked: totalStaked,
		Stakers:     stakers,
	}

	epoch := k.epochKeeper.GetEpochInfo(ctx, k.GetParams(ctx).EpochIdentifier).CurrentEpoch
	key := sdk.Uint64ToBigEndian(uint64(epoch))
	store.Set(key, k.cdc.MustMarshal(&snapshot))
}

//func (k Keeper) IterateActiveStakers(ctx sdk.Context, cb func(addr sdk.AccAddress, stake types.StakeInfo)) {
//	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.ActiveStakerPrefix))
//
//	iterator := sdk.KVStorePrefixIterator(store, nil)
//	defer iterator.Close()
//
//	for ; iterator.Valid(); iterator.Next() {
//		var stakeInfo types.StakeInfo
//		k.cdc.MustUnmarshal(iterator.Value(), &stakeInfo)
//
//		addr, err := sdk.AccAddressFromBech32(string(iterator.Key()))
//		if err != nil {
//			panic(fmt.Sprintf("invalid address in active staker store: %s", err))
//		}
//
//		cb(addr, stakeInfo)
//	}
//}

//func (k Keeper) DistributeRewardsToLastEpochStakers(ctx sdk.Context, totalReward sdk.Int) {
//	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.SnapshotPrefix))
//
//	currentEpoch := k.epochKeeper.GetEpochInfo(ctx, k.GetParams(ctx).EpochIdentifier).CurrentEpoch
//	lastEpochKey := sdk.Uint64ToBigEndian(uint64(currentEpoch - 1))
//
//	bz := store.Get(lastEpochKey)
//	if bz == nil {
//		return // нема знімку — нема винагороди
//	}
//
//	var snapshot types.EpochSnapshot
//	k.cdc.MustUnmarshal(bz, &snapshot)
//
//	for addrStr, stakeInfo := range snapshot.Stakers {
//		if snapshot.TotalShares.IsZero() {
//			continue
//		}
//
//		reward := stakeInfo.Shares.Quo(snapshot.TotalShares).MulInt(totalReward).TruncateInt()
//		addr, _ := sdk.AccAddressFromBech32(addrStr)
//
//		// нарахування токенів
//		rewardCoin := sdk.NewCoin(k.GetParams(ctx).RewardDenom, reward)
//		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.NewCoins(rewardCoin))
//		if err != nil {
//			panic(err)
//		}
//	}
//}
