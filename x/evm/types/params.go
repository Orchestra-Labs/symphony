package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	ParamStoreKeyEnableCreate = []byte("EnableCreate")
	ParamStoreKeyEnableCall   = []byte("EnableCall")
	ParamStoreKeyExtraEIPs    = []byte("ExtraEIPs")
	ParamStoreKeyChainID      = []byte("ChainID")
)

// ParamKeyTable returns the parameter key table for the EVM module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// Params defines the parameters for the EVM module.
type Params struct {
	// EnableCreate enables contract creation (CREATE/CREATE2 opcodes).
	EnableCreate bool `protobuf:"varint,1,opt,name=enable_create,json=enableCreate,proto3" json:"enable_create,omitempty"`
	// EnableCall enables contract calls (CALL/CALLCODE/DELEGATECALL/STATICCALL opcodes).
	EnableCall bool `protobuf:"varint,2,opt,name=enable_call,json=enableCall,proto3" json:"enable_call,omitempty"`
	// ExtraEIPs defines the additional EIPs to enable beyond the default set.
	ExtraEIPs []int64 `protobuf:"varint,3,rep,packed,name=extra_eips,json=extraEips,proto3" json:"extra_eips,omitempty"`
	// ChainID is the EVM chain ID (used for EIP-155 replay protection).
	ChainID string `protobuf:"bytes,4,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
}

// DefaultParams returns default EVM parameters.
func DefaultParams() Params {
	return Params{
		EnableCreate: true,
		EnableCall:   true,
		ExtraEIPs:    []int64{},
		ChainID:      "9999", // Default to 9999, should be configured before mainnet
	}
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyEnableCreate, &p.EnableCreate, validateBool),
		paramtypes.NewParamSetPair(ParamStoreKeyEnableCall, &p.EnableCall, validateBool),
		paramtypes.NewParamSetPair(ParamStoreKeyExtraEIPs, &p.ExtraEIPs, validateEIPs),
		paramtypes.NewParamSetPair(ParamStoreKeyChainID, &p.ChainID, validateChainID),
	}
}

// Validate performs basic validation on EVM parameters.
func (p Params) Validate() error {
	if err := validateChainID(p.ChainID); err != nil {
		return err
	}
	if err := validateEIPs(p.ExtraEIPs); err != nil {
		return err
	}
	return nil
}

func validateBool(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateEIPs(i interface{}) error {
	eips, ok := i.([]int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	// Validate that EIP numbers are positive
	for _, eip := range eips {
		if eip <= 0 {
			return fmt.Errorf("invalid EIP number: %d", eip)
		}
	}
	return nil
}

func validateChainID(i interface{}) error {
	chainID, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if chainID == "" {
		return fmt.Errorf("chain ID cannot be empty")
	}
	return nil
}

// String implements fmt.Stringer
func (p Params) String() string {
	return fmt.Sprintf(`EVM Params:
  EnableCreate: %t
  EnableCall:   %t
  ExtraEIPs:    %v
  ChainID:      %s
`, p.EnableCreate, p.EnableCall, p.ExtraEIPs, p.ChainID)
}
