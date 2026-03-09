package keeper

import (
	"context"
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/osmosis-labs/osmosis/v27/x/evm/types"
)

var _ types.QueryServer = Keeper{}

// Params returns the EVM module parameters.
func (k Keeper) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidState
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{
		Params: params,
	}, nil
}

// Balance returns the EVM balance for an address.
func (k Keeper) Balance(goCtx context.Context, req *types.QueryBalanceRequest) (*types.QueryBalanceResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidState
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Decode the address
	addr, err := hex.DecodeString(req.Address)
	if err != nil {
		return nil, types.ErrInvalidAddress
	}

	balance := k.GetBalance(ctx, sdk.AccAddress(addr))

	return &types.QueryBalanceResponse{
		Balance: balance.String(),
	}, nil
}

// Code returns the contract code for an address.
func (k Keeper) Code(goCtx context.Context, req *types.QueryCodeRequest) (*types.QueryCodeResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidState
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Decode the address
	addr, err := hex.DecodeString(req.Address)
	if err != nil {
		return nil, types.ErrInvalidAddress
	}

	// Get the code hash
	codeHash := k.GetCodeHash(ctx, sdk.AccAddress(addr))
	if codeHash == nil {
		return &types.QueryCodeResponse{
			Code: []byte{},
		}, nil
	}

	// Get the code
	code := k.GetCode(ctx, codeHash)

	return &types.QueryCodeResponse{
		Code: code,
	}, nil
}
