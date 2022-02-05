package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	"github.com/cosmos/ibc-go/v3/modules/core/exported"
)

// VerifyHeader checks if the provided header is valid and updates
// the consensus state if appropriate. It returns an error if:
// - the header provided is not parseable to a solo machine header
// - the header sequence does not match the current sequence
// - the header timestamp is less than the consensus state timestamp
// - the currently registered public key did not provide the update signature
func (cs ClientState) VerifyHeader(
	ctx sdk.Context, cdc codec.BinaryCodec, clientStore sdk.KVStore,
	header exported.Header,
) error {
	smHeader, ok := header.(*Header)
	if !ok {
		return sdkerrors.Wrapf(
			clienttypes.ErrInvalidHeader, "header type %T, expected  %T", header, &Header{},
		)
	}

	return checkHeader(cdc, &cs, smHeader)
}

// CheckForMisbehaviour will check for misbehaviour.
func (cs ClientState) CheckForMisbehaviour(
	ctx sdk.Context, cdc codec.BinaryCodec, clientStore sdk.KVStore,
	header exported.Header,
) (bool, error) {
	misbehaviour, ok := header.(*Misbehaviour)
	if ok {
		return cs.checkMisbehaviourHeader(ctx, cdc, clientStore, misbehaviour)
	}
	// No handler for *Header
	return false, nil
}

// UpdateStateOnMisbehaviour freezes client state on misbehaviour.
func (cs ClientState) UpdateStateOnMisbehaviour(
	_ sdk.Context, cdc codec.BinaryCodec, clientStore sdk.KVStore,
) {
	cs.IsFrozen = true
	clientStore.Set(host.ClientStateKey(), clienttypes.MustMarshalClientState(cdc, &cs))
}

// UpdateStateFromHeader updates the consensus state.
func (cs ClientState) UpdateStateFromHeader(
	ctx sdk.Context, cdc codec.BinaryCodec, clientStore sdk.KVStore,
	header exported.Header,
) error {
	tmHeader := header.(*Header)
	clientState, consensusState := update(&cs, tmHeader)
	clientStore.Set(host.ConsensusStateKey(tmHeader.GetHeight()), clienttypes.MustMarshalConsensusState(cdc, consensusState))
	clientStore.Set(host.ClientStateKey(), clienttypes.MustMarshalClientState(cdc, clientState))
	return nil
}

// checkHeader checks if the Solo Machine update signature is valid.
func checkHeader(cdc codec.BinaryCodec, clientState *ClientState, header *Header) error {
	// assert update sequence is current sequence
	if header.Sequence != clientState.Sequence {
		return sdkerrors.Wrapf(
			clienttypes.ErrInvalidHeader,
			"header sequence does not match the client state sequence (%d != %d)", header.Sequence, clientState.Sequence,
		)
	}

	// assert update timestamp is not less than current consensus state timestamp
	if header.Timestamp < clientState.ConsensusState.Timestamp {
		return sdkerrors.Wrapf(
			clienttypes.ErrInvalidHeader,
			"header timestamp is less than to the consensus state timestamp (%d < %d)", header.Timestamp, clientState.ConsensusState.Timestamp,
		)
	}

	// assert currently registered public key signed over the new public key with correct sequence
	data, err := HeaderSignBytes(cdc, header)
	if err != nil {
		return err
	}

	sigData, err := UnmarshalSignatureData(cdc, header.Signature)
	if err != nil {
		return err
	}

	publicKey, err := clientState.ConsensusState.GetPubKey()
	if err != nil {
		return err
	}

	if err := VerifySignature(publicKey, data, sigData); err != nil {
		return sdkerrors.Wrap(ErrInvalidHeader, err.Error())
	}

	return nil
}

// update the consensus state to the new public key and an incremented sequence
func update(clientState *ClientState, header *Header) (*ClientState, *ConsensusState) {
	consensusState := &ConsensusState{
		PublicKey:   header.NewPublicKey,
		Diversifier: header.NewDiversifier,
		Timestamp:   header.Timestamp,
	}

	// increment sequence number
	clientState.Sequence++
	clientState.ConsensusState = consensusState
	return clientState, consensusState
}
