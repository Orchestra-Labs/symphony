package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/osmosis-labs/osmosis/v26/x/oracle/types"
)

// HandleUpdateSellOnlyProposal updates the sell only list and enable buy list based on the proposal content
func (k Keeper) HandleUpdateSellOnlyProposal(ctx sdk.Context, p *types.UpdateSellOnlyProposal) error {
	for _, denom := range p.SellOnly {
		k.SetSellOnly(ctx, denom, true)
	}
	for _, denom := range p.EnableBuy {
		k.SetSellOnly(ctx, denom, false)
	}
	return nil
}

// NewOracleProposalHandler creates a new oracle proposal handler which will handle gov proposals for oracle module.
func NewOracleProposalHandler(k Keeper) govtypesv1.Handler {
	return func(ctx sdk.Context, content govtypesv1.Content) error {
		switch c := content.(type) {
		case *types.UpdateSellOnlyProposal:
			return k.HandleUpdateSellOnlyProposal(ctx, c)

		default:
			return fmt.Errorf("unrecognized incentives proposal content type: %T", c)
		}
	}
}
