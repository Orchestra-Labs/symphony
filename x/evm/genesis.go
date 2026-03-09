package evm

import (
	"encoding/hex"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/osmosis-labs/osmosis/v27/x/evm/keeper"
	"github.com/osmosis-labs/osmosis/v27/x/evm/types"
)

// InitGenesis initializes the module's state from a genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set module parameters
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}

	// Initialize accounts and their storage
	for _, account := range genState.Accounts {
		// Decode the address
		addr, err := hex.DecodeString(account.Address)
		if err != nil {
			panic(err)
		}
		accAddr := sdk.AccAddress(addr)

		// Set balance if present
		if account.Balance != "" {
			// Balance is in wei (string format big.Int)
			// SetBalance will convert wei to note
			balanceInt, ok := math.NewIntFromString(account.Balance)
			if !ok {
				panic(fmt.Errorf("invalid balance for account %s: %s", account.Address, account.Balance))
			}
			if err := k.SetBalance(ctx, accAddr, balanceInt.BigInt()); err != nil {
				panic(err)
			}
		}

		// Set nonce
		if account.Nonce > 0 {
			if err := k.SetNonce(ctx, accAddr, account.Nonce); err != nil {
				panic(err)
			}
		}

		// Set code hash if present (contract account)
		if account.CodeHash != "" {
			codeHash, err := hex.DecodeString(account.CodeHash)
			if err != nil {
				panic(err)
			}
			k.SetCodeHash(ctx, accAddr, codeHash)
		}

		// Initialize storage if present
		if account.Storage != nil {
			for key, value := range account.Storage {
				keyBz, err := hex.DecodeString(key)
				if err != nil {
					panic(err)
				}
				valueBz, err := hex.DecodeString(value)
				if err != nil {
					panic(err)
				}
				k.SetStorage(ctx, accAddr, keyBz, valueBz)
			}
		}
	}
}

// ExportGenesis exports the module's state to a genesis state.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	params := k.GetParams(ctx)

	// In a full implementation, we would iterate through all accounts
	// and export their state. For now, we return the params only.
	return &types.GenesisState{
		Params:   params,
		Accounts: []types.State{},
	}
}
