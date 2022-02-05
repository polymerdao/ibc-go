package keeper

import (
	"encoding/hex"

	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v3/modules/core/exported"
)

// CreateClient creates a new client state and populates it with a given consensus
// state as defined in https://github.com/cosmos/ibc/tree/master/spec/core/ics-002-client-semantics#create
func (k Keeper) CreateClient(
	ctx sdk.Context, clientState exported.ClientState, consensusState exported.ConsensusState,
) (string, error) {
	params := k.GetParams(ctx)
	if !params.IsAllowedClient(clientState.ClientType()) {
		return "", sdkerrors.Wrapf(
			types.ErrInvalidClientType,
			"client state type %s is not registered in the allowlist", clientState.ClientType(),
		)
	}

	clientID := k.GenerateClientIdentifier(ctx, clientState.ClientType())

	k.SetClientState(ctx, clientID, clientState)
	k.Logger(ctx).Info("client created at height", "client-id", clientID, "height", clientState.GetLatestHeight().String())

	// verifies initial consensus state against client state and initializes client store with any client-specific metadata
	// e.g. set ProcessedTime in Tendermint clients
	if err := clientState.Initialize(ctx, k.cdc, k.ClientStore(ctx, clientID), consensusState); err != nil {
		return "", err
	}

	k.Logger(ctx).Info("client created at height", "client-id", clientID, "height", clientState.GetLatestHeight().String())

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{"ibc", "client", "create"},
			1,
			[]metrics.Label{telemetry.NewLabel(types.LabelClientType, clientState.ClientType())},
		)
	}()

	return clientID, nil
}

// UpdateClient updates the consensus state and the state root from a provided header.
func (k Keeper) UpdateClient(ctx sdk.Context, clientID string, header exported.Header) error {
	clientState, found := k.GetClientState(ctx, clientID)
	if !found {
		return sdkerrors.Wrapf(types.ErrClientNotFound, "cannot update client with ID %s", clientID)
	}

	clientStore := k.ClientStore(ctx, clientID)

	if status := clientState.Status(ctx, clientStore, k.cdc); status != exported.Active {
		return sdkerrors.Wrapf(types.ErrClientNotActive, "cannot update client (%s) with status %s", clientID, status)
	}

	eventType := types.EventTypeUpdateClient

	if err := clientState.VerifyHeader(ctx, k.cdc, clientStore, header); err != nil {
		return sdkerrors.Wrapf(err, "header verification failed for client (%s)", clientID)
	}

	misbehaving, err := clientState.CheckForMisbehaviour(ctx, k.cdc, clientStore, header)
	if err != nil {
		return sdkerrors.Wrapf(err, "misbehaviour check failed for client (%s)", clientID)
	}
	if misbehaving {
		// Freeze client.
		clientState.UpdateStateOnMisbehaviour(ctx, k.cdc, clientStore)

		// set eventType to SubmitMisbehaviour
		eventType = types.EventTypeSubmitMisbehaviour
		k.Logger(ctx).Info("client frozen due to misbehaviour", "client-id", clientID)
		defer func() {
			telemetry.IncrCounterWithLabels(
				[]string{"ibc", "client", "misbehaviour"},
				1,
				[]metrics.Label{
					telemetry.NewLabel(types.LabelClientType, clientState.ClientType()),
					telemetry.NewLabel(types.LabelClientID, clientID),
					telemetry.NewLabel(types.LabelMsgType, "update"),
				},
			)
		}()
		return nil
	}

	// Any writes made in updates are persisted on valid updates.
	// Light client implementations are responsible for writing the correct metadata.
	if err := clientState.UpdateStateFromHeader(ctx, k.cdc, clientStore, header); err != nil {
		return sdkerrors.Wrapf(err, "cannot update client with ID %s", clientID)
	}

	// emit the full header in events
	var headerStr string
	if header != nil {
		// Marshal the Header as an Any and encode the resulting bytes to hex.
		// This prevents the event value from containing invalid UTF-8 characters
		// which may cause data to be lost when JSON encoding/decoding.
		headerStr = hex.EncodeToString(types.MustMarshalHeader(k.cdc, header))

	}

	// TODO: Do we want to use a self reported height from clientstate update?
	k.Logger(ctx).Info("client state updated", "client-id", clientID)

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{"ibc", "client", "update"},
			1,
			[]metrics.Label{
				telemetry.NewLabel(types.LabelClientType, clientState.ClientType()),
				telemetry.NewLabel(types.LabelClientID, clientID),
				telemetry.NewLabel(types.LabelUpdateType, "msg"),
			},
		)
	}()

	// emitting events in the keeper emits for both begin block and handler client updates
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			eventType,
			sdk.NewAttribute(types.AttributeKeyClientID, clientID),
			sdk.NewAttribute(types.AttributeKeyClientType, clientState.ClientType()),
			sdk.NewAttribute(types.AttributeKeyHeader, headerStr),
		),
	)
	return nil
}

// UpgradeClient upgrades the client to a new client state if this new client was committed to
// by the old client at the specified upgrade height
func (k Keeper) UpgradeClient(ctx sdk.Context, clientID string, upgradedClient exported.ClientState, upgradedConsState exported.ConsensusState,
	proofUpgradeClient, proofUpgradeConsState []byte) error {
	clientState, found := k.GetClientState(ctx, clientID)
	if !found {
		return sdkerrors.Wrapf(types.ErrClientNotFound, "cannot update client with ID %s", clientID)
	}

	clientStore := k.ClientStore(ctx, clientID)

	if status := clientState.Status(ctx, clientStore, k.cdc); status != exported.Active {
		return sdkerrors.Wrapf(types.ErrClientNotActive, "cannot upgrade client (%s) with status %s", clientID, status)
	}

	if err := clientState.VerifyUpgradeAndUpdateState(ctx, k.cdc, clientStore,
		upgradedClient, upgradedConsState, proofUpgradeClient, proofUpgradeConsState); err != nil {
		return sdkerrors.Wrapf(err, "cannot upgrade client with ID %s", clientID)
	}

	k.Logger(ctx).Info("client state upgraded", "client-id", clientID, "height", upgradedClient.GetLatestHeight().String())

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{"ibc", "client", "upgrade"},
			1,
			[]metrics.Label{
				telemetry.NewLabel(types.LabelClientType, upgradedClient.ClientType()),
				telemetry.NewLabel(types.LabelClientID, clientID),
			},
		)
	}()

	// emitting events in the keeper emits for client upgrades
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpgradeClient,
			sdk.NewAttribute(types.AttributeKeyClientID, clientID),
			sdk.NewAttribute(types.AttributeKeyClientType, upgradedClient.ClientType()),
			sdk.NewAttribute(types.AttributeKeyConsensusHeight, upgradedClient.GetLatestHeight().String()),
		),
	)

	return nil
}

// CheckMisbehaviourAndUpdateState checks for client misbehaviour and freezes the
// client if so.
func (k Keeper) CheckMisbehaviourAndUpdateState(ctx sdk.Context, clientID string, misbehaviour exported.Header) error {
	clientState, found := k.GetClientState(ctx, clientID)
	if !found {
		return sdkerrors.Wrapf(types.ErrClientNotFound, "cannot check misbehaviour for client with ID %s", clientID)
	}

	clientStore := k.ClientStore(ctx, clientID)

	if status := clientState.Status(ctx, clientStore, k.cdc); status != exported.Active {
		return sdkerrors.Wrapf(types.ErrClientNotActive, "cannot process misbehaviour for client (%s) with status %s", clientID, status)
	}

	if err := misbehaviour.ValidateBasic(); err != nil {
		return err
	}

	misbehaving, err := clientState.CheckForMisbehaviour(ctx, k.cdc, clientStore, misbehaviour)
	if err != nil {
		return sdkerrors.Wrapf(err, "misbehaviour check failed for client (%s)", clientID)
	}
	if !misbehaving {
		// no-op if not misbehaving.
		return nil
	}

	// Freeze client.
	clientState.UpdateStateOnMisbehaviour(ctx, k.cdc, clientStore)
	k.Logger(ctx).Info("client frozen due to misbehaviour", "client-id", clientID)

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{"ibc", "client", "misbehaviour"},
			1,
			[]metrics.Label{
				telemetry.NewLabel(types.LabelClientType, misbehaviour.ClientType()),
				telemetry.NewLabel(types.LabelClientID, clientID),
			},
		)
	}()

	return nil
}
