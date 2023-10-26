package keeper

import (
	"context"

	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/alice/checkers/x/checkers/rules"
)

func (k Keeper) PlayerInfoAllWithNotStartedGames(goCtx context.Context, req *types.QueryPlayerInfoAllWithNotStartedGamesRequest) (*types.QueryPlayerInfoAllWithNotStartedGamesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Process the query
	var notStartedGameBoard = rules.New().String()
	var storedGames []types.StoredGame = k.GetAllStoredGame((ctx))
	var notStartedStoredGames []types.StoredGame
	for _, storedGame := range storedGames {
		if storedGame.Board == notStartedGameBoard {
			notStartedStoredGames = append(notStartedStoredGames, storedGame)
		}
	}

	var playersInfo []types.PlayerInfo = k.GetAllPlayerInfo(ctx)
	var filteredPlayersInfo []types.PlayerInfo
	for _, playerInfo := range playersInfo {
		for _, notStartedStoredGame := range notStartedStoredGames {
			if playerInfo.Account == notStartedStoredGame.Black || playerInfo.Account == notStartedStoredGame.Red {
				filteredPlayersInfo = append(filteredPlayersInfo, playerInfo)
				break
			}
		}
	}

	return &types.QueryPlayerInfoAllWithNotStartedGamesResponse{PlayerInfoAllWithNotStartedGames: filteredPlayersInfo}, nil
}
