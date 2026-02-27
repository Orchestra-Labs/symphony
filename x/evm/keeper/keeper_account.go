package keeper

import (
	"encoding/hex"
	"math/big"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/osmosis-labs/osmosis/v27/x/evm/types"
)

// GetEVMDenom returns the denomination used for EVM transactions.
// This is "note" with 10^12 scaling to match EVM's 18 decimals.
const EVMDenom = "note"

// Wei per note: 1 note = 10^12 wei (to bridge the 6 decimal -> 18 decimal gap)
var WeiPerNote = new(big.Int).Exp(big.NewInt(10), big.NewInt(12), nil)

// GetOrCreateAccount retrieves an account or creates a new one if it doesn't exist.
func (k Keeper) GetOrCreateAccount(ctx sdk.Context, addr sdk.AccAddress) sdk.AccountI {
	account := k.accountKeeper.GetAccount(ctx, addr)
	if account == nil {
		account = k.accountKeeper.NewAccountWithAddress(ctx, addr)
		k.accountKeeper.SetAccount(ctx, account)
	}
	return account
}

// GetBalance returns the EVM balance (in wei) for an account.
// Converts from note (6 decimals) to wei (18 decimals).
func (k Keeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress) *big.Int {
	coin := k.bankKeeper.GetBalance(ctx, addr, EVMDenom)

	// Convert note to wei: note_amount * 10^12
	noteAmount := coin.Amount.BigInt()
	weiAmount := new(big.Int).Mul(noteAmount, WeiPerNote)

	return weiAmount
}

// SetBalance sets the EVM balance (in wei) for an account.
// Converts from wei (18 decimals) to note (6 decimals).
func (k Keeper) SetBalance(ctx sdk.Context, addr sdk.AccAddress, amount *big.Int) error {
	// Convert wei to note: wei_amount / 10^12
	noteAmount := new(big.Int).Div(amount, WeiPerNote)

	currentBalance := k.bankKeeper.GetBalance(ctx, addr, EVMDenom)
	newBalance := sdk.NewCoin(EVMDenom, math.NewIntFromBigInt(noteAmount))

	if newBalance.Amount.GT(currentBalance.Amount) {
		// Need to mint the difference
		diff := newBalance.Amount.Sub(currentBalance.Amount)
		if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(EVMDenom, diff))); err != nil {
			return err
		}
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.NewCoins(sdk.NewCoin(EVMDenom, diff))); err != nil {
			return err
		}
	} else if newBalance.Amount.LT(currentBalance.Amount) {
		// Need to burn the difference
		diff := currentBalance.Amount.Sub(newBalance.Amount)
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, sdk.NewCoins(sdk.NewCoin(EVMDenom, diff))); err != nil {
			return err
		}
		if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(EVMDenom, diff))); err != nil {
			return err
		}
	}

	return nil
}

// AddBalance adds to the EVM balance (in wei) for an account.
func (k Keeper) AddBalance(ctx sdk.Context, addr sdk.AccAddress, amount *big.Int) error {
	currentBalance := k.GetBalance(ctx, addr)
	newBalance := new(big.Int).Add(currentBalance, amount)
	return k.SetBalance(ctx, addr, newBalance)
}

// SubBalance subtracts from the EVM balance (in wei) for an account.
func (k Keeper) SubBalance(ctx sdk.Context, addr sdk.AccAddress, amount *big.Int) error {
	currentBalance := k.GetBalance(ctx, addr)
	if currentBalance.Cmp(amount) < 0 {
		return types.ErrInsufficientFunds
	}
	newBalance := new(big.Int).Sub(currentBalance, amount)
	return k.SetBalance(ctx, addr, newBalance)
}

// GetNonce returns the nonce for an account.
func (k Keeper) GetNonce(ctx sdk.Context, addr sdk.AccAddress) uint64 {
	account := k.accountKeeper.GetAccount(ctx, addr)
	if account == nil {
		return 0
	}
	return account.GetSequence()
}

// SetNonce sets the nonce for an account.
func (k Keeper) SetNonce(ctx sdk.Context, addr sdk.AccAddress, nonce uint64) error {
	account := k.GetOrCreateAccount(ctx, addr)
	if err := account.SetSequence(nonce); err != nil {
		return err
	}
	k.accountKeeper.SetAccount(ctx, account)
	return nil
}

// GetCodeHash returns the code hash for a contract account.
func (k Keeper) GetCodeHash(ctx sdk.Context, addr sdk.AccAddress) []byte {
	account := k.accountKeeper.GetAccount(ctx, addr)
	if account == nil {
		return nil
	}

	// Check if it's an EthAccount with a code hash
	if ethAcc, ok := account.(*authtypes.EthAccount); ok {
		codeHash := ethAcc.GetCodeHash()
		return codeHash
	}

	return nil
}

// SetCodeHash sets the code hash for a contract account.
func (k Keeper) SetCodeHash(ctx sdk.Context, addr sdk.AccAddress, codeHash []byte) {
	account := k.GetOrCreateAccount(ctx, addr)

	// If it's not an EthAccount, we need to convert it
	ethAcc, ok := account.(*authtypes.EthAccount)
	if !ok {
		// Create a new EthAccount from the existing account
		baseAcc, ok := account.(*authtypes.BaseAccount)
		if !ok {
			// If it's not even a BaseAccount, create a new one
			baseAcc = authtypes.NewBaseAccount(addr, nil, 0, 0)
		}
		ethAcc = authtypes.NewEthAccountWithCode(baseAcc, codeHash)
	} else {
		ethAcc.CodeHash = hex.EncodeToString(codeHash)
	}

	k.accountKeeper.SetAccount(ctx, ethAcc)
}

// HasCode returns true if the account has contract code.
func (k Keeper) HasCode(ctx sdk.Context, addr sdk.AccAddress) bool {
	codeHash := k.GetCodeHash(ctx, addr)
	if codeHash == nil || len(codeHash) == 0 {
		return false
	}

	// Check if it's the empty code hash
	emptyHash, _ := hex.DecodeString(authtypes.EmptyCodeHash)
	if hex.EncodeToString(codeHash) == hex.EncodeToString(emptyHash) {
		return false
	}

	return true
}

// TransferBalance transfers balance from one account to another.
func (k Keeper) TransferBalance(ctx sdk.Context, from, to sdk.AccAddress, amount *big.Int) error {
	if amount.Sign() == 0 {
		return nil
	}

	if amount.Sign() < 0 {
		return types.ErrInvalidValue
	}

	// Convert wei to note for the transfer
	noteAmount := new(big.Int).Div(amount, WeiPerNote)
	coins := sdk.NewCoins(sdk.NewCoin(EVMDenom, math.NewIntFromBigInt(noteAmount)))

	return k.bankKeeper.SendCoins(ctx, from, to, coins)
}
