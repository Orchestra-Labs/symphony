package types

// EVM module event types
const (
	// EventTypeEthereumTx defines the event type for Ethereum transactions.
	EventTypeEthereumTx = "ethereum_tx"

	// EventTypeBlockBloom defines the event type for block bloom filters.
	EventTypeBlockBloom = "block_bloom"

	// AttributeKeyContractAddress defines the attribute key for contract addresses.
	AttributeKeyContractAddress = "contract_address"

	// AttributeKeyRecipient defines the attribute key for transaction recipients.
	AttributeKeyRecipient = "recipient"

	// AttributeKeyTxHash defines the attribute key for transaction hashes.
	AttributeKeyTxHash = "tx_hash"

	// AttributeKeyEthereumTxHash defines the attribute key for Ethereum transaction hashes.
	AttributeKeyEthereumTxHash = "ethereum_tx_hash"

	// AttributeKeyTxIndex defines the attribute key for transaction index.
	AttributeKeyTxIndex = "tx_index"

	// AttributeKeyTxGasUsed defines the attribute key for gas used by a transaction.
	AttributeKeyTxGasUsed = "tx_gas_used"

	// AttributeKeyTxType defines the attribute key for transaction type.
	AttributeKeyTxType = "tx_type"

	// AttributeKeyTxLog defines the attribute key for transaction logs (JSON).
	AttributeKeyTxLog = "tx_log"

	// AttributeValueCategory defines the value for the module category.
	AttributeValueCategory = ModuleName

	// MetricKeyTransitionDB defines the metric key for StateDB transitions.
	MetricKeyTransitionDB = "transition_db"

	// MetricKeyStaticCall defines the metric key for static calls.
	MetricKeyStaticCall = "static_call"
)
