package keeper_test

import (
	"testing"

	testifysuite "github.com/stretchr/testify/suite"

	ibctesting "github.com/cosmos/ibc-go/v8/testing"
)

// TestMultihopMultihopTestSuite runs all multihop related tests.
func TestMultihopTestSuite(t *testing.T) {
	testifysuite.Run(t, new(MultihopTestSuite))
}

// MultihopTestSuite is a testing suite to test keeper functions.
type MultihopTestSuite struct {
	testifysuite.Suite
	// multihop channel path
	chanPath *ibctesting.PathM
	coord    *ibctesting.CoordinatorM
}

// SetupTest is run before each test method in the suite
// No IBC connections or channels are created.
func (suite *MultihopTestSuite) SetupTest(numChains int) {
	coord, paths := ibctesting.CreateLinkedChains(&suite.Suite, numChains)
	suite.chanPath = paths.ToPathM()
	suite.coord = &ibctesting.CoordinatorM{Coordinator: coord}
}

// SetupConnections creates connections between each pair of chains in the multihop path.
func (suite *MultihopTestSuite) SetupConnections() {
	suite.coord.SetupConnections(suite.chanPath)
}

// SetupConnections creates connections between each pair of chains in the multihop path.
func (suite *MultihopTestSuite) SetupAllButTheSpecifiedConnection(index int) {
	err := suite.coord.SetupAllButTheSpecifiedConnection(suite.chanPath, index)
	suite.Require().NoError(err)
}

// SetupChannels create a multihop channel after creating all its preprequisites in order, ie. clients, connections.
func (suite *MultihopTestSuite) SetupChannels() {
	suite.coord.SetupChannels(suite.chanPath)
}

// A returns the one endpoint of the multihop channel.
func (suite *MultihopTestSuite) A() *ibctesting.EndpointM {
	return suite.chanPath.EndpointA
}

// Z returns the other endpoint of the multihop channel.
func (suite *MultihopTestSuite) Z() *ibctesting.EndpointM {
	return suite.chanPath.EndpointZ
}
