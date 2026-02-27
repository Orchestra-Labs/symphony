package types

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"golang.org/x/crypto/sha3"
)

var (
	_ sdk.Msg = &MsgEthereumTx{}
)

// MsgEthereumTx represents an Ethereum transaction.
// It wraps the Ethereum transaction data and provides Cosmos SDK Msg interface.
type MsgEthereumTx struct {
	// Data is the Ethereum transaction payload.
	// For legacy transactions: RLP([nonce, gasPrice, gas, to, value, data, v, r, s])
	// For EIP-1559: RLP([chainId, nonce, maxPriorityFeePerGas, maxFeePerGas, gas, to, value, data, accessList, v, r, s])
	Data *EthereumTxData `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`

	// Hash is the Ethereum transaction hash (keccak256 of RLP-encoded tx).
	// This is computed on the fly and not stored.
	Hash string `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`

	// From is the Ethereum sender address (20 bytes, hex-encoded).
	// Recovered from the transaction signature.
	From string `protobuf:"bytes,3,opt,name=from,proto3" json:"from,omitempty"`
}

// EthereumTxData contains the Ethereum transaction data fields.
type EthereumTxData struct {
	// ChainID is the EIP-155 chain ID (for replay protection).
	ChainID string `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`

	// Nonce is the sender's transaction count.
	Nonce uint64 `protobuf:"varint,2,opt,name=nonce,proto3" json:"nonce,omitempty"`

	// GasPrice is the gas price in wei (legacy transactions).
	GasPrice string `protobuf:"bytes,3,opt,name=gas_price,json=gasPrice,proto3" json:"gas_price,omitempty"`

	// Gas is the gas limit for the transaction.
	Gas uint64 `protobuf:"varint,4,opt,name=gas,proto3" json:"gas,omitempty"`

	// To is the recipient address (nil for contract creation).
	To string `protobuf:"bytes,5,opt,name=to,proto3" json:"to,omitempty"`

	// Value is the amount to transfer in wei.
	Value string `protobuf:"bytes,6,opt,name=value,proto3" json:"value,omitempty"`

	// Data is the contract call data or contract creation code.
	Input []byte `protobuf:"bytes,7,opt,name=input,proto3" json:"input,omitempty"`

	// V, R, S are the ECDSA signature values.
	V []byte `protobuf:"bytes,8,opt,name=v,proto3" json:"v,omitempty"`
	R []byte `protobuf:"bytes,9,opt,name=r,proto3" json:"r,omitempty"`
	S []byte `protobuf:"bytes,10,opt,name=s,proto3" json:"s,omitempty"`
}

// NewMsgEthereumTx creates a new MsgEthereumTx.
func NewMsgEthereumTx(
	chainID *big.Int,
	nonce uint64,
	to string,
	amount *big.Int,
	gasLimit uint64,
	gasPrice *big.Int,
	input []byte,
	v, r, s []byte,
) *MsgEthereumTx {
	return &MsgEthereumTx{
		Data: &EthereumTxData{
			ChainID:  chainID.String(),
			Nonce:    nonce,
			GasPrice: gasPrice.String(),
			Gas:      gasLimit,
			To:       to,
			Value:    amount.String(),
			Input:    input,
			V:        v,
			R:        r,
			S:        s,
		},
	}
}

// Route implements sdk.Msg
func (msg MsgEthereumTx) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg MsgEthereumTx) Type() string {
	return EventTypeEthereumTx
}

// ValidateBasic implements sdk.Msg
func (msg MsgEthereumTx) ValidateBasic() error {
	if msg.Data == nil {
		return fmt.Errorf("transaction data cannot be nil")
	}

	// Validate chain ID
	if msg.Data.ChainID == "" {
		return ErrInvalidChainID
	}

	// Validate gas
	if msg.Data.Gas == 0 {
		return ErrInvalidGasLimit
	}

	// Validate gas price
	gasPrice := new(big.Int)
	if _, ok := gasPrice.SetString(msg.Data.GasPrice, 10); !ok {
		return ErrInvalidGasPrice
	}
	if gasPrice.Sign() < 0 {
		return ErrInvalidGasPrice
	}

	// Validate value
	value := new(big.Int)
	if _, ok := value.SetString(msg.Data.Value, 10); !ok {
		return ErrInvalidValue
	}
	if value.Sign() < 0 {
		return ErrInvalidValue
	}

	// Validate signature
	if len(msg.Data.V) == 0 || len(msg.Data.R) == 0 || len(msg.Data.S) == 0 {
		return ErrInvalidSignature
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgEthereumTx) GetSignBytes() []byte {
	// Ethereum transactions are signed differently than Cosmos transactions.
	// We return the keccak256 hash of the RLP-encoded transaction data.
	bz := msg.Data.SignBytes()
	hash := sha3.NewLegacyKeccak256()
	hash.Write(bz)
	return hash.Sum(nil)
}

// GetSigners implements sdk.Msg
func (msg MsgEthereumTx) GetSigners() []sdk.AccAddress {
	// The signer is the 'from' address recovered from the signature.
	// This is populated during signature verification.
	if msg.From == "" {
		return []sdk.AccAddress{}
	}

	fromBz, err := hex.DecodeString(msg.From)
	if err != nil {
		return []sdk.AccAddress{}
	}

	return []sdk.AccAddress{sdk.AccAddress(fromBz)}
}

// GetFrom returns the sender address as a byte slice.
func (msg *MsgEthereumTx) GetFrom() []byte {
	if msg.From == "" {
		return nil
	}
	fromBz, _ := hex.DecodeString(msg.From)
	return fromBz
}

// GetTo returns the recipient address as a byte slice (nil for contract creation).
func (msg *MsgEthereumTx) GetTo() []byte {
	if msg.Data.To == "" {
		return nil
	}
	toBz, _ := hex.DecodeString(msg.Data.To)
	return toBz
}

// GetGasPrice returns the gas price as a big.Int.
func (msg *MsgEthereumTx) GetGasPrice() *big.Int {
	gasPrice := new(big.Int)
	gasPrice.SetString(msg.Data.GasPrice, 10)
	return gasPrice
}

// GetValue returns the transaction value as a big.Int.
func (msg *MsgEthereumTx) GetValue() *big.Int {
	value := new(big.Int)
	value.SetString(msg.Data.Value, 10)
	return value
}

// IsContractCreation returns true if the transaction creates a contract.
func (msg *MsgEthereumTx) IsContractCreation() bool {
	return msg.Data.To == ""
}

// ComputeHash computes the Ethereum transaction hash.
func (msg *MsgEthereumTx) ComputeHash() string {
	bz := msg.Data.SignBytes()
	hash := sha3.NewLegacyKeccak256()
	hash.Write(bz)
	return hex.EncodeToString(hash.Sum(nil))
}

// SignBytes returns the bytes to be signed for this transaction data.
// This is a simplified version - full implementation would handle EIP-155, EIP-1559, etc.
func (data *EthereumTxData) SignBytes() []byte {
	// For simplicity, we're hashing the concatenation of fields.
	// Real implementation would use proper RLP encoding.
	h := sha256.New()
	h.Write([]byte(data.ChainID))
	h.Write([]byte(fmt.Sprintf("%d", data.Nonce)))
	h.Write([]byte(data.GasPrice))
	h.Write([]byte(fmt.Sprintf("%d", data.Gas)))
	h.Write([]byte(data.To))
	h.Write([]byte(data.Value))
	h.Write(data.Input)
	return h.Sum(nil)
}

// String implements fmt.Stringer
func (msg MsgEthereumTx) String() string {
	return fmt.Sprintf(`MsgEthereumTx{
  Hash: %s
  From: %s
  To: %s
  Value: %s
  Gas: %d
  GasPrice: %s
  Nonce: %d
}`, msg.Hash, msg.From, msg.Data.To, msg.Data.Value, msg.Data.Gas, msg.Data.GasPrice, msg.Data.Nonce)
}
