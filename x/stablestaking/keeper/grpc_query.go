package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/v27/x/stablestaking/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Querier{}

type Querier struct {
	Keeper
}

func (q Querier) StablePool(c context.Context, request *types.QueryPoolRequest) (*types.QueryPoolResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
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

	if err := sdk.ValidateDenom(request.Token); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid stake denom")
	}

	ctx := sdk.UnwrapSDKContext(c)
	userStake, found := q.GetUserStake(ctx, sdk.AccAddress(request.Address), request.Token)
	if !found {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("User stake not found with %s", request.Token))
	}
	return &types.QueryUserStakeResponse{Stakes: &userStake}, nil
}

func (q Querier) UserTotalStake(c context.Context, request *types.QueryUserTotalStakeRequest) (*types.QueryUserTotalStakeResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	userTotalStake := q.GetUserTotalStake(ctx, sdk.AccAddress(request.Address))

	return &types.QueryUserTotalStakeResponse{Stakes: userTotalStake}, nil
}

func (q Querier) UserUnbonding(ctx context.Context, request *types.QueryUserUnbondingRequest) (*types.QueryUserUnbondingResponse, error) {
	//TODO implement me
	return &types.QueryUserUnbondingResponse{}, nil
}

func (q Querier) UserTotalUnbonding(ctx context.Context, request *types.QueryUserTotalUnbondingRequest) (*types.QueryUserTotalUnbondingResponse, error) {
	//TODO implement me
	return &types.QueryUserTotalUnbondingResponse{}, nil
}
