package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary x/evm interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgEthereumTx{}, "evm/MsgEthereumTx", nil)
}

// RegisterInterfaces registers the x/evm interfaces types with the interface registry
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgEthereumTx{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	amino.Seal()
}

// Placeholder for the gRPC service descriptor (will be auto-generated from proto in production)
var _Msg_serviceDesc = struct {
	ServiceName string
	HandlerType interface{}
	Methods     []interface{}
	Streams     []interface{}
	Metadata    string
}{
	ServiceName: "symphony.evm.v1.Msg",
	Metadata:    "symphony/evm/v1/tx.proto",
}
