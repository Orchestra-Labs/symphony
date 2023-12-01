package authenticator

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/CosmWasm/wasmd/x/wasm/keeper"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	txsigning "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/osmosis-labs/osmosis/v20/x/authenticator/iface"
)

type CosmwasmAuthenticator struct {
	contractKeeper *keeper.PermissionedKeeper
	ak             *authkeeper.AccountKeeper
	cdc            codectypes.AnyUnpacker
	sigModeHandler authsigning.SignModeHandler

	contractAddr sdk.AccAddress
}

var (
	_ iface.Authenticator = &CosmwasmAuthenticator{}
)

func NewCosmwasmAuthenticator(contractKeeper *keeper.PermissionedKeeper, accountKeeper *authkeeper.AccountKeeper, sigModeHandler authsigning.SignModeHandler, cdc codectypes.AnyUnpacker) CosmwasmAuthenticator {
	return CosmwasmAuthenticator{
		contractKeeper: contractKeeper,
		sigModeHandler: sigModeHandler,
		ak:             accountKeeper,
		cdc:            cdc,
	}
}

func (cwa CosmwasmAuthenticator) Type() string {
	return "CosmwasmAuthenticator"
}

func (cwa CosmwasmAuthenticator) StaticGas() uint64 {
	return 0
}

type CosmwasmAuthenticatorInitData struct {
	Contract string `json:"contract"`
}

func (cwa CosmwasmAuthenticator) Initialize(data []byte) (iface.Authenticator, error) {
	var initData CosmwasmAuthenticatorInitData
	err := json.Unmarshal(data, &initData)
	if err != nil {
		return nil, err
	}
	if len(initData.Contract) == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "missing contract address")
	}
	contractAddr, err := sdk.AccAddressFromBech32(initData.Contract)
	if err != nil {
		return nil, err
	}
	cwa.contractAddr = contractAddr
	return cwa, nil
}

func (cwa CosmwasmAuthenticator) GetAuthenticationData(
	ctx sdk.Context,
	tx sdk.Tx,
	messageIndex int,
	simulate bool,
) (iface.AuthenticatorData, error) {
	signers, signatures, signingTx, err := GetCommonAuthenticationData(ctx, tx, messageIndex, simulate)
	if err != nil {
		return SignatureData{}, err
	}

	// Get the signature for the message at msgIndex
	return SignatureData{
		Signers:    signers,
		Signatures: signatures,
		Tx:         signingTx,
		Simulate:   simulate,
	}, nil
}

type SignModeData struct {
	Direct  []byte `json:"sign_mode_direct"`
	Textual string `json:"sign_mode_textual"`
}

type LocalAny struct {
	TypeURL string `json:"type_url"`
	Value   []byte `json:"value"`
}

type ExplicitTxData struct {
	ChainID       string     `json:"chain_id"`
	AccountNumber uint64     `json:"account_number"`
	Sequence      uint64     `json:"sequence"`
	TimeoutHeight uint64     `json:"timeout_height"`
	Msgs          []LocalAny `json:"msgs"`
	Memo          string     `json:"memo"`
	Simulate      bool       `json:"simulate"`
}

type simplifiedSignatureData struct {
	Signers    []sdk.AccAddress `json:"signers"`
	Signatures [][]byte         `json:"signatures"`
}

type AuthenticateSudoMsg struct {
	Authenticate AuthenticateMsg `json:"authenticate"`
}

type AuthenticateMsg struct {
	Account        sdk.AccAddress          `json:"account"`
	Msg            LocalAny                `json:"msg"`
	Signature      []byte                  `json:"signature"` // Only allowing messages with a single signer
	SignModeTxData SignModeData            `json:"sign_mode_tx_data"`
	TxData         ExplicitTxData          `json:"tx_data"`
	SignatureData  simplifiedSignatureData `json:"signature_data"`
}

type AuthenticationResult struct {
}

// TODO: decide if we want to reject or just not authenticate
func (cwa CosmwasmAuthenticator) Authenticate(ctx sdk.Context, account sdk.AccAddress, msg sdk.Msg, authenticationData iface.AuthenticatorData) iface.AuthenticationResult {
	if len(msg.GetSigners()) != 1 {
		return iface.Rejected("only messages with a single signer are supported", sdkerrors.ErrInvalidType)
	}

	signatureData, ok := authenticationData.(SignatureData)
	if !ok {
		return iface.Rejected("invalid signature verification data", sdkerrors.ErrInvalidType)
	}

	// Retrieve and build the signer data struct
	genesis := ctx.IsGenesis() || ctx.BlockHeight() == 0
	chainID := ctx.ChainID()
	var accNum uint64
	baseAccount := cwa.ak.GetAccount(ctx, account)
	if !genesis {
		accNum = baseAccount.GetAccountNumber()
	}

	signerData := authsigning.SignerData{
		ChainID:       chainID,
		AccountNumber: accNum,
		Sequence:      baseAccount.GetSequence(),
	}

	signBytes, err := cwa.sigModeHandler.GetSignBytes(txsigning.SignMode_SIGN_MODE_DIRECT, signerData, signatureData.Tx)
	if err != nil {
		return iface.Rejected("failed to get signBytes", err)
	}

	encodedMsg, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return iface.Rejected("failed to encode msg", err)
	}

	timeoutTx, ok := signatureData.Tx.(sdk.TxWithTimeoutHeight)
	if !ok {
		return iface.Rejected("failed to cast tx to TxWithTimeoutHeight", sdkerrors.ErrInvalidType)
	}
	memoTx, ok := signatureData.Tx.(sdk.TxWithMemo)
	if !ok {
		return iface.Rejected("failed to cast tx to TxWithMemo", sdkerrors.ErrInvalidType)
	}

	msgs := make([]LocalAny, len(signatureData.Tx.GetMsgs()))
	for i, msg := range signatureData.Tx.GetMsgs() {
		encodedMsg, err := codectypes.NewAnyWithValue(msg)
		if err != nil {
			return iface.Rejected("failed to encode msg", err)
		}
		msgs[i] = LocalAny{
			TypeURL: encodedMsg.TypeUrl,
			Value:   encodedMsg.Value,
		}
	}

	txData := ExplicitTxData{
		ChainID:       chainID,
		AccountNumber: accNum,
		Sequence:      baseAccount.GetSequence(),
		TimeoutHeight: timeoutTx.GetTimeoutHeight(),
		Msgs:          msgs,
		Memo:          memoTx.GetMemo(),
		Simulate:      signatureData.Simulate,
	}

	signer := msg.GetSigners()[0]
	var signatures [][]byte
	var msgSignature []byte
	for i, signature := range signatureData.Signatures {
		single, ok := signature.Data.(*txsigning.SingleSignatureData)
		if !ok {
			return iface.Rejected("failed to cast signature to SingleSignatureData", sdkerrors.ErrInvalidType)
		}
		signatures = append(signatures, single.Signature)
		if signatureData.Signers[i].Equals(signer) {
			msgSignature = single.Signature
		}
	}

	authMsg := AuthenticateMsg{
		Account: account,
		Msg: LocalAny{
			TypeURL: encodedMsg.TypeUrl,
			Value:   encodedMsg.Value,
		},
		Signature: msgSignature, // TODO: currently only allowing one signer per message.
		TxData:    txData,
		SignModeTxData: SignModeData{ // TODO: Add other sign modes. Specifically textual when it becomes available
			Direct: signBytes,
		},
		SignatureData: simplifiedSignatureData{
			Signers:    signatureData.Signers,
			Signatures: signatures,
		},
	}
	bz, err := json.Marshal(AuthenticateSudoMsg{authMsg})
	fmt.Println(string(bz))
	if err != nil {
		return iface.Rejected("failed to marshall AuthenticateMsg", err)
	}

	// sha256 of signbytes
	hash := sha256.New()
	hash.Write(signBytes)
	hashBz := hash.Sum(nil)

	fmt.Println("msg hash", hashBz)
	signatureHex := fmt.Sprintf("%x", msgSignature)
	fmt.Println("signature", signatureHex)

	result, err := cwa.contractKeeper.Sudo(ctx, cwa.contractAddr, bz)
	if err != nil {
		return iface.Rejected("failed to sudo", err)
	}

	authResult, err := UnmarshalAuthenticationResult(result)
	if err != nil {
		return iface.Rejected("failed to unmarshal authentication result", err)
	}
	return authResult
}

func (cwa CosmwasmAuthenticator) Track(ctx sdk.Context, account sdk.AccAddress, msg sdk.Msg) error {
	return nil
}

func (cwa CosmwasmAuthenticator) ConfirmExecution(ctx sdk.Context, account sdk.AccAddress, msg sdk.Msg, authenticationData iface.AuthenticatorData) iface.ConfirmationResult {
	return iface.Confirm()
}

func (cwa CosmwasmAuthenticator) OnAuthenticatorAdded(ctx sdk.Context, account sdk.AccAddress, data []byte) error {
	_, err := sdk.AccAddressFromBech32(string(data))
	if err != nil {
		return err
	}
	// TODO: check contract address length. Check contract exists?
	return nil
}

func (cwa CosmwasmAuthenticator) OnAuthenticatorRemoved(ctx sdk.Context, account sdk.AccAddress, data []byte) error {
	return nil
}

func UnmarshalAuthenticationResult(data []byte) (iface.AuthenticationResult, error) {
	// Unmarshal type field
	var rawType struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &rawType); err != nil {
		return nil, err
	}

	switch rawType.Type { // usign snake case here because that's what cosmwasm defaults to
	case "authenticated":
		return iface.Authenticated(), nil
	case "not_authenticated":
		return iface.NotAuthenticated(), nil
	case "rejected":
		var content struct {
			Msg string `json:"msg"`
		}
		if err := json.Unmarshal(data, &content); err != nil {
			return nil, err
		}
		return iface.Rejected(content.Msg, fmt.Errorf("cosmwasm contract error")), nil
	default:
		return nil, fmt.Errorf("invalid authentication result type: %s", rawType.Type)
	}
}