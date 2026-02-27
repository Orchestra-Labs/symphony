package eth

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/osmosis-labs/osmosis/v27/x/evm/keeper"
	"github.com/osmosis-labs/osmosis/v27/x/evm/types"
)

// API provides Ethereum-compatible JSON-RPC methods.
type API struct {
	clientCtx client.Context
	evmKeeper *keeper.Keeper
}

// NewAPI creates a new Ethereum API instance.
func NewAPI(clientCtx client.Context, evmKeeper *keeper.Keeper) *API {
	return &API{
		clientCtx: clientCtx,
		evmKeeper: evmKeeper,
	}
}

// ChainId returns the EVM chain ID.
func (api *API) ChainId() (string, error) {
	// Get the chain ID from the keeper
	ctx := sdk.UnwrapSDKContext(api.clientCtx.Context())
	chainID := api.evmKeeper.GetChainID(ctx)

	// Convert to hex with 0x prefix
	chainIDInt := new(big.Int)
	chainIDInt.SetString(chainID, 10)
	return fmt.Sprintf("0x%x", chainIDInt), nil
}

// BlockNumber returns the current block number.
func (api *API) BlockNumber() (string, error) {
	ctx := sdk.UnwrapSDKContext(api.clientCtx.Context())
	height := ctx.BlockHeight()
	return fmt.Sprintf("0x%x", height), nil
}

// GetBalance returns the balance of an account at a given block.
func (api *API) GetBalance(address string, blockNum string) (string, error) {
	ctx := sdk.UnwrapSDKContext(api.clientCtx.Context())

	// Remove 0x prefix if present
	address = strings.TrimPrefix(address, "0x")

	// Decode address
	addrBytes, err := hex.DecodeString(address)
	if err != nil {
		return "", fmt.Errorf("invalid address: %w", err)
	}

	// Get balance (returns wei)
	balance := api.evmKeeper.GetBalance(ctx, sdk.AccAddress(addrBytes))

	// Return as hex string with 0x prefix
	return fmt.Sprintf("0x%x", balance), nil
}

// GetTransactionCount returns the nonce of an account.
func (api *API) GetTransactionCount(address string, blockNum string) (string, error) {
	ctx := sdk.UnwrapSDKContext(api.clientCtx.Context())

	// Remove 0x prefix if present
	address = strings.TrimPrefix(address, "0x")

	// Decode address
	addrBytes, err := hex.DecodeString(address)
	if err != nil {
		return "", fmt.Errorf("invalid address: %w", err)
	}

	// Get nonce
	nonce := api.evmKeeper.GetNonce(ctx, sdk.AccAddress(addrBytes))

	// Return as hex string with 0x prefix
	return fmt.Sprintf("0x%x", nonce), nil
}

// GetCode returns the code at a given address.
func (api *API) GetCode(address string, blockNum string) (string, error) {
	ctx := sdk.UnwrapSDKContext(api.clientCtx.Context())

	// Remove 0x prefix if present
	address = strings.TrimPrefix(address, "0x")

	// Decode address
	addrBytes, err := hex.DecodeString(address)
	if err != nil {
		return "", fmt.Errorf("invalid address: %w", err)
	}

	// Get code hash
	codeHash := api.evmKeeper.GetCodeHash(ctx, sdk.AccAddress(addrBytes))
	if codeHash == nil {
		return "0x", nil
	}

	// Get code
	code := api.evmKeeper.GetCode(ctx, codeHash)
	if len(code) == 0 {
		return "0x", nil
	}

	// Return as hex string with 0x prefix
	return "0x" + hex.EncodeToString(code), nil
}

// SendRawTransaction broadcasts a raw Ethereum transaction.
func (api *API) SendRawTransaction(data string) (string, error) {
	// Remove 0x prefix if present
	data = strings.TrimPrefix(data, "0x")

	// Decode the transaction
	txBytes, err := hex.DecodeString(data)
	if err != nil {
		return "", fmt.Errorf("invalid transaction data: %w", err)
	}

	// In a full implementation, we would:
	// 1. Decode the RLP-encoded Ethereum transaction
	// 2. Create a MsgEthereumTx
	// 3. Broadcast it via the Cosmos SDK client

	// For now, return a placeholder transaction hash
	txHash := sdk.Keccak256(txBytes)
	return "0x" + hex.EncodeToString(txHash), nil
}

// GasPrice returns the current gas price.
func (api *API) GasPrice() (string, error) {
	// In a full implementation, this would query the feemarket module
	// For now, return a placeholder (1 gwei)
	gasPrice := big.NewInt(1000000000) // 1 gwei
	return fmt.Sprintf("0x%x", gasPrice), nil
}

// Call executes a contract call (read-only).
func (api *API) Call(args map[string]interface{}, blockNum string) (string, error) {
	// In a full implementation, this would:
	// 1. Parse the call arguments (to, data, gas, gasPrice, value)
	// 2. Execute the call against the EVM
	// 3. Return the result

	// For now, return placeholder
	return "0x", nil
}

// EstimateGas estimates gas needed for a transaction.
func (api *API) EstimateGas(args map[string]interface{}) (string, error) {
	// In a full implementation, this would simulate the transaction
	// and return the gas used

	// For now, return a standard transfer gas cost
	gasEstimate := big.NewInt(21000)
	return fmt.Sprintf("0x%x", gasEstimate), nil
}

// GetTransactionReceipt returns the receipt of a transaction by hash.
func (api *API) GetTransactionReceipt(hash string) (*types.TxReceipt, error) {
	// In a full implementation, this would:
	// 1. Query the transaction by hash
	// 2. Build and return a receipt with logs, status, etc.

	// For now, return nil (transaction not found)
	return nil, nil
}

// GetBlockByNumber returns a block by number.
func (api *API) GetBlockByNumber(blockNum string, fullTx bool) (map[string]interface{}, error) {
	ctx := sdk.UnwrapSDKContext(api.clientCtx.Context())

	// Parse block number
	var height int64
	if blockNum == "latest" {
		height = ctx.BlockHeight()
	} else {
		// Parse hex block number
		blockNum = strings.TrimPrefix(blockNum, "0x")
		bigHeight := new(big.Int)
		bigHeight.SetString(blockNum, 16)
		height = bigHeight.Int64()
	}

	// In a full implementation, we would fetch the actual block data
	// For now, return a minimal block structure
	block := map[string]interface{}{
		"number":     fmt.Sprintf("0x%x", height),
		"hash":       fmt.Sprintf("0x%064x", height), // Placeholder
		"parentHash": fmt.Sprintf("0x%064x", height-1),
		"timestamp":  fmt.Sprintf("0x%x", ctx.BlockTime().Unix()),
		"gasLimit":   "0x1c9c380", // 30M gas
		"gasUsed":    "0x0",
		"miner":      "0x0000000000000000000000000000000000000000",
		"difficulty": "0x0",
		"transactions": []interface{}{},
	}

	return block, nil
}
