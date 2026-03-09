package simulation

import (
	abci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
)

// RequestBeginBlock and ResponseEndBlock are compatibility types
// for the simulation framework. In CometBFT v0.38+, these were replaced
// with RequestFinalizeBlock/ResponseFinalizeBlock as part of ABCI++.
// These types allow simulation code to continue working during the transition.

// RequestBeginBlock is a compatibility type that maps to RequestFinalizeBlock
type RequestBeginBlock = abci.RequestFinalizeBlock

// ResponseEndBlock is a compatibility type
type ResponseEndBlock struct {
	ValidatorUpdates      []abci.ValidatorUpdate
	ConsensusParamUpdates *cmtproto.ConsensusParams
	Events                []abci.Event
}
