package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	abci "github.com/tendermint/tendermint/abci/types"
	//sdk "github.com/cosmos/cosmos-sdk/types"
	//stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	//icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
)

func main() {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	/*
		data, err := icatypes.SerializeCosmosTx(cdc, []sdk.Msg{
			stakingtypes.QueryUnbondingDelegationRequest{
				ValidatorAddr: "someaddr",
			},
		})
		if err != nil {
			log.Fatal(err)
		}

		icaPacket := icatypes.InterchainAccountPacketData{
			Data: data,
		}
		log.Println(data)
		log.Println(icaPacket)
	*/
	log.Println(cdc)

	req := []abci.RequestQuery{
		{
			Path:   fmt.Sprintf("store/%s/key", host.StoreKey),
			Height: 0,
			Data:   []byte("key1"),
			Prove:  true,
		},
		{
			Path:   fmt.Sprintf("store/%s/key", host.StoreKey),
			Height: 1,
			Data:   []byte("key2"),
			Prove:  true,
		},
	}
	b, _ := json.Marshal(req)

	s := string(b)
	log.Println(s)

	log.Println(req)
}
