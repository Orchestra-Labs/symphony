package ante_test

import (
	"github.com/osmosis-labs/osmosis/v27/app/apptesting/assets"
	"testing"

	"github.com/osmosis-labs/osmosis/osmomath"

	"github.com/stretchr/testify/suite"

	"cosmossdk.io/log"
	dbm "github.com/cosmos/cosmos-db"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/osmosis-labs/osmosis/v27/app"
	treasurytypes "github.com/osmosis-labs/osmosis/v27/x/treasury/types"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
)

// AnteTestSuite is a test suite to be used with ante handler tests.
type AnteTestSuite struct {
	suite.Suite

	app *app.SymphonyApp
	// anteHandler sdk.AnteHandler
	ctx       sdk.Context
	clientCtx client.Context
	txBuilder client.TxBuilder
}

// returns context and app with params set on account keeper
func createTestApp(isCheckTx bool, tempDir string) (*app.SymphonyApp, sdk.Context) {
	// TODO: we need to feed in custom binding here?
	var wasmOpts []wasmkeeper.Option
	app := app.NewSymphonyApp(
		log.NewNopLogger(), dbm.NewMemDB(), nil, true, map[int64]bool{},
		tempDir, 0,
		simtestutil.EmptyAppOptions{}, wasmOpts,
	)
	ctx := app.BaseApp.NewContext(isCheckTx)
	app.TreasuryKeeper.SetParams(ctx, treasurytypes.DefaultParams())
	app.TreasuryKeeper.SetTaxRate(ctx, osmomath.NewDecWithPrec(1, 2)) // 0.01%
	app.OracleKeeper.SetTobinTax(ctx, assets.MicroSDRDenom, osmomath.NewDecWithPrec(1, 2))
	app.OracleKeeper.SetTobinTax(ctx, assets.MicroUSDDenom, osmomath.NewDecWithPrec(1, 2))
	app.OracleKeeper.SetTobinTax(ctx, assets.MicroHKDDenom, osmomath.NewDecWithPrec(1, 2))
	err := app.DistrKeeper.FeePool.Set(ctx, distributiontypes.InitialFeePool())
	if err != nil {
		return nil, sdk.Context{}
	}

	return app, ctx
}

// SetupTest setups a new test, with new app, context, and anteHandler.
func (suite *AnteTestSuite) SetupTest(isCheckTx bool) {
	tempDir := suite.T().TempDir()
	suite.app, suite.ctx = createTestApp(isCheckTx, tempDir)
	suite.ctx = suite.ctx.WithBlockHeight(1)

	// Set up TxConfig.
	encodingConfig := suite.SetupEncoding()

	suite.clientCtx = client.Context{}.
		WithTxConfig(encodingConfig.TxConfig)
}

func (suite *AnteTestSuite) SetupEncoding() testutil.TestEncodingConfig {
	encodingConfig := testutil.MakeTestEncodingConfig()
	// We're using TestMsg encoding in some tests, so register it here.
	encodingConfig.Amino.RegisterConcrete(&testdata.TestMsg{}, "testdata.TestMsg", nil)
	testdata.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	return encodingConfig
}

// CreateTestTx is a helper function to create a tx given multiple inputs.
func (suite *AnteTestSuite) CreateTestTx(privs []cryptotypes.PrivKey, accNums []uint64, accSeqs []uint64, chainID string) (xauthsigning.Tx, error) {
	// First round: we gather all the signer infos. We use the "set empty
	// signature" hack to do that.
	var sigsV2 []signing.SignatureV2
	for i, priv := range privs {
		sigV2 := signing.SignatureV2{
			PubKey: priv.PubKey(),
			Data: &signing.SingleSignatureData{
				SignMode:  signing.SignMode(suite.clientCtx.TxConfig.SignModeHandler().DefaultMode()),
				Signature: nil,
			},
			Sequence: accSeqs[i],
		}

		sigsV2 = append(sigsV2, sigV2)
	}
	err := suite.txBuilder.SetSignatures(sigsV2...)
	if err != nil {
		return nil, err
	}

	// Second round: all signer infos are set, so each signer can sign.
	sigsV2 = []signing.SignatureV2{}
	for i, priv := range privs {
		signerData := xauthsigning.SignerData{
			ChainID:       chainID,
			AccountNumber: accNums[i],
			Sequence:      accSeqs[i],
		}
		sigV2, err := tx.SignWithPrivKey(sdk.Context{},
			signing.SignMode(suite.clientCtx.TxConfig.SignModeHandler().DefaultMode()),
			signerData,
			suite.txBuilder,
			priv,
			suite.clientCtx.TxConfig,
			accSeqs[i])
		if err != nil {
			return nil, err
		}

		sigsV2 = append(sigsV2, sigV2)
	}
	err = suite.txBuilder.SetSignatures(sigsV2...)
	if err != nil {
		return nil, err
	}

	return suite.txBuilder.GetTx(), nil
}

func TestAnteTestSuite(t *testing.T) {
	suite.Run(t, new(AnteTestSuite))
}

// func generatePubKeysAndSignatures(n int, msg []byte, _ bool) (pubkeys []cryptotypes.PubKey, signatures [][]byte) {
// 	pubkeys = make([]cryptotypes.PubKey, n)
// 	signatures = make([][]byte, n)
// 	for i := 0; i < n; i++ {
// 		var privkey cryptotypes.PrivKey = secp256k1.GenPrivKey()

// 		// TODO: also generate ed25519 keys as below when ed25519 keys are
// 		//  actually supported, https://github.com/cosmos/cosmos-sdk/issues/4789
// 		// for now this fails:
// 		// if rand.Int63()%2 == 0 {
// 		//	privkey = ed25519.GenPrivKey()
// 		// } else {
// 		//	privkey = secp256k1.GenPrivKey()
// 		// }

// 		pubkeys[i] = privkey.PubKey()
// 		signatures[i], _ = privkey.Sign(msg)
// 	}
// 	return
// }

// func expectedGasCostByKeys(pubkeys []cryptotypes.PubKey) uint64 {
// 	cost := uint64(0)
// 	for _, pubkey := range pubkeys {
// 		pubkeyType := strings.ToLower(fmt.Sprintf("%T", pubkey))
// 		switch {
// 		case strings.Contains(pubkeyType, "ed25519"):
// 			cost += authtypes.DefaultParams().SigVerifyCostED25519
// 		case strings.Contains(pubkeyType, "secp256k1"):
// 			cost += authtypes.DefaultParams().SigVerifyCostSecp256k1
// 		default:
// 			panic("unexpected key type")
// 		}
// 	}
// 	return cost
// }
