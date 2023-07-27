package multihop_test

import (
	"github.com/cosmos/cosmos-sdk/codec"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	commitmenttypes "github.com/cosmos/ibc-go/v7/modules/core/23-commitment/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	"github.com/cosmos/ibc-go/v7/modules/core/multihop"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ChanPathTestSuite struct {
	suite.Suite
}

type mockEndpoint struct {
	connectionID string
}

func (m mockEndpoint) ChainID() string {
	//TODO implement me
	panic("implement me")
}

func (m mockEndpoint) Codec() codec.BinaryCodec {
	//TODO implement me
	panic("implement me")
}

func (m mockEndpoint) ClientID() string {
	//TODO implement me
	panic("implement me")
}

func (m mockEndpoint) GetKeyValueProofHeight() exported.Height {
	//TODO implement me
	panic("implement me")
}

func (m mockEndpoint) GetConsensusHeight() exported.Height {
	//TODO implement me
	panic("implement me")
}

func (m mockEndpoint) GetConsensusState(height exported.Height) (exported.ConsensusState, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockEndpoint) ConnectionID() string {
	return m.connectionID
}

func (m mockEndpoint) GetConnection() (*types.ConnectionEnd, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockEndpoint) QueryProofAtHeight(key []byte, height int64) ([]byte, clienttypes.Height, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockEndpoint) GetMerklePath(path string) (commitmenttypes.MerklePath, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockEndpoint) UpdateClient() error {
	//TODO implement me
	panic("implement me")
}

func (m mockEndpoint) Counterparty() multihop.Endpoint {
	//TODO implement me
	panic("implement me")
}

var _ multihop.Endpoint = (*mockEndpoint)(nil)

func (suite *ChanPathTestSuite) TestCounterparty() {
	//       A      <->             B              <->      C
	// connection-0 <-> connection-2, connection-1 <-> connection-3
	epA := mockEndpoint{connectionID: "connection-0"}
	epAB := mockEndpoint{connectionID: "connection-2"}
	epBC := mockEndpoint{connectionID: "connection-1"}
	epC := mockEndpoint{connectionID: "connection-3"}
	pathAB := &multihop.Path{epA, epAB}
	pathBC := &multihop.Path{epBC, epC}
	chanPath := multihop.NewChanPath([]*multihop.Path{pathAB, pathBC})
	hops := chanPath.GetConnectionHops()
	suite.Require().Equal([]string{"connection-0", "connection-1"}, hops)
	suite.Require().Equal([]string{"connection-3", "connection-2"}, chanPath.Counterparty().GetConnectionHops())
}

func TestChanPathTestSuite(t *testing.T) {
	suite.Run(t, new(ChanPathTestSuite))
}
