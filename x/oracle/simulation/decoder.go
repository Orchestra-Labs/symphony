package simulation

import (
	"bytes"
	"fmt"

	gogotypes "github.com/gogo/protobuf/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/osmosis-labs/osmosis/v27/x/oracle/types"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding oracle type.
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.ExchangeRateKey):
			var exchangeRateA, exchangeRateB sdk.DecProto
			cdc.MustUnmarshal(kvA.Value, &exchangeRateA)
			cdc.MustUnmarshal(kvB.Value, &exchangeRateB)
			return fmt.Sprintf("%v\n%v", exchangeRateA, exchangeRateB)
		case bytes.Equal(kvA.Key[:1], types.FeederDelegationKey):
			return fmt.Sprintf("%v\n%v", sdk.AccAddress(kvA.Value), sdk.AccAddress(kvB.Value))
		case bytes.Equal(kvA.Key[:1], types.MissCounterKey):
			var counterA, counterB gogotypes.UInt64Value
			cdc.MustUnmarshal(kvA.Value, &counterA)
			cdc.MustUnmarshal(kvB.Value, &counterB)
			return fmt.Sprintf("%v\n%v", counterA.Value, counterB.Value)
		case bytes.Equal(kvA.Key[:1], types.AggregateExchangeRatePrevoteKey):
			var prevoteA, prevoteB types.AggregateExchangeRatePrevote
			cdc.MustUnmarshal(kvA.Value, &prevoteA)
			cdc.MustUnmarshal(kvB.Value, &prevoteB)
			return fmt.Sprintf("%v\n%v", prevoteA, prevoteB)
		case bytes.Equal(kvA.Key[:1], types.AggregateExchangeRateVoteKey):
			var voteA, voteB types.AggregateExchangeRateVote
			cdc.MustUnmarshal(kvA.Value, &voteA)
			cdc.MustUnmarshal(kvB.Value, &voteB)
			return fmt.Sprintf("%v\n%v", voteA, voteB)
		case bytes.Equal(kvA.Key[:1], types.TobinTaxKey):
			var tobinTaxA, tobinTaxB sdk.DecProto
			cdc.MustUnmarshal(kvA.Value, &tobinTaxA)
			cdc.MustUnmarshal(kvB.Value, &tobinTaxB)
			return fmt.Sprintf("%v\n%v", tobinTaxA, tobinTaxB)
		default:
			panic(fmt.Sprintf("invalid oracle key prefix %X", kvA.Key[:1]))
		}
	}
}
