package types_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	clienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
	ibctesting "github.com/cosmos/ibc-go/v5/testing"
)

var (
	clientHeight = clienttypes.NewHeight(0, 4)
)

type LocalhostTestSuite struct {
	suite.Suite

	coordinator ibctesting.Coordinator
	chain       *ibctesting.TestChain
}

func (suite *LocalhostTestSuite) SetupTest() {
	suite.coordinator = *ibctesting.NewCoordinator(suite.T(), 1)
	suite.chain = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	// commit some blocks so that QueryProof returns valid proof (cannot return valid query if height <= 1)
	suite.coordinator.CommitNBlocks(suite.chain, 2)
}

func TestLocalhostTestSuite(t *testing.T) {
	suite.Run(t, new(LocalhostTestSuite))
}
