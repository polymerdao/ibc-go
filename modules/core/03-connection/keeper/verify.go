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

	// get time and block delays
	timeDelay := connection.GetDelayPeriod()
	blockDelay := k.getBlockDelay(ctx, connection)

	merklePath := commitmenttypes.NewMerklePath(host.PacketCommitmentPath(portID, channelID, sequence))
	merklePath, err = commitmenttypes.ApplyPrefix(connection.GetCounterparty().GetPrefix(), merklePath)
	if err != nil {
		return err
	}

	if err := clientState.VerifyMembership(
		ctx, clientStore, k.cdc, height,
		timeDelay, blockDelay,
		proof, merklePath, commitmentBytes,
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
	var multihopProof channeltypes.MsgMultihopProofs
	if err := k.cdc.Unmarshal(proof, &multihopProof); err != nil {
		return err
	}

	// get the last hop connection on the other side of the multihop channel
	// the last hop connection is the connection end on the chain before the counterparty multihop chain
	lastHopConnectionEnd, err := multihopProof.GetLastHopConnectionEnd(k.cdc, connection)
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

	delayPeriod, err := multihopProof.GetMaximumDelayPeriod(k.cdc, connection)
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

	return mh.VerifyMultihopMembership(k.cdc, consensusState, connectionHops, &multihopProof, path, value)
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
	var mProof channeltypes.MsgMultihopProofs
	if err := k.cdc.Unmarshal(proof, &mProof); err != nil {
		return err
	}

	lastHopConnectionEnd, err := mProof.GetLastHopConnectionEnd(k.cdc, connection)
	if err != nil {
		return err
	}

	// generate the key on the counterparty end that needs to be proven
	key, err := kvGenerator(&mProof, lastHopConnectionEnd)
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

	delayPeriod, err := mProof.GetMaximumDelayPeriod(k.cdc, connection)
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

	return mh.VerifyMultihopNonMembership(k.cdc, consensusState, connectionHops, &mProof, path)
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
