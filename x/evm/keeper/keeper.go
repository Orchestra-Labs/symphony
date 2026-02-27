package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/osmosis-labs/osmosis/v27/x/evm/types"
)

// Keeper manages the EVM module state.
type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	paramSpace paramtypes.Subspace

	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	stakingKeeper types.StakingKeeper

	// hooks for other modules to listen to EVM events
	hooks types.EvmHooks
}

// NewKeeper creates a new EVM Keeper instance.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	paramSpace paramtypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		paramSpace:    paramSpace,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		stakingKeeper: stakingKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// SetHooks sets the EVM hooks.
func (k *Keeper) SetHooks(hooks types.EvmHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set EVM hooks twice")
	}
	k.hooks = hooks
	return k
}

// GetCode retrieves the contract bytecode by its code hash.
func (k Keeper) GetCode(ctx sdk.Context, codeHash []byte) []byte {
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixCode, codeHash...)
	return store.Get(key)
}

// SetCode stores contract bytecode, keyed by its keccak256 hash.
func (k Keeper) SetCode(ctx sdk.Context, codeHash []byte, code []byte) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixCode, codeHash...)
	store.Set(key, code)
}

// GetStorage retrieves a storage value for a given contract address and key.
func (k Keeper) GetStorage(ctx sdk.Context, address sdk.AccAddress, key []byte) []byte {
	store := ctx.KVStore(k.storeKey)
	storageKey := append(types.AddressStoragePrefix(address), key...)
	return store.Get(storageKey)
}

// SetStorage stores a storage value for a given contract address and key.
func (k Keeper) SetStorage(ctx sdk.Context, address sdk.AccAddress, key []byte, value []byte) {
	store := ctx.KVStore(k.storeKey)
	storageKey := append(types.AddressStoragePrefix(address), key...)
	if len(value) == 0 {
		store.Delete(storageKey)
	} else {
		store.Set(storageKey, value)
	}
}

// DeleteStorage removes a storage entry for a given contract address and key.
func (k Keeper) DeleteStorage(ctx sdk.Context, address sdk.AccAddress, key []byte) {
	store := ctx.KVStore(k.storeKey)
	storageKey := append(types.AddressStoragePrefix(address), key...)
	store.Delete(storageKey)
}

// IterateStorage iterates over all storage entries for a given address.
func (k Keeper) IterateStorage(ctx sdk.Context, address sdk.AccAddress, cb func(key, value []byte) bool) {
	store := ctx.KVStore(k.storeKey)
	prefix := types.AddressStoragePrefix(address)

	iterator := storetypes.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		// Remove the prefix to get the actual storage key
		key := iterator.Key()[len(prefix):]
		value := iterator.Value()

		if cb(key, value) {
			break
		}
	}
}
