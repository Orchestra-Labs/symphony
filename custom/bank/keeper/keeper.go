package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/osmosis-labs/osmosis/osmomath"
	appparams "github.com/osmosis-labs/osmosis/v27/app/params"
	treasurytypes "github.com/osmosis-labs/osmosis/v27/x/treasury/types"
	"golang.org/x/exp/slices"
)

var _ bankkeeper.Keeper = (*CustomKeeper)(nil)

type CustomKeeper struct {
	*bankkeeper.BaseKeeper
	ak AccountKeeper
}

func (k *CustomKeeper) AddSupplyOffset(ctx context.Context, denom string, offsetAmount osmomath.Int) {
	//TODO implement me
	panic("implement me")
}

func (k *CustomKeeper) SendCoinsFromModuleToManyAccounts(ctx context.Context, senderModule string, recipientAddrs []sdk.AccAddress, amts []sdk.Coins) error {
	if len(recipientAddrs) != len(amts) {
		panic(fmt.Errorf("addresses and amounts numbers does not match"))
	}

	senderAddr := k.ak.GetModuleAddress(senderModule)
	if senderAddr == nil {
		panic(errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", senderModule))
	}

	for _, recipientAddr := range recipientAddrs {
		if k.BlockedAddr(recipientAddr) {
			return errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", recipientAddr)
		}
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	for i, toAddr := range recipientAddrs {
		err := k.SendCoins(ctx, senderAddr, toAddr, amts[i])
		if err != nil {
			return err
		}
		sdkCtx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				EventTypeTransfer,
				sdk.NewAttribute(AttributeKeyRecipient, toAddr.String()),
				sdk.NewAttribute(AttributeKeySender, senderAddr.String()),
				sdk.NewAttribute(sdk.AttributeKeyAmount, amts[i].String()),
			),
		})
	}
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(AttributeKeySender, senderAddr.String()),
	))

	return nil
}

func NewCustomKeeper(baseBankKeeper *bankkeeper.BaseKeeper, ak AccountKeeper) CustomKeeper {
	return CustomKeeper{
		BaseKeeper: baseBankKeeper,
		ak:         ak,
	}
}

func (k *CustomKeeper) BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error {
	acc := k.ak.GetModuleAccount(ctx, moduleName)
	if acc == nil {
		panic(errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", moduleName))
	}

	if !acc.HasPermission(authtypes.Burner) {
		panic(errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "module account %s does not have permissions to burn tokens", moduleName))
	}
	index := slices.IndexFunc(amt, func(coin sdk.Coin) bool {
		return coin.Denom == appparams.BaseCoinUnit
	})
	// instead of burning, we would like to send the coins to the reserve.
	if index >= 0 {
		nativeCoin := amt[index]

		err := k.SendCoinsFromModuleToModule(ctx, moduleName, treasurytypes.ModuleName, sdk.NewCoins(nativeCoin))
		if err != nil {
			return fmt.Errorf("failed to send coins to the reserve on burn: %w", err)
		}

		// proceed with base logic but exclude native coin
		amt = slices.Delete(amt, index, index+1)
	}
	return k.BaseKeeper.BurnCoins(ctx, moduleName, amt)
}

func (k *CustomKeeper) BurnCoinsEnable(ctx context.Context, moduleName string, amt sdk.Coins) error {
	acc := k.ak.GetModuleAccount(ctx, moduleName)
	if acc == nil {
		panic(errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", moduleName))
	}

	if !acc.HasPermission(authtypes.Burner) {
		panic(errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "module account %s does not have permissions to burn tokens", moduleName))
	}

	return k.BaseKeeper.BurnCoins(ctx, moduleName, amt)
}
