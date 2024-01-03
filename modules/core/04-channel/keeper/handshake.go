package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	"github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	commitmenttypes "github.com/cosmos/ibc-go/v8/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	tmclient "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
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
		return "", nil, errorsmod.Wrap(connectiontypes.ErrConnectionNotFound, connectionHops[0])
	}

	// ******************************************************************************************
	// TODO: This check does not apply to multihop channels. For multihop we need the connectionEnd
	// corresponding to the connection to chain Z for the logic to be meaningful.
	// ******************************************************************************************

	// getVersions := connectionEnd.GetVersions()
	// if len(getVersions) != 1 {
	// 	return "", nil, errorsmod.Wrapf(
	// 		connectiontypes.ErrInvalidVersion,
	// 		"single version must be negotiated on connection before opening channel, got: %v",
	// 		getVersions,
	// 	)
	// }

	getVersions := connectionEnd.GetVersions()
	if len(getVersions) != 1 {
		return "", nil, errorsmod.Wrapf(
			connectiontypes.ErrInvalidVersion,
			"single version must be negotiated on connection before opening channel, got: %v",
			getVersions,
		)
	}

	if !connectiontypes.VerifySupportedFeature(getVersions[0], order.String()) {
		return "", nil, errorsmod.Wrapf(
			connectiontypes.ErrInvalidVersion,
			"connection version %s does not support channel ordering: %s",
			getVersions[0], order.String(),
		)
	}

	clientState, found := k.clientKeeper.GetClientState(ctx, connectionEnd.ClientId)
	if !found {
		return "", nil, errorsmod.Wrapf(clienttypes.ErrClientNotFound, "clientID (%s)", connectionEnd.ClientId)
	}

	if status := k.clientKeeper.GetClientStatus(ctx, clientState, connectionEnd.ClientId); status != exported.Active {
		return "", nil, errorsmod.Wrapf(clienttypes.ErrClientNotActive, "client (%s) status is %s", connectionEnd.ClientId, status)
	}

	if !k.portKeeper.Authenticate(ctx, portCap, portID) {
		return "", nil, errorsmod.Wrapf(porttypes.ErrInvalidPort, "caller does not own port capability for port ID %s", portID)
	}

	channelID := k.GenerateChannelIdentifier(ctx)

	capKey, err := k.scopedKeeper.NewCapability(ctx, host.ChannelCapabilityPath(portID, channelID))
	if err != nil {
		return "", nil, errorsmod.Wrapf(err, "could not create channel capability for port ID %s and channel ID %s", portID, channelID)
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

	k.Logger(ctx).Info("channel state updated", "port-id", portID, "channel-id", channelID, "previous-state", types.UNINITIALIZED.String(), "new-state", types.INIT.String())

	defer telemetry.IncrCounter(1, "ibc", "channel", "open-init")

	emitChannelOpenInitEvent(ctx, portID, channelID, channel)
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
		return "", nil, errorsmod.Wrapf(porttypes.ErrInvalidPort, "caller does not own port capability for port ID %s", portID)
	}

	// Directly verify the last connectionHop. In a multihop hop scenario only the final
	// connection hop can be verified directly. The remaining connections are verified below.
	connectionEnd, found := k.connectionKeeper.GetConnection(ctx, connectionHops[0])
	if !found {
		return "", nil, errorsmod.Wrap(connectiontypes.ErrConnectionNotFound, connectionHops[0])
	}

	if connectionEnd.GetState() != int32(connectiontypes.OPEN) {
		return "", nil, errorsmod.Wrapf(
			connectiontypes.ErrInvalidConnectionState,
			"connection state is not OPEN (got %s)", connectiontypes.State(connectionEnd.GetState()).String(),
		)
	}
	// determine counterparty hops
	counterpartyHops := []string{connectionEnd.GetCounterparty().GetConnectionID()}

	// check version support
	if err := checkVersion(&connectionEnd, order); err != nil {
		return "", nil, err
	}

	// expectedCounterparty is the counterparty of the counterparty's channel end (i.e self)
	expectedCounterparty := types.NewCounterparty(portID, "")
	expectedChannel := types.NewChannel(
		types.INIT, order, expectedCounterparty,
		counterpartyHops, counterpartyVersion,
	)

	if err := k.connectionKeeper.VerifyChannelState(
		ctx, connectionHops, proofHeight, proofInit,
		counterparty.PortId, counterparty.ChannelId, expectedChannel,
	); err != nil {
		return "", nil, err
	}

	var (
		capKey *capabilitytypes.Capability
		err    error
	)

	capKey, err = k.scopedKeeper.NewCapability(ctx, host.ChannelCapabilityPath(portID, channelID))
	if err != nil {
		return "", nil, errorsmod.Wrapf(err, "could not create channel capability for port ID %s and channel ID %s", portID, channelID)
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

	k.Logger(ctx).Info("channel state updated", "port-id", portID, "channel-id", channelID, "previous-state", types.UNINITIALIZED.String(), "new-state", types.TRYOPEN.String())

	defer telemetry.IncrCounter(1, "ibc", "channel", "open-try")

	emitChannelOpenTryEvent(ctx, portID, channelID, channel)
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
		return errorsmod.Wrapf(types.ErrChannelNotFound, "port ID (%s) channel ID (%s)", portID, channelID)
	}

	if channel.State != types.INIT {
		return errorsmod.Wrapf(types.ErrInvalidChannelState, "channel state should be INIT (got %s)", channel.State.String())
	}

	if !k.scopedKeeper.AuthenticateCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)) {
		return errorsmod.Wrapf(types.ErrChannelCapabilityNotFound, "caller does not own capability for channel, port ID (%s) channel ID (%s)", portID, channelID)
	}

	connectionEnd, found := k.connectionKeeper.GetConnection(ctx, channel.ConnectionHops[0])
	if !found {
		return errorsmod.Wrap(connectiontypes.ErrConnectionNotFound, channel.ConnectionHops[0])
	}

	if connectionEnd.GetState() != int32(connectiontypes.OPEN) {
		return errorsmod.Wrapf(
			connectiontypes.ErrInvalidConnectionState,
			"connection state is not OPEN (got %s)", connectiontypes.State(connectionEnd.GetState()).String(),
		)
	}
	counterpartyHops := []string{connectionEnd.GetCounterparty().GetConnectionID()}
	expectedCounterparty := types.NewCounterparty(portID, channelID)
	expectedChannel := types.NewChannel(
		types.TRYOPEN, channel.Ordering, expectedCounterparty,
		counterpartyHops, counterpartyVersion,
	)

	return k.connectionKeeper.VerifyChannelState(
		ctx, channel.ConnectionHops, proofHeight, proofTry,
		channel.Counterparty.PortId, counterpartyChannelID,
		expectedChannel,
	)
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
		panic(fmt.Errorf("could not find existing channel when updating channel state in successful ChanOpenAck step, channelID: %s, portID: %s", channelID, portID))
	}

	channel.State = types.OPEN
	channel.Version = counterpartyVersion
	channel.Counterparty.ChannelId = counterpartyChannelID
	k.SetChannel(ctx, portID, channelID, channel)

	k.Logger(ctx).Info("channel state updated", "port-id", portID, "channel-id", channelID, "previous-state", types.INIT.String(), "new-state", types.OPEN.String())

	defer telemetry.IncrCounter(1, "ibc", "channel", "open-ack")

	emitChannelOpenAckEvent(ctx, portID, channelID, channel)
}

// ChanOpenConfirm is called by the handshake-accepting module to confirm the acknowledgement
// of the handshake-originating module on the other chain and finish the channel opening handshake.
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
		return errorsmod.Wrapf(types.ErrChannelNotFound, "port ID (%s) channel ID (%s)", portID, channelID)
	}

	if channel.State != types.TRYOPEN {
		return errorsmod.Wrapf(
			types.ErrInvalidChannelState,
			"channel state is not TRYOPEN (got %s)", channel.State.String(),
		)
	}

	if !k.scopedKeeper.AuthenticateCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)) {
		return errorsmod.Wrapf(types.ErrChannelCapabilityNotFound, "caller does not own capability for channel, port ID (%s) channel ID (%s)", portID, channelID)
	}

	connectionEnd, found := k.connectionKeeper.GetConnection(ctx, channel.ConnectionHops[0])
	if !found {
		return errorsmod.Wrap(connectiontypes.ErrConnectionNotFound, channel.ConnectionHops[0])
	}

	if connectionEnd.GetState() != int32(connectiontypes.OPEN) {
		return errorsmod.Wrapf(
			connectiontypes.ErrInvalidConnectionState,
			"connection state is not OPEN (got %s)", connectiontypes.State(connectionEnd.GetState()).String(),
		)
	}

	counterpartyHops := []string{connectionEnd.GetCounterparty().GetConnectionID()}
	counterparty := types.NewCounterparty(portID, channelID)
	expectedChannel := types.NewChannel(
		types.OPEN, channel.Ordering, counterparty,
		counterpartyHops, channel.Version)

	return k.connectionKeeper.VerifyChannelState(
		ctx, channel.ConnectionHops, proofHeight, proofAck,
		channel.Counterparty.PortId, channel.Counterparty.ChannelId,
		expectedChannel)
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
		panic(fmt.Errorf("could not find existing channel when updating channel state in successful ChanOpenConfirm step, channelID: %s, portID: %s", channelID, portID))
	}

	channel.State = types.OPEN
	k.SetChannel(ctx, portID, channelID, channel)
	k.Logger(ctx).Info("channel state updated", "port-id", portID, "channel-id", channelID, "previous-state", types.TRYOPEN.String(), "new-state", types.OPEN.String())

	defer telemetry.IncrCounter(1, "ibc", "channel", "open-confirm")

	emitChannelOpenConfirmEvent(ctx, portID, channelID, channel)
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
		return errorsmod.Wrapf(types.ErrChannelCapabilityNotFound, "caller does not own capability for channel, port ID (%s) channel ID (%s)", portID, channelID)
	}

	channel, found := k.GetChannel(ctx, portID, channelID)
	if !found {
		return errorsmod.Wrapf(types.ErrChannelNotFound, "port ID (%s) channel ID (%s)", portID, channelID)
	}

	if channel.IsClosed() {
		return errorsmod.Wrap(types.ErrInvalidChannelState, "channel is already CLOSED")
	}

	connectionEnd, found := k.connectionKeeper.GetConnection(ctx, channel.ConnectionHops[0])
	if !found {
		return errorsmod.Wrap(connectiontypes.ErrConnectionNotFound, channel.ConnectionHops[0])
	}

	clientState, found := k.clientKeeper.GetClientState(ctx, connectionEnd.ClientId)
	if !found {
		return errorsmod.Wrapf(clienttypes.ErrClientNotFound, "clientID (%s)", connectionEnd.ClientId)
	}

	if status := k.clientKeeper.GetClientStatus(ctx, clientState, connectionEnd.ClientId); status != exported.Active {
		return errorsmod.Wrapf(clienttypes.ErrClientNotActive, "client (%s) status is %s", connectionEnd.ClientId, status)
	}

	if connectionEnd.GetState() != int32(connectiontypes.OPEN) {
		return errorsmod.Wrapf(
			connectiontypes.ErrInvalidConnectionState,
			"connection state is not OPEN (got %s)", connectiontypes.State(connectionEnd.GetState()).String(),
		)
	}

	k.Logger(ctx).Info("channel state updated", "port-id", portID, "channel-id", channelID, "previous-state", channel.State.String(), "new-state", types.CLOSED.String())

	defer telemetry.IncrCounter(1, "ibc", "channel", "close-init")

	channel.State = types.CLOSED
	k.SetChannel(ctx, portID, channelID, channel)

	emitChannelCloseInitEvent(ctx, portID, channelID, channel)

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
		return errorsmod.Wrap(types.ErrChannelCapabilityNotFound, "caller does not own capability for channel, port ID (%s) channel ID (%s)")
	}

	channel, found := k.GetChannel(ctx, portID, channelID)
	if !found {
		return errorsmod.Wrapf(types.ErrChannelNotFound, "port ID (%s) channel ID (%s)", portID, channelID)
	}

	if channel.IsClosed() {
		return errorsmod.Wrap(types.ErrInvalidChannelState, "channel is already CLOSED")
	}

	connectionEnd, found := k.connectionKeeper.GetConnection(ctx, channel.ConnectionHops[0])
	if !found {
		return errorsmod.Wrap(connectiontypes.ErrConnectionNotFound, channel.ConnectionHops[0])
	}

	if connectionEnd.GetState() != int32(connectiontypes.OPEN) {
		return errorsmod.Wrapf(
			connectiontypes.ErrInvalidConnectionState,
			"connection state is not OPEN (got %s)", connectiontypes.State(connectionEnd.GetState()).String(),
		)
	}

	counterpartyHops := []string{connectionEnd.GetCounterparty().GetConnectionID()}
	counterparty := types.NewCounterparty(portID, channelID)
	expectedChannel := types.NewChannel(
		types.CLOSED, channel.Ordering, counterparty,
		counterpartyHops, channel.Version)

	if err := k.connectionKeeper.VerifyChannelState(
		ctx, channel.ConnectionHops, proofHeight, proofInit,
		channel.Counterparty.PortId, channel.Counterparty.ChannelId,
		expectedChannel); err != nil {
		return err
	}

	k.Logger(ctx).Info("channel state updated", "port-id", portID, "channel-id", channelID, "previous-state", channel.State.String(), "new-state", types.CLOSED.String())

	defer telemetry.IncrCounter(1, "ibc", "channel", "close-confirm")

	channel.State = types.CLOSED
	k.SetChannel(ctx, portID, channelID, channel)

	emitChannelCloseConfirmEvent(ctx, portID, channelID, channel)

	return nil
}

// ChanCloseFrozen is called by the counterparty module to close their end of the
// channel due to a frozen client in the channel path. The proof arguments are
// expected to have been queried from the chain with the frozen client along the
// channel path.
//
// proofConnection - proof of the connectionEnd for the frozen chain (contains the clientID)
// proofClientState - proof of the frozen clientState (clientID determined from the connectionEnd)
func (k Keeper) ChanCloseFrozen(
	ctx sdk.Context,
	portID,
	channelID string,
	chanCap *capabilitytypes.Capability,
	proofConnection []byte,
	proofClientState []byte,
	proofHeight exported.Height,
) error {
	if !k.scopedKeeper.AuthenticateCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)) {
		return errorsmod.Wrap(types.ErrChannelCapabilityNotFound, "caller does not own capability for channel, port ID (%s) channel ID (%s)")
	}

	channel, found := k.GetChannel(ctx, portID, channelID)
	if !found {
		return errorsmod.Wrapf(types.ErrChannelNotFound, "port ID (%s) channel ID (%s)", portID, channelID)
	}

	// ChanCloseFrozen only applies to multi-hop channels
	if len(channel.ConnectionHops) <= 1 {
		return types.ErrInvalidHopLength
	}

	if channel.State == types.CLOSED {
		return errorsmod.Wrap(types.ErrInvalidChannelState, "channel is already CLOSED")
	}

	connectionEnd, found := k.connectionKeeper.GetConnection(ctx, channel.ConnectionHops[0])
	if !found {
		return errorsmod.Wrap(
			connectiontypes.ErrConnectionNotFound,
			channel.ConnectionHops[len(channel.ConnectionHops)-1],
		)
	}

	if connectionEnd.GetState() != int32(connectiontypes.OPEN) {
		return errorsmod.Wrapf(
			connectiontypes.ErrInvalidConnectionState,
			"connection state is not OPEN (got %s)", connectiontypes.State(connectionEnd.GetState()).String(),
		)
	}

	var mProof types.MsgMultihopProofs
	if err := k.cdc.Unmarshal(proofConnection, &mProof); err != nil {
		return fmt.Errorf("cannot unmarshal proof: %v", err)
	}

	connectionIdx := len(mProof.ConsensusProofs) + 1
	if connectionIdx > len(channel.ConnectionHops)-1 {
		return fmt.Errorf("connectionKVGenerator: connectionHops index out of range (%d > %d)", connectionIdx, len(channel.ConnectionHops)-1)
	}
	connectionID := channel.ConnectionHops[connectionIdx]
	connKey := commitmenttypes.NewMerklePath(host.ConnectionPath(connectionID))
	connValue := mProof.KeyProof.Value
	var frozenConnection connectiontypes.ConnectionEnd
	if err := k.cdc.Unmarshal(connValue, &connectionEnd); err != nil {
		return fmt.Errorf("connectionKVGenerator: %s", err)
	}
	clientID := frozenConnection.ClientId

	// truncated connectionHops (e.g. [0, 1, 2, 3] -> [0, 1])
	connectionHops := channel.ConnectionHops[:len(mProof.ConnectionProofs)+1]

	// prove connectionEnd
	if err := k.connectionKeeper.VerifyMultihopMembership(
		ctx, connectionHops, proofHeight, proofConnection,
		connKey, connValue); err != nil {
		return err
	}

	if err := k.cdc.Unmarshal(proofClientState, &mProof); err != nil {
		return fmt.Errorf("cannot unmarshal proof: %v", err)
	}

	clientKey := commitmenttypes.NewMerklePath(host.FullClientStatePath(clientID))
	clientValue := mProof.KeyProof.Value // client proof

	// unmarshal to client state interface
	var exportedClientState exported.ClientState
	if err := k.cdc.UnmarshalInterface(mProof.KeyProof.Value, &exportedClientState); err != nil {
		return err
	}

	// try to cast to tendermint client state
	cs, ok := exportedClientState.(*tmclient.ClientState)
	if !ok {
		return fmt.Errorf("cannot cast exported client state to tendermint client state")
	}

	// check client is frozen
	if cs.FrozenHeight.RevisionHeight == 0 {
		return fmt.Errorf("cannot close channel, client is not frozen")
	}

	// prove frozen client
	if err := k.connectionKeeper.VerifyMultihopMembership(
		ctx, connectionHops, proofHeight, proofClientState,
		clientKey, clientValue); err != nil {
		return err
	}

	k.Logger(ctx).Info("channel state updated", "port-id", portID, "channel-id", channelID, "previous-state", channel.State.String(), "new-state", "FROZEN")

	defer telemetry.IncrCounter(1, "ibc", "channel", "close-frozen")

	channel.State = types.FROZEN
	k.SetChannel(ctx, portID, channelID, channel)

	emitChannelCloseFrozenEvent(ctx, portID, channelID, channel)
	return nil
}

// chekcVersion checks that the connection's features include support for order.
func checkVersion(connection *connectiontypes.ConnectionEnd, order types.Order) error {
	getVersions := connection.GetVersions()
	if len(getVersions) != 1 {
		return errorsmod.Wrapf(
			connectiontypes.ErrInvalidVersion,
			"single version must be negotiated on connection before opening channel, got: %v",
			getVersions,
		)
	}

	// verify chain supports the requested ordering.
	if !connectiontypes.VerifySupportedFeature(getVersions[0], order.String()) {
		return errorsmod.Wrapf(
			connectiontypes.ErrInvalidVersion,
			"connection version %s does not support channel ordering: %s",
			getVersions[0], order.String(),
		)
	}

	return nil
}
