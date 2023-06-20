package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
)

type ConnectionEnd = connectiontypes.ConnectionEnd

// TODO: Create multihop proof struct to serialize multihop proofs into.

// GetMultihopConnectionEnd returns the final connectionEnd from the counterparty perspective
func (m *MsgMultihopProofs) GetMultihopConnectionEnd(cdc codec.BinaryCodec, connection exported.ConnectionI) (*ConnectionEnd, error) {
	var connectionEnd ConnectionEnd
	if len(m.ConnectionProofs) > 0 {
		if err := cdc.Unmarshal(m.ConnectionProofs[len(m.ConnectionProofs)-1].Value, &connectionEnd); err != nil {
			return nil, err
		}
	} else {
		var ok bool
		connectionEnd, ok = connection.(connectiontypes.ConnectionEnd)
		if !ok {
			return nil, fmt.Errorf("failed to cast connection interface. expected type 'connectiontypes.ConnectionEnd'")
		}
	}
	return &connectionEnd, nil
}

// GetMultihopCounterpartyConsensus returns the final consensusState from the counterparty perspective (e.g. the source chain state).
func (m *MsgMultihopProofs) GetMultihopCounterpartyConsensus(cdc codec.BinaryCodec) (consensusState exported.ConsensusState, err error) {
	err = cdc.UnmarshalInterface(m.ConsensusProofs[len(m.ConsensusProofs)-1].Value, &consensusState)
	return
}

// GetMaximumDelayPeriod returns the maximum delay period over all connections in the multi-hop channel path.
func (m *MsgMultihopProofs) GetMaximumDelayPeriod(
	cdc codec.BinaryCodec,
	lastConnection exported.ConnectionI,
) (uint64, error) {
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

// GetCounterpartyHops returns the counter party connectionHops. Connection proofs are ordered from receiving chain to sending chain
// so in order to get the counterparty connection hops we need to reverse iterate through the proofs and then add the final counterparty
// connection id for the receiving chain.
func (m *MsgMultihopProofs) GetCounterpartyHops(
	cdc codec.BinaryCodec,
	lastConnection *ConnectionEnd,
) (counterpartyHops []string, err error) {
	var connectionEnd ConnectionEnd
	for _, connData := range m.ConnectionProofs {
		if err = cdc.Unmarshal(connData.Value, &connectionEnd); err != nil {
			return nil, err
		}
		counterpartyHops = append([]string{connectionEnd.GetCounterparty().GetConnectionID()}, counterpartyHops...)
	}

	counterpartyHops = append(counterpartyHops, lastConnection.GetCounterparty().GetConnectionID())

	return counterpartyHops, nil
}
