package ante_test

import (
	"encoding/json"
	"fmt"
	"github.com/osmosis-labs/osmosis/osmomath"
	"github.com/osmosis-labs/osmosis/v27/ante"
	"github.com/osmosis-labs/osmosis/v27/app/apptesting/assets"
	"os"
	"time"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authz "github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	markettypes "github.com/osmosis-labs/osmosis/v27/x/market/types"
	oracletypes "github.com/osmosis-labs/osmosis/v27/x/oracle/types"
)

func (suite *AnteTestSuite) TestDeductFeeDecorator_ZeroGas() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	mfd := ante.NewDeductFeeDecorator(suite.app.TxFeesKeeper, suite.app.AccountKeeper, suite.app.BankKeeper, nil, suite.app.TreasuryKeeper, suite.app.OracleKeeper)
	antehandler := sdk.ChainAnteDecorators(mfd)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()
	coins := sdk.NewCoins(sdk.NewCoin("atom", osmomath.NewInt(300)))
	testutil.FundAccount(suite.ctx, suite.app.BankKeeper, addr1, coins)

	// msg and signatures
	msg := testdata.NewTestMsg(addr1)
	suite.Require().NoError(suite.txBuilder.SetMsgs(msg))

	// set zero gas
	suite.txBuilder.SetGasLimit(0)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// Set IsCheckTx to true
	suite.ctx = suite.ctx.WithIsCheckTx(true)

	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().Error(err)

	// zero gas is accepted in simulation mode
	_, err = antehandler(suite.ctx, tx, true)
	suite.Require().NoError(err)
}

func (suite *AnteTestSuite) TestEnsureMempoolFees() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	mfd := ante.NewDeductFeeDecorator(suite.app.TxFeesKeeper, suite.app.AccountKeeper, suite.app.BankKeeper, nil, suite.app.TreasuryKeeper, suite.app.OracleKeeper)
	antehandler := sdk.ChainAnteDecorators(mfd)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()
	coins := sdk.NewCoins(sdk.NewCoin("atom", osmomath.NewInt(300)))
	testutil.FundAccount(suite.ctx, suite.app.BankKeeper, addr1, coins)

	// msg and signatures
	msg := testdata.NewTestMsg(addr1)
	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := uint64(15)
	suite.Require().NoError(suite.txBuilder.SetMsgs(msg))
	suite.txBuilder.SetFeeAmount(feeAmount)
	suite.txBuilder.SetGasLimit(gasLimit)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// Set high gas price so standard test fee fails
	atomPrice := sdk.NewDecCoinFromDec("atom", osmomath.NewDec(20))
	highGasPrice := []sdk.DecCoin{atomPrice}
	suite.ctx = suite.ctx.WithMinGasPrices(highGasPrice)

	// Set IsCheckTx to true
	suite.ctx = suite.ctx.WithIsCheckTx(true)

	// antehandler errors with insufficient fees
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NotNil(err, "Decorator should have errored on too low fee for local gasPrice")

	// antehandler should not error since we do not check minGasPrice in simulation mode
	cacheCtx, _ := suite.ctx.CacheContext()
	_, err = antehandler(cacheCtx, tx, true)
	suite.Require().Nil(err, "Decorator should not have errored in simulation mode")

	// Set IsCheckTx to false
	suite.ctx = suite.ctx.WithIsCheckTx(false)

	// antehandler should not error since we do not check minGasPrice in DeliverTx
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().Nil(err, "MempoolFeeDecorator returned error in DeliverTx")

	// Set IsCheckTx back to true for testing sufficient mempool fee
	suite.ctx = suite.ctx.WithIsCheckTx(true)

	atomPrice = sdk.NewDecCoinFromDec("atom", osmomath.NewDec(0).Quo(osmomath.NewDec(100000)))
	lowGasPrice := []sdk.DecCoin{atomPrice}
	suite.ctx = suite.ctx.WithMinGasPrices(lowGasPrice)

	newCtx, err := antehandler(suite.ctx, tx, false)
	suite.Require().Nil(err, "Decorator should not have errored on fee higher than local gasPrice")
	// Priority is the smallest gas price amount in any denom. Since we have only 1 gas price
	// of 10atom, the priority here is 10.
	suite.Require().Equal(int64(10), newCtx.Priority())
}

func (suite *AnteTestSuite) TestDeductFees() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()

	// msg and signatures
	msg := testdata.NewTestMsg(addr1)
	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()
	suite.Require().NoError(suite.txBuilder.SetMsgs(msg))
	suite.txBuilder.SetFeeAmount(feeAmount)
	suite.txBuilder.SetGasLimit(gasLimit)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// Set account with insufficient funds
	acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addr1)
	suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
	coins := sdk.NewCoins(sdk.NewCoin("atom", osmomath.NewInt(10)))
	err = testutil.FundAccount(suite.ctx, suite.app.BankKeeper, addr1, coins)
	suite.Require().NoError(err)

	dfd := ante.NewDeductFeeDecorator(suite.app.TxFeesKeeper, suite.app.AccountKeeper, suite.app.BankKeeper, nil, suite.app.TreasuryKeeper, suite.app.OracleKeeper)
	antehandler := sdk.ChainAnteDecorators(dfd)

	_, err = antehandler(suite.ctx, tx, false)

	suite.Require().NotNil(err, "Tx did not error when fee payer had insufficient funds")

	// Set account with sufficient funds
	suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
	err = testutil.FundAccount(suite.ctx, suite.app.BankKeeper, addr1, sdk.NewCoins(sdk.NewCoin("atom", osmomath.NewInt(200))))
	suite.Require().NoError(err)

	_, err = antehandler(suite.ctx, tx, false)

	suite.Require().Nil(err, "Tx errored after account has been set with sufficient funds")
}

func (suite *AnteTestSuite) TestEnsureMempoolFeesSend() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	mfd := ante.NewDeductFeeDecorator(suite.app.TxFeesKeeper, suite.app.AccountKeeper, suite.app.BankKeeper, nil, suite.app.TreasuryKeeper, suite.app.OracleKeeper)
	antehandler := sdk.ChainAnteDecorators(mfd)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()
	coins := sdk.NewCoins(sdk.NewCoin(assets.MicroSDRDenom, osmomath.NewInt(1000000)))
	testutil.FundAccount(suite.ctx, suite.app.BankKeeper, addr1, coins)

	// msg and signatures
	sendAmount := int64(1000000)
	sendCoins := sdk.NewCoins(sdk.NewInt64Coin(assets.MicroSDRDenom, sendAmount))
	msg := banktypes.NewMsgSend(addr1, addr1, sendCoins)

	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()
	suite.Require().NoError(suite.txBuilder.SetMsgs(msg))
	suite.txBuilder.SetFeeAmount(feeAmount)
	suite.txBuilder.SetGasLimit(gasLimit)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// set zero gas prices
	suite.ctx = suite.ctx.WithMinGasPrices(sdk.NewDecCoins())

	// Set IsCheckTx to true
	suite.ctx = suite.ctx.WithIsCheckTx(true)

	// antehandler errors with insufficient fees due to tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().Error(err, "Decorator should errored on low fee for local gasPrice + tax")

	tk := suite.app.TreasuryKeeper
	expectedTax := tk.GetTaxRate(suite.ctx).MulInt64(sendAmount).TruncateInt()

	// set tax amount
	feeAmount = sdk.NewCoins(sdk.NewCoin(assets.MicroSDRDenom, expectedTax))
	suite.txBuilder.SetFeeAmount(feeAmount)
	tx, err = suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// must pass with tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "Decorator should not have errored on fee higher than local gasPrice")
}

func (suite *AnteTestSuite) TestEnsureMempoolFeesSwapSend() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	mfd := ante.NewDeductFeeDecorator(suite.app.TxFeesKeeper, suite.app.AccountKeeper, suite.app.BankKeeper, nil, suite.app.TreasuryKeeper, suite.app.OracleKeeper)
	antehandler := sdk.ChainAnteDecorators(mfd)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()
	coins := sdk.NewCoins(sdk.NewCoin(assets.MicroSDRDenom, osmomath.NewInt(1000000)))
	testutil.FundAccount(suite.ctx, suite.app.BankKeeper, addr1, coins)

	// msg and signatures
	sendAmount := int64(1000000)
	sendCoin := sdk.NewInt64Coin(assets.MicroSDRDenom, sendAmount)
	msg := markettypes.NewMsgSwapSend(addr1, addr1, sendCoin, assets.MicroKRWDenom)

	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()
	suite.Require().NoError(suite.txBuilder.SetMsgs(msg))
	suite.txBuilder.SetFeeAmount(feeAmount)
	suite.txBuilder.SetGasLimit(gasLimit)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// set zero gas prices
	suite.ctx = suite.ctx.WithMinGasPrices(sdk.NewDecCoins())

	// Set IsCheckTx to true
	suite.ctx = suite.ctx.WithIsCheckTx(true)

	// antehandler errors with insufficient fees due to tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().Error(err, "Decorator should errored on low fee for local gasPrice + tax")

	tk := suite.app.TreasuryKeeper
	expectedTax := tk.GetTaxRate(suite.ctx).MulInt64(sendAmount).TruncateInt()

	// set tax amount
	suite.txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(assets.MicroSDRDenom, expectedTax)))
	tx, err = suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// must pass with tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "Decorator should not have errored on fee higher than local gasPrice")
}

func (suite *AnteTestSuite) TestEnsureMempoolFeesMultiSend() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	mfd := ante.NewDeductFeeDecorator(suite.app.TxFeesKeeper, suite.app.AccountKeeper, suite.app.BankKeeper, nil, suite.app.TreasuryKeeper, suite.app.OracleKeeper)
	antehandler := sdk.ChainAnteDecorators(mfd)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()
	coins := sdk.NewCoins(sdk.NewCoin(assets.MicroSDRDenom, osmomath.NewInt(1000000)))
	testutil.FundAccount(suite.ctx, suite.app.BankKeeper, addr1, coins)

	// msg and signatures
	sendAmount := int64(1000000)
	sendCoins := sdk.NewCoins(sdk.NewInt64Coin(assets.MicroSDRDenom, sendAmount))
	msg := banktypes.NewMsgMultiSend(
		banktypes.NewInput(addr1, sendCoins),
		[]banktypes.Output{
			banktypes.NewOutput(addr1, sendCoins),
			banktypes.NewOutput(addr1, sendCoins),
		},
	)

	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()
	suite.Require().NoError(suite.txBuilder.SetMsgs(msg))
	suite.txBuilder.SetFeeAmount(feeAmount)
	suite.txBuilder.SetGasLimit(gasLimit)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// set zero gas prices
	suite.ctx = suite.ctx.WithMinGasPrices(sdk.NewDecCoins())

	// Set IsCheckTx to true
	suite.ctx = suite.ctx.WithIsCheckTx(true)

	// antehandler errors with insufficient fees due to tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().Error(err, "Decorator should errored on low fee for local gasPrice + tax")

	tk := suite.app.TreasuryKeeper
	expectedTax := tk.GetTaxRate(suite.ctx).MulInt64(sendAmount).TruncateInt()

	// set tax amount
	suite.txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(assets.MicroSDRDenom, expectedTax)))
	tx, err = suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().Error(err, "Decorator should errored on low fee for local gasPrice + tax")

	// must pass with tax
	suite.txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(assets.MicroSDRDenom, expectedTax.Add(expectedTax))))
	tx, err = suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "Decorator should not have errored on fee higher than local gasPrice")
}

func (suite *AnteTestSuite) TestEnsureMempoolFeesInstantiateContract() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	mfd := ante.NewDeductFeeDecorator(suite.app.TxFeesKeeper, suite.app.AccountKeeper, suite.app.BankKeeper, nil, suite.app.TreasuryKeeper, suite.app.OracleKeeper)
	antehandler := sdk.ChainAnteDecorators(mfd)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()
	coins := sdk.NewCoins(sdk.NewCoin(assets.MicroSDRDenom, osmomath.NewInt(1000000)))
	testutil.FundAccount(suite.ctx, suite.app.BankKeeper, addr1, coins)

	// msg and signatures
	sendAmount := int64(1000000)
	sendCoins := sdk.NewCoins(sdk.NewInt64Coin(assets.MicroSDRDenom, sendAmount))
	msg := &wasmtypes.MsgInstantiateContract{
		Sender: addr1.String(),
		Admin:  addr1.String(),
		CodeID: 0,
		Msg:    []byte{},
		Funds:  sendCoins,
	}

	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()
	suite.Require().NoError(suite.txBuilder.SetMsgs(msg))
	suite.txBuilder.SetFeeAmount(feeAmount)
	suite.txBuilder.SetGasLimit(gasLimit)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// set zero gas prices
	suite.ctx = suite.ctx.WithMinGasPrices(sdk.NewDecCoins())

	// Set IsCheckTx to true
	suite.ctx = suite.ctx.WithIsCheckTx(true)

	// antehandler errors with insufficient fees due to tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().Error(err, "Decorator should errored on low fee for local gasPrice + tax")

	tk := suite.app.TreasuryKeeper
	expectedTax := tk.GetTaxRate(suite.ctx).MulInt64(sendAmount).TruncateInt()

	// set tax amount
	suite.txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(assets.MicroSDRDenom, expectedTax)))
	tx, err = suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// must pass with tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "Decorator should not have errored on fee higher than local gasPrice")
}

func (suite *AnteTestSuite) TestEnsureMempoolFeesExecuteContract() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	mfd := ante.NewDeductFeeDecorator(suite.app.TxFeesKeeper, suite.app.AccountKeeper, suite.app.BankKeeper, nil, suite.app.TreasuryKeeper, suite.app.OracleKeeper)
	antehandler := sdk.ChainAnteDecorators(mfd)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()
	coins := sdk.NewCoins(sdk.NewCoin(assets.MicroSDRDenom, osmomath.NewInt(1000000)))
	testutil.FundAccount(suite.ctx, suite.app.BankKeeper, addr1, coins)

	// msg and signatures
	sendAmount := int64(1000000)
	sendCoins := sdk.NewCoins(sdk.NewInt64Coin(assets.MicroSDRDenom, sendAmount))
	msg := &wasmtypes.MsgExecuteContract{
		Sender:   addr1.String(),
		Contract: addr1.String(),
		Msg:      []byte{},
		Funds:    sendCoins,
	}

	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()
	suite.Require().NoError(suite.txBuilder.SetMsgs(msg))
	suite.txBuilder.SetFeeAmount(feeAmount)
	suite.txBuilder.SetGasLimit(gasLimit)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// set zero gas prices
	suite.ctx = suite.ctx.WithMinGasPrices(sdk.NewDecCoins())

	// Set IsCheckTx to true
	suite.ctx = suite.ctx.WithIsCheckTx(true)

	// antehandler errors with insufficient fees due to tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().Error(err, "Decorator should errored on low fee for local gasPrice + tax")

	tk := suite.app.TreasuryKeeper
	expectedTax := tk.GetTaxRate(suite.ctx).MulInt64(sendAmount).TruncateInt()

	// set tax amount
	suite.txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(assets.MicroSDRDenom, expectedTax)))
	tx, err = suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// must pass with tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "Decorator should not have errored on fee higher than local gasPrice")
}

func (suite *AnteTestSuite) TestEnsureMempoolFeesAuthzExec() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	mfd := ante.NewDeductFeeDecorator(suite.app.TxFeesKeeper, suite.app.AccountKeeper, suite.app.BankKeeper, nil, suite.app.TreasuryKeeper, suite.app.OracleKeeper)
	antehandler := sdk.ChainAnteDecorators(mfd)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()
	coins := sdk.NewCoins(sdk.NewCoin(assets.MicroSDRDenom, osmomath.NewInt(1000000)))
	testutil.FundAccount(suite.ctx, suite.app.BankKeeper, addr1, coins)

	// msg and signatures
	sendAmount := int64(1000000)
	sendCoins := sdk.NewCoins(sdk.NewInt64Coin(assets.MicroSDRDenom, sendAmount))
	msg := authz.NewMsgExec(addr1, []sdk.Msg{banktypes.NewMsgSend(addr1, addr1, sendCoins)})

	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()
	suite.Require().NoError(suite.txBuilder.SetMsgs(&msg))
	suite.txBuilder.SetFeeAmount(feeAmount)
	suite.txBuilder.SetGasLimit(gasLimit)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// set zero gas prices
	suite.ctx = suite.ctx.WithMinGasPrices(sdk.NewDecCoins())

	// Set IsCheckTx to true
	suite.ctx = suite.ctx.WithIsCheckTx(true)

	// antehandler errors with insufficient fees due to tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().Error(err, "Decorator should errored on low fee for local gasPrice + tax")

	tk := suite.app.TreasuryKeeper
	expectedTax := tk.GetTaxRate(suite.ctx).MulInt64(sendAmount).TruncateInt()

	// set tax amount
	suite.txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(assets.MicroSDRDenom, expectedTax)))
	tx, err = suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// must pass with tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "could not retrieve exchange rate for usdr: usdr: unknown denom")
}

func (suite *AnteTestSuite) TestTaxExemption() {
	// keys and addresses
	var privs []cryptotypes.PrivKey
	var addrs []sdk.AccAddress

	// 0, 1: exemption
	// 2, 3: normal
	for i := 0; i < 4; i++ {
		priv, _, addr := testdata.KeyTestPubAddr()
		privs = append(privs, priv)
		addrs = append(addrs, addr)
	}

	// set send amount
	sendAmt := int64(1000000)
	sendCoin := sdk.NewInt64Coin(assets.MicroSDRDenom, sendAmt)
	feeAmt := int64(1000)

	cases := []struct {
		name           string
		msgSigner      cryptotypes.PrivKey
		msgCreator     func() []sdk.Msg
		minFeeAmount   int64
		expectProceeds int64
	}{
		{
			name:      "MsgSend(exemption -> exemption)",
			msgSigner: privs[0],
			msgCreator: func() []sdk.Msg {
				var msgs []sdk.Msg

				msg1 := banktypes.NewMsgSend(addrs[0], addrs[1], sdk.NewCoins(sendCoin))
				msgs = append(msgs, msg1)

				return msgs
			},
			minFeeAmount:   0,
			expectProceeds: 0,
		}, {
			name:      "MsgSend(normal -> normal)",
			msgSigner: privs[2],
			msgCreator: func() []sdk.Msg {
				var msgs []sdk.Msg

				msg1 := banktypes.NewMsgSend(addrs[2], addrs[3], sdk.NewCoins(sendCoin))
				msgs = append(msgs, msg1)

				return msgs
			},
			// tax this one hence burn amount is fee amount
			minFeeAmount:   feeAmt,
			expectProceeds: feeAmt,
		}, {
			name:      "MsgExec(MsgSend(normal -> normal))",
			msgSigner: privs[2],
			msgCreator: func() []sdk.Msg {
				var msgs []sdk.Msg

				msg1 := authz.NewMsgExec(addrs[1], []sdk.Msg{banktypes.NewMsgSend(addrs[2], addrs[3], sdk.NewCoins(sendCoin))})
				msgs = append(msgs, &msg1)

				return msgs
			},
			// tax this one hence burn amount is fee amount
			minFeeAmount:   feeAmt,
			expectProceeds: feeAmt,
		}, {
			name:      "MsgSend(exemption -> normal), MsgSend(exemption -> exemption)",
			msgSigner: privs[0],
			msgCreator: func() []sdk.Msg {
				var msgs []sdk.Msg

				msg1 := banktypes.NewMsgSend(addrs[0], addrs[2], sdk.NewCoins(sendCoin))
				msgs = append(msgs, msg1)
				msg2 := banktypes.NewMsgSend(addrs[0], addrs[1], sdk.NewCoins(sendCoin))
				msgs = append(msgs, msg2)

				return msgs
			},
			// tax this one hence burn amount is fee amount
			minFeeAmount:   feeAmt,
			expectProceeds: feeAmt,
		}, {
			name:      "MsgSend(exemption -> exemption), MsgMultiSend(exemption -> normal, exemption -> exemption)",
			msgSigner: privs[0],
			msgCreator: func() []sdk.Msg {
				var msgs []sdk.Msg

				msg1 := banktypes.NewMsgSend(addrs[0], addrs[1], sdk.NewCoins(sendCoin))
				msgs = append(msgs, msg1)
				msg2 := banktypes.NewMsgMultiSend(
					banktypes.NewInput(addrs[0], sdk.NewCoins(sendCoin)),
					[]banktypes.Output{
						{
							Address: addrs[2].String(),
							Coins:   sdk.NewCoins(sendCoin),
						},
						{
							Address: addrs[1].String(),
							Coins:   sdk.NewCoins(sendCoin),
						},
					},
				)
				msgs = append(msgs, msg2)

				return msgs
			},
			minFeeAmount:   feeAmt * 2,
			expectProceeds: feeAmt * 2,
		}, {
			name:      "MsgExecuteContract(exemption), MsgExecuteContract(normal)",
			msgSigner: privs[3],
			msgCreator: func() []sdk.Msg {
				sendAmount := int64(1000000)
				sendCoins := sdk.NewCoins(sdk.NewInt64Coin(assets.MicroSDRDenom, sendAmount))
				// get wasm code for wasm contract create and instantiate
				wasmCode, err := os.ReadFile("./testdata/hackatom.wasm")
				suite.Require().NoError(err)
				per := wasmkeeper.NewDefaultPermissionKeeper(suite.app.WasmKeeper)
				// set wasm default params
				suite.app.WasmKeeper.SetParams(suite.ctx, wasmtypes.DefaultParams())
				// wasm create
				CodeID, _, err := per.Create(suite.ctx, addrs[0], wasmCode, nil)
				suite.Require().NoError(err)
				// params for contract init
				r := wasmkeeper.HackatomExampleInitMsg{Verifier: addrs[0], Beneficiary: addrs[0]}
				bz, err := json.Marshal(r)
				suite.Require().NoError(err)
				// change block time for contract instantiate
				suite.ctx = suite.ctx.WithBlockTime(time.Date(2020, time.April, 22, 12, 0, 0, 0, time.UTC))
				// instantiate contract then set the contract address to tax exemption
				addr, _, err := per.Instantiate(suite.ctx, CodeID, addrs[0], nil, bz, "my label", nil)
				suite.Require().NoError(err)
				// instantiate contract then not set to tax exemption
				addr1, _, err := per.Instantiate(suite.ctx, CodeID, addrs[0], nil, bz, "my label", nil)
				suite.Require().NoError(err)

				var msgs []sdk.Msg
				// msg and signatures
				msg1 := &wasmtypes.MsgExecuteContract{
					Sender:   addrs[0].String(),
					Contract: addr.String(),
					Msg:      []byte{},
					Funds:    sendCoins,
				}
				msgs = append(msgs, msg1)

				msg2 := &wasmtypes.MsgExecuteContract{
					Sender:   addrs[3].String(),
					Contract: addr1.String(),
					Msg:      []byte{},
					Funds:    sendCoins,
				}
				msgs = append(msgs, msg2)
				return msgs
			},
			minFeeAmount:   feeAmt,
			expectProceeds: feeAmt,
		},
	}

	// there should be no coin in burn module
	for _, c := range cases {
		suite.SetupTest(true) // setup
		require := suite.Require()
		ak := suite.app.AccountKeeper
		bk := suite.app.BankKeeper

		fmt.Printf("CASE = %s \n", c.name)
		suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

		mfd := ante.NewDeductFeeDecorator(suite.app.TxFeesKeeper, suite.app.AccountKeeper, suite.app.BankKeeper, nil, suite.app.TreasuryKeeper, suite.app.OracleKeeper)
		antehandler := sdk.ChainAnteDecorators(mfd)

		for i := 0; i < 4; i++ {
			coins := sdk.NewCoins(sdk.NewCoin(assets.MicroSDRDenom, osmomath.NewInt(10000000)))
			testutil.FundAccount(suite.ctx, suite.app.BankKeeper, addrs[i], coins)
		}

		// msg and signatures
		feeAmount := sdk.NewCoins(sdk.NewInt64Coin(assets.MicroSDRDenom, c.minFeeAmount))
		gasLimit := testdata.NewTestGasLimit()
		require.NoError(suite.txBuilder.SetMsgs(c.msgCreator()...))
		suite.txBuilder.SetFeeAmount(feeAmount)
		suite.txBuilder.SetGasLimit(gasLimit)

		privs, accNums, accSeqs := []cryptotypes.PrivKey{c.msgSigner}, []uint64{0}, []uint64{0}
		tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
		require.NoError(err)

		_, err = antehandler(suite.ctx, tx, false)
		require.NoError(err)

		// check fee collector and treasury
		feeCollector := ak.GetModuleAccount(suite.ctx, authtypes.FeeCollectorName)
		//treasuryCollector := ak.GetModuleAccount(s.ctx, authtypes.ModuleName)
		amountFee := bk.GetBalance(suite.ctx, feeCollector.GetAddress(), assets.MicroSDRDenom)
		require.Equal(amountFee, sdk.NewCoin(assets.MicroSDRDenom, osmomath.NewDec(c.minFeeAmount).TruncateInt()))

		// check tax proceeds
		//taxProceeds := s.app.TreasuryKeeper.PeekEpochTaxProceeds(s.ctx)
		//require.Equal(taxProceeds, sdk.NewCoins(sdk.NewCoin(assets.MicroSDRDenom, osmomath.NewInt(c.expectProceeds))))
	}
}

func (suite *AnteTestSuite) TestOracleZeroFee() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	mfd := ante.NewDeductFeeDecorator(suite.app.TxFeesKeeper, suite.app.AccountKeeper, suite.app.BankKeeper, nil, suite.app.TreasuryKeeper, suite.app.OracleKeeper)
	antehandler := sdk.ChainAnteDecorators(mfd)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()
	account := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addr1)
	suite.app.AccountKeeper.SetAccount(suite.ctx, account)
	testutil.FundAccount(suite.ctx, suite.app.BankKeeper, addr1, sdk.NewCoins(sdk.NewInt64Coin(assets.MicroSDRDenom, 1_000_000_000)))

	// new val
	val, err := stakingtypes.NewValidator(sdk.ValAddress(addr1).String(), priv1.PubKey(), stakingtypes.Description{})
	suite.Require().NoError(err)
	suite.app.StakingKeeper.SetValidator(suite.ctx, val)

	// msg and signatures

	// MsgAggregateExchangeRatePrevote
	msg := oracletypes.NewMsgAggregateExchangeRatePrevote(
		oracletypes.GetAggregateVoteHash("salt", "exchange rates",
			sdk.ValAddress(val.GetOperator())), addr1, sdk.ValAddress(val.GetOperator()))
	suite.txBuilder.SetMsgs(msg)
	suite.txBuilder.SetGasLimit(testdata.NewTestGasLimit())
	suite.txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewInt64Coin(assets.MicroSDRDenom, 0)))
	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err)

	// check fee collector empty
	balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, suite.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName))
	suite.Require().Equal(sdk.Coins{}, balances)

	// MsgAggregateExchangeRateVote
	msg1 := oracletypes.NewMsgAggregateExchangeRateVote("salt", "exchange rates", addr1, sdk.ValAddress(val.GetOperator()))
	suite.txBuilder.SetMsgs(msg1)
	tx, err = suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err)

	// check fee collector empty
	balances = suite.app.BankKeeper.GetAllBalances(suite.ctx, suite.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName))
	suite.Require().Equal(sdk.Coins{}, balances)
}
