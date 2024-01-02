package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"

	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	mh "github.com/cosmos/ibc-go/v8/modules/core/33-multihop"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

// VerifyMultihopMembership verifies a multi-hop membership proof.
func (k Keeper) VerifyMultihopMembership(
	ctx sdk.Context,
	connection exported.ConnectionI,
	height exported.Height,
	proof []byte,
	connectionHops []string,
	path exported.Path,
	value []byte,
) error {
	var multihopProof channeltypes.MsgMultihopProofs
	if err := k.cdc.Unmarshal(proof, &multihopProof); err != nil {
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
	path exported.Path,
) error {
	var mProof channeltypes.MsgMultihopProofs
	if err := k.cdc.Unmarshal(proof, &mProof); err != nil {
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
