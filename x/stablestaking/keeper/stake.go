package keeper

import (
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/v27/x/stablestaking/types"
)

func (k Keeper) StakeTokens(ctx sdk.Context, staker sdk.AccAddress, amount sdk.Coin) (*types.MsgStakeTokensResponse, error) {
	pool, found := k.GetPool(ctx, amount.Denom)
	if !found {
		pool = types.StakingPool{
			Token:       amount.Denom,
			TotalStaked: math.LegacyZeroDec(),
			TotalShares: math.LegacyZeroDec(),
		}
	}

	userShares := sdk.NewDecCoin(amount.Denom, amount.Amount)
	if pool.TotalStaked.IsPositive() {
		userShares.Amount = userShares.Amount.Mul(pool.TotalShares).Quo(pool.TotalStaked)
	}

	pool.TotalStaked = pool.TotalStaked.Add(math.LegacyDec(amount.Amount))
	pool.TotalShares = pool.TotalShares.Add(userShares.Amount)
	k.SetPool(ctx, pool)

	userStake, found := k.GetUserStake(ctx, staker, amount.Denom)
	if !found {
		userStake = types.UserStake{
			Address: staker.String(),
			Shares:  userShares.Amount,
			Epoch:   0,
		}
	} else {
		userStake.Shares = userStake.Shares.Add(userShares.Amount)
		//TODO: how do not override previous stake and correct distribute rewards?
	}

	k.SetUserStake(ctx, userStake, amount.Denom)

	err := k.BankKeeper.SendCoinsFromAccountToModule(ctx, staker, types.ModuleName, sdk.NewCoins(amount))
	if err != nil {
		return nil, err
	}

	return &types.MsgStakeTokensResponse{}, nil
}

func (k Keeper) UnStakeTokens(ctx sdk.Context, staker sdk.AccAddress, amount sdk.Coin) (*types.MsgUnstakeTokensResponse, error) {
	pool, _ := k.GetPool(ctx, amount.Denom)
	if pool.TotalStaked.IsZero() {
		return nil, fmt.Errorf("total staked is zero")
	}

	stakedBalance, found := k.GetUserStake(ctx, staker, amount.Denom)
	if !found {
		return nil, fmt.Errorf("not found staked amount for user %s and denom %s", staker.String(), amount.Denom)
	}

	sharesToRemove := math.LegacyDec(amount.Amount).Mul(pool.TotalShares).Quo(pool.TotalStaked)
	if stakedBalance.Shares.LT(sharesToRemove) {
		return nil, fmt.Errorf("unstake amount exceeds user's share: %s", stakedBalance.Shares.String())
	}

	stakedBalance.Shares = stakedBalance.Shares.Sub(sharesToRemove)
	k.SetUserStake(ctx, stakedBalance, amount.Denom)
	k.AddUnbondingRequest(ctx, staker, amount)
	//TODO: send coins to user after unbonding period
	//err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, address, sdk.NewCoins(unstakeAmount))
	//if err != nil {
	//	return err
	//}

	pool.TotalStaked = pool.TotalStaked.Sub(math.LegacyDec(amount.Amount))
	pool.TotalShares = pool.TotalShares.Add(sharesToRemove)
	k.SetPool(ctx, pool)

	return &types.MsgUnstakeTokensResponse{}, nil
}

func (k Keeper) AddUnbondingRequest(ctx sdk.Context, staker sdk.AccAddress, amount sdk.Coin) {
	store := ctx.KVStore(k.storeKey)
	key := k.GetUnbondingKey(staker)

	unbondingTime := ctx.BlockTime().Add(k.GetParams(ctx).UnbondingTime)
	unbondingRequest := types.UnbondingRequest{
		Address:    staker.String(),
		Shares:     math.LegacyDec(amount.Amount),
		UnbondTime: unbondingTime.UnixMicro(), //TODO: check if it's correct
	}

	store.Set(key, k.cdc.MustMarshal(&unbondingRequest))
}

func (k Keeper) GetUnbondingInfo(ctx sdk.Context, staker sdk.AccAddress) (types.UnbondingRequest, bool) {
	store := ctx.KVStore(k.storeKey)
	key := k.GetUnbondingKey(staker)
	bz := store.Get(key)
	if bz == nil {
		return types.UnbondingRequest{}, false
	}

	var info types.UnbondingRequest
	k.cdc.MustUnmarshal(bz, &info)
	return info, true
}

func (k Keeper) GetUnbondingKey(staker sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("unbonding:%s", staker.String()))
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

func (k Keeper) SetPool(ctx sdk.Context, pool types.StakingPool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PoolKey))
	bz := k.cdc.MustMarshal(&pool)
	store.Set([]byte(pool.Token), bz)
}

func (k Keeper) GetUserStake(ctx sdk.Context, address sdk.AccAddress, token string) (types.UserStake, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PoolKey))
	bz := store.Get([]byte(address.String() + token))
	if bz == nil {
		return types.UserStake{}, false
	}

	var stake types.UserStake
	k.cdc.MustUnmarshal(bz, &stake)
	return stake, true
}

func (k Keeper) SetUserStake(ctx sdk.Context, stake types.UserStake, token string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PoolKey))
	bz := k.cdc.MustMarshal(&stake)
	store.Set([]byte(stake.Address+token), bz)
}
