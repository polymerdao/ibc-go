package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"

	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	"github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	mh "github.com/cosmos/ibc-go/v7/modules/core/multihop"
)

// ChanOpenInit is called by a module to initiate a channel opening handshake with
// a module on another chain. The counterparty channel identifier is validated to be
// empty in msg validation.
func (k Keeper) ChanOpenInit(
	ctx sdk.Context,
	order types.Order,
	connectionHops []string,
	portID string,
	portCap *capabilitytypes.Capability,
	counterparty types.Counterparty,
	version string,
) (string, *capabilitytypes.Capability, error) {
	// connection hop length checked on msg.ValidateBasic()
	connectionEnd, found := k.connectionKeeper.GetConnection(ctx, connectionHops[0])
	if !found {
		return "", nil, sdkerrors.Wrap(connectiontypes.ErrConnectionNotFound, connectionHops[0])
	}

	// ******************************************************************************************
	// TODO: This is a bug for a multihop channels. For multihop we need the connectionEnd
	// corresponding to the connection to chain Z for the logic to be meaningful.
	// ******************************************************************************************

	if len(connectionHops) == 1 {
		getVersions := connectionEnd.GetVersions()
		if len(getVersions) != 1 {
			return "", nil, sdkerrors.Wrapf(
				connectiontypes.ErrInvalidVersion,
				"single version must be negotiated on connection before opening channel, got: %v",
				getVersions,
			)
		}

		if !connectiontypes.VerifySupportedFeature(getVersions[0], order.String()) {
			return "", nil, sdkerrors.Wrapf(
				connectiontypes.ErrInvalidVersion,
				"connection version %s does not support channel ordering: %s",
				getVersions[0], order.String(),
			)
		}
	}

	if !k.portKeeper.Authenticate(ctx, portCap, portID) {
		return "", nil, sdkerrors.Wrapf(porttypes.ErrInvalidPort, "caller does not own port capability for port ID %s", portID)
	}

	channelID := k.GenerateChannelIdentifier(ctx)

	capKey, err := k.scopedKeeper.NewCapability(ctx, host.ChannelCapabilityPath(portID, channelID))
	if err != nil {
		return "", nil, sdkerrors.Wrapf(err, "could not create channel capability for port ID %s and channel ID %s", portID, channelID)
	}

	return channelID, capKey, nil
}

// WriteOpenInitChannel writes a channel which has successfully passed the OpenInit handshake step.
// The channel is set in state and all the associated Send and Recv sequences are set to 1.
// An event is emitted for the handshake step.
func (k Keeper) WriteOpenInitChannel(
	ctx sdk.Context,
	portID,
	channelID string,
	order types.Order,
	connectionHops []string,
	counterparty types.Counterparty,
	version string,
) {
	channel := types.NewChannel(types.INIT, order, counterparty, connectionHops, version)
	k.SetChannel(ctx, portID, channelID, channel)

	k.SetNextSequenceSend(ctx, portID, channelID, 1)
	k.SetNextSequenceRecv(ctx, portID, channelID, 1)
	k.SetNextSequenceAck(ctx, portID, channelID, 1)

	k.Logger(ctx).Info("channel state updated", "port-id", portID, "channel-id", channelID, "previous-state", "NONE", "new-state", "INIT")

	defer func() {
		telemetry.IncrCounter(1, "ibc", "channel", "open-init")
	}()

	EmitChannelOpenInitEvent(ctx, portID, channelID, channel)
}

// ChanOpenTry is called by a module to accept the first step of a channel opening
// handshake initiated by a module on another chain.
func (k Keeper) ChanOpenTry(
	ctx sdk.Context,
	order types.Order,
	connectionHops []string,
	portID string,
	portCap *capabilitytypes.Capability,
	counterparty types.Counterparty,
	counterpartyVersion string,
	proofInit []byte,
	proofHeight exported.Height,
) (string, *capabilitytypes.Capability, error) {

	// generate a new channel
	channelID := k.GenerateChannelIdentifier(ctx)

	if !k.portKeeper.Authenticate(ctx, portCap, portID) {
		return "", nil, sdkerrors.Wrapf(porttypes.ErrInvalidPort, "caller does not own port capability for port ID %s", portID)
	}

	// Directly verify the last connectionHop. In a multihop hop scenario only the final
	// connection hop can be verified directly. The remaining connections are verified below.
	connectionEnd, found := k.connectionKeeper.GetConnection(ctx, connectionHops[0])
	if !found {
		return "", nil, sdkerrors.Wrap(connectiontypes.ErrConnectionNotFound, connectionHops[0])
	}

	if connectionEnd.GetState() != int32(connectiontypes.OPEN) {
		return "", nil, sdkerrors.Wrapf(
			connectiontypes.ErrInvalidConnectionState,
			"connection state is not OPEN (got %s)", connectiontypes.State(connectionEnd.GetState()).String(),
		)
	}

	// connectionHops A --> Z
	var counterpartyHops []string
	var checkEnd *connectiontypes.ConnectionEnd
	if len(connectionHops) > 1 {
		var err error
		// get the connectionEnd associated with chain A
		checkEnd, err = mh.GetMultihopConnectionEnd(k.cdc, proofInit)
		if err != nil {
			return "", nil, err
		}

		counterpartyHops, err = mh.GetCounterPartyHops(k.cdc, proofInit, &connectionEnd)
		if err != nil {
			return "", nil, err
		}

	} else {
		counterpartyHops = []string{connectionEnd.GetCounterparty().GetConnectionID()}
		checkEnd = &connectionEnd
	}

	// get chain A's version
	getVersions := checkEnd.GetVersions()
	if len(getVersions) != 1 {
		return "", nil, sdkerrors.Wrapf(
			connectiontypes.ErrInvalidVersion,
			"single version must be negotiated on connection before opening channel, got: %v",
			getVersions,
		)
	}

	// verify chain A supports the requested ordering.
	if !connectiontypes.VerifySupportedFeature(getVersions[0], order.String()) {
		return "", nil, sdkerrors.Wrapf(
			connectiontypes.ErrInvalidVersion,
			"connection version %s does not support channel ordering: %s",
			getVersions[0], order.String(),
		)
	}

	// expectedCounterpaty is the counterparty of the counterparty's channel end
	// (i.e self)
	expectedCounterparty := types.NewCounterparty(portID, "")
	expectedChannel := types.NewChannel(
		types.INIT, order, expectedCounterparty,
		counterpartyHops, counterpartyVersion,
	)

	// handle multihop case
	if len(connectionHops) > 1 {

		// expected value bytes
		value, err := expectedChannel.Marshal()
		if err != nil {
			return "", nil, err
		}

		// verify multihop proof
		consensusState, found := k.clientKeeper.GetClientConsensusState(ctx, connectionEnd.ClientId, proofHeight)
		if !found {
			return "", nil, sdkerrors.Wrapf(clienttypes.ErrConsensusStateNotFound,
				"consensus state not found for client id: %s", connectionEnd.ClientId)
		}

		// connectionHops is ZA, we want AZ
		if err := mh.VerifyMultihopProof(k.cdc, consensusState, connectionHops, proofInit, value); err != nil {
			return "", nil, err
		}
	} else {

		if err := k.connectionKeeper.VerifyChannelState(
			ctx, connectionEnd, proofHeight, proofInit,
			counterparty.PortId, counterparty.ChannelId, expectedChannel,
		); err != nil {
			return "", nil, err
		}
	}

	var (
		capKey *capabilitytypes.Capability
		err    error
	)

	capKey, err = k.scopedKeeper.NewCapability(ctx, host.ChannelCapabilityPath(portID, channelID))
	if err != nil {
		return "", nil, sdkerrors.Wrapf(err, "could not create channel capability for port ID %s and channel ID %s", portID, channelID)
	}

	return channelID, capKey, nil
}

// WriteOpenTryChannel writes a channel which has successfully passed the OpenTry handshake step.
// The channel is set in state. If a previous channel state did not exist, all the Send and Recv
// sequences are set to 1. An event is emitted for the handshake step.
func (k Keeper) WriteOpenTryChannel(
	ctx sdk.Context,
	portID,
	channelID string,
	order types.Order,
	connectionHops []string,
	counterparty types.Counterparty,
	version string,
) {
	k.SetNextSequenceSend(ctx, portID, channelID, 1)
	k.SetNextSequenceRecv(ctx, portID, channelID, 1)
	k.SetNextSequenceAck(ctx, portID, channelID, 1)

	channel := types.NewChannel(types.TRYOPEN, order, counterparty, connectionHops, version)

	k.SetChannel(ctx, portID, channelID, channel)

	k.Logger(ctx).Info("channel state updated", "port-id", portID, "channel-id", channelID, "previous-state", "NONE", "new-state", "TRYOPEN")

	defer func() {
		telemetry.IncrCounter(1, "ibc", "channel", "open-try")
	}()

	EmitChannelOpenTryEvent(ctx, portID, channelID, channel)
}

// ChanOpenAck is called by the handshake-originating module to acknowledge the
// acceptance of the initial request by the counterparty module on the other chain.
func (k Keeper) ChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterpartyVersion,
	counterpartyChannelID string,
	proofTry []byte,
	proofHeight exported.Height,
) error {
	channel, found := k.GetChannel(ctx, portID, channelID)
	if !found {
		return sdkerrors.Wrapf(types.ErrChannelNotFound, "port ID (%s) channel ID (%s)", portID, channelID)
	}

	if channel.State != types.INIT {
		return sdkerrors.Wrapf(types.ErrInvalidChannelState, "channel state should be INIT (got %s)", channel.State.String())
	}

	if !k.scopedKeeper.AuthenticateCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)) {
		return sdkerrors.Wrapf(types.ErrChannelCapabilityNotFound, "caller does not own capability for channel, port ID (%s) channel ID (%s)", portID, channelID)
	}

	connectionEnd, found := k.connectionKeeper.GetConnection(ctx, channel.ConnectionHops[0])
	if !found {
		return sdkerrors.Wrap(connectiontypes.ErrConnectionNotFound, channel.ConnectionHops[0])
	}

	if connectionEnd.GetState() != int32(connectiontypes.OPEN) {
		return sdkerrors.Wrapf(
			connectiontypes.ErrInvalidConnectionState,
			"connection state is not OPEN (got %s)", connectiontypes.State(connectionEnd.GetState()).String(),
		)
	}

	// connectionHops Z --> A
	var counterpartyHops []string
	if len(channel.ConnectionHops) > 1 {
		var err error
		counterpartyHops, err = mh.GetCounterPartyHops(k.cdc, proofTry, &connectionEnd)
		if err != nil {
			return err
		}
	} else {
		counterpartyHops = []string{connectionEnd.GetCounterparty().GetConnectionID()}
	}

	expectedCounterparty := types.NewCounterparty(portID, channelID)
	expectedChannel := types.NewChannel(
		types.TRYOPEN, channel.Ordering, expectedCounterparty,
		counterpartyHops, counterpartyVersion,
	)

	// verify multihop proof
	if len(channel.ConnectionHops) > 1 {
		value, err := expectedChannel.Marshal()
		if err != nil {
			return err
		}

		consensusState, found := k.clientKeeper.GetClientConsensusState(ctx, connectionEnd.ClientId, proofHeight)
		if !found {
			return sdkerrors.Wrapf(clienttypes.ErrConsensusStateNotFound,
				"consensus state not found for client id: %s", connectionEnd.ClientId)
		}

		if err := mh.VerifyMultihopProof(k.cdc, consensusState, channel.ConnectionHops, proofTry, value); err != nil {
			return err
		}
	} else {
		if err := k.connectionKeeper.VerifyChannelState(
			ctx, connectionEnd, proofHeight, proofTry,
			channel.Counterparty.PortId, counterpartyChannelID, expectedChannel,
		); err != nil {
			return err
		}
	}

	return nil
}

// WriteOpenAckChannel writes an updated channel state for the successful OpenAck handshake step.
// An event is emitted for the handshake step.
func (k Keeper) WriteOpenAckChannel(
	ctx sdk.Context,
	portID,
	channelID,
	counterpartyVersion,
	counterpartyChannelID string,
) {
	channel, found := k.GetChannel(ctx, portID, channelID)
	if !found {
		panic(fmt.Sprintf("could not find existing channel when updating channel state in successful ChanOpenAck step, channelID: %s, portID: %s", channelID, portID))
	}

	channel.State = types.OPEN
	channel.Version = counterpartyVersion
	channel.Counterparty.ChannelId = counterpartyChannelID
	k.SetChannel(ctx, portID, channelID, channel)

	k.Logger(ctx).Info("channel state updated", "port-id", portID, "channel-id", channelID, "previous-state", channel.State.String(), "new-state", "OPEN")

	defer func() {
		telemetry.IncrCounter(1, "ibc", "channel", "open-ack")
	}()

	EmitChannelOpenAckEvent(ctx, portID, channelID, channel)
}

// ChanOpenConfirm is called by the counterparty module to close their end of the
// channel, since the other end has been closed.
func (k Keeper) ChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
	chanCap *capabilitytypes.Capability,
	proofAck []byte,
	proofHeight exported.Height,
) error {
	channel, found := k.GetChannel(ctx, portID, channelID)
	if !found {
		return sdkerrors.Wrapf(types.ErrChannelNotFound, "port ID (%s) channel ID (%s)", portID, channelID)
	}

	if channel.State != types.TRYOPEN {
		return sdkerrors.Wrapf(
			types.ErrInvalidChannelState,
			"channel state is not TRYOPEN (got %s)", channel.State.String(),
		)
	}

	if !k.scopedKeeper.AuthenticateCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)) {
		return sdkerrors.Wrapf(types.ErrChannelCapabilityNotFound, "caller does not own capability for channel, port ID (%s) channel ID (%s)", portID, channelID)
	}

	connectionEnd, found := k.connectionKeeper.GetConnection(ctx, channel.ConnectionHops[0])
	if !found {
		return sdkerrors.Wrap(connectiontypes.ErrConnectionNotFound, channel.ConnectionHops[0])
	}

	if connectionEnd.GetState() != int32(connectiontypes.OPEN) {
		return sdkerrors.Wrapf(
			connectiontypes.ErrInvalidConnectionState,
			"connection state is not OPEN (got %s)", connectiontypes.State(connectionEnd.GetState()).String(),
		)
	}

	// connectionHops A --> Z
	var counterpartyHops []string
	if len(channel.ConnectionHops) > 1 {
		var err error
		counterpartyHops, err = mh.GetCounterPartyHops(k.cdc, proofAck, &connectionEnd)
		if err != nil {
			return err
		}
	} else {
		counterpartyHops = []string{connectionEnd.GetCounterparty().GetConnectionID()}
	}

	counterparty := types.NewCounterparty(portID, channelID)
	expectedChannel := types.NewChannel(
		types.OPEN, channel.Ordering, counterparty,
		counterpartyHops, channel.Version,
	)

	// verify multihop proof or standard proof
	if len(channel.ConnectionHops) > 1 {
		value, err := expectedChannel.Marshal()
		if err != nil {
			return err
		}

		// verify multihop proof
		consensusState, found := k.clientKeeper.GetClientConsensusState(ctx, connectionEnd.ClientId, proofHeight)
		if !found {
			return sdkerrors.Wrapf(clienttypes.ErrConsensusStateNotFound,
				"consensus state not found for client id: %s", connectionEnd.ClientId)
		}

		if err := mh.VerifyMultihopProof(k.cdc, consensusState, channel.ConnectionHops, proofAck, value); err != nil {
			return err
		}
	} else {
		if err := k.connectionKeeper.VerifyChannelState(
			ctx, connectionEnd, proofHeight, proofAck,
			channel.Counterparty.PortId, channel.Counterparty.ChannelId,
			expectedChannel,
		); err != nil {
			return err
		}
	}

	return nil
}

// WriteOpenConfirmChannel writes an updated channel state for the successful OpenConfirm handshake step.
// An event is emitted for the handshake step.
func (k Keeper) WriteOpenConfirmChannel(
	ctx sdk.Context,
	portID,
	channelID string,
) {
	channel, found := k.GetChannel(ctx, portID, channelID)
	if !found {
		panic(fmt.Sprintf("could not find existing channel when updating channel state in successful ChanOpenConfirm step, channelID: %s, portID: %s", channelID, portID))
	}

	channel.State = types.OPEN
	k.SetChannel(ctx, portID, channelID, channel)
	k.Logger(ctx).Info("channel state updated", "port-id", portID, "channel-id", channelID, "previous-state", "TRYOPEN", "new-state", "OPEN")

	defer func() {
		telemetry.IncrCounter(1, "ibc", "channel", "open-confirm")
	}()

	EmitChannelOpenConfirmEvent(ctx, portID, channelID, channel)
}

// Closing Handshake
//
// This section defines the set of functions required to close a channel handshake
// as defined in https://github.com/cosmos/ibc/tree/master/spec/core/ics-004-channel-and-packet-semantics#closing-handshake
//
// ChanCloseInit is called by either module to close their end of the channel. Once
// closed, channels cannot be reopened.
func (k Keeper) ChanCloseInit(
	ctx sdk.Context,
	portID,
	channelID string,
	chanCap *capabilitytypes.Capability,
) error {
	if !k.scopedKeeper.AuthenticateCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)) {
		return sdkerrors.Wrapf(types.ErrChannelCapabilityNotFound, "caller does not own capability for channel, port ID (%s) channel ID (%s)", portID, channelID)
	}

	channel, found := k.GetChannel(ctx, portID, channelID)
	if !found {
		return sdkerrors.Wrapf(types.ErrChannelNotFound, "port ID (%s) channel ID (%s)", portID, channelID)
	}

	if channel.State == types.CLOSED {
		return sdkerrors.Wrap(types.ErrInvalidChannelState, "channel is already CLOSED")
	}

	// TODO: skip this if connectionHops > 1 (alternative would be to pass connection state along with proof)

	connectionEnd, found := k.connectionKeeper.GetConnection(ctx, channel.ConnectionHops[0])
	if !found {
		return sdkerrors.Wrap(connectiontypes.ErrConnectionNotFound, channel.ConnectionHops[0])
	}

	if connectionEnd.GetState() != int32(connectiontypes.OPEN) {
		return sdkerrors.Wrapf(
			connectiontypes.ErrInvalidConnectionState,
			"connection state is not OPEN (got %s)", connectiontypes.State(connectionEnd.GetState()).String(),
		)
	}

	k.Logger(ctx).Info("channel state updated", "port-id", portID, "channel-id", channelID, "previous-state", channel.State.String(), "new-state", "CLOSED")

	defer func() {
		telemetry.IncrCounter(1, "ibc", "channel", "close-init")
	}()

	channel.State = types.CLOSED
	k.SetChannel(ctx, portID, channelID, channel)

	EmitChannelCloseInitEvent(ctx, portID, channelID, channel)

	return nil
}

// ChanCloseConfirm is called by the counterparty module to close their end of the
// channel, since the other end has been closed.
func (k Keeper) ChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
	chanCap *capabilitytypes.Capability,
	proofInit []byte,
	proofHeight exported.Height,
) error {
	if !k.scopedKeeper.AuthenticateCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)) {
		return sdkerrors.Wrap(types.ErrChannelCapabilityNotFound, "caller does not own capability for channel, port ID (%s) channel ID (%s)")
	}

	channel, found := k.GetChannel(ctx, portID, channelID)
	if !found {
		return sdkerrors.Wrapf(types.ErrChannelNotFound, "port ID (%s) channel ID (%s)", portID, channelID)
	}

	if channel.State == types.CLOSED {
		return sdkerrors.Wrap(types.ErrInvalidChannelState, "channel is already CLOSED")
	}

	connectionEnd, found := k.connectionKeeper.GetConnection(ctx, channel.ConnectionHops[len(channel.ConnectionHops)-1])
	if !found {
		return sdkerrors.Wrap(connectiontypes.ErrConnectionNotFound, channel.ConnectionHops[len(channel.ConnectionHops)-1])
	}

	if connectionEnd.GetState() != int32(connectiontypes.OPEN) {
		return sdkerrors.Wrapf(
			connectiontypes.ErrInvalidConnectionState,
			"connection state is not OPEN (got %s)", connectiontypes.State(connectionEnd.GetState()).String(),
		)
	}

	// determine counterparty hops
	var counterpartyHops []string
	if len(channel.ConnectionHops) > 1 {
		// get the connectionHops A --> Z
		var err error
		counterpartyHops, err = mh.GetCounterPartyHops(k.cdc, proofInit, &connectionEnd)
		if err != nil {
			return err
		}
	} else {
		counterpartyHops = []string{connectionEnd.GetCounterparty().GetConnectionID()}
	}

	counterparty := types.NewCounterparty(portID, channelID)
	expectedChannel := types.NewChannel(
		types.CLOSED, channel.Ordering, counterparty,
		counterpartyHops, channel.Version,
	)

	// verify multihop proof
	if len(channel.ConnectionHops) > 1 {

		value, err := expectedChannel.Marshal()
		if err != nil {
			return err
		}

		consensusState, found := k.clientKeeper.GetClientConsensusState(ctx, connectionEnd.ClientId, proofHeight)
		if !found {
			return sdkerrors.Wrapf(clienttypes.ErrConsensusStateNotFound,
				"consensus state not found for client id: %s", connectionEnd.ClientId)
		}

		if err := mh.VerifyMultihopProof(k.cdc, consensusState, channel.ConnectionHops, proofInit, value); err != nil {
			return err
		}

	} else { // standard proof
		if err := k.connectionKeeper.VerifyChannelState(
			ctx, connectionEnd, proofHeight, proofInit,
			channel.Counterparty.PortId, channel.Counterparty.ChannelId,
			expectedChannel,
		); err != nil {
			return err
		}
	}

	k.Logger(ctx).Info("channel state updated", "port-id", portID, "channel-id", channelID, "previous-state", channel.State.String(), "new-state", "CLOSED")

	defer func() {
		telemetry.IncrCounter(1, "ibc", "channel", "close-confirm")
	}()

	channel.State = types.CLOSED
	k.SetChannel(ctx, portID, channelID, channel)

	EmitChannelCloseConfirmEvent(ctx, portID, channelID, channel)

	return nil
}
