package keeper

import (
	"context"
	"strconv"

	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/alice/checkers/x/checkers/rules"
)

func (k msgServer) CreateGame(goCtx context.Context, msg *types.MsgCreateGame) (*types.MsgCreateGameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// DONE: Handling the message (start)
	systemInfo, found := k.Keeper.GetSystemInfo(ctx)
	if !found {
		panic("SystemInfo not found")
	}
	newIndex := strconv.FormatUint(systemInfo.NextId, 10)

	var newGame *rules.Game = rules.New()
	var newStoredGame types.StoredGame = types.StoredGame{
		Index: newIndex,
		Board: rules.New().String(),
		Turn:  rules.PieceStrings[newGame.Turn],
		Black: msg.Black,
		Red:   msg.Red,
	}
	err := newStoredGame.Validate()
	if err != nil {
		return nil, err
	}
	k.Keeper.SetStoredGame(ctx, newStoredGame)

	playersName := []string{msg.Black, msg.Red}
	for _, playerName := range playersName {
		player, found := k.Keeper.GetPlayerInfo(ctx, playerName)
		if !found {
			player = types.PlayerInfo{
				Account:    playerName,
				TotalGames: 1,
			}
		} else {
			player.TotalGames += 1
		}
		k.Keeper.SetPlayerInfo(ctx, player)
	}

	systemInfo.NextId++
	k.Keeper.SetSystemInfo(ctx, systemInfo)
	// DONE: Handling the message (end)

	_ = ctx

	return &types.MsgCreateGameResponse{
		GameIndex: newIndex,
	}, nil
}
