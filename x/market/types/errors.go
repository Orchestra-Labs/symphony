package types

import (
	errorsmod "cosmossdk.io/errors"
)

// Market errors
var (
	ErrRecursiveSwap                  = errorsmod.Register(ModuleName, 2, "recursive swap")
	ErrNoEffectivePrice               = errorsmod.Register(ModuleName, 3, "no price registered with oracle")
	ErrZeroSwapCoin                   = errorsmod.Register(ModuleName, 4, "zero swap coin")
	ErrNotEnoughBalanceOnMarketVaults = errorsmod.Register(ModuleName, 5, "not enough balance on market vaults")
	ErrSellOnlyDenom                  = errorsmod.Register(ModuleName, 6, "this denom cannot be bought at this time")
)
