package keeper

import (
	"context"
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/osmosis-labs/osmosis/v27/x/evm/types"
)

var _ types.QueryServer = Keeper{}

// Params returns the EVM module parameters.
func (k Keeper) Params(goCtx context.Context, req *QueryParamsRequest) (*QueryParamsResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidState
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)

	return &QueryParamsResponse{
		Params: params,
	}, nil
}

// Balance returns the EVM balance for an address.
func (k Keeper) Balance(goCtx context.Context, req *QueryBalanceRequest) (*QueryBalanceResponse, error) {
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

	return &QueryBalanceResponse{
		Balance: balance.String(),
	}, nil
}

// Code returns the contract code for an address.
func (k Keeper) Code(goCtx context.Context, req *QueryCodeRequest) (*QueryCodeRequest, error) {
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
		return &QueryCodeRequest{
			Address: req.Address,
			Code:    []byte{},
		}, nil
	}

	// Get the code
	code := k.GetCode(ctx, codeHash)

	return &QueryCodeRequest{
		Address: req.Address,
		Code:    code,
	}, nil
}

// Storage returns a storage value for a contract.
func (k Keeper) Storage(goCtx context.Context, req *QueryStorageRequest) (*QueryStorageResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidState
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Decode the address
	addr, err := hex.DecodeString(req.Address)
	if err != nil {
		return nil, types.ErrInvalidAddress
	}

	// Decode the storage key
	key, err := hex.DecodeString(req.Key)
	if err != nil {
		return nil, types.ErrInvalidState
	}

	// Get the storage value
	value := k.GetStorage(ctx, sdk.AccAddress(addr), key)

	return &QueryStorageResponse{
		Value: hex.EncodeToString(value),
	}, nil
}

// Placeholder query request/response types (would be auto-generated from proto in production)
type QueryParamsRequest struct{}
type QueryParamsResponse struct {
	Params types.Params
}

type QueryBalanceRequest struct {
	Address string
}
type QueryBalanceResponse struct {
	Balance string
}

type QueryCodeRequest struct {
	Address string
	Code    []byte
}

type QueryStorageRequest struct {
	Address string
	Key     string
}
type QueryStorageResponse struct {
	Value string
}
