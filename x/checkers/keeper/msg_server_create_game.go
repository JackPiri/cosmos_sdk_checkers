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

	systemInfo.NextId++
	k.Keeper.SetSystemInfo(ctx, systemInfo)
	// DONE: Handling the message (start)

	_ = ctx

	return &types.MsgCreateGameResponse{
		GameIndex: newIndex,
	}, nil
}
