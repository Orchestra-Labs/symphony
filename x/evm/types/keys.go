package types

const (
	// ModuleName defines the module name for the EVM module.
	ModuleName = "evm"

	// StoreKey defines the primary module store key for the EVM module.
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key.
	RouterKey = ModuleName

	// TransientKey is the transient store key for the EVM module.
	// Used for temporary state during transaction execution.
	TransientKey = "transient_evm"
)

// KVStore key prefixes
var (
	// KeyPrefixCode stores contract bytecode, keyed by keccak256(code).
	// Key format: KeyPrefixCode + keccak256(code) → bytecode
	KeyPrefixCode = []byte{0x01}

	// KeyPrefixStorage stores contract storage, keyed by address + storage key.
	// Key format: KeyPrefixStorage + address(20 bytes) + storageKey(32 bytes) → value(32 bytes)
	KeyPrefixStorage = []byte{0x02}

	// KeyPrefixParams stores module parameters.
	KeyPrefixParams = []byte{0x03}
)

// AddressStoragePrefix returns the prefix for a contract's storage.
// All storage keys for a given address share this prefix.
func AddressStoragePrefix(address []byte) []byte {
	return append(KeyPrefixStorage, address...)
}
