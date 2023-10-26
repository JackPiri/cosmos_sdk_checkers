package keeper

import (
	"context"

	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) PlayerInfoAllWithNotStartedGames(goCtx context.Context, req *types.QueryPlayerInfoAllWithNotStartedGamesRequest) (*types.QueryPlayerInfoAllWithNotStartedGamesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Process the query
	_ = ctx

	return &types.QueryPlayerInfoAllWithNotStartedGamesResponse{}, nil
}
