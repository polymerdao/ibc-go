package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// checkMisbehaviourHeader determines whether or not the currently registered
// public key signed over two different messages with the same sequence. If this is true
// the client state is updated to a frozen status.
// NOTE: Misbehaviour is not tracked for previous public keys, a solo machine may update to
// a new public key before the misbehaviour is processed. Therefore, misbehaviour is data
// order processing dependent.
func (cs ClientState) checkMisbehaviourHeader(
	ctx sdk.Context,
	cdc codec.BinaryCodec,
	clientStore sdk.KVStore,
	misbehaviour *Misbehaviour,
) (bool, error) {
	// NOTE: a check that the misbehaviour message data are not equal is done by
	// misbehaviour.ValidateBasic which is called by the 02-client keeper.

	// verify first signature
	if err := verifySignatureAndData(cdc, cs, misbehaviour, misbehaviour.SignatureOne); err != nil {
		return false, sdkerrors.Wrap(err, "failed to verify signature one")
	}

	// verify second signature
	if err := verifySignatureAndData(cdc, cs, misbehaviour, misbehaviour.SignatureTwo); err != nil {
		return false, sdkerrors.Wrap(err, "failed to verify signature two")
	}

	return true, nil
}

// verifySignatureAndData verifies that the currently registered public key has signed
// over the provided data and that the data is valid. The data is valid if it can be
// unmarshaled into the specified data type.
func verifySignatureAndData(cdc codec.BinaryCodec, clientState ClientState, misbehaviour *Misbehaviour, sigAndData *SignatureAndData) error {

	// do not check misbehaviour timestamp since we want to allow processing of past misbehaviour

	// ensure data can be unmarshaled to the specified data type
	if _, err := UnmarshalDataByType(cdc, sigAndData.DataType, sigAndData.Data); err != nil {
		return err
	}

	data, err := MisbehaviourSignBytes(
		cdc,
		misbehaviour.Sequence, sigAndData.Timestamp,
		clientState.ConsensusState.Diversifier,
		sigAndData.DataType,
		sigAndData.Data,
	)
	if err != nil {
		return err
	}

	sigData, err := UnmarshalSignatureData(cdc, sigAndData.Signature)
	if err != nil {
		return err
	}

	publicKey, err := clientState.ConsensusState.GetPubKey()
	if err != nil {
		return err
	}

	if err := VerifySignature(publicKey, data, sigData); err != nil {
		return err
	}

	return nil

}
