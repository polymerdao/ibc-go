package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/codec"
	clienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v5/modules/core/03-connection/types"

	//channeltypes "github.com/cosmos/ibc-go/v5/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/v5/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/v5/modules/core/24-host"
	"github.com/cosmos/ibc-go/v5/modules/core/exported"
	ibctmtypes "github.com/cosmos/ibc-go/v5/modules/light-clients/07-tendermint/types"
	"github.com/cosmos/ibc-go/v5/modules/light-clients/09-localhost/types"
	ibctesting "github.com/cosmos/ibc-go/v5/testing"
)

const (
	testConnectionID = "connectionid"
	testPortID       = "testportid"
	testChannelID    = "testchannelid"
	testSequence     = 1
)

func (suite *LocalhostTestSuite) TestStatus() {
	ctx := suite.chain.GetContext()
	clientState := types.NewClientState("chainID", clienttypes.NewHeight(3, 10))

	// localhost should always return active
	status := clientState.Status(ctx, nil, nil)
	suite.Require().Equal(exported.Active, status)
}

func (suite *LocalhostTestSuite) TestValidate() {
	testCases := []struct {
		name        string
		clientState *types.ClientState
		expPass     bool
	}{
		{
			name:        "valid client",
			clientState: types.NewClientState("chainID", clienttypes.NewHeight(3, 10)),
			expPass:     true,
		},
		{
			name:        "invalid chain id",
			clientState: types.NewClientState(" ", clienttypes.NewHeight(3, 10)),
			expPass:     false,
		},
		{
			name:        "invalid height",
			clientState: types.NewClientState("chainID", clienttypes.ZeroHeight()),
			expPass:     false,
		},
	}

	for _, tc := range testCases {
		err := tc.clientState.Validate()
		if tc.expPass {
			suite.Require().NoError(err, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
	}
}

func (suite *LocalhostTestSuite) TestInitialize() {
	testCases := []struct {
		name      string
		consState exported.ConsensusState
		expPass   bool
	}{
		{
			"valid initialization",
			nil,
			true,
		},
		{
			"invalid consenus state",
			&ibctmtypes.ConsensusState{},
			false,
		},
	}

	clientState := types.NewClientState("chainID", clienttypes.NewHeight(3, 10))

	for _, tc := range testCases {
		err := clientState.Initialize(suite.chain.GetContext(), suite.chain.Codec, nil, tc.consState)

		if tc.expPass {
			suite.Require().NoError(err, "valid testcase: %s failed", tc.name)
		} else {
			suite.Require().Error(err, "invalid testcase: %s passed", tc.name)
		}
	}
}

func (suite *LocalhostTestSuite) TestVerifyClientState() {
	clientState := types.NewClientState("chainID", clientHeight)
	invalidClient := types.NewClientState("chainID", clienttypes.NewHeight(0, 12))
	testCases := []struct {
		name         string
		clientState  *types.ClientState
		malleate     func(codec.BinaryCodec, sdk.KVStore)
		counterparty *types.ClientState
		expPass      bool
	}{
		{
			name:        "proof verification success",
			clientState: clientState,
			malleate: func(cdc codec.BinaryCodec, store sdk.KVStore) {
				bz := clienttypes.MustMarshalClientState(cdc, clientState)
				store.Set(host.ClientStateKey(), bz)
			},
			counterparty: clientState,
			expPass:      true,
		},
		{
			name:        "proof verification failed: invalid client",
			clientState: clientState,
			malleate: func(cdc codec.BinaryCodec, store sdk.KVStore) {
				bz := clienttypes.MustMarshalClientState(cdc, clientState)
				store.Set(host.ClientStateKey(), bz)
			},
			counterparty: invalidClient,
			expPass:      false,
		},
		{
			name:         "proof verification failed: client not stored",
			clientState:  clientState,
			malleate:     func(cdc codec.BinaryCodec, store sdk.KVStore) {},
			counterparty: clientState,
			expPass:      false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest()
			cdc := suite.chain.Codec
			store := suite.chain.GetContext().KVStore(suite.chain.App.GetKey(host.StoreKey))
			tc.malleate(cdc, store)

			err := tc.clientState.VerifyClientState(
				store, cdc, clienttypes.NewHeight(0, 10), nil, "", []byte{}, tc.counterparty,
			)

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *LocalhostTestSuite) TestVerifyClientConsensusState() {
	clientState := types.NewClientState("chainID", clientHeight)
	err := clientState.VerifyClientConsensusState(
		nil, nil, nil, "", nil, nil, nil, nil,
	)
	suite.Require().NoError(err)
}

func (suite *LocalhostTestSuite) TestCheckHeaderAndUpdateState() {
	ctx := suite.chain.GetContext()
	clientState := types.NewClientState("chainID", clientHeight)
	cs, _, err := clientState.CheckHeaderAndUpdateState(ctx, nil, nil, nil)
	suite.Require().NoError(err)
	suite.Require().Equal(uint64(0), cs.GetLatestHeight().GetRevisionNumber())
	suite.Require().Equal(ctx.BlockHeight(), int64(cs.GetLatestHeight().GetRevisionHeight()))
	suite.Require().Equal(ctx.BlockHeader().ChainID, clientState.ChainId)
}

func (suite *LocalhostTestSuite) TestMisbehaviourAndUpdateState() {
	ctx := suite.chain.GetContext()
	clientState := types.NewClientState("chainID", clientHeight)
	cs, err := clientState.CheckMisbehaviourAndUpdateState(ctx, nil, nil, nil)
	suite.Require().Error(err)
	suite.Require().Nil(cs)
}

func (suite *LocalhostTestSuite) TestProposedHeaderAndUpdateState() {
	ctx := suite.chain.GetContext()
	clientState := types.NewClientState("chainID", clientHeight)
	cs, err := clientState.CheckSubstituteAndUpdateState(ctx, nil, nil, nil, nil)
	suite.Require().Error(err)
	suite.Require().Nil(cs)
}

func (suite *LocalhostTestSuite) TestVerifyConnectionState() {

	testCases := []struct {
		name                   string
		clientState            *types.ClientState
		getTestConnectionAndID func(*ibctesting.Path) (string, connectiontypes.ConnectionEnd)
		expPass                bool
	}{
		{
			name:        "proof verification success",
			clientState: types.NewClientState("chainID", clientHeight),
			getTestConnectionAndID: func(path *ibctesting.Path) (string, connectiontypes.ConnectionEnd) {
				conn := path.EndpointB.GetConnection()
				return path.EndpointB.ConnectionID, conn
			},
			expPass: true,
		},
		{
			name:        "proof verification failed: connection not stored",
			clientState: types.NewClientState("chainID", clientHeight),
			getTestConnectionAndID: func(path *ibctesting.Path) (string, connectiontypes.ConnectionEnd) {
				counterparty := connectiontypes.NewCounterparty("clientB", testConnectionID, commitmenttypes.NewMerklePrefix([]byte("ibc")))
				conn := connectiontypes.NewConnectionEnd(connectiontypes.OPEN, "clientA", counterparty, []*connectiontypes.Version{connectiontypes.NewVersion("2", nil)}, 0)
				return path.EndpointB.ConnectionID, conn
			},
			expPass: false,
		},
		{
			name:        "proof verification failed: different connection stored",
			clientState: types.NewClientState("chainID", clientHeight),
			getTestConnectionAndID: func(path *ibctesting.Path) (string, connectiontypes.ConnectionEnd) {
				counterparty := connectiontypes.NewCounterparty(path.EndpointB.ClientID, path.EndpointB.ConnectionID, commitmenttypes.NewMerklePrefix([]byte("ibc")))
				conn := connectiontypes.NewConnectionEnd(connectiontypes.OPEN, path.EndpointA.ClientID, counterparty, []*connectiontypes.Version{connectiontypes.NewVersion("2", nil)}, 0)
				return conn.Counterparty.ConnectionId, conn
			},
			expPass: false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest()
			path := ibctesting.NewLocalPath(suite.chain)
			suite.coordinator.Setup(path)

			connID, conn := tc.getTestConnectionAndID(path)
			store := suite.chain.GetContext().KVStore(suite.chain.App.GetKey(host.StoreKey))
			err := tc.clientState.VerifyConnectionState(
				store, suite.chain.Codec, clientHeight, nil, []byte{}, connID, conn,
			)

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

//func (suite *LocalhostTestSuite) TestVerifyChannelState() {
//	counterparty := channeltypes.NewCounterparty(testPortID, testChannelID)
//	ch1 := channeltypes.NewChannel(channeltypes.OPEN, channeltypes.ORDERED, counterparty, []string{testConnectionID}, "1.0.0")
//	ch2 := channeltypes.NewChannel(channeltypes.OPEN, channeltypes.ORDERED, counterparty, []string{testConnectionID}, "2.0.0")
//
//	testCases := []struct {
//		name        string
//		clientState *types.ClientState
//		malleate    func()
//		channel     channeltypes.Channel
//		expPass     bool
//	}{
//		{
//			name:        "proof verification success",
//			clientState: types.NewClientState("chainID", clientHeight),
//			malleate: func() {
//				bz, err := suite.cdc.Marshal(&ch1)
//				suite.Require().NoError(err)
//				suite.store.Set(host.ChannelKey(testPortID, testChannelID), bz)
//			},
//			channel: ch1,
//			expPass: true,
//		},
//		{
//			name:        "proof verification failed: channel not stored",
//			clientState: types.NewClientState("chainID", clientHeight),
//			malleate:    func() {},
//			channel:     ch1,
//			expPass:     false,
//		},
//		{
//			name:        "proof verification failed: unmarshal failed",
//			clientState: types.NewClientState("chainID", clientHeight),
//			malleate: func() {
//				suite.store.Set(host.ChannelKey(testPortID, testChannelID), []byte("channel"))
//			},
//			channel: ch1,
//			expPass: false,
//		},
//		{
//			name:        "proof verification failed: different channel stored",
//			clientState: types.NewClientState("chainID", clientHeight),
//			malleate: func() {
//				bz, err := suite.cdc.Marshal(&ch2)
//				suite.Require().NoError(err)
//				suite.store.Set(host.ChannelKey(testPortID, testChannelID), bz)
//			},
//			channel: ch1,
//			expPass: false,
//		},
//	}
//
//	for _, tc := range testCases {
//		tc := tc
//
//		suite.Run(tc.name, func() {
//			suite.SetupTest()
//			tc.malleate()
//
//			err := tc.clientState.VerifyChannelState(
//				suite.store, suite.cdc, clientHeight, nil, []byte{}, testPortID, testChannelID, &tc.channel,
//			)
//
//			if tc.expPass {
//				suite.Require().NoError(err)
//			} else {
//				suite.Require().Error(err)
//			}
//		})
//	}
//}
//
//func (suite *LocalhostTestSuite) TestVerifyPacketCommitment() {
//	testCases := []struct {
//		name        string
//		clientState *types.ClientState
//		malleate    func()
//		commitment  []byte
//		expPass     bool
//	}{
//		{
//			name:        "proof verification success",
//			clientState: types.NewClientState("chainID", clientHeight),
//			malleate: func() {
//				suite.store.Set(
//					host.PacketCommitmentKey(testPortID, testChannelID, testSequence), []byte("commitment"),
//				)
//			},
//			commitment: []byte("commitment"),
//			expPass:    true,
//		},
//		{
//			name:        "proof verification failed: different commitment stored",
//			clientState: types.NewClientState("chainID", clientHeight),
//			malleate: func() {
//				suite.store.Set(
//					host.PacketCommitmentKey(testPortID, testChannelID, testSequence), []byte("different"),
//				)
//			},
//			commitment: []byte("commitment"),
//			expPass:    false,
//		},
//		{
//			name:        "proof verification failed: no commitment stored",
//			clientState: types.NewClientState("chainID", clientHeight),
//			malleate:    func() {},
//			commitment:  []byte{},
//			expPass:     false,
//		},
//	}
//
//	for _, tc := range testCases {
//		tc := tc
//
//		suite.Run(tc.name, func() {
//			suite.SetupTest()
//			tc.malleate()
//
//			err := tc.clientState.VerifyPacketCommitment(
//				suite.ctx, suite.store, suite.cdc, clientHeight, 0, 0, nil, []byte{}, testPortID, testChannelID, testSequence, tc.commitment,
//			)
//
//			if tc.expPass {
//				suite.Require().NoError(err)
//			} else {
//				suite.Require().Error(err)
//			}
//		})
//	}
//}
//
//func (suite *LocalhostTestSuite) TestVerifyPacketAcknowledgement() {
//	testCases := []struct {
//		name        string
//		clientState *types.ClientState
//		malleate    func()
//		ack         []byte
//		expPass     bool
//	}{
//		{
//			name:        "proof verification success",
//			clientState: types.NewClientState("chainID", clientHeight),
//			malleate: func() {
//				suite.store.Set(
//					host.PacketAcknowledgementKey(testPortID, testChannelID, testSequence), []byte("acknowledgement"),
//				)
//			},
//			ack:     []byte("acknowledgement"),
//			expPass: true,
//		},
//		{
//			name:        "proof verification failed: different ack stored",
//			clientState: types.NewClientState("chainID", clientHeight),
//			malleate: func() {
//				suite.store.Set(
//					host.PacketAcknowledgementKey(testPortID, testChannelID, testSequence), []byte("different"),
//				)
//			},
//			ack:     []byte("acknowledgement"),
//			expPass: false,
//		},
//		{
//			name:        "proof verification failed: no commitment stored",
//			clientState: types.NewClientState("chainID", clientHeight),
//			malleate:    func() {},
//			ack:         []byte{},
//			expPass:     false,
//		},
//	}
//
//	for _, tc := range testCases {
//		tc := tc
//
//		suite.Run(tc.name, func() {
//			suite.SetupTest()
//			tc.malleate()
//
//			err := tc.clientState.VerifyPacketAcknowledgement(
//				suite.ctx, suite.store, suite.cdc, clientHeight, 0, 0, nil, []byte{}, testPortID, testChannelID, testSequence, tc.ack,
//			)
//
//			if tc.expPass {
//				suite.Require().NoError(err)
//			} else {
//				suite.Require().Error(err)
//			}
//		})
//	}
//}
//
//func (suite *LocalhostTestSuite) TestVerifyPacketReceiptAbsence() {
//	clientState := types.NewClientState("chainID", clientHeight)
//
//	err := clientState.VerifyPacketReceiptAbsence(
//		suite.ctx, suite.store, suite.cdc, clientHeight, 0, 0, nil, nil, testPortID, testChannelID, testSequence,
//	)
//
//	suite.Require().NoError(err, "receipt absence failed")
//
//	suite.store.Set(host.PacketReceiptKey(testPortID, testChannelID, testSequence), []byte("receipt"))
//
//	err = clientState.VerifyPacketReceiptAbsence(
//		suite.ctx, suite.store, suite.cdc, clientHeight, 0, 0, nil, nil, testPortID, testChannelID, testSequence,
//	)
//	suite.Require().Error(err, "receipt exists in store")
//}
//
//func (suite *LocalhostTestSuite) TestVerifyNextSeqRecv() {
//	nextSeqRecv := uint64(5)
//
//	testCases := []struct {
//		name        string
//		clientState *types.ClientState
//		malleate    func()
//		nextSeqRecv uint64
//		expPass     bool
//	}{
//		{
//			name:        "proof verification success",
//			clientState: types.NewClientState("chainID", clientHeight),
//			malleate: func() {
//				suite.store.Set(
//					host.NextSequenceRecvKey(testPortID, testChannelID),
//					sdk.Uint64ToBigEndian(nextSeqRecv),
//				)
//			},
//			nextSeqRecv: nextSeqRecv,
//			expPass:     true,
//		},
//		{
//			name:        "proof verification failed: different nextSeqRecv stored",
//			clientState: types.NewClientState("chainID", clientHeight),
//			malleate: func() {
//				suite.store.Set(
//					host.NextSequenceRecvKey(testPortID, testChannelID),
//					sdk.Uint64ToBigEndian(3),
//				)
//			},
//			nextSeqRecv: nextSeqRecv,
//			expPass:     false,
//		},
//		{
//			name:        "proof verification failed: no nextSeqRecv stored",
//			clientState: types.NewClientState("chainID", clientHeight),
//			malleate:    func() {},
//			nextSeqRecv: nextSeqRecv,
//			expPass:     false,
//		},
//	}
//
//	for _, tc := range testCases {
//		tc := tc
//
//		suite.Run(tc.name, func() {
//			suite.SetupTest()
//			tc.malleate()
//
//			err := tc.clientState.VerifyNextSequenceRecv(
//				suite.ctx, suite.store, suite.cdc, clientHeight, 0, 0, nil, []byte{}, testPortID, testChannelID, nextSeqRecv,
//			)
//
//			if tc.expPass {
//				suite.Require().NoError(err)
//			} else {
//				suite.Require().Error(err)
//			}
//		})
//	}
//}
