package keeper_test

import (
	"github.com/osmosis-labs/osmosis/v27/x/autopilot/types"
)

func (s *KeeperTestSuite) TestGenesis() {
	expectedGenesisState := types.GenesisState{
		Params: types.Params{
			StakeibcActive: true,
			ClaimActive:    true,
		},
	}

	s.App.AutopilotKeeper.InitGenesis(s.Ctx, expectedGenesisState)

	actualGenesisState := s.App.AutopilotKeeper.ExportGenesis(s.Ctx)
	s.Require().NotNil(actualGenesisState)
	s.Require().Equal(expectedGenesisState.Params, actualGenesisState.Params)
}
