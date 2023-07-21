//go:build !noproof

package keeper

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	"github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
)

// TimeoutPacket is called by a module which originally attempted to send a
// packet to a counterparty module, where the timeout height has passed on the
// counterparty chain without the packet being committed, to prove that the
// packet can no longer be executed and to allow the calling module to safely
// perform appropriate state transitions. Its intended usage is within the
// ante handler.
func (k Keeper) TimeoutPacket(
	ctx sdk.Context,
	packet exported.PacketI,
	proof []byte,
	proofHeight exported.Height,
	nextSequenceRecv uint64,
) error {
	channel, found := k.GetChannel(ctx, packet.GetSourcePort(), packet.GetSourceChannel())
	if !found {
		return sdkerrors.Wrapf(
			types.ErrChannelNotFound,
			"port ID (%s) channel ID (%s)", packet.GetSourcePort(), packet.GetSourceChannel(),
		)
	}

	// NOTE: TimeoutPacket is called by the AnteHandler which acts upon the packet.Route(),
	// so the capability authentication can be omitted here

	if packet.GetDestPort() != channel.Counterparty.PortId {
		return sdkerrors.Wrapf(
			types.ErrInvalidPacket,
			"packet destination port doesn't match the counterparty's port (%s ≠ %s)", packet.GetDestPort(), channel.Counterparty.PortId,
		)
	}

	if packet.GetDestChannel() != channel.Counterparty.ChannelId {
		return sdkerrors.Wrapf(
			types.ErrInvalidPacket,
			"packet destination channel doesn't match the counterparty's channel (%s ≠ %s)", packet.GetDestChannel(), channel.Counterparty.ChannelId,
		)
	}

	connectionEnd, found := k.connectionKeeper.GetConnection(ctx, channel.ConnectionHops[0])
	if !found {
		return sdkerrors.Wrap(
			connectiontypes.ErrConnectionNotFound,
			channel.ConnectionHops[0],
		)
	}

	var mProof types.MsgMultihopProofs
	var proofTimestamp uint64
	var err error
	if len(channel.ConnectionHops) > 1 {
		err := k.cdc.Unmarshal(proof, &mProof)
		if err != nil {
			return err
		}

		consensusState, err := mProof.GetMultihopCounterpartyConsensus(k.cdc)
		if err != nil {
			return err
		}
		proofTimestamp = consensusState.GetTimestamp()
	} else {
		// check that timeout height or timeout timestamp has passed on the other end
		var err error
		proofTimestamp, err = k.connectionKeeper.GetTimestampAtHeight(ctx, connectionEnd, proofHeight)
		if err != nil {
			return err
		}
	}

	timeoutHeight := packet.GetTimeoutHeight()
	if (timeoutHeight.IsZero() || proofHeight.LT(timeoutHeight)) &&
		(packet.GetTimeoutTimestamp() == 0 || proofTimestamp < packet.GetTimeoutTimestamp()) {
		return sdkerrors.Wrap(types.ErrPacketTimeout, "packet timeout has not been reached for height or timestamp")
	}

	commitment := k.GetPacketCommitment(ctx, packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetSequence())

	if len(commitment) == 0 {
		EmitTimeoutPacketEvent(ctx, packet, channel)
		// This error indicates that the timeout has already been relayed
		// or there is a misconfigured relayer attempting to prove a timeout
		// for a packet never sent. Core IBC will treat this error as a no-op in order to
		// prevent an entire relay transaction from failing and consuming unnecessary fees.
		return types.ErrNoOpMsg
	}

	if channel.State != types.OPEN {
		return sdkerrors.Wrapf(
			types.ErrInvalidChannelState,
			"channel state is not OPEN (got %s)", channel.State.String(),
		)
	}

	packetCommitment := types.CommitPacket(k.cdc, packet)

	// verify we sent the packet and haven't cleared it out yet
	if !bytes.Equal(commitment, packetCommitment) {
		return sdkerrors.Wrapf(types.ErrInvalidPacket, "packet commitment bytes are not equal: got (%v), expected (%v)", commitment, packetCommitment)
	}

	switch channel.Ordering {
	case types.ORDERED:
		// check that packet has not been received
		if nextSequenceRecv > packet.GetSequence() {
			return sdkerrors.Wrapf(
				types.ErrPacketReceived,
				"packet already received, next sequence receive > packet sequence (%d > %d)", nextSequenceRecv, packet.GetSequence(),
			)
		}

		// check that the recv sequence is as claimed
		if len(channel.ConnectionHops) > 1 {
			// verify multihop proof
			kvGenerator := func(_ *types.MsgMultihopProofs, _ *connectiontypes.ConnectionEnd) (string, []byte, error) {
				key := host.NextSequenceRecvPath(packet.GetSourcePort(), packet.GetSourceChannel())
				value := sdk.Uint64ToBigEndian(nextSequenceRecv)
				return key, value, nil
			}

			err = k.connectionKeeper.VerifyMultihopMembership(
				ctx, connectionEnd, proofHeight, proof,
				channel.ConnectionHops, kvGenerator)
		} else {
			err = k.connectionKeeper.VerifyNextSequenceRecv(
				ctx, connectionEnd, proofHeight, proof,
				packet.GetDestPort(), packet.GetDestChannel(), nextSequenceRecv,
			)
		}
	case types.UNORDERED:
		if len(channel.ConnectionHops) > 1 {
			// verify multihop proof
			keyGenerator := func(_ *types.MsgMultihopProofs, _ *connectiontypes.ConnectionEnd) (string, error) {
				key := host.PacketReceiptPath(
					packet.GetSourcePort(),
					packet.GetSourceChannel(),
					packet.GetSequence(),
				)
				return key, nil
			}

			err = k.connectionKeeper.VerifyMultihopNonMembership(
				ctx, connectionEnd, proofHeight, proof,
				channel.ConnectionHops, keyGenerator)
		} else {
			err = k.connectionKeeper.VerifyPacketReceiptAbsence(
				ctx, connectionEnd, proofHeight, proof,
				packet.GetDestPort(), packet.GetDestChannel(), packet.GetSequence(),
			)
		}
	default:
		panic(sdkerrors.Wrapf(types.ErrInvalidChannelOrdering, channel.Ordering.String()))
	}

	if err != nil {
		return err
	}

	// NOTE: the remaining code is located in the TimeoutExecuted function
	return nil
}
