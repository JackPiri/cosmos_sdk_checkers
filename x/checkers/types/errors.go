package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/checkers module sentinel errors
var (
	ErrInvalidBlack         = sdkerrors.Register(ModuleName, 1100, "black address is invalid: %s")
	ErrInvalidRed           = sdkerrors.Register(ModuleName, 1101, "red address is invalid: %s")
	ErrGameNotParseable     = sdkerrors.Register(ModuleName, 1102, "game cannot be parsed")
	ErrInvalidGameIndex     = sdkerrors.Register(ModuleName, 1103, "game index is invalid")
	ErrInvalidPiecePosition = sdkerrors.Register(ModuleName, 1104, "piece position is invalid")
	ErrNotExistingGame      = sdkerrors.Register(ModuleName, 1105, "game does not exist")
	ErrInvalidPlayer        = sdkerrors.Register(ModuleName, 1106, "creator address is not among players")
	ErrTurnMismatch         = sdkerrors.Register(ModuleName, 1107, "turn mismatch")
	ErrWrongMove            = sdkerrors.Register(ModuleName, 1108, "wrong move")
)
