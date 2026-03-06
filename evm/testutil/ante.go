package testutil

import sdk "github.com/cosmos/cosmos-sdk/types"

// NoOpNextFn is a no-op function that returns the context and no error in order to mock
// the next function in the AnteHandler chain.
//
// It can be used in unit tests when calling a decorator's AnteHandle method, e.g.
// `dec.AnteHandle(ctx, tx, false, NoOpNextFn)`
func NoOpNextFn(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) {
	return ctx, nil
}
