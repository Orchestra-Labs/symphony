package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
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
	// ensure stable staking module account is set
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

	// Store both by epoch and as latest
	epoch := k.epochKeeper.GetEpochInfo(ctx, k.GetParams(ctx).EpochIdentifier).CurrentEpoch
	epochKey := sdk.Uint64ToBigEndian(uint64(epoch))
	store.Set(epochKey, bz)

	// Also store as latest for quick access
	store.Set([]byte(fmt.Sprintf("latest:%s", denom)), bz)
}

func (k Keeper) GetEpochSnapshot(ctx sdk.Context, denom string) types.EpochSnapshot {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.SnapshotKey))

	// Try to get the latest snapshot first
	bz := store.Get([]byte(fmt.Sprintf("latest:%s", denom)))
	if bz == nil {
		return types.EpochSnapshot{}
	}

	var snapshot types.EpochSnapshot
	k.cdc.MustUnmarshal(bz, &snapshot)
	return snapshot
}

func (k Keeper) SnapshotCurrentEpoch(ctx sdk.Context) {
	params := k.GetParams(ctx)
	if len(params.SupportedTokens) == 0 {
		return
	}

	var totalShares math.LegacyDec
	var totalStaked math.LegacyDec
	var stakers []*types.UserStake

	// Iterate through all stakers and collect their stakes
	k.IterateActiveStakers(ctx, func(addr sdk.AccAddress, stake types.UserStake) {
		// Get the denom from the key
		key := addr.String() + params.SupportedTokens[0]
		store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.UserStakeKey))
		if store.Has([]byte(key)) {
			stakers = append(stakers, &stake)
			totalShares = totalShares.Add(stake.Shares)
			totalStaked = totalStaked.Add(stake.Shares)
		}
	})

	snapshot := types.EpochSnapshot{
		TotalShares: totalShares,
		TotalStaked: totalStaked,
		Stakers:     stakers,
	}

	// Store the snapshot
	k.SetEpochSnapshot(ctx, snapshot, params.SupportedTokens[0])
}

func (k Keeper) IterateActiveStakers(ctx sdk.Context, cb func(addr sdk.AccAddress, stake types.UserStake)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.UserStakeKey))

	iterator := storetypes.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var stake types.UserStake
		k.cdc.MustUnmarshal(iterator.Value(), &stake)

		addr, err := sdk.AccAddressFromBech32(stake.Address)
		if err != nil {
			panic(fmt.Sprintf("invalid address in active staker store: %s", err))
		}

		cb(addr, stake)
	}
}

func (k Keeper) DistributeRewardsToLastEpochStakers(ctx sdk.Context, totalReward math.Int) {
	params := k.GetParams(ctx)
	if len(params.SupportedTokens) == 0 {
		return
	}

	snapshot := k.GetEpochSnapshot(ctx, params.SupportedTokens[0])
	if snapshot.TotalShares.IsZero() {
		return // No snapshot - no rewards
	}

	// Verify we have enough balances in the module account
	moduleAddr := k.AccountKeeper.GetModuleAddress(types.ModuleName)
	balance := k.BankKeeper.GetBalance(ctx, moduleAddr, params.SupportedTokens[0])
	if balance.Amount.LT(totalReward) {
		panic(fmt.Sprintf("insufficient balance in module account: %s < %s", balance.Amount, totalReward))
	}

	for _, staker := range snapshot.Stakers {
		if staker.Shares.IsZero() {
			continue
		}

		// Calculate reward based on a share ratio
		reward := staker.Shares.Quo(snapshot.TotalShares).MulInt(totalReward).TruncateInt()
		if reward.IsZero() {
			continue
		}

		addr, err := sdk.AccAddressFromBech32(staker.Address)
		if err != nil {
			panic(fmt.Sprintf("invalid address in snapshot: %s", err))
		}

		// Send reward tokens
		rewardCoin := sdk.NewCoin(params.SupportedTokens[0], reward)
		err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.NewCoins(rewardCoin))
		if err != nil {
			panic(fmt.Sprintf("failed to send rewards: %s", err))
		}
	}
}

func (k Keeper) GetEpochReward(ctx sdk.Context) math.Int {
	params := k.GetParams(ctx)
	if len(params.SupportedTokens) == 0 {
		return math.ZeroInt()
	}

	// Get the total staked amount for the first supported token
	pool, found := k.GetPool(ctx, params.SupportedTokens[0])
	if !found || pool.TotalStaked.IsZero() {
		return math.ZeroInt()
	}

	// Parse reward rate from params
	rewardRate, err := math.LegacyNewDecFromStr(params.RewardRate)
	if err != nil {
		panic(fmt.Sprintf("invalid reward rate: %s", err))
	}

	// Calculate reward: total_staked * reward_rate
	reward := pool.TotalStaked.Mul(rewardRate).TruncateInt()
	return reward
}
