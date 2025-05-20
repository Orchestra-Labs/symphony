package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	epochstypes "github.com/osmosis-labs/osmosis/v27/x/epochs/types"
	txfeestypes "github.com/osmosis-labs/osmosis/v27/x/txfees/types"
)

type Hooks struct {
	k Keeper
}

var (
	_ epochstypes.EpochHooks = Hooks{}
)

// Hooks creates new pool incentives hooks.
func (k Keeper) Hooks() Hooks { return Hooks{k} }

// GetModuleName implements types.EpochHooks.
func (Hooks) GetModuleName() string {
	return txfeestypes.ModuleName
}

func (h Hooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return h.k.BeforeEpochStart(ctx, epochIdentifier, epochNumber)
}

func (h Hooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return h.k.AfterEpochEnd(ctx, epochIdentifier, epochNumber)
}

func (k Keeper) BeforeEpochStart(_ctx sdk.Context, _epochIdentifier string, _epochNumber int64) error {
	return nil
}

// AfterEpochEnd at the end of each epoch, take snapshot and distribute rewards to Stakers for previous epoch
func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, _epochNumber int64) error {
	if epochIdentifier != k.GetParams(ctx).EpochIdentifier {
		return nil
	}

	//TODO:
	// 1 - take snapshot
	// 2 - distribute rewards to Stakers
	//defaultFeesDenom, _ := k.GetBaseDenom(ctx)
	//nonNativefeeTokenCollectorAddress := k.accountKeeper.GetModuleAddress(txfeestypes.NonNativeTxFeeCollectorName)
	//
	//// Non-native fee token collector for staking rewards get swapped entirely into base denom.
	//k.swapNonNativeFeeToDenom(ctx, defaultFeesDenom, nonNativefeeTokenCollectorAddress)
	//
	//// Now that the rewards have been swapped, transfer any base denom existing in the non-native tx fee collector to the auth fee token collector (indirectly distributing to stakers)
	//baseDenomCoins := sdk.NewCoins(k.bankKeeper.GetBalance(ctx, nonNativefeeTokenCollectorAddress, defaultFeesDenom))
	//err := osmoutils.ApplyFuncIfNoError(ctx, func(cacheCtx sdk.Context) error {
	//	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, txfeestypes.NonNativeTxFeeCollectorName, authtypes.FeeCollectorName, baseDenomCoins)
	//	return err
	//})
	//if err != nil {
	//	incTelementryCounter(txfeestypes.TakerFeeFailedNativeRewardUpdateMetricName, baseDenomCoins.String(), err.Error())
	//}

	return nil
}
