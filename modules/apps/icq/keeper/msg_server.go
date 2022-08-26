package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/ibc-go/v3/modules/apps/icq/types"
)

var _ types.MsgServer = Keeper{}

// Query defines a rpc handler method for MsgQuery.
func (k Keeper) Query(goCtx context.Context, msg *types.MsgQuery) (*types.MsgQueryResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.SendQuery(
		ctx, msg.SourcePort, msg.SourceChannel, msg.Requests, msg.TimeoutHeight, msg.TimeoutTimestamp,
	); err != nil {
		return nil, err
	}

	k.Logger(ctx).Info("IBC interchain query", "num_requests", len(msg.Requests))

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeQuery,
			sdk.NewAttribute(types.AttributeKeyNumRequests, len(msg.Requests)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	})

	return &types.MsgQueryResponse{}, nil
}
