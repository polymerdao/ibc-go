package types

import (
	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/codec"

	connectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	ibcerrors "github.com/cosmos/ibc-go/v8/modules/core/errors"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

type ConnectionEnd = connectiontypes.ConnectionEnd

// TODO: Create multihop proof struct to serialize multihop proofs into.

// GetLastHopConnectionEnd returns the last hop connectionEnd from the perpespective of the executing chain.
// The last hop connection is the connection end on the chain before the counterparty multihop chain (i.e. the chain on the other end of the multihop channel).
func (m *MsgMultihopProofs) GetLastHopConnectionEnd(cdc codec.BinaryCodec, connection exported.ConnectionI) (*ConnectionEnd, error) {
	var connectionEnd ConnectionEnd
	if len(m.ConnectionProofs) > 0 {
		// connection proofs are ordered from executing chain to counterparty,
		// so the last hop connection end is the value of last connection proof
		if err := cdc.Unmarshal(m.ConnectionProofs[len(m.ConnectionProofs)-1].Value, &connectionEnd); err != nil {
			return nil, err
		}
	} else {
		var ok bool
		connectionEnd, ok = connection.(connectiontypes.ConnectionEnd)
		if !ok {
			return nil, errorsmod.Wrapf(ibcerrors.ErrInvalidType, "expected %T, got %T", connectiontypes.ConnectionEnd{}, connection)
		}
	}
	return &connectionEnd, nil
}

// GeLastHopConsensusState returns the last hop consensusState from the perspective of the executing chain.
// The last hop connection is the connection end on the chain before the counterparty multihop chain (i.e. the chain on the other end of the multihop channel).
func (m *MsgMultihopProofs) GeLastHopConsensusState(cdc codec.BinaryCodec) (exported.ConsensusState, error) {
	var consensusState exported.ConsensusState
	if len(m.ConnectionProofs) > 0 {
		if err := cdc.UnmarshalInterface(m.ConsensusProofs[len(m.ConsensusProofs)-1].Value, &consensusState); err != nil {
			return nil, err
		}
	} else {
		panic("") // TODO: is this reachable?
	}
	return consensusState, nil
}

// GetMaximumDelayPeriod returns the maximum delay period over all connections in the multi-hop channel path.
func (m *MsgMultihopProofs) GetMaximumDelayPeriod(cdc codec.BinaryCodec, lastConnection exported.ConnectionI) (uint64, error) {
	delayPeriod := lastConnection.GetDelayPeriod()
	for _, connData := range m.ConnectionProofs {
		var connectionEnd ConnectionEnd
		if err := cdc.Unmarshal(connData.Value, &connectionEnd); err != nil {
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
func (m *MsgMultihopProofs) GetCounterpartyConnectionHops(cdc codec.BinaryCodec, connection *ConnectionEnd) (counterpartyHops []string, err error) {
	var connectionEnd ConnectionEnd
	for _, connectionProof := range m.ConnectionProofs {
		if err = cdc.Unmarshal(connectionProof.Value, &connectionEnd); err != nil {
			return nil, err
		}
		counterpartyHops = append([]string{connectionEnd.GetCounterparty().GetConnectionID()}, counterpartyHops...)
	}

	// the last hop is the counterparty connection ID of the connection on the other end of the multihop channel
	counterpartyHops = append(counterpartyHops, connection.GetCounterparty().GetConnectionID())

	return counterpartyHops, nil
}
