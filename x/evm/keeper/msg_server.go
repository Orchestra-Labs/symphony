package keeper

import (
	"context"
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/osmosis-labs/osmosis/v27/x/evm/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// EthereumTx processes an Ethereum transaction.
func (m msgServer) EthereumTx(goCtx context.Context, msg *types.MsgEthereumTx) (*MsgEthereumTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate the message
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Check if contract creation is enabled
	if msg.IsContractCreation() && !m.GetEnableCreate(ctx) {
		return nil, types.ErrCreateDisabled
	}

	// Check if contract calls are enabled
	if !msg.IsContractCreation() && !m.GetEnableCall(ctx) {
		return nil, types.ErrCallDisabled
	}

	// Decode the sender address
	fromAddr, err := hex.DecodeString(msg.From)
	if err != nil {
		return nil, types.ErrInvalidAddress
	}
	from := sdk.AccAddress(fromAddr)

	// Check sender balance
	balance := m.GetBalance(ctx, from)
	totalCost := msg.GetValue()
	// Add gas cost: gasLimit * gasPrice
	gasCost := msg.GetGasPrice()
	gasCost.Mul(gasCost, sdk.NewIntFromUint64(msg.Data.Gas).BigInt())
	totalCost.Add(totalCost, gasCost)

	if balance.Cmp(totalCost) < 0 {
		return nil, types.ErrInsufficientFunds
	}

	// Check and increment nonce
	currentNonce := m.GetNonce(ctx, from)
	if msg.Data.Nonce != currentNonce {
		if msg.Data.Nonce < currentNonce {
			return nil, types.ErrNonceTooLow
		}
		return nil, types.ErrNonceTooHigh
	}

	// Increment nonce
	if err := m.SetNonce(ctx, from, currentNonce+1); err != nil {
		return nil, err
	}

	// Deduct gas cost from sender
	if err := m.SubBalance(ctx, from, gasCost); err != nil {
		return nil, err
	}

	var contractAddress string
	var gasUsed uint64 = msg.Data.Gas // Simplified: assume all gas is used

	if msg.IsContractCreation() {
		// Contract creation
		contractAddress, err = m.createContract(ctx, from, msg)
		if err != nil {
			// Refund gas on error (simplified)
			_ = m.AddBalance(ctx, from, gasCost)
			return nil, err
		}
	} else {
		// Contract call or transfer
		toAddr := msg.GetTo()
		if err := m.executeCall(ctx, from, toAddr, msg); err != nil {
			// Refund gas on error (simplified)
			_ = m.AddBalance(ctx, from, gasCost)
			return nil, err
		}
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEthereumTx,
			sdk.NewAttribute(types.AttributeKeyTxHash, msg.ComputeHash()),
			sdk.NewAttribute(types.AttributeKeyRecipient, msg.Data.To),
			sdk.NewAttribute(types.AttributeKeyContractAddress, contractAddress),
			sdk.NewAttribute(types.AttributeKeyTxGasUsed, fmt.Sprintf("%d", gasUsed)),
		),
	)

	return &MsgEthereumTxResponse{
		Hash:            msg.ComputeHash(),
		ContractAddress: contractAddress,
		GasUsed:         gasUsed,
	}, nil
}

// createContract handles contract creation.
// This is a simplified placeholder - full implementation would use go-ethereum's EVM.
func (m msgServer) createContract(ctx sdk.Context, from sdk.AccAddress, msg *types.MsgEthereumTx) (string, error) {
	// Generate contract address (simplified - should use CREATE or CREATE2 logic)
	// For CREATE: keccak256(rlp([sender, nonce]))[12:]
	nonce := m.GetNonce(ctx, from) - 1 // We already incremented it
	contractAddr := sdk.AccAddress(append(from, byte(nonce)))[:20]

	// Create the account
	_ = m.GetOrCreateAccount(ctx, contractAddr)

	// Store the contract code
	// In reality, we would execute the init code and store the deployed code
	code := msg.Data.Input
	codeHash := sdk.Keccak256(code)

	m.SetCode(ctx, codeHash, code)
	m.SetCodeHash(ctx, contractAddr, codeHash)

	// Transfer value if any
	value := msg.GetValue()
	if value.Sign() > 0 {
		if err := m.TransferBalance(ctx, from, contractAddr, value); err != nil {
			return "", err
		}
	}

	return hex.EncodeToString(contractAddr), nil
}

// executeCall handles contract calls and transfers.
// This is a simplified placeholder - full implementation would use go-ethereum's EVM.
func (m msgServer) executeCall(ctx sdk.Context, from sdk.AccAddress, to []byte, msg *types.MsgEthereumTx) error {
	toAddr := sdk.AccAddress(to)

	// Transfer value if any
	value := msg.GetValue()
	if value.Sign() > 0 {
		if err := m.TransferBalance(ctx, from, toAddr, value); err != nil {
			return err
		}
	}

	// If there's input data and the recipient has code, this is a contract call
	if len(msg.Data.Input) > 0 && m.HasCode(ctx, toAddr) {
		// In a full implementation, we would:
		// 1. Create an EVM instance
		// 2. Execute the contract code with the input data
		// 3. Handle state changes, events, and return data
		// For now, this is a placeholder
		m.Logger(ctx).Info("Contract call executed", "to", hex.EncodeToString(to), "data_len", len(msg.Data.Input))
	}

	return nil
}

// Placeholder response type (would be auto-generated from proto in production)
type MsgEthereumTxResponse struct {
	Hash            string
	ContractAddress string
	GasUsed         uint64
}
