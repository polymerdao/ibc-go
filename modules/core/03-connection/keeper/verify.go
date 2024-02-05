package keeper

import (
	"math"

	errorsmod "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/v8/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
	mh "github.com/cosmos/ibc-go/v8/modules/core/33-multihop"
	ibcerrors "github.com/cosmos/ibc-go/v8/modules/core/errors"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

// VerifyClientState verifies a proof of a client state of the running machine
// stored on the target machine
func (k Keeper) VerifyClientState(
	ctx sdk.Context,
	connection exported.ConnectionI,
	height exported.Height,
	proof []byte,
	clientState exported.ClientState,
) error {
	clientID := connection.GetClientID()
	targetClient, clientStore, err := k.getClientStateAndVerificationStore(ctx, clientID)
	if err != nil {
		return err
	}

	if status := k.clientKeeper.GetClientStatus(ctx, targetClient, clientID); status != exported.Active {
		return errorsmod.Wrapf(clienttypes.ErrClientNotActive, "client (%s) status is %s", clientID, status)
	}

	merklePath := commitmenttypes.NewMerklePath(host.FullClientStatePath(connection.GetCounterparty().GetClientID()))
	merklePath, err = commitmenttypes.ApplyPrefix(connection.GetCounterparty().GetPrefix(), merklePath)
	if err != nil {
		return err
	}

	bz, err := k.cdc.MarshalInterface(clientState)
	if err != nil {
		return err
	}

	if err := targetClient.VerifyMembership(
		ctx, clientStore, k.cdc, height,
		0, 0, // skip delay period checks for non-packet processing verification
		proof, merklePath, bz,
	); err != nil {
		return errorsmod.Wrapf(err, "failed client state verification for target client: %s", clientID)
	}

	return nil
}

// VerifyClientConsensusState verifies a proof of the consensus state of the
// specified client stored on the target machine.
func (k Keeper) VerifyClientConsensusState(
	ctx sdk.Context,
	connection exported.ConnectionI,
	height exported.Height,
	consensusHeight exported.Height,
	proof []byte,
	consensusState exported.ConsensusState,
) error {
	clientID := connection.GetClientID()
	clientState, clientStore, err := k.getClientStateAndVerificationStore(ctx, clientID)
	if err != nil {
		return err
	}

	if status := k.clientKeeper.GetClientStatus(ctx, clientState, clientID); status != exported.Active {
		return errorsmod.Wrapf(clienttypes.ErrClientNotActive, "client (%s) status is %s", clientID, status)
	}

	merklePath := commitmenttypes.NewMerklePath(host.FullConsensusStatePath(connection.GetCounterparty().GetClientID(), consensusHeight))
	merklePath, err = commitmenttypes.ApplyPrefix(connection.GetCounterparty().GetPrefix(), merklePath)
	if err != nil {
		return err
	}

	bz, err := k.cdc.MarshalInterface(consensusState)
	if err != nil {
		return err
	}

	if err := clientState.VerifyMembership(
		ctx, clientStore, k.cdc, height,
		0, 0, // skip delay period checks for non-packet processing verification
		proof, merklePath, bz,
	); err != nil {
		return errorsmod.Wrapf(err, "failed consensus state verification for client (%s)", clientID)
	}

	return nil
}

// VerifyConnectionState verifies a proof of the connection state of the
// specified connection end stored on the target machine.
func (k Keeper) VerifyConnectionState(
	ctx sdk.Context,
	connection exported.ConnectionI,
	height exported.Height,
	proof []byte,
	connectionID string,
	counterpartyConnection exported.ConnectionI, // opposite connection
) error {
	clientID := connection.GetClientID()
	clientState, clientStore, err := k.getClientStateAndVerificationStore(ctx, clientID)
	if err != nil {
		return err
	}

	if status := k.clientKeeper.GetClientStatus(ctx, clientState, clientID); status != exported.Active {
		return errorsmod.Wrapf(clienttypes.ErrClientNotActive, "client (%s) status is %s", clientID, status)
	}

	merklePath := commitmenttypes.NewMerklePath(host.ConnectionPath(connectionID))
	merklePath, err = commitmenttypes.ApplyPrefix(connection.GetCounterparty().GetPrefix(), merklePath)
	if err != nil {
		return err
	}

	connectionEnd, ok := counterpartyConnection.(connectiontypes.ConnectionEnd)
	if !ok {
		return errorsmod.Wrapf(ibcerrors.ErrInvalidType, "invalid connection type %T", counterpartyConnection)
	}

	bz, err := k.cdc.Marshal(&connectionEnd)
	if err != nil {
		return err
	}

	if err := clientState.VerifyMembership(
		ctx, clientStore, k.cdc, height,
		0, 0, // skip delay period checks for non-packet processing verification
		proof, merklePath, bz,
	); err != nil {
		return errorsmod.Wrapf(err, "failed connection state verification for client (%s)", clientID)
	}

	return nil
}

// VerifyChannelState verifies a proof of the channel state of the specified
// channel end, under the specified port, stored on the target machine.
func (k Keeper) VerifyChannelState(
	ctx sdk.Context,
	connection exported.ConnectionI,
	height exported.Height,
	proof []byte,
	portID,
	channelID string,
	channel exported.ChannelI,
) error {
	clientID := connection.GetClientID()
	clientState, clientStore, err := k.getClientStateAndVerificationStore(ctx, clientID)
	if err != nil {
		return err
	}

	if status := k.clientKeeper.GetClientStatus(ctx, clientState, clientID); status != exported.Active {
		return errorsmod.Wrapf(clienttypes.ErrClientNotActive, "client (%s) status is %s", clientID, status)
	}

	merklePath := commitmenttypes.NewMerklePath(host.ChannelPath(portID, channelID))
	merklePath, err = commitmenttypes.ApplyPrefix(connection.GetCounterparty().GetPrefix(), merklePath)
	if err != nil {
		return err
	}

	channelEnd, ok := channel.(channeltypes.Channel)
	if !ok {
		return errorsmod.Wrapf(ibcerrors.ErrInvalidType, "invalid channel type %T", channel)
	}

	bz, err := k.cdc.Marshal(&channelEnd)
	if err != nil {
		return err
	}

	if err := clientState.VerifyMembership(
		ctx, clientStore, k.cdc, height,
		0, 0, // skip delay period checks for non-packet processing verification
		proof, merklePath, bz,
	); err != nil {
		return errorsmod.Wrapf(err, "failed channel state verification for client (%s)", clientID)
	}

	return nil
}

// VerifyPacketCommitment verifies a proof of an outgoing packet commitment at
// the specified port, specified channel, and specified sequence.
func (k Keeper) VerifyPacketCommitment(
	ctx sdk.Context,
	connection exported.ConnectionI,
	height exported.Height,
	proof []byte,
	portID,
	channelID string,
	sequence uint64,
	connectionHops []string,
	commitmentBytes []byte,
) error {
	clientID := connection.GetClientID()
	clientState, clientStore, err := k.getClientStateAndVerificationStore(ctx, clientID)
	if err != nil {
		return err
	}

	if status := k.clientKeeper.GetClientStatus(ctx, clientState, clientID); status != exported.Active {
		return errorsmod.Wrapf(clienttypes.ErrClientNotActive, "client (%s) status is %s", clientID, status)
	}

	var multihopProof commitmenttypes.MsgMultihopProofs
	if len(connectionHops) > 1 {
		if err := k.cdc.Unmarshal(proof, &multihopProof); err != nil {
			return errorsmod.Wrap(err, "failed to unmarshal multi-hop proof")
		}
	} else {
		multihopProof.KeyProof = &commitmenttypes.MultihopProof{}
		multihopProof.KeyProof.Proof = proof
	}

	// get maxmimum delay for any hop in the channel path
	delayTimePeriod, err := k.GetMaximumDelayPeriod(connection, &multihopProof)
	if err != nil {
		return err
	}
	// FIXME: this may not be accurate for multi-hop channels
	expectedTimePerBlock := k.GetParams(ctx).MaxExpectedTimePerBlock

	// get the last hop connection on the other side of the multihop channel
	// the last hop connection is the connection end on the chain before the counterparty multihop chain
	connectionEnd, err := k.GetLastHopConnectionEnd(connection, &multihopProof)
	if err != nil {
		return errorsmod.Wrap(err, "failed to get connection end")
	}

	prefix := connectionEnd.GetCounterparty().GetPrefix()
	merklePath := commitmenttypes.NewMerklePath(host.PacketCommitmentPath(portID, channelID, sequence))
	merklePath, err = commitmenttypes.ApplyPrefix(prefix, merklePath)
	if err != nil {
		return err
	}

	if err := clientState.VerifyMultihopMembership(
		ctx,
		clientStore,
		k.cdc,
		height,
		delayTimePeriod,
		expectedTimePerBlock,
		proof, // TODO: pass in multihop proof, need to fix import cycle
		connectionHops,
		merklePath,
		commitmentBytes,
	); err != nil {
		return errorsmod.Wrapf(err, "failed packet commitment verification for client (%s)", clientID)
	}

	return nil
}

// VerifyPacketAcknowledgement verifies a proof of an incoming packet
// acknowledgement at the specified port, specified channel, and specified sequence.
func (k Keeper) VerifyPacketAcknowledgement(
	ctx sdk.Context,
	connection exported.ConnectionI,
	height exported.Height,
	proof []byte,
	portID,
	channelID string,
	sequence uint64,
	acknowledgement []byte,
) error {
	clientID := connection.GetClientID()
	clientState, clientStore, err := k.getClientStateAndVerificationStore(ctx, clientID)
	if err != nil {
		return err
	}

	if status := k.clientKeeper.GetClientStatus(ctx, clientState, clientID); status != exported.Active {
		return errorsmod.Wrapf(clienttypes.ErrClientNotActive, "client (%s) status is %s", clientID, status)
	}

	// get time and block delays
	timeDelay := connection.GetDelayPeriod()
	blockDelay := k.getBlockDelay(ctx, connection)

	merklePath := commitmenttypes.NewMerklePath(host.PacketAcknowledgementPath(portID, channelID, sequence))
	merklePath, err = commitmenttypes.ApplyPrefix(connection.GetCounterparty().GetPrefix(), merklePath)
	if err != nil {
		return err
	}

	if err := clientState.VerifyMembership(
		ctx, clientStore, k.cdc, height,
		timeDelay, blockDelay,
		proof, merklePath, channeltypes.CommitAcknowledgement(acknowledgement),
	); err != nil {
		return errorsmod.Wrapf(err, "failed packet acknowledgement verification for client (%s)", clientID)
	}

	return nil
}

// VerifyPacketReceiptAbsence verifies a proof of the absence of an
// incoming packet receipt at the specified port, specified channel, and
// specified sequence.
func (k Keeper) VerifyPacketReceiptAbsence(
	ctx sdk.Context,
	connection exported.ConnectionI,
	height exported.Height,
	proof []byte,
	portID,
	channelID string,
	sequence uint64,
) error {
	clientID := connection.GetClientID()
	clientState, clientStore, err := k.getClientStateAndVerificationStore(ctx, clientID)
	if err != nil {
		return err
	}

	if status := k.clientKeeper.GetClientStatus(ctx, clientState, clientID); status != exported.Active {
		return errorsmod.Wrapf(clienttypes.ErrClientNotActive, "client (%s) status is %s", clientID, status)
	}

	// get time and block delays
	timeDelay := connection.GetDelayPeriod()
	blockDelay := k.getBlockDelay(ctx, connection)

	merklePath := commitmenttypes.NewMerklePath(host.PacketReceiptPath(portID, channelID, sequence))
	merklePath, err = commitmenttypes.ApplyPrefix(connection.GetCounterparty().GetPrefix(), merklePath)
	if err != nil {
		return err
	}

	if err := clientState.VerifyNonMembership(
		ctx, clientStore, k.cdc, height,
		timeDelay, blockDelay,
		proof, merklePath,
	); err != nil {
		return errorsmod.Wrapf(err, "failed packet receipt absence verification for client (%s)", clientID)
	}

	return nil
}

// VerifyNextSequenceRecv verifies a proof of the next sequence number to be
// received of the specified channel at the specified port.
func (k Keeper) VerifyNextSequenceRecv(
	ctx sdk.Context,
	connection exported.ConnectionI,
	height exported.Height,
	proof []byte,
	portID,
	channelID string,
	nextSequenceRecv uint64,
) error {
	clientID := connection.GetClientID()
	clientState, clientStore, err := k.getClientStateAndVerificationStore(ctx, clientID)
	if err != nil {
		return err
	}

	if status := k.clientKeeper.GetClientStatus(ctx, clientState, clientID); status != exported.Active {
		return errorsmod.Wrapf(clienttypes.ErrClientNotActive, "client (%s) status is %s", clientID, status)
	}

	// get time and block delays
	timeDelay := connection.GetDelayPeriod()
	blockDelay := k.getBlockDelay(ctx, connection)

	merklePath := commitmenttypes.NewMerklePath(host.NextSequenceRecvPath(portID, channelID))
	merklePath, err = commitmenttypes.ApplyPrefix(connection.GetCounterparty().GetPrefix(), merklePath)
	if err != nil {
		return err
	}

	if err := clientState.VerifyMembership(
		ctx, clientStore, k.cdc, height,
		timeDelay, blockDelay,
		proof, merklePath, sdk.Uint64ToBigEndian(nextSequenceRecv),
	); err != nil {
		return errorsmod.Wrapf(err, "failed next sequence receive verification for client (%s)", clientID)
	}

	return nil
}

// VerifyMultihopMembership verifies a multi-hop membership proof.
func (k Keeper) VerifyMultihopMembership(
	ctx sdk.Context,
	connection exported.ConnectionI,
	height exported.Height,
	proof []byte,
	connectionHops []string,
	kvGenerator channeltypes.KeyValueGenFunc,
) error {
	var multihopProof commitmenttypes.MsgMultihopProofs
	if err := k.cdc.Unmarshal(proof, &multihopProof); err != nil {
		return err
	}

	// get the last hop connection on the other side of the multihop channel
	// the last hop connection is the connection end on the chain before the counterparty multihop chain
	lastHopConnectionEnd, err := k.GetLastHopConnectionEnd(connection, &multihopProof)
	if err != nil {
		return err
	}

	// generate the key/value for the expected value on the counterparty end that needs to be proven
	key, value, err := kvGenerator(&multihopProof, lastHopConnectionEnd)
	if err != nil {
		return errorsmod.Wrap(err, "failed to generate key and value to prove")
	}

	// the counterparty of the last hop connection end is the connection end on the other end of the multihop channel
	prefix := lastHopConnectionEnd.GetCounterparty().GetPrefix()
	merklePath, err := commitmenttypes.ApplyPrefix(prefix, commitmenttypes.NewMerklePath(key))
	if err != nil {
		return err
	}

	clientID := connection.GetClientID()
	clientStore := k.clientKeeper.ClientStore(ctx, clientID)

	clientState, found := k.clientKeeper.GetClientState(ctx, clientID)
	if !found {
		return errorsmod.Wrap(clienttypes.ErrClientNotFound, clientID)
	}

	// check client associated with connection on this end of the multihop channel is active
	if status := clientState.Status(ctx, clientStore, k.cdc); status != exported.Active {
		return errorsmod.Wrapf(clienttypes.ErrClientNotActive, "client (%s) status is %s", clientID, status)
	}

	if clientState.GetLatestHeight().LT(height) {
		return errorsmod.Wrapf(
			clienttypes.ErrInvalidHeight,
			"client state height < proof height (%d < %d), please ensure the client has been updated",
			clientState.GetLatestHeight(), height,
		)
	}

	delayPeriod, err := k.GetMaximumDelayPeriod(connection, &multihopProof)
	if err != nil {
		return err
	}

	expectedTimePerBlock := k.GetParams(ctx).MaxExpectedTimePerBlock

	// ensure the delayPeriod passed
	if err := mh.VerifyDelayPeriodPassed(ctx, clientStore, height, delayPeriod, expectedTimePerBlock); err != nil {
		return err
	}

	consensusState, found := k.clientKeeper.GetClientConsensusState(ctx, clientID, height)
	if !found {
		return errorsmod.Wrapf(clienttypes.ErrConsensusStateNotFound, "consensus state not found for client ID (%s) at height (%s)", clientID, height)
	}

	return mh.VerifyMultihopMembership(k.cdc, consensusState, connectionHops, &multihopProof, merklePath, value)
}

// VerifyMultihopNonMembership verifies a multi-hop non-membership proof.
func (k Keeper) VerifyMultihopNonMembership(
	ctx sdk.Context,
	connection exported.ConnectionI,
	height exported.Height,
	proof []byte,
	connectionHops []string,
	kvGenerator channeltypes.KeyGenFunc,
) error {
	var multihopProof commitmenttypes.MsgMultihopProofs
	if err := k.cdc.Unmarshal(proof, &multihopProof); err != nil {
		return err
	}

	lastHopConnectionEnd, err := k.GetLastHopConnectionEnd(connection, &multihopProof)
	if err != nil {
		return err
	}

	// generate the key on the counterparty end that needs to be proven
	key, err := kvGenerator(&multihopProof, lastHopConnectionEnd)
	if err != nil {
		return errorsmod.Wrap(err, "failed to generate key")
	}

	// the counterparty of the last hop connection end is the connection end on the other end of the multihop channel
	prefix := lastHopConnectionEnd.GetCounterparty().GetPrefix()
	path, err := commitmenttypes.ApplyPrefix(prefix, commitmenttypes.NewMerklePath(key))
	if err != nil {
		return err
	}

	clientID := connection.GetClientID()
	clientStore := k.clientKeeper.ClientStore(ctx, clientID)

	clientState, found := k.clientKeeper.GetClientState(ctx, clientID)
	if !found {
		return errorsmod.Wrap(clienttypes.ErrClientNotFound, clientID)
	}

	// check client associated with connection on this end of the multihop channel is active
	if status := clientState.Status(ctx, clientStore, k.cdc); status != exported.Active {
		return errorsmod.Wrapf(clienttypes.ErrClientNotActive, "client (%s) status is %s", clientID, status)
	}

	if clientState.GetLatestHeight().LT(height) {
		return errorsmod.Wrapf(
			clienttypes.ErrInvalidHeight,
			"client state height < proof height (%d < %d), please ensure the client has been updated",
			clientState.GetLatestHeight(), height,
		)
	}

	delayPeriod, err := k.GetMaximumDelayPeriod(connection, &multihopProof)
	if err != nil {
		return err
	}

	expectedTimePerBlock := k.GetParams(ctx).MaxExpectedTimePerBlock

	// ensure the delayPeriod passed
	if err := mh.VerifyDelayPeriodPassed(ctx, clientStore, height, delayPeriod, expectedTimePerBlock); err != nil {
		return err
	}

	consensusState, found := k.clientKeeper.GetClientConsensusState(ctx, clientID, height)
	if !found {
		return errorsmod.Wrapf(clienttypes.ErrConsensusStateNotFound, "consensus state not found for client ID (%s) at height (%s)", clientID, height)
	}

	return mh.VerifyMultihopNonMembership(k.cdc, consensusState, connectionHops, &multihopProof, path)
}

// getBlockDelay calculates the block delay period from the time delay of the connection
// and the maximum expected time per block.
func (k Keeper) getBlockDelay(ctx sdk.Context, connection exported.ConnectionI) uint64 {
	// expectedTimePerBlock should never be zero, however if it is then return a 0 blcok delay for safety
	// as the expectedTimePerBlock parameter was not set.
	expectedTimePerBlock := k.GetParams(ctx).MaxExpectedTimePerBlock
	if expectedTimePerBlock == 0 {
		return 0
	}
	// calculate minimum block delay by dividing time delay period
	// by the expected time per block. Round up the block delay.
	timeDelay := connection.GetDelayPeriod()
	return uint64(math.Ceil(float64(timeDelay) / float64(expectedTimePerBlock)))
}

// getClientStateAndVerificationStore returns the client state and associated KVStore for the provided client identifier.
// If the client type is localhost then the core IBC KVStore is returned, otherwise the client prefixed store is returned.
func (k Keeper) getClientStateAndVerificationStore(ctx sdk.Context, clientID string) (exported.ClientState, storetypes.KVStore, error) {
	clientState, found := k.clientKeeper.GetClientState(ctx, clientID)
	if !found {
		return nil, nil, errorsmod.Wrap(clienttypes.ErrClientNotFound, clientID)
	}

	store := k.clientKeeper.ClientStore(ctx, clientID)
	if clientID == exported.LocalhostClientID {
		store = ctx.KVStore(k.storeKey)
	}

	return clientState, store, nil
}

// GetLastHopConnectionEnd returns the last hop connectionEnd from the perpespective of the executing chain.
// The last hop connection is the connection end on the chain before the counterparty multihop chain (i.e. the chain on the other end of the multihop channel).
func (k Keeper) GetLastHopConnectionEnd(c exported.ConnectionI, m *commitmenttypes.MsgMultihopProofs) (exported.ConnectionI, error) {
	var connectionEnd connectiontypes.ConnectionEnd
	if len(m.ConnectionProofs) > 0 {
		// connection proofs are ordered from executing chain to counterparty,
		// so the last hop connection end is the value of last connection proof
		if err := k.cdc.Unmarshal(m.ConnectionProofs[len(m.ConnectionProofs)-1].Value, &connectionEnd); err != nil {
			return nil, err
		}
		return &connectionEnd, nil
	}

	return c, nil
}

// GetMaximumDelayPeriod returns the maximum delay period over all connections in the multi-hop channel path.
func (k Keeper) GetMaximumDelayPeriod(c exported.ConnectionI, m *commitmenttypes.MsgMultihopProofs) (uint64, error) {
	delayPeriod := c.GetDelayPeriod()
	for _, connData := range m.ConnectionProofs {
		var connectionEnd connectiontypes.ConnectionEnd
		if err := k.cdc.Unmarshal(connData.Value, &connectionEnd); err != nil {
			return 0, err
		}
		if connectionEnd.DelayPeriod > delayPeriod {
			delayPeriod = connectionEnd.DelayPeriod
		}
	}
	return delayPeriod, nil
}

// GetCounterpartyConnectionHops returns the counterparty connectionHops.
// connection is the connection end on one of the ends of the multihop channel, and this function returns
// the connection hops for the connection on the other end of the multihop channel.
// Since connection proofs are ordered from the perspective of the connection parameter, in order to get the
// counterparty connection hops we need to reverse iterate through the proofs and then add the final
// counterparty connection ID for connection.
func (k Keeper) GetCounterpartyConnectionHops(c exported.ConnectionI, m *commitmenttypes.MsgMultihopProofs) (counterpartyHops []string, err error) {
	var connectionEnd connectiontypes.ConnectionEnd
	for _, connectionProof := range m.ConnectionProofs {
		if err = k.cdc.Unmarshal(connectionProof.Value, &connectionEnd); err != nil {
			return nil, err
		}
		counterpartyHops = append([]string{connectionEnd.GetCounterparty().GetConnectionID()}, counterpartyHops...)
	}

	// the last hop is the counterparty connection ID of the connection on the other end of the multihop channel
	counterpartyHops = append(counterpartyHops, c.GetCounterparty().GetConnectionID())

	return counterpartyHops, nil
}

// GeLastHopConsensusState returns the last hop consensusState from the perspective of the executing chain.
// The last hop connection is the connection end on the chain before the counterparty multihop chain (i.e. the chain on the other end of the multihop channel).
func (k Keeper) GetLastHopConsensusState(m *commitmenttypes.MsgMultihopProofs) (exported.ConsensusState, error) {
	var consensusState exported.ConsensusState
	if len(m.ConnectionProofs) > 0 {
		if err := k.cdc.UnmarshalInterface(m.ConsensusProofs[len(m.ConsensusProofs)-1].Value, &consensusState); err != nil {
			return nil, err
		}
	} else {
		panic("") // TODO: is this reachable?
	}
	return consensusState, nil
}
