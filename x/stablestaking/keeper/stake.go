package keeper

import (
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
	"strings"

	"github.com/osmosis-labs/osmosis/v27/x/stablestaking/types"
)

func (k Keeper) StakeTokens(ctx sdk.Context, staker sdk.AccAddress, amount sdk.Coin) (*types.MsgStakeTokensResponse, error) {
	if !types.IsAllowedToken(amount.Denom) {
		return nil, fmt.Errorf("unsupported token: %s", amount.Denom)
	}

	pool, found := k.GetPool(ctx, amount.Denom)
	if !found {
		pool = types.StakingPool{
			Denom:       amount.Denom,
			TotalStaked: math.LegacyZeroDec(),
			TotalShares: math.LegacyZeroDec(),
		}
	}

	userShares := sdk.NewDecCoin(amount.Denom, amount.Amount)
	if pool.TotalStaked.IsPositive() {
		userShares.Amount = userShares.Amount.Mul(pool.TotalShares).Quo(pool.TotalStaked)
	}

	pool.TotalStaked = pool.TotalStaked.Add(math.LegacyNewDecFromInt(amount.Amount))
	pool.TotalShares = pool.TotalShares.Add(userShares.Amount)
	k.SetPool(ctx, pool)

	currentEpoch := k.epochKeeper.GetEpochInfo(ctx, "week")
	userStake, found := k.GetUserStake(ctx, staker, amount.Denom)
	if !found {
		userStake = types.UserStake{
			Address: staker.String(),
			Shares:  userShares.Amount,
			Epoch:   currentEpoch.CurrentEpoch,
		}
	} else {
		userStake.Shares = userStake.Shares.Add(userShares.Amount)
		userStake.Epoch = currentEpoch.CurrentEpoch
	}

	k.SetUserStake(ctx, userStake, amount.Denom)

	err := k.BankKeeper.SendCoinsFromAccountToModule(ctx, staker, types.ModuleName, sdk.NewCoins(amount))
	if err != nil {
		return nil, err
	}

	return &types.MsgStakeTokensResponse{}, nil
}

func (k Keeper) UnStakeTokens(ctx sdk.Context, staker sdk.AccAddress, amount sdk.Coin) (*types.MsgUnstakeTokensResponse, error) {
	if !types.IsAllowedToken(amount.Denom) {
		return nil, fmt.Errorf("unsupported token: %s", amount.Denom)
	}

	pool, find := k.GetPool(ctx, amount.Denom)
	if !find {
		return nil, fmt.Errorf("not found pool for denom %s", amount.Denom)
	}

	if pool.TotalStaked.IsZero() {
		return nil, fmt.Errorf("total staked is zero")
	}

	stakedBalance, found := k.GetUserStake(ctx, staker, amount.Denom)
	if !found {
		return nil, fmt.Errorf("not found staked amount for user %s and denom %s", staker.String(), amount.Denom)
	}

	sharesToRemove := math.LegacyNewDecFromInt(amount.Amount).Mul(pool.TotalShares).Quo(pool.TotalStaked)
	if stakedBalance.Shares.LT(sharesToRemove) {
		return nil, fmt.Errorf("unstake amount exceeds user's share: %s", stakedBalance.Shares.String())
	}

	stakedBalance.Shares = stakedBalance.Shares.Sub(sharesToRemove)
	k.SetUserStake(ctx, stakedBalance, amount.Denom)
	k.AddUnbondingRequest(ctx, staker, amount)

	pool.TotalStaked = pool.TotalStaked.Sub(math.LegacyNewDecFromInt(amount.Amount))
	pool.TotalShares = pool.TotalShares.Sub(sharesToRemove)
	k.SetPool(ctx, pool)

	return &types.MsgUnstakeTokensResponse{
		Staker: staker.String(),
		Amount: sdk.DecCoin{
			Denom:  amount.Denom,
			Amount: math.LegacyNewDecFromInt(amount.Amount),
		},
		TotalShares: pool.TotalShares,
		TotalStaked: pool.TotalStaked,
	}, nil
}

func (k Keeper) AddUnbondingRequest(ctx sdk.Context, staker sdk.AccAddress, amount sdk.Coin) {
	store := ctx.KVStore(k.storeKey)
	key := k.GetUnbondingKey(staker, amount.Denom)

	currentEpoch := k.epochKeeper.GetEpochInfo(ctx, "day")
	countDays := k.GetParams(ctx).UnbondingDuration.Milliseconds() / 1000 / 60 / 60 / 24 // milliseconds, minutes, hours, days
	unbondingEpoch := currentEpoch.CurrentEpoch + countDays

	unbondingInfo, found := k.GetUnbondingInfo(ctx, staker, amount.Denom)
	if found {
		unbondingInfo.Amount = unbondingInfo.Amount.Add(math.LegacyNewDecFromInt(amount.Amount))
		unbondingInfo.UnbondEpoch = unbondingEpoch
	} else {
		unbondingInfo.UnbondEpoch = unbondingEpoch
		unbondingInfo.Denom = amount.Denom
		unbondingInfo.Amount = math.LegacyNewDecFromInt(amount.Amount)
		unbondingInfo.Address = staker.String()
	}

	store.Set(key, k.cdc.MustMarshal(&unbondingInfo))
}

func (k Keeper) GetUnbondingInfo(ctx sdk.Context, staker sdk.AccAddress, denom string) (types.UnbondingInfo, bool) {
	store := ctx.KVStore(k.storeKey)
	key := k.GetUnbondingKey(staker, denom)
	bz := store.Get(key)
	if bz == nil {
		return types.UnbondingInfo{}, false
	}

	var info types.UnbondingInfo
	k.cdc.MustUnmarshal(bz, &info)
	return info, true
}

func (k Keeper) GetUnbondingTotalInfo(ctx sdk.Context, staker sdk.AccAddress) []types.UnbondingInfo {
	store := ctx.KVStore(k.storeKey)
	var totalUnbondingInfo []types.UnbondingInfo

	key := k.GetUnbondingKey(staker, "")
	iterator := storetypes.KVStorePrefixIterator(store, key)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var unbondingInfo types.UnbondingInfo
		k.cdc.MustUnmarshal(iterator.Value(), &unbondingInfo)
		totalUnbondingInfo = append(totalUnbondingInfo, unbondingInfo)
	}

	return totalUnbondingInfo
}

func (k Keeper) GetUnbondingKey(staker sdk.AccAddress, denom string) []byte {
	return []byte(fmt.Sprintf("%s:%s%s", types.UnbondingKey, staker.String(), denom))
}

func (k Keeper) GetPool(ctx sdk.Context, token string) (types.StakingPool, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PoolKey))
	bz := store.Get([]byte(token))
	if bz == nil {
		return types.StakingPool{}, false
	}

	var pool types.StakingPool
	k.cdc.MustUnmarshal(bz, &pool)
	return pool, true
}

func (k Keeper) GetPools(ctx sdk.Context) []types.StakingPool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PoolKey))
	var pools []types.StakingPool

	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var pool types.StakingPool
		k.cdc.MustUnmarshal(iterator.Value(), &pool)

		pools = append(pools, pool)
	}

	return pools
}

func (k Keeper) SetPool(ctx sdk.Context, pool types.StakingPool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PoolKey))
	bz := k.cdc.MustMarshal(&pool)
	store.Set([]byte(pool.Denom), bz)
}

func (k Keeper) GetUserStakeKey(staker string, denom string) []byte {
	return []byte(fmt.Sprintf("%s:%s:%s", types.UserStakeKey, staker, denom))
}

func (k Keeper) GetUserStake(ctx sdk.Context, address sdk.AccAddress, token string) (types.UserStake, bool) {
	store := ctx.KVStore(k.storeKey)
	keyVal := k.GetUserStakeKey(address.String(), token)
	bz := store.Get(keyVal)
	if bz == nil {
		return types.UserStake{}, false
	}

	var stake types.UserStake
	k.cdc.MustUnmarshal(bz, &stake)
	return stake, true
}

func (k Keeper) SetUserStake(ctx sdk.Context, stake types.UserStake, token string) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&stake)
	keyVal := k.GetUserStakeKey(stake.Address, token)
	store.Set(keyVal, bz)
}

func (k Keeper) GetUserTotalStake(ctx sdk.Context, address sdk.AccAddress) sdk.DecCoins {
	store := ctx.KVStore(k.storeKey)
	var stakes sdk.DecCoins

	keyPrefix := k.GetUserStakeKey(address.String(), "")
	iterator := storetypes.KVStorePrefixIterator(store, keyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var stake types.UserStake
		k.cdc.MustUnmarshal(iterator.Value(), &stake)

		key := string(iterator.Key())
		token := key[len(keyPrefix):]

		stakes = append(stakes, sdk.DecCoin{
			Denom:  token,
			Amount: stake.Shares,
		})
	}

	return stakes
}

func (k Keeper) GetTotalStakersPerPool(ctx sdk.Context, token string) (int32, error) {
	store := ctx.KVStore(k.storeKey)

	_, found := k.GetPool(ctx, token)
	if !found {
		return 0, fmt.Errorf("not found pool with denom %s", token)
	}

	keyPrefix := k.GetUserStakeKey("", "")
	iterator := storetypes.KVStorePrefixIterator(store, keyPrefix)
	defer iterator.Close()

	var stakers int32
	for ; iterator.Valid(); iterator.Next() {
		var stake types.UserStake
		k.cdc.MustUnmarshal(iterator.Value(), &stake)

		key := string(iterator.Key())
		parts := strings.Split(key, ":")

		if len(parts) >= 3 {
			denom := parts[2]
			if denom == token {
				stakers += 1
			}
		}
	}

	return stakers, nil
}

func (k Keeper) GetTotalStakers(ctx sdk.Context) []*types.TotalStakers {
	store := ctx.KVStore(k.storeKey)
	counts := make(map[string]int)

	keyPrefix := k.GetUserStakeKey("", "")
	iterator := storetypes.KVStorePrefixIterator(store, keyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var stake types.UserStake
		k.cdc.MustUnmarshal(iterator.Value(), &stake)

		key := string(iterator.Key())
		parts := strings.Split(key, ":")

		if len(parts) >= 3 {
			denom := parts[2]
			counts[denom] = counts[denom] + 1
		}
	}

	kvList := make([]*types.TotalStakers, 0, len(counts))
	for k, v := range counts {
		kvList = append(kvList, &types.TotalStakers{Denom: k, Count: strconv.Itoa(v)})
	}

	return kvList
}
