package simulation

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding market type.
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		// The market module stores only params in the params subspace,
		// and doesn't have its own KV store entries to decode.
		// Return a default message for any keys encountered.
		return fmt.Sprintf("market key: %X\nvalueA: %X\nvalueB: %X", kvA.Key, kvA.Value, kvB.Value)
	}
}
