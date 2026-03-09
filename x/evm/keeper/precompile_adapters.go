package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/v27/x/evm/types"
)

// bankKeeperAdapter adapts types.BankKeeper (which uses context.Context) to the interface
// expected by precompiles (which uses sdk.Context).
type bankKeeperAdapter struct {
	keeper types.BankKeeper
}

func (a *bankKeeperAdapter) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	return a.keeper.GetBalance(ctx, addr, denom)
}

func (a *bankKeeperAdapter) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	return a.keeper.SendCoins(ctx, fromAddr, toAddr, amt)
}

// accountKeeperAdapter adapts types.AccountKeeper (which uses context.Context) to the interface
// expected by precompiles (which uses sdk.Context).
type accountKeeperAdapter struct {
	keeper types.AccountKeeper
}

func (a *accountKeeperAdapter) GetAccount(ctx sdk.Context, addr sdk.AccAddress) sdk.AccountI {
	return a.keeper.GetAccount(ctx, addr)
}

// stakingKeeperAdapter adapts types.StakingKeeper to the interface expected by staking precompile.
type stakingKeeperAdapter struct {
	// For now, we don't wrap StakingKeeper because the staking precompile
	// already defines its own interface that doesn't match types.StakingKeeper.
	// This will be a placeholder for future implementation.
}
