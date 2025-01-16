package types

import (
	"fmt"
	"golang.org/x/exp/slices"
	"strings"

	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	ProposalTypeUpdateSellOnly = "UpdateSellOnly"
)

func init() {
	govtypesv1.RegisterProposalType(ProposalTypeUpdateSellOnly)
}

var (
	_ govtypesv1.Content = &UpdateSellOnlyProposal{}
)

// NewUpdateSellOnlyProposal returns a new instance of a proposal struct.
func NewUpdateSellOnlyProposal(title, description string, sellOnlyDenoms, enableBuyDenoms []string) govtypesv1.Content {
	return &UpdateSellOnlyProposal{
		Title:       title,
		Description: description,
		EnableBuy:   enableBuyDenoms,
		SellOnly:    sellOnlyDenoms,
	}
}

func (p *UpdateSellOnlyProposal) GetTitle() string { return p.Title }

// GetDescription gets the description of the proposal
func (p *UpdateSellOnlyProposal) GetDescription() string { return p.Description }

// ProposalRoute returns the router key for the proposal
func (p *UpdateSellOnlyProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of the proposal
func (p *UpdateSellOnlyProposal) ProposalType() string {
	return ProposalTypeUpdateSellOnly
}

// ValidateBasic validates a governance proposal's abstract and basic contents
func (p *UpdateSellOnlyProposal) ValidateBasic() error {
	if len(p.EnableBuy) == 0 && len(p.SellOnly) == 0 {
		return fmt.Errorf("must provide at item to update")
	}

	for i, denom := range p.EnableBuy {
		if slices.Contains(p.EnableBuy[i+1:], denom) {
			return fmt.Errorf("%s is duplicated", denom)
		}
		if slices.Contains(p.SellOnly, denom) {
			return fmt.Errorf("%s is in both SellOnly and EnableBuy", denom)
		}
	}

	for i, denom := range p.SellOnly {
		if slices.Contains(p.SellOnly[i+1:], denom) {
			return fmt.Errorf("%s is duplicated", denom)
		}
	}

	return nil
}

// String returns a string to display the proposal.
func (p UpdateSellOnlyProposal) String() string {
	sellOnly := ""
	for _, denom := range p.SellOnly {
		sellOnly = sellOnly + fmt.Sprintf("(SellOnly: %s) ", denom)
	}

	sellBuy := ""
	for _, denom := range p.EnableBuy {
		sellBuy = sellBuy + fmt.Sprintf("(SellOnly: %s) ", denom)
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Update SellOnly Proposal:
Title:       %s
Description: %s
SellOnly:    %s
SellBuy:	 %s
`, p.Title, p.Description, sellOnly, sellBuy))
	return b.String()
}
