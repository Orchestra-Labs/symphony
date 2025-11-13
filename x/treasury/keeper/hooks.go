package keeper

import (
	appparams "github.com/osmosis-labs/osmosis/v27/app/params"
	"time"

	epochstypes "github.com/osmosis-labs/osmosis/v27/x/epochs/types"
	"github.com/osmosis-labs/osmosis/v27/x/treasury/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeforeEpochStart is the epoch start hook.
func (k Keeper) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return nil
}

// AfterEpochEnd is the epoch end hook.
func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	// Try to run burn if amount enough once per day
	if epochIdentifier == "day" {
		moduleBurnAddr := k.accountKeeper.GetModuleAddress(types.NativeBurnCollectorName)
		totalAmount := k.BankKeeper.GetBalance(ctx, moduleBurnAddr, appparams.BaseCoinUnit)

		if !totalAmount.IsZero() {
			err := k.BankKeeper.SendCoinsFromModuleToModule(ctx, types.NativeBurnCollectorName, types.ModuleName, sdk.NewCoins(totalAmount))
			if err != nil {
				ctx.Logger().Error("SendCoinsFromModuleToModule", "module", types.NativeBurnCollectorName, "denom", totalAmount.Denom, "amount", totalAmount.Amount, "height", ctx.BlockHeight())
			} else {
				err = k.BankKeeper.BurnCoinsEnable(ctx, types.ModuleName, sdk.Coins{totalAmount})
				if err != nil {
					ctx.Logger().Error("BurnCoinsEnable", "module", types.ModuleName, "denom", totalAmount.Denom, "amount", totalAmount.Amount, "height", ctx.BlockHeight())
				}
			}
		}
	}

	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)
	params := k.GetParams(ctx)

	if params.UpdateTreasuryEpochIdentifier != epochIdentifier {
		return nil
	}

	refillAmount := k.RefillExchangePool(ctx)
	oldTaxRate := k.GetTaxRate(ctx)
	newTaxRate := k.UpdateReserveFee(ctx)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeTaxRateUpdate,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyExchangePoolRefillAmount, refillAmount.String()),
			sdk.NewAttribute(types.AttributeKeyOldTaxRate, oldTaxRate.String()),
			sdk.NewAttribute(types.AttributeKeyNewTaxRate, newTaxRate.String()),
		),
	)

	return nil
}

// ___________________________________________________________________________________________________

// Hooks is the wrapper struct for the incentives keeper.
type Hooks struct {
	k Keeper
}

var _ epochstypes.EpochHooks = Hooks{}

// Hooks returns the hook wrapper struct.
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// GetModuleName implements types.EpochHooks.
func (Hooks) GetModuleName() string {
	return types.ModuleName
}

// BeforeEpochStart is the epoch start hook.
func (h Hooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return h.k.BeforeEpochStart(ctx, epochIdentifier, epochNumber)
}

// AfterEpochEnd is the epoch end hook.
func (h Hooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return h.k.AfterEpochEnd(ctx, epochIdentifier, epochNumber)
}
