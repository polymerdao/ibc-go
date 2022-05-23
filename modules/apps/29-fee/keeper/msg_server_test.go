package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/ibc-go/v3/modules/apps/29-fee/types"
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v3/testing"
)

func (suite *KeeperTestSuite) TestRegisterCounterpartyAddress() {
	var (
		sender       string
		counterparty string
		channelID    string
		ctx          sdk.Context
	)

	testCases := []struct {
		name     string
		expPass  bool
		malleate func()
	}{
		{
			"success",
			true,
			func() {},
		},
		{
			"counterparty is an arbitrary string",
			true,
			func() { counterparty = "arbitrary-string" },
		},
		{
			"channel does not exist",
			false,
			func() { channelID = "channel-22" },
		},
		{
			"channel is not fee enabled",
			false,
			func() {
				suite.chainA.GetSimApp().IBCFeeKeeper.DeleteFeeEnabled(ctx, suite.path.EndpointA.ChannelConfig.PortID, channelID)
			},
		},
	}

	for _, tc := range testCases {
		suite.SetupTest()
		ctx = suite.chainA.GetContext()
		suite.coordinator.Setup(suite.path) // setup channel

		sender = suite.chainA.SenderAccount.GetAddress().String()
		counterparty = suite.chainB.SenderAccount.GetAddress().String()
		channelID = suite.path.EndpointA.ChannelID

		tc.malleate()

		msg := types.NewMsgRegisterCounterpartyAddress(sender, counterparty, suite.path.EndpointA.ChannelConfig.PortID, channelID)
		_, err := suite.chainA.GetSimApp().IBCFeeKeeper.RegisterCounterpartyAddress(sdk.WrapSDKContext(ctx), msg)

		if tc.expPass {
			suite.Require().NoError(err) // message committed

			counterpartyAddress, _ := suite.chainA.GetSimApp().IBCFeeKeeper.GetCounterpartyAddress(ctx, suite.chainA.SenderAccount.GetAddress().String(), ibctesting.FirstChannelID)
			suite.Require().Equal(counterparty, counterpartyAddress)
		} else {
			suite.Require().Error(err)
		}
	}
}

func (suite *KeeperTestSuite) TestPayPacketFee() {
	var (
		expEscrowBalance sdk.Coins
		expFeesInEscrow  []types.PacketFee
		msg              *types.MsgPayPacketFee
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"success",
			func() {},
			true,
		},
		{
			"success with existing packet fees in escrow",
			func() {
				fee := types.NewFee(defaultRecvFee, defaultAckFee, defaultTimeoutFee)

				packetID := channeltypes.NewPacketId(suite.path.EndpointA.ChannelConfig.PortID, suite.path.EndpointA.ChannelID, 1)
				packetFee := types.NewPacketFee(fee, suite.chainA.SenderAccount.GetAddress().String(), nil)
				feesInEscrow := types.NewPacketFees([]types.PacketFee{packetFee})

				suite.chainA.GetSimApp().IBCFeeKeeper.SetFeesInEscrow(suite.chainA.GetContext(), packetID, feesInEscrow)
				err := suite.chainA.GetSimApp().BankKeeper.SendCoinsFromAccountToModule(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), types.ModuleName, fee.Total())
				suite.Require().NoError(err)

				expEscrowBalance = expEscrowBalance.Add(fee.Total()...)
				expFeesInEscrow = append(expFeesInEscrow, packetFee)
			},
			true,
		},
		{
			"fee module is locked",
			func() {
				lockFeeModule(suite.chainA)
			},
			false,
		},
		{
			"fee module disabled on channel",
			func() {
				msg.SourcePortId = "invalid-port"
				msg.SourceChannelId = "invalid-channel"
			},
			false,
		},
		{
			"invalid refund address",
			func() {
				msg.Signer = "invalid-address"
			},
			false,
		},
		{
			"refund account does not exist",
			func() {
				msg.Signer = suite.chainB.SenderAccount.GetAddress().String()
			},
			false,
		},
		{
			"acknowledgement fee balance not found",
			func() {
				msg.Fee.AckFee = invalidCoins
			},
			false,
		},
		{
			"receive fee balance not found",
			func() {
				msg.Fee.RecvFee = invalidCoins
			},
			false,
		},
		{
			"timeout fee balance not found",
			func() {
				msg.Fee.TimeoutFee = invalidCoins
			},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest()
			suite.coordinator.Setup(suite.path) // setup channel

			fee := types.NewFee(defaultRecvFee, defaultAckFee, defaultTimeoutFee)
			msg = types.NewMsgPayPacketFee(
				fee,
				suite.path.EndpointA.ChannelConfig.PortID,
				suite.path.EndpointA.ChannelID,
				suite.chainA.SenderAccount.GetAddress().String(),
				nil,
			)

			expEscrowBalance = fee.Total()
			expPacketFee := types.NewPacketFee(fee, suite.chainA.SenderAccount.GetAddress().String(), nil)
			expFeesInEscrow = []types.PacketFee{expPacketFee}

			tc.malleate()

			_, err := suite.chainA.GetSimApp().IBCFeeKeeper.PayPacketFee(sdk.WrapSDKContext(suite.chainA.GetContext()), msg)

			if tc.expPass {
				suite.Require().NoError(err) // message committed

				packetID := channeltypes.NewPacketId(suite.path.EndpointA.ChannelConfig.PortID, suite.path.EndpointA.ChannelID, 1)
				feesInEscrow, found := suite.chainA.GetSimApp().IBCFeeKeeper.GetFeesInEscrow(suite.chainA.GetContext(), packetID)
				suite.Require().True(found)
				suite.Require().Equal(expFeesInEscrow, feesInEscrow.PacketFees)

				escrowBalance := suite.chainA.GetSimApp().BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.GetSimApp().IBCFeeKeeper.GetFeeModuleAddress(), sdk.DefaultBondDenom)
				suite.Require().Equal(expEscrowBalance.AmountOf(sdk.DefaultBondDenom), escrowBalance.Amount)
			} else {
				suite.Require().Error(err)

				escrowBalance := suite.chainA.GetSimApp().BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.GetSimApp().IBCFeeKeeper.GetFeeModuleAddress(), sdk.DefaultBondDenom)
				suite.Require().Equal(sdk.NewInt(0), escrowBalance.Amount)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestPayPacketFeeAsync() {
	var (
		packet           channeltypes.Packet
		expEscrowBalance sdk.Coins
		expFeesInEscrow  []types.PacketFee
		msg              *types.MsgPayPacketFeeAsync
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"success",
			func() {},
			true,
		},
		{
			"success with existing packet fees in escrow",
			func() {
				fee := types.NewFee(defaultRecvFee, defaultAckFee, defaultTimeoutFee)

				packetID := channeltypes.NewPacketId(suite.path.EndpointA.ChannelConfig.PortID, suite.path.EndpointA.ChannelID, 1)
				packetFee := types.NewPacketFee(fee, suite.chainA.SenderAccount.GetAddress().String(), nil)
				feesInEscrow := types.NewPacketFees([]types.PacketFee{packetFee})

				suite.chainA.GetSimApp().IBCFeeKeeper.SetFeesInEscrow(suite.chainA.GetContext(), packetID, feesInEscrow)
				err := suite.chainA.GetSimApp().BankKeeper.SendCoinsFromAccountToModule(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), types.ModuleName, fee.Total())
				suite.Require().NoError(err)

				expEscrowBalance = expEscrowBalance.Add(fee.Total()...)
				expFeesInEscrow = append(expFeesInEscrow, packetFee)
			},
			true,
		},
		{
			"fee module is locked",
			func() {
				lockFeeModule(suite.chainA)
			},
			false,
		},
		{
			"fee module disabled on channel",
			func() {
				msg.PacketId.PortId = "invalid-port"
				msg.PacketId.ChannelId = "invalid-channel"
			},
			false,
		},
		{
			"channel does not exist",
			func() {
				msg.PacketId.ChannelId = "channel-100"

				// to test this functionality, we must set the fee to enabled for this non existent channel
				// NOTE: the channel doesn't exist in 04-channel keeper, but we will add a mapping within ics29 anyways
				suite.chainA.GetSimApp().IBCFeeKeeper.SetFeeEnabled(suite.chainA.GetContext(), msg.PacketId.PortId, msg.PacketId.ChannelId)
			},
			false,
		},
		{
			"packet not sent",
			func() {
				msg.PacketId.Sequence = msg.PacketId.Sequence + 1
			},
			false,
		},
		{
			"packet already acknowledged",
			func() {
				err := suite.path.RelayPacket(packet)
				suite.Require().NoError(err)
			},
			false,
		},
		{
			"packet already timed out",
			func() {
				// try to incentivze a packet which is timed out
				packetID := channeltypes.NewPacketId(suite.path.EndpointA.ChannelConfig.PortID, suite.path.EndpointA.ChannelID, msg.PacketId.Sequence+1)
				packet = channeltypes.NewPacket(ibctesting.MockPacketData, packetID.Sequence, packetID.PortId, packetID.ChannelId, suite.path.EndpointB.ChannelConfig.PortID, suite.path.EndpointB.ChannelID, clienttypes.GetSelfHeight(suite.chainB.GetContext()), 0)

				err := suite.path.EndpointA.SendPacket(packet)
				suite.Require().NoError(err)

				// need to update chainA's client representing chainB to prove missing ack
				err = suite.path.EndpointA.UpdateClient()
				suite.Require().NoError(err)

				err = suite.path.EndpointA.TimeoutPacket(packet)
				suite.Require().NoError(err)

				msg.PacketId = packetID
			},
			false,
		},
		{
			"invalid refund address",
			func() {
				msg.PacketFee.RefundAddress = "invalid-address"
			},
			false,
		},
		{
			"refund account does not exist",
			func() {
				msg.PacketFee.RefundAddress = suite.chainB.SenderAccount.GetAddress().String()
			},
			false,
		},
		{
			"acknowledgement fee balance not found",
			func() {
				msg.PacketFee.Fee.AckFee = invalidCoins
			},
			false,
		},
		{
			"receive fee balance not found",
			func() {
				msg.PacketFee.Fee.RecvFee = invalidCoins
			},
			false,
		},
		{
			"timeout fee balance not found",
			func() {
				msg.PacketFee.Fee.TimeoutFee = invalidCoins
			},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest()
			suite.coordinator.Setup(suite.path) // setup channel

			// send a packet to incentivize
			packetID := channeltypes.NewPacketId(suite.path.EndpointA.ChannelConfig.PortID, suite.path.EndpointA.ChannelID, 1)
			packet = channeltypes.NewPacket(ibctesting.MockPacketData, packetID.Sequence, packetID.PortId, packetID.ChannelId, suite.path.EndpointB.ChannelConfig.PortID, suite.path.EndpointB.ChannelID, clienttypes.NewHeight(clienttypes.ParseChainID(suite.chainB.ChainID), 100), 0)
			err := suite.path.EndpointA.SendPacket(packet)
			suite.Require().NoError(err)

			fee := types.NewFee(defaultRecvFee, defaultAckFee, defaultTimeoutFee)
			packetFee := types.NewPacketFee(fee, suite.chainA.SenderAccount.GetAddress().String(), nil)

			expEscrowBalance = fee.Total()
			expFeesInEscrow = []types.PacketFee{packetFee}
			msg = types.NewMsgPayPacketFeeAsync(packetID, packetFee)

			tc.malleate()

			_, err = suite.chainA.GetSimApp().IBCFeeKeeper.PayPacketFeeAsync(sdk.WrapSDKContext(suite.chainA.GetContext()), msg)

			if tc.expPass {
				suite.Require().NoError(err) // message committed

				feesInEscrow, found := suite.chainA.GetSimApp().IBCFeeKeeper.GetFeesInEscrow(suite.chainA.GetContext(), packetID)
				suite.Require().True(found)
				suite.Require().Equal(expFeesInEscrow, feesInEscrow.PacketFees)

				escrowBalance := suite.chainA.GetSimApp().BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.GetSimApp().IBCFeeKeeper.GetFeeModuleAddress(), sdk.DefaultBondDenom)
				suite.Require().Equal(expEscrowBalance.AmountOf(sdk.DefaultBondDenom), escrowBalance.Amount)
			} else {
				suite.Require().Error(err)

				escrowBalance := suite.chainA.GetSimApp().BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.GetSimApp().IBCFeeKeeper.GetFeeModuleAddress(), sdk.DefaultBondDenom)
				suite.Require().Equal(sdk.NewInt(0), escrowBalance.Amount)
			}
		})
	}
}
