package keeper_test

import (
	"testing"

	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	"github.com/stretchr/testify/suite"
)

// TestMultihopMultihopTestSuite runs all multihop related tests.
func TestMultihopTestSuite(t *testing.T) {
	suite.Run(t, new(MultihopTestSuite))
}

// MultihopTestSuite is a testing suite to test keeper functions.
type MultihopTestSuite struct {
	suite.Suite
	// multihop channel path
	chanPath *ibctesting.PathM
	coord    *ibctesting.CoordinatorM
	// testing chains used for convenience and readability
	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain
}

// SetupTest is run before each test method in the suite
// No IBC connections or channels are created.
func (suite *MultihopTestSuite) SetupTest() {
	coord, paths := ibctesting.CreateLinkedChains(&suite.Suite, 5)
	suite.chanPath = paths.ToPathM()
	suite.coord = &ibctesting.CoordinatorM{Coordinator: coord}
	suite.chainA = suite.coord.GetChain(ibctesting.GetChainID(1))
	suite.chainB = suite.coord.GetChain(ibctesting.GetChainID(2))
	// commit some blocks so that QueryProof returns valid proof (cannot return valid query if height <= 1)
	suite.coord.CommitNBlocks(suite.chainA, 2)
	suite.coord.CommitNBlocks(suite.chainB, 2)
}

// SetupConnections creates connections between each pair of chains in the multihop path.
func (s *MultihopTestSuite) SetupConnections() {
	s.coord.SetupConnections(s.chanPath)
}

// SetupChannels create a multihop channel after creating all its preprequisites in order, ie. clients, connections.
func (s *MultihopTestSuite) SetupChannels() {
	s.coord.SetupChannels(s.chanPath)
}

// A returns the one endpoint of the multihop channel.
func (s *MultihopTestSuite) A() *ibctesting.EndpointM {
	return s.chanPath.EndpointA
}

// Z returns the other endpoint of the multihop channel.
func (s *MultihopTestSuite) Z() *ibctesting.EndpointM {
	return s.chanPath.EndpointZ
}
