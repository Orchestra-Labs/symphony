package keeper

import (
	"context"
	"fmt"
	appparams "github.com/osmosis-labs/osmosis/v27/app/params"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/osmosis-labs/osmosis/v27/x/stablestaking/types"
)

var _ types.QueryServer = Querier{}

type Querier struct {
	Keeper
}

func (q Querier) StablePool(c context.Context, request *types.QueryPoolRequest) (*types.QueryPoolResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if err := sdk.ValidateDenom(request.Denom); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid stake denom")
	}

	ctx := sdk.UnwrapSDKContext(c)
	pool, find := q.GetPool(ctx, request.Denom)
	if !find {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Pool not found with %s", request.Denom))
	}

	return &types.QueryPoolResponse{Pool: &pool}, nil
}

func (q Querier) StablePools(c context.Context, request *types.QueryPoolsRequest) (*types.QueryPoolsResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	pools := q.GetPools(ctx)

	pointerPools := make([]*types.StakingPool, len(pools))
	for i := range pools {
		pointerPools[i] = &pools[i]
	}

	return &types.QueryPoolsResponse{Pools: pointerPools}, nil
}

// NewQueryServerImpl returns an implementation of the stablestaking QueryServer interface
// for the provided Keeper.
func NewQueryServerImpl(keeper Keeper) types.QueryServer {
	return &Querier{Keeper: keeper}
}

// Params queries params of stablestaking module
func (q Querier) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	return &types.QueryParamsResponse{Params: q.GetParams(ctx)}, nil
}

func (q Querier) UserStake(c context.Context, request *types.QueryUserStakeRequest) (*types.QueryUserStakeResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	reqAddress, err := sdk.AccAddressFromBech32(request.Address)
	if err != nil {
		panic(fmt.Sprintf("invalid staker address : %s", err))
	}

	if err := sdk.ValidateDenom(request.Denom); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid stake denom")
	}

	ctx := sdk.UnwrapSDKContext(c)
	userStake, found := q.GetUserStake(ctx, reqAddress, request.Denom)
	if !found {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("User stake not found with %s", request.Denom))
	}
	return &types.QueryUserStakeResponse{Stakes: &userStake}, nil
}

func (q Querier) UserTotalStake(c context.Context, request *types.QueryUserTotalStakeRequest) (*types.QueryUserTotalStakeResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	reqAddress, err := sdk.AccAddressFromBech32(request.Address)
	if err != nil {
		panic(fmt.Sprintf("invalid staker address : %s", err))
	}

	ctx := sdk.UnwrapSDKContext(c)
	userTotalStake := q.GetUserTotalStake(ctx, reqAddress)

	return &types.QueryUserTotalStakeResponse{Stakes: userTotalStake}, nil
}

func (q Querier) UserUnbonding(ctx context.Context, request *types.QueryUserUnbondingRequest) (*types.QueryUserUnbondingResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	reqAddress, err := sdk.AccAddressFromBech32(request.Address)
	if err != nil {
		panic(fmt.Sprintf("invalid staker address : %s", err))
	}

	if err := sdk.ValidateDenom(request.Denom); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid denom")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	unbondingInfo, found := q.GetUnbondingInfo(sdkCtx, reqAddress, request.Denom)
	if !found {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Unbonding info not found for address %s and denom %s", request.Address, request.Denom))
	}

	return &types.QueryUserUnbondingResponse{
		Info: &unbondingInfo,
	}, nil
}

func (q Querier) UserTotalUnbonding(ctx context.Context, request *types.QueryUserTotalUnbondingRequest) (*types.QueryUserTotalUnbondingResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	reqAddress, err := sdk.AccAddressFromBech32(request.Address)
	if err != nil {
		panic(fmt.Sprintf("invalid staker address : %s", err))
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	unbondingInfos := q.GetUnbondingTotalInfo(sdkCtx, reqAddress)

	var unbondInfos []*types.UnbondingInfo
	for _, info := range unbondingInfos {
		unbondInfos = append(unbondInfos, &info)
	}
	return &types.QueryUserTotalUnbondingResponse{
		Info: unbondInfos,
	}, nil
}

func (q Querier) TotalStakersPerPool(ctx context.Context, request *types.QueryPoolRequest) (*types.QueryTotalStakersPerPoolResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	totalStakers, err := q.GetTotalStakersPerPool(sdkCtx, request.GetDenom(), request.GetLimit())
	if err != nil {
		return nil, err
	}

	return &types.QueryTotalStakersPerPoolResponse{
		Stakers: &types.TotalStakers{
			Denom:   request.Denom,
			Stakers: totalStakers,
		},
	}, nil
}

func (q Querier) TotalStakers(ctx context.Context, request *types.QueryPoolsRequest) (*types.QueryTotalStakersResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	totalStakers := q.GetTotalStakers(sdkCtx)

	return &types.QueryTotalStakersResponse{
		Stakers: totalStakers,
	}, nil
}

func (q Querier) RewardAmountPerPool(ctx context.Context, request *types.QueryPoolRequest) (*types.QueryRewardPerPoolResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// Get total rewards available
	moduleRewardsAddr := q.AccountKeeper.GetModuleAddress(types.NativeRewardsCollectorName)
	totalReward := q.BankKeeper.GetBalance(ctx, moduleRewardsAddr, appparams.BaseCoinUnit)
	if totalReward.IsZero() {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Available rewards: %s", totalReward))
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	snapshotData, err := q.CalculateStakedPools(sdkCtx)
	if err != nil {
		return nil, err
	}

	if snapshotData.TotalStakedAcrossPools.IsZero() {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("No stake found for pool %s. Available rewards: %s", request.Denom, totalReward))
	}

	poolSnapshot, ok := snapshotData.PoolSnapshots[request.Denom]
	if !ok {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("No stake found for pool %s or reward distribution has not started yet", request.Denom))
	}
	poolRewardShare := poolSnapshot.TotalStaked.Quo(snapshotData.TotalStakedAcrossPools)
	poolReward := poolRewardShare.MulInt(totalReward.Amount).TruncateInt()

	if poolReward.IsZero() {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Available rewards for pool %s: %s", request.Denom, poolReward))
	}

	return &types.QueryRewardPerPoolResponse{
		Pool: &types.StakingPool{
			Denom:       request.Denom,
			TotalStaked: poolSnapshot.TotalStaked,
			TotalShares: poolSnapshot.TotalShares,
		},
		Reward: &sdk.Coin{
			Denom:  appparams.BaseCoinUnit,
			Amount: poolReward,
		},
	}, nil
}
