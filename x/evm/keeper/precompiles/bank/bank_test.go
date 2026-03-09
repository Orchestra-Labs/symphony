package bank_test

import (
	"math/big"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"

	"github.com/osmosis-labs/osmosis/v27/x/evm/keeper/precompiles/bank"
)

type BankPrecompileTestSuite struct {
	suite.Suite

	ctx        sdk.Context
	precompile *bank.BankPrecompile
}

func TestBankPrecompileTestSuite(t *testing.T) {
	suite.Run(t, new(BankPrecompileTestSuite))
}

func (suite *BankPrecompileTestSuite) SetupTest() {
	suite.T().Skip("Setup test app infrastructure needed")
}

func (suite *BankPrecompileTestSuite) TestBalanceOf() {
	// Test address
	testAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")

	// Encode balanceOf(address) call
	// Method ID: 0x70a08231 (first 4 bytes of keccak256("balanceOf(address)"))
	methodID := []byte{0x70, 0xa0, 0x82, 0x31}

	// Pad address to 32 bytes
	paddedAddr := common.LeftPadBytes(testAddr.Bytes(), 32)

	input := append(methodID, paddedAddr...)

	// Call the precompile
	output, gas, err := suite.precompile.Run(suite.ctx, common.Address{}, input, 100000, true)
	suite.Require().NoError(err)
	suite.Require().NotNil(output)
	suite.Require().Greater(gas, uint64(0))

	// Decode output (should be uint256)
	balance := new(big.Int).SetBytes(output)
	suite.Require().NotNil(balance)
}

func (suite *BankPrecompileTestSuite) TestTransfer() {
	// Test addresses
	fromAddr := common.HexToAddress("0x1111111111111111111111111111111111111111")
	toAddr := common.HexToAddress("0x2222222222222222222222222222222222222222")

	// Encode transfer(address,uint256) call
	// Method ID: 0xa9059cbb (first 4 bytes of keccak256("transfer(address,uint256)"))
	methodID := []byte{0xa9, 0x05, 0x9c, 0xbb}

	// Pad recipient address to 32 bytes
	paddedTo := common.LeftPadBytes(toAddr.Bytes(), 32)

	// Amount: 1 NOTE = 1e18 wei
	amount := new(big.Int).Mul(big.NewInt(1), big.NewInt(1e18))
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	input := append(methodID, paddedTo...)
	input = append(input, paddedAmount...)

	// Call the precompile (not read-only since it's a transfer)
	output, gas, err := suite.precompile.Run(suite.ctx, fromAddr, input, 100000, false)
	suite.Require().NoError(err)
	suite.Require().NotNil(output)
	suite.Require().Greater(gas, uint64(0))

	// Output should be true (success) = 1
	success := new(big.Int).SetBytes(output)
	suite.Require().Equal(big.NewInt(1), success)
}

func (suite *BankPrecompileTestSuite) TestRequiredGas() {
	input := []byte{0x70, 0xa0, 0x82, 0x31} // balanceOf method ID
	gas := suite.precompile.RequiredGas(input)

	// Should require some gas
	suite.Require().Greater(gas, uint64(0))
}

func (suite *BankPrecompileTestSuite) TestAddress() {
	addr := suite.precompile.Address()
	expectedAddr := common.HexToAddress("0x0000000000000000000000000000000000000801")
	suite.Require().Equal(expectedAddr, addr)
}
