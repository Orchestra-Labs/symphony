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
	params := k.GetParams(ctx)
	if epochIdentifier != params.EpochIdentifier {
		return nil
	}

	// 1. Take a snapshot of current active stakers
	k.SnapshotCurrentEpoch(ctx)

	// 2. Distribute rewards to stakers from last snapshot
	rewardPerEpoch := k.GetEpochReward(ctx)
	if !rewardPerEpoch.IsZero() {
		k.DistributeRewardsToLastEpochStakers(ctx, rewardPerEpoch)
	}

	return nil
}
