package types

import (
	"context"
)

// QueryServer defines the gRPC query service for the EVM module.
// This is a simplified stub for non-protobuf implementation.
type QueryServer interface {
	// Params returns the EVM module parameters.
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	// Balance returns the EVM balance for an address.
	Balance(context.Context, *QueryBalanceRequest) (*QueryBalanceResponse, error)
	// Code returns the contract code for an address.
	Code(context.Context, *QueryCodeRequest) (*QueryCodeResponse, error)
}

// MsgServer defines the gRPC msg service for the EVM module.
// This is a simplified stub for non-protobuf implementation.
type MsgServer interface {
	// EthereumTx processes an Ethereum transaction.
	EthereumTx(context.Context, *MsgEthereumTx) (*MsgEthereumTxResponse, error)
}

// Query request/response types (simplified stubs)

type QueryParamsRequest struct{}

type QueryParamsResponse struct {
	Params Params
}

type QueryBalanceRequest struct {
	Address string
}

type QueryBalanceResponse struct {
	Balance string
}

type QueryCodeRequest struct {
	Address string
}

type QueryCodeResponse struct {
	Code []byte
}

// Msg response types (simplified stubs)

type MsgEthereumTxResponse struct {
	// Hash is the Ethereum transaction hash.
	Hash string
	// ContractAddress is the address of the created contract (if any).
	ContractAddress string
	// Logs contains the event logs emitted by the transaction.
	Logs []*Log
}

// RegisterMsgServer registers the msg server with the provided registrar.
// This is a stub for simplified implementation without gRPC.
func RegisterMsgServer(registrar interface{}, server MsgServer) {
	// Stub implementation - would register gRPC server in full implementation
}

// RegisterQueryServer registers the query server with the provided registrar.
// This is a stub for simplified implementation without gRPC.
func RegisterQueryServer(registrar interface{}, server QueryServer) {
	// Stub implementation - would register gRPC server in full implementation
}
