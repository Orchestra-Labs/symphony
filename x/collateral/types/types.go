package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

type CollateralRecord struct {
	User      sdk.AccAddress `json:"user" yaml:"user"`
	Amount    sdk.DecCoin    `json:"amount" yaml:"amount"`
	Deposited time.Time      `json:"deposited" yaml:"deposited"`
}

// DebtRecord records how many tracking units minted for user
type DebtRecord struct {
	User   sdk.AccAddress `json:"user" yaml:"user"`
	Amount sdk.DecCoin    `json:"amount" yaml:"amount"`
}
