package types

import (
	"encoding/json"
	"fmt"
	"github.com/osmosis-labs/osmosis/osmomath"

	"github.com/cosmos/cosmos-sdk/codec"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, taxRate osmomath.Dec) *GenesisState {
	return &GenesisState{
		Params:  params,
		TaxRate: taxRate,
	}
}

// DefaultGenesisState gets raw genesis raw message for testing
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:  DefaultParams(),
		TaxRate: DefaultTaxRate,
	}
}

// ValidateGenesis validates the provided oracle genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds, no duplicate validators)
func ValidateGenesis(data *GenesisState) error {
	if data.TaxRate.LT(osmomath.ZeroDec()) {
		return fmt.Errorf("tax_rate must be positive, is %s", data.TaxRate)
	}
	if data.TaxRate.GT(data.Params.MaxFeeMultiplier) {
		return fmt.Errorf("tax_rate must less than RateMax(%s)", data.Params.MaxFeeMultiplier)
	}

	return data.Params.Validate()
}

// GetGenesisStateFromAppState returns x/market GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(cdc codec.JSONCodec, appState map[string]json.RawMessage) *GenesisState {
	var genesisState GenesisState

	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}

	return &genesisState
}
