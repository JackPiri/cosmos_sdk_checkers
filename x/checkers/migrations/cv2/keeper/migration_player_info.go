package keeper

import (
	"context"

	"github.com/alice/checkers/x/checkers/keeper"
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
)

func getOrNewPlayerInfoInMap(infoSoFar *map[string]*types.PlayerInfo, playerIndex string) (playerInfo *types.PlayerInfo) {
	playerInfo, found := (*infoSoFar)[playerIndex]
	if !found {
		playerInfo = &types.PlayerInfo{
			Account:    playerIndex,
			TotalGames: 0,
		}
		(*infoSoFar)[playerIndex] = playerInfo
	}
	return playerInfo
}

type storedGamesChunk struct {
	StoredGames []types.StoredGame
	Error       error
}

func loadStoredGames(context context.Context,
	k keeper.Keeper,
	gamesChannel chan<- storedGamesChunk,
	chunk uint64) {
	defer func() { close(gamesChannel) }()
	response, err := k.StoredGameAll(context, &types.QueryAllStoredGameRequest{
		Pagination: &query.PageRequest{Limit: chunk},
	})
	if err != nil {
		gamesChannel <- storedGamesChunk{Error: err}
		return
	}
	gamesChannel <- storedGamesChunk{StoredGames: response.StoredGame}
	for response.Pagination.NextKey != nil {
		response, err = k.StoredGameAll(context, &types.QueryAllStoredGameRequest{
			Pagination: &query.PageRequest{
				Key:   response.Pagination.NextKey,
				Limit: chunk,
			},
		})
		if err != nil {
			gamesChannel <- storedGamesChunk{Error: err}
			return
		}
		gamesChannel <- storedGamesChunk{StoredGames: response.StoredGame}
	}
}

type playerInfoTuple struct {
	PlayerInfo types.PlayerInfo
	Error      error
}

func handleStoredGameChannel(k keeper.Keeper,
	gamesChannel <-chan storedGamesChunk,
	playerInfoChannel chan<- playerInfoTuple) {
	defer func() { close(playerInfoChannel) }()
	for games := range gamesChannel {
		if games.Error != nil {
			playerInfoChannel <- playerInfoTuple{Error: games.Error}
			return
		}
		playerInfos := make(map[string]*types.PlayerInfo, len(games.StoredGames))
		for _, game := range games.StoredGames {
			getOrNewPlayerInfoInMap(&playerInfos, game.Black).TotalGames++
			getOrNewPlayerInfoInMap(&playerInfos, game.Red).TotalGames++
		}
		for _, playerInfo := range playerInfos {
			if playerInfo != nil {
				playerInfoChannel <- playerInfoTuple{PlayerInfo: *playerInfo}
			}
		}
	}
}

func handlePlayerInfoChannel(ctx sdk.Context, k keeper.Keeper,
	playerInfoChannel <-chan playerInfoTuple,
	done chan<- error) {
	defer func() { close(done) }()
	for receivedInfo := range playerInfoChannel {
		if receivedInfo.Error != nil {
			done <- receivedInfo.Error
			return
		}
		existingInfo, found := k.GetPlayerInfo(ctx, receivedInfo.PlayerInfo.Account)
		if found {
			existingInfo.TotalGames += receivedInfo.PlayerInfo.TotalGames
		} else {
			existingInfo = receivedInfo.PlayerInfo
		}
		k.SetPlayerInfo(ctx, existingInfo)
	}
	done <- nil
}

func MapStoredGamesReduceToPlayerInfo(ctx sdk.Context, k keeper.Keeper, chunk uint64) error {
	context := sdk.WrapSDKContext(ctx)
	gamesChannel := make(chan storedGamesChunk)
	playerInfoChannel := make(chan playerInfoTuple)
	done := make(chan error)

	go handlePlayerInfoChannel(ctx, k, playerInfoChannel, done)
	go handleStoredGameChannel(k, gamesChannel, playerInfoChannel)
	go loadStoredGames(context, k, gamesChannel, chunk)

	return <-done
}
