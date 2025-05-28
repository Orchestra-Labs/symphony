package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/osmosis-labs/osmosis/osmomath"
	"github.com/osmosis-labs/osmosis/v27/app/apptesting"
	"github.com/osmosis-labs/osmosis/v27/app/apptesting/assets"
)

type UnbondingTestSuite struct {
	apptesting.KeeperTestHelper
}

func TestUnbondingTestSuite(t *testing.T) {
	suite.Run(t, new(UnbondingTestSuite))
}

func (s *UnbondingTestSuite) SetupTest() {
	s.Setup()
	// Set Oracle Price
	sdrPriceInMelody := osmomath.NewDecWithPrec(17, 1)
	s.App.OracleKeeper.SetMelodyExchangeRate(s.Ctx, assets.MicroUSDDenom, sdrPriceInMelody)

	// Mint initial tokens
	totalUsddSupply := sdk.NewCoins(sdk.NewCoin(assets.MicroUSDDenom, InitTokens.MulRaw(int64(len(Addr)*10))))
	err := s.App.BankKeeper.MintCoins(s.Ctx, FaucetAccountName, totalUsddSupply)
	s.Require().NoError(err)

	// Fund test accounts
	staker := sdk.AccAddress("staker1")
	err = s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, FaucetAccountName, staker, InitUSDDCoins)
	s.Require().NoError(err)
}

func (s *UnbondingTestSuite) TestUnbondTokens() {
	staker := sdk.AccAddress("staker1")

	s.Run("fail on non-existent stake", func() {
		_, err := s.App.StableStakingKeeper.UnStakeTokens(s.Ctx, staker, sdk.NewCoin(assets.MicroUSDDenom, math.NewInt(100)))
		require.Error(s.T(), err)
		require.Contains(s.T(), err.Error(), "stake not found")
	})

	s.Run("fail on insufficient stake", func() {
		// First stake some tokens
		token := sdk.NewCoin(assets.MicroUSDDenom, math.NewInt(100))
		_, err := s.App.StableStakingKeeper.StakeTokens(s.Ctx, staker, token)
		require.NoError(s.T(), err)

		// Try to unbond more than staked
		_, err = s.App.StableStakingKeeper.UnStakeTokens(s.Ctx, staker, sdk.NewCoin(assets.MicroUSDDenom, math.NewInt(200)))
		require.Error(s.T(), err)
		require.Contains(s.T(), err.Error(), "insufficient stake")
	})

	s.Run("successful unbonding", func() {
		// First stake some tokens
		token := sdk.NewCoin(assets.MicroUSDDenom, math.NewInt(100))
		_, err := s.App.StableStakingKeeper.StakeTokens(s.Ctx, staker, token)
		require.NoError(s.T(), err)

		// Unbond half of the stake
		unbondAmount := math.NewInt(50)
		unbondingInfo, err := s.App.StableStakingKeeper.UnStakeTokens(s.Ctx, staker, sdk.NewCoin(assets.MicroUSDDenom, unbondAmount))
		require.NoError(s.T(), err)
		require.NotNil(s.T(), unbondingInfo)

		// Verify unbonding info
		require.Equal(s.T(), staker.String(), unbondingInfo.Staker)
		require.Equal(s.T(), assets.MicroUSDDenom, unbondingInfo.Amount.Denom)
		require.Equal(s.T(), unbondAmount.String(), unbondingInfo.Amount.Amount.TruncateInt().String())

		// Verify user stake is reduced
		userStake, found := s.App.StableStakingKeeper.GetUserStake(s.Ctx, staker, assets.MicroUSDDenom)
		require.True(s.T(), found)
		require.Equal(s.T(), math.LegacyNewDecFromInt(math.NewInt(50)), userStake.Shares)

		// Verify pool is updated
		pool, found := s.App.StableStakingKeeper.GetPool(s.Ctx, assets.MicroUSDDenom)
		require.True(s.T(), found)
		require.Equal(s.T(), math.LegacyNewDecFromInt(math.NewInt(50)), pool.TotalShares)
		require.Equal(s.T(), math.LegacyNewDecFromInt(math.NewInt(50)), pool.TotalStaked)
	})
}

func (s *UnbondingTestSuite) TestCompleteUnbonding() {
	staker := sdk.AccAddress("staker1")

	s.Run("complete unbonding after period", func() {
		// First stake some tokens
		token := sdk.NewCoin(assets.MicroUSDDenom, math.NewInt(100))
		_, err := s.App.StableStakingKeeper.StakeTokens(s.Ctx, staker, token)
		require.NoError(s.T(), err)

		// Unbond tokens
		unbondAmount := math.NewInt(50)
		_, err = s.App.StableStakingKeeper.UnStakeTokens(s.Ctx, staker, sdk.NewCoin(assets.MicroUSDDenom, unbondAmount))
		require.NoError(s.T(), err)

		// Move time forward past an unbonding period
		unbondingDuration := s.App.StableStakingKeeper.GetParams(s.Ctx).UnbondingDuration
		s.Ctx = s.Ctx.WithBlockTime(s.Ctx.BlockTime().Add(unbondingDuration + time.Hour))

		// Verify unbonding info is removed
		unbondInfo, found := s.App.StableStakingKeeper.GetUnbondingInfo(s.Ctx, staker, assets.MicroUSDDenom)
		require.False(s.T(), found)

		// Verify tokens are returned to user
		require.Equal(s.T(), unbondAmount.String(), unbondInfo.Amount.String())
	})

	s.Run("fail to complete unbonding before period", func() {
		// First stake some tokens
		token := sdk.NewCoin(assets.MicroUSDDenom, math.NewInt(100))
		_, err := s.App.StableStakingKeeper.StakeTokens(s.Ctx, staker, token)
		require.NoError(s.T(), err)

		// Unbond tokens
		unbondAmount := math.NewInt(50)
		_, err = s.App.StableStakingKeeper.UnStakeTokens(s.Ctx, staker, sdk.NewCoin(assets.MicroUSDDenom, unbondAmount))
		require.NoError(s.T(), err)
	})
}
