package keeper

import (
	"context"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/osmosis-labs/osmosis/osmomath"
	"github.com/osmosis-labs/osmosis/v27/x/collateral/types"
)

type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Context, addr sdk.AccAddress, moduleName string, coins sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, moduleName string, addr sdk.AccAddress, coins sdk.Coins) error
	SendCoinsFromModuleToModule(ctx context.Context, senderModule string, recipientModule string, amt sdk.Coins) error

	MintCoins(ctx sdk.Context, moduleName string, coins sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, coins sdk.Coins) error

	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetSupply(ctx context.Context, denom string) sdk.Coin
}

type OracleKeeper interface {
	GetPrice(ctx sdk.Context, denom string) (sdk.DecCoin, error)
}

// AccountKeeper is expected keeper for auth module
type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx context.Context, moduleName string) sdk.ModuleAccountI
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI
}

type Keeper struct {
	storeKey   storetypes.StoreKey
	cdc        codec.Codec
	paramSpace paramstypes.Subspace

	accKeeper    AccountKeeper
	bankKeeper   BankKeeper
	oracleKeeper OracleKeeper
}

func NewKeeper(
	cdc codec.Codec,
	key storetypes.StoreKey,
	paramSpace paramstypes.Subspace,
	accKeeper AccountKeeper,
	bankKeeper BankKeeper,
	oracleKeeper OracleKeeper,
) Keeper {
	// ensure stable staking module account is set
	if addr := accKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:     key,
		cdc:          cdc,
		paramSpace:   paramSpace,
		bankKeeper:   bankKeeper,
		oracleKeeper: oracleKeeper,
		accKeeper:    accKeeper,
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

func (k Keeper) GetCollateralKey(denom string) []byte {
	if denom == "" {
		return []byte(fmt.Sprintf("%s:", types.PoolKey))
	}
	return []byte(fmt.Sprintf("%s:%s", types.PoolKey, denom))
}

func (k Keeper) GetTotalCollateral(ctx sdk.Context, denom string) (sdk.Coin, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PoolKey))

	bz := store.Get(k.GetCollateralKey(denom))
	if bz == nil {
		return sdk.Coin{}, fmt.Errorf("total collateral does not exist for denom: %s", denom)
	}

	var result sdk.Coin
	k.cdc.MustUnmarshal(bz, &result)

	return result, nil
}

func (k Keeper) setTotalCollateral(ctx sdk.Context, amount sdk.Coin) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&amount)
	store.Set(k.GetCollateralKey(amount.Denom), bz)
}

func (k Keeper) addCollateral(ctx sdk.Context, amount sdk.Coin) error {
	current, err := k.GetTotalCollateral(ctx, amount.Denom)
	if err != nil {
		return fmt.Errorf("add collateral error: %s", err)
	}

	k.setTotalCollateral(ctx, current.Add(amount))
	return nil
}

func (k Keeper) subCollateral(ctx sdk.Context, amount sdk.Coin) error {
	current, err := k.GetTotalCollateral(ctx, amount.Denom)
	if err != nil {
		return fmt.Errorf("sub collateral error: %s", err)
	}

	if current.IsLT(amount) {
		return fmt.Errorf("insufficient total collateral: %s < %s", current.String(), amount.String())
	}

	k.setTotalCollateral(ctx, current.Sub(amount))
	return nil
}

func (k Keeper) Deposit(ctx sdk.Context, depositor sdk.AccAddress, amount sdk.Coin) (*types.MsgDepositResponse, error) {
	params := k.GetParams(ctx)
	if !types.IsAllowedToken(amount.Denom, params.AllowedAssets) {
		return nil, fmt.Errorf("unsupported token: %s", amount.Denom)
	}

	if amount.IsZero() || amount.IsNegative() {
		return nil, fmt.Errorf("invalid amount")
	}

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, depositor, types.ModuleName, sdk.NewCoins(amount)); err != nil {
		return nil, err
	}

	trackDenom := params.TrackTokenPrefix + amount.Denom
	totalCollateral, _ := k.GetTotalCollateral(ctx, amount.Denom)
	totalLpSupply := k.bankKeeper.GetSupply(ctx, trackDenom)

	var mintAmt osmomath.Int
	if totalCollateral.IsZero() || totalLpSupply.IsZero() {
		mintAmt = amount.Amount
	} else {
		// mint = deposit_amount * total_lp_supply / total_collateral
		mintAmt = amount.Amount.Mul(totalLpSupply.Amount).Quo(totalCollateral.Amount)
		if mintAmt.IsZero() {
			return nil, fmt.Errorf("invalid amount of denom: %s", amount.Denom)
		}
	}

	lpCoin := sdk.NewCoin(trackDenom, mintAmt)

	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(lpCoin)); err != nil {
		return nil, fmt.Errorf("failed to mint tracking coins: %s", err)
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositor, sdk.NewCoins(lpCoin)); err != nil {
		// rollback: try to burn minted coins to avoid inconsistent state
		_ = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(lpCoin))
		return nil, fmt.Errorf("failed to send tracking coins to depositor: %s. Error: %s", depositor, err)
	}

	if totalCollateral.IsZero() || totalLpSupply.IsZero() {
		k.setTotalCollateral(ctx, amount)
	} else {
		if err := k.addCollateral(ctx, amount); err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeDeposit,
		sdk.NewAttribute(types.AttributeKeyDepositor, depositor.String()),
		sdk.NewAttribute(types.AttributeKeyUnderlyingDenom, amount.Denom),
		sdk.NewAttribute(types.AttributeKeyUnderlyingAmount, amount.Amount.String()),
		sdk.NewAttribute(types.AttributeKeyTrackingDenom, trackDenom),
		sdk.NewAttribute(types.AttributeKeyTrackingAmount, mintAmt.String()),
	))

	return &types.MsgDepositResponse{
		Receiver:       depositor.String(),
		TrackingAmount: lpCoin,
	}, nil
}

func (k Keeper) Withdraw(ctx sdk.Context, depositor sdk.AccAddress, amount sdk.Coin) (*types.MsgWithdrawResponse, error) {
	params := k.GetParams(ctx)
	underlyingDenom := amount.Denom[len(params.TrackTokenPrefix):]

	if !types.IsAllowedToken(underlyingDenom, params.AllowedAssets) {
		return nil, fmt.Errorf("unsupported token: %s", amount.Denom)
	}

	if amount.IsZero() || amount.IsNegative() {
		return nil, fmt.Errorf("withdraw amount must be positive")
	}

	if !params.WithdrawalsEnabled {
		return nil, fmt.Errorf("withdrawals are globally disabled")
	}

	lpBalance := k.bankKeeper.GetBalance(ctx, depositor, amount.Denom)

	if lpBalance.Amount.LT(amount.Amount) {
		return nil, fmt.Errorf("insufficient LP tokens: balance %s, required %s", lpBalance.Amount, amount.Amount)
	}

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, depositor, types.ModuleName, sdk.NewCoins(amount)); err != nil {
		return nil, fmt.Errorf("failed moving LP tokens to module: %s", err)
	}

	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(amount)); err != nil {
		return nil, fmt.Errorf("failed burning LP tokens: %s", err)
	}

	totalUnderlying := k.bankKeeper.GetBalance(ctx, k.accKeeper.GetModuleAddress(types.ModuleName), underlyingDenom)
	if totalUnderlying.IsZero() {
		return nil, fmt.Errorf("underlying denom balance is zero")
	}

	totalLPSupply := k.bankKeeper.GetSupply(ctx, amount.Denom)

	if totalLPSupply.IsZero() {
		return nil, fmt.Errorf("lp supply is zero")
	}

	// withdrawAmount = lpAmount * totalUnderlying / totalLPSupply
	withdrawUnderlying := amount.Amount.Mul(totalUnderlying.Amount).Quo(totalLPSupply.Amount)
	withdrawUnderlyingCoin := sdk.NewCoin(underlyingDenom, withdrawUnderlying)

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositor, sdk.NewCoins(withdrawUnderlyingCoin)); err != nil {
		return nil, fmt.Errorf("failed sending underlying asset: %s", err)
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeWithdraw,
		sdk.NewAttribute(types.AttributeKeyDepositor, depositor.String()),
		sdk.NewAttribute(types.AttributeKeyUnderlyingDenom, withdrawUnderlyingCoin.Denom),
		sdk.NewAttribute(types.AttributeKeyUnderlyingAmount, withdrawUnderlyingCoin.Amount.String()),
		sdk.NewAttribute(types.AttributeKeyTrackingDenom, amount.Denom),
		sdk.NewAttribute(types.AttributeKeyTrackingAmount, amount.Amount.String()),
	))

	return &types.MsgWithdrawResponse{
		Receiver:       depositor.String(),
		ReturnedAmount: &withdrawUnderlyingCoin,
	}, nil
}
