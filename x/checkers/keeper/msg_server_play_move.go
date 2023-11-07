package keeper

import (
	"context"
	"strconv"

	"github.com/alice/checkers/x/checkers/rules"
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) PlayMove(goCtx context.Context, msg *types.MsgPlayMove) (*types.MsgPlayMoveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// DONE: Handling the message (start)
	storedGame, found := k.Keeper.GetStoredGame(ctx, msg.GameIndex)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrNotExistingGame, "missing game (%s)", msg.GameIndex)
	}

	if msg.Creator != storedGame.Black && msg.Creator != storedGame.Red {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPlayer, "invalid player (%s)", msg.Creator)
	}
	if (msg.Creator == storedGame.Black && storedGame.Turn != rules.PieceStrings[rules.BLACK_PLAYER]) ||
		(msg.Creator == storedGame.Red && storedGame.Turn != rules.PieceStrings[rules.RED_PLAYER]) {
		return nil, sdkerrors.Wrapf(types.ErrTurnMismatch, "turn is (%s)", storedGame.Turn)
	}
	game, err := storedGame.ParseGame()
	if err != nil {
		panic(err.Error())
	}
	captured, moveErr := game.Move(
		rules.Pos{
			X: int(msg.FromX),
			Y: int(msg.FromY),
		},
		rules.Pos{
			X: int(msg.ToX),
			Y: int(msg.ToY),
		},
	)
	if moveErr != nil {
		return nil, sdkerrors.Wrapf(types.ErrWrongMove, moveErr.Error())
	}

	storedGame.Board = game.String()
	storedGame.Turn = rules.PieceStrings[game.Turn]
	k.Keeper.SetStoredGame(ctx, storedGame)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.MovePlayedEventType,
			sdk.NewAttribute(types.MovePlayedEventCreator, msg.Creator),
			sdk.NewAttribute(types.MovePlayedEventGameIndex, msg.GameIndex),
			sdk.NewAttribute(types.MovePlayedEventCapturedX, strconv.FormatInt(int64(captured.X), 10)),
			sdk.NewAttribute(types.MovePlayedEventCapturedY, strconv.FormatInt(int64(captured.Y), 10)),
			sdk.NewAttribute(types.MovePlayedEventWinner, rules.PieceStrings[game.Winner()]),
		),
	)
	// DONE: Handling the message (end)

	return &types.MsgPlayMoveResponse{
		CapturedX: int32(captured.X),
		CapturedY: int32(captured.Y),
		Winner:    rules.PieceStrings[game.Winner()],
	}, nil
}
