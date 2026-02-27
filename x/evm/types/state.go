package types

import (
	"encoding/hex"
	"fmt"
)

// State represents the EVM state for an account.
type State struct {
	// Address is the account address (20 bytes).
	Address string `json:"address"`
	// Balance is the account balance in wei.
	Balance string `json:"balance"`
	// Nonce is the account nonce.
	Nonce uint64 `json:"nonce"`
	// CodeHash is the hash of the contract code (if any).
	CodeHash string `json:"code_hash"`
	// Storage is the contract storage (key-value pairs).
	Storage map[string]string `json:"storage,omitempty"`
}

// NewState creates a new State for an account.
func NewState(address string, balance string, nonce uint64) *State {
	return &State{
		Address:  address,
		Balance:  balance,
		Nonce:    nonce,
		CodeHash: "",
		Storage:  make(map[string]string),
	}
}

// Validate performs basic validation of the state.
func (s State) Validate() error {
	// Validate address format (40 hex characters = 20 bytes)
	if s.Address != "" {
		if len(s.Address) != 40 {
			return fmt.Errorf("invalid address length: expected 40, got %d", len(s.Address))
		}
		if _, err := hex.DecodeString(s.Address); err != nil {
			return fmt.Errorf("invalid address hex: %w", err)
		}
	}

	// Validate code hash if present (64 hex characters = 32 bytes)
	if s.CodeHash != "" {
		if len(s.CodeHash) != 64 {
			return fmt.Errorf("invalid code hash length: expected 64, got %d", len(s.CodeHash))
		}
		if _, err := hex.DecodeString(s.CodeHash); err != nil {
			return fmt.Errorf("invalid code hash hex: %w", err)
		}
	}

	return nil
}

// GenesisState defines the evm module's genesis state.
type GenesisState struct {
	// Params are the module parameters.
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	// Accounts contains the initial EVM account states.
	Accounts []State `protobuf:"bytes,2,rep,name=accounts,proto3" json:"accounts"`
}

// DefaultGenesisState returns the default genesis state for the evm module.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:   DefaultParams(),
		Accounts: []State{},
	}
}

// Validate performs basic validation of genesis data.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}

	// Validate all account states
	addresses := make(map[string]bool)
	for i, account := range gs.Accounts {
		if err := account.Validate(); err != nil {
			return fmt.Errorf("invalid account at index %d: %w", i, err)
		}

		// Check for duplicate addresses
		if addresses[account.Address] {
			return fmt.Errorf("duplicate account address: %s", account.Address)
		}
		addresses[account.Address] = true
	}

	return nil
}

// ProtoMessage implements proto.Message interface (simplified stub for non-protobuf implementation)
func (gs *GenesisState) ProtoMessage() {}

// Reset implements proto.Message interface (simplified stub for non-protobuf implementation)
func (gs *GenesisState) Reset() {
	*gs = GenesisState{}
}

// String implements fmt.Stringer
func (gs GenesisState) String() string {
	return fmt.Sprintf("GenesisState:\n  Params: %s\n  Accounts: %d", gs.Params.String(), len(gs.Accounts))
}

// StorageEntry represents a single storage entry for a contract.
type StorageEntry struct {
	// Key is the storage key (32 bytes, hex-encoded).
	Key string `json:"key"`
	// Value is the storage value (32 bytes, hex-encoded).
	Value string `json:"value"`
}

// ValidateStorageKey validates a storage key format.
func ValidateStorageKey(key string) error {
	if len(key) != 64 {
		return fmt.Errorf("invalid storage key length: expected 64, got %d", len(key))
	}
	if _, err := hex.DecodeString(key); err != nil {
		return fmt.Errorf("invalid storage key hex: %w", err)
	}
	return nil
}

// ValidateStorageValue validates a storage value format.
func ValidateStorageValue(value string) error {
	if len(value) != 64 {
		return fmt.Errorf("invalid storage value length: expected 64, got %d", len(value))
	}
	if _, err := hex.DecodeString(value); err != nil {
		return fmt.Errorf("invalid storage value hex: %w", err)
	}
	return nil
}
