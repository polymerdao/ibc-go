package multihop

import (
	"fmt"
	"math"
	"strings"

	errorsmod "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	connectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/v8/modules/core/23-commitment/types"
	ibcerrors "github.com/cosmos/ibc-go/v8/modules/core/errors"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	tmclient "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
)

// parseID is a helper function that parses a client or connection id from a key path
func parseID(keyPath []string) (string, error) {
	if len(keyPath) < 2 {
		return "", fmt.Errorf(
			"invalid consensus proof key path length: %d",
			len(keyPath),
		)
	}
	parts := strings.Split(keyPath[1], "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid consensus proof key path: %s", keyPath)
	}
	return parts[1], nil
}

// parseClientIDFromKey parses the client id from a consensus state key path and returns it.
func parseClientIDFromKey(keyPath []string) (string, error) {
	return parseID(keyPath)
}

// parseConnectionIDFromKey parses the connectionID from the connection proof key and returns it.
func parseConnectionIDFromKey(keyPath []string) (string, error) {
	return parseID(keyPath)
}

// VerifyDelayPeriodPassed will ensure that at least delayTimePeriod amount of time and delayBlockPeriod number of blocks have passed
// since consensus state was submitted before allowing verification to continue.
func VerifyDelayPeriodPassed(
	ctx sdk.Context,
	store storetypes.KVStore,
	proofHeight exported.Height,
	timeDelay uint64,
	expectedTimePerBlock uint64,
) error {
	// get time and block delays
	blockDelay := getBlockDelay(timeDelay, expectedTimePerBlock)
	return tmclient.VerifyDelayPeriodPassed(ctx, store, proofHeight, timeDelay, blockDelay)
}

// getBlockDelay calculates the block delay period from the time delay of the connection
// and the maximum expected time per block.
func getBlockDelay(timeDelay uint64, expectedTimePerBlock uint64) uint64 {
	// expectedTimePerBlock should never be zero, however if it is then return a 0 block delay for safety
	// as the expectedTimePerBlock parameter was not set.
	if expectedTimePerBlock == 0 {
		return 0
	}
	return uint64(math.Ceil(float64(timeDelay) / float64(expectedTimePerBlock)))
}

// VerifyMultihopMembership verifies a multihop proof. A nil value indicates a non-inclusion proof (proof of absence).
func VerifyMultihopMembership(
	cdc codec.BinaryCodec,
	consensusState exported.ConsensusState,
	connectionHops []string,
	proofs *channeltypes.MsgMultihopProofs,
	path exported.Path,
	value []byte,
) error {
	if len(proofs.ConsensusProofs) != len(proofs.ConnectionProofs) {
		return fmt.Errorf("the number of connection (%d) and consensus (%d) proofs must be equal", len(proofs.ConnectionProofs), len(proofs.ConsensusProofs))
	}

	// verify connection states and ordering
	if err := validateIntermediateStates(cdc, proofs.ConsensusProofs, proofs.ConnectionProofs, connectionHops[1:]); err != nil {
		return err
	}

	// verify intermediate consensus and connection states from destination --> source
	if err := verifyIntermediateStates(cdc, consensusState, proofs.ConsensusProofs, proofs.ConnectionProofs); err != nil {
		return err
	}

	// verify the key proof on counterparty chain's consensus state.
	return verifyKeyValueMembership(cdc, consensusState, proofs, path, value)
}

// VerifyMultihopNonMembership verifies a multihop proof.
func VerifyMultihopNonMembership(
	cdc codec.BinaryCodec,
	consensusState exported.ConsensusState,
	connectionHops []string,
	proofs *channeltypes.MsgMultihopProofs,
	path exported.Path,
) error {
	if len(proofs.ConsensusProofs) != len(proofs.ConnectionProofs) {
		return fmt.Errorf("the number of connection (%d) and consensus (%d) proofs must be equal", len(proofs.ConnectionProofs), len(proofs.ConsensusProofs))
	}

	// validate connection states and ordering
	if err := validateIntermediateStates(cdc, proofs.ConsensusProofs, proofs.ConnectionProofs, connectionHops[1:]); err != nil {
		return err
	}

	// verify intermediate consensus and connection states from destination --> source
	if err := verifyIntermediateStates(cdc, consensusState, proofs.ConsensusProofs, proofs.ConnectionProofs); err != nil {
		return err
	}

	// verify the keyproof on source chain's consensus state.
	return verifyKeyNonMembership(cdc, consensusState, proofs, path)
}

// validateIntermediateStates validates the following:
// 1. That the connections in the values of connectionProofs match the connectionHops of the channel and are in OPEN state.
// 2. That the client IDs in the keys of consensusProofs matches the client ID in the connection end on the same hop.
// - consensusProofs is a list of proofs, with the corresponding key/value pair, for the consensus statesfrom the executing chain
// (but excluding it) to the counterparty.
// - connectionProofs is a list of proofs, with the corresponding key/value pair, for the connection ends that represent the connection hops
// from the executing chain (but excluding it) to the counterparty.
// - connectionHops contains the connection IDs for the hops from the executing chain to the counterparty, but it does not include the
// connection ID for the first hop on the executing chain.
func validateIntermediateStates(
	cdc codec.BinaryCodec,
	consensusProofs []*channeltypes.MultihopProof,
	connectionProofs []*channeltypes.MultihopProof,
	connectionHops []string,
) error {
	// check that the are as many proofs as connection hops
	if len(connectionProofs) != len(connectionHops) {
		return errorsmod.Wrapf(connectiontypes.ErrInvalidLengthConnection,
			"connectionHops length (%d) must match the connectionProofData length (%d)",
			len(connectionHops), len(connectionProofs))
	}

	// check all connections are in OPEN state and that the connection IDs match and are in the right order.
	for i := 0; i < len(consensusProofs); i++ {
		consensusProof := consensusProofs[i]
		connectionProof := connectionProofs[i]

		var connectionEnd connectiontypes.ConnectionEnd
		if err := cdc.Unmarshal(connectionProof.Value, &connectionEnd); err != nil {
			return errorsmod.Wrapf(err, "failed to unmarshal connection end for hop %d", i)
		}

		// parse connection ID from the connectionEnd key path
		connectionID, err := parseConnectionIDFromKey(connectionProof.PrefixedKey.KeyPath)
		if err != nil {
			return err
		}

		// 1. check the connectionHop value matches the connection ID from the proof and are in the same order.
		if connectionID != connectionHops[i] {
			return errorsmod.Wrapf(
				connectiontypes.ErrConnectionPath,
				"connection hop (%s) does not match connection proof hop (%s) for hop %d",
				connectionHops[i], connectionID, i)
		}

		// 2. check that the connectionEnd is in the OPEN state.
		if connectionEnd.GetState() != int32(connectiontypes.OPEN) {
			return errorsmod.Wrapf(
				connectiontypes.ErrInvalidConnectionState,
				"connection state is not OPEN (got %s) for connection ID (%s) for hop %d",
				connectiontypes.State(connectionEnd.GetState()).String(),
				connectionEnd.Counterparty.ConnectionId, i)
		}

		// parse the client ID from the the consensusState key path
		clientID, err := parseClientIDFromKey(consensusProof.PrefixedKey.KeyPath)
		if err != nil {
			return err
		}

		// check that the client ID matches the client ID in the connectionEnd on the same hop.
		if connectionEnd.ClientId != clientID {
			return fmt.Errorf("consensus state client id (%s) does not match expected client id (%s)",
				connectionEnd.ClientId,
				clientID,
			)
		}
	}
	return nil
}

// verifyIntermediateStates verifies the intermediate consensus and connection states in the multi-hop proof.
// - consensusState
// - consensusProofs is a list of proofs, with the corresponding key/value pair, for the consensus statesfrom the executing chain
// (but excluding it) to the counterparty.
// - connectionProofs is a list of proofs, with the corresponding key/value pair, for the connection ends that represent the connection hops
// from the executing chain (but excluding it) to the counterparty.
func verifyIntermediateStates(
	cdc codec.BinaryCodec,
	consensusState exported.ConsensusState,
	consensusProofs []*channeltypes.MultihopProof,
	connectionProofs []*channeltypes.MultihopProof,
) error {
	for i := 0; i < len(consensusProofs); i++ {
		consensusProof := consensusProofs[i]
		connectionProof := connectionProofs[i]

		cs, ok := consensusState.(*tmclient.ConsensusState)
		if !ok {
			return errorsmod.Wrapf(ibcerrors.ErrInvalidType, "expected %T, got: %T", tmclient.ConsensusState{}, consensusState)
		}

		// prove consensus state
		var proof commitmenttypes.MerkleProof
		if err := cdc.Unmarshal(consensusProof.Proof, &proof); err != nil {
			return errorsmod.Wrapf(err, "failed to unmarshal consensus state proof")
		}
		if err := proof.VerifyMembership(
			commitmenttypes.GetSDKSpecs(),
			cs.GetRoot(),
			*consensusProof.PrefixedKey,
			consensusProof.Value,
		); err != nil {
			return errorsmod.Wrapf(err, "failed to verify consensus state proof for hop %d", i)
		}

		// prove connection state
		proof.Reset()
		if err := cdc.Unmarshal(connectionProof.Proof, &proof); err != nil {
			return errorsmod.Wrapf(err, "failed to unmarshal connection state proof for hop %d", i)
		}
		if err := proof.VerifyMembership(
			commitmenttypes.GetSDKSpecs(),
			cs.GetRoot(),
			*connectionProof.PrefixedKey,
			connectionProof.Value,
		); err != nil {
			return errorsmod.Wrapf(err, "failed to verify connection proof for hop %d", i)
		}

		// set the next consensus state
		if err := cdc.UnmarshalInterface(consensusProof.Value, &consensusState); err != nil {
			return errorsmod.Wrapf(err, "failed to unmarshal consensus state for hop %d", i)
		}
	}
	return nil
}

// verifyKeyValueMembership verifies a multihop membership proof including all intermediate state proofs.
func verifyKeyValueMembership(
	cdc codec.BinaryCodec,
	consensusStateI exported.ConsensusState,
	proofs *channeltypes.MsgMultihopProofs,
	path exported.Path,
	value []byte,
) error {
	// no keyproof provided, nothing to verify
	if proofs.KeyProof == nil {
		return nil
	}

	var merkleProof commitmenttypes.MerkleProof
	if err := cdc.Unmarshal(proofs.KeyProof.Proof, &merkleProof); err != nil {
		return errorsmod.Wrap(commitmenttypes.ErrInvalidProof, "failed to unmarshal proof into ICS 23 commitment merkle proof")
	}

	if len(proofs.ConsensusProofs) > 0 {
		index := uint32(len(proofs.ConsensusProofs)) - 1
		if err := cdc.UnmarshalInterface(proofs.ConsensusProofs[index].Value, &consensusStateI); err != nil {
			return errorsmod.Wrap(err, "failed to unmarshall consensus state")
		}
	}
	consensusState, ok := consensusStateI.(*tmclient.ConsensusState)
	if !ok {
		return errorsmod.Wrapf(ibcerrors.ErrInvalidType, "expected %T, got: %T", tmclient.ConsensusState{}, consensusState)
	}

	return merkleProof.VerifyMembership(commitmenttypes.GetSDKSpecs(), consensusState.GetRoot(), path, value)
}

// verifyKeyNonMembership verifies a multihop non-membership proof including all intermediate state proofs.
func verifyKeyNonMembership(
	cdc codec.BinaryCodec,
	consensusStateI exported.ConsensusState,
	proofs *channeltypes.MsgMultihopProofs,
	path exported.Path,
) error {
	// no keyproof provided, nothing to verify
	if proofs.KeyProof == nil {
		return nil
	}

	var keyProof commitmenttypes.MerkleProof
	if err := cdc.Unmarshal(proofs.KeyProof.Proof, &keyProof); err != nil {
		return fmt.Errorf("failed to unmarshal key proof: %w", err)
	}

	if len(proofs.ConsensusProofs) > 0 {
		index := uint32(len(proofs.ConsensusProofs)) - 1
		if err := cdc.UnmarshalInterface(proofs.ConsensusProofs[index].Value, &consensusStateI); err != nil {
			return fmt.Errorf("failed to unpack consensus state: %w", err)
		}
	}
	consensusState, ok := consensusStateI.(*tmclient.ConsensusState)
	if !ok {
		return fmt.Errorf("expected consensus state to be tendermint consensus state, got: %T", consensusStateI)
	}

	return keyProof.VerifyNonMembership(commitmenttypes.GetSDKSpecs(), consensusState.GetRoot(), path)
}
