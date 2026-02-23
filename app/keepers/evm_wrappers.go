package keepers

import (
	"context"
	"time"

	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// EVMAccountKeeper is a wrapper for AccountKeeper to satisfy the VM module's AccountKeeper interface
type EVMAccountKeeper struct {
	*authkeeper.AccountKeeper
}

func (k EVMAccountKeeper) UnorderedTransactionsEnabled() bool {
	return false
}

func (k EVMAccountKeeper) RemoveExpiredUnorderedNonces(ctx sdk.Context) error {
	return nil
}

func (k EVMAccountKeeper) TryAddUnorderedNonce(ctx sdk.Context, sender []byte, timestamp time.Time) error {
	return nil
}

func (k EVMAccountKeeper) AddressCodec() address.Codec {
	// Standard AccountKeeper in v0.50 has an AddressCodec or we can get it from the config
	// However, we can just return the bech32 codec here if needed.
	// But let's check if we can get it from the embedded keeper.
    // In Symphony's SDK version, it might be available directly.
	return k.AccountKeeper.AddressCodec()
}

func (k EVMAccountKeeper) GetParams(ctx context.Context) authtypes.Params {
    return k.AccountKeeper.GetParams(ctx)
}
