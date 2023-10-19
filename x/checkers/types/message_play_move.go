package types

import (
	"strconv"

	"github.com/alice/checkers/x/checkers/rules"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgPlayMove = "play_move"

var _ sdk.Msg = &MsgPlayMove{}

func NewMsgPlayMove(creator string, gameIndex string, fromX uint64, fromY uint64, toX uint64, toY uint64) *MsgPlayMove {
	return &MsgPlayMove{
		Creator:   creator,
		GameIndex: gameIndex,
		FromX:     fromX,
		FromY:     fromY,
		ToX:       toX,
		ToY:       toY,
	}
}

func (msg *MsgPlayMove) Route() string {
	return RouterKey
}

func (msg *MsgPlayMove) Type() string {
	return TypeMsgPlayMove
}

func (msg *MsgPlayMove) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgPlayMove) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgPlayMove) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	gameIndex, err := strconv.ParseUint(msg.GameIndex, 10, 64)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidGameIndex, "not parseable (%s)", err)
	}
	if gameIndex < DefaultIndex {
		return sdkerrors.Wrapf(ErrInvalidGameIndex, "value too low (%d)", msg.GameIndex)
	}
	if msg.FromX >= rules.BOARD_DIM || msg.FromY >= rules.BOARD_DIM ||
		msg.ToX >= rules.BOARD_DIM || msg.ToY >= rules.BOARD_DIM ||
		(msg.FromX+msg.FromY)%2 == 0 ||
		(msg.ToX+msg.ToY)%2 == 0 { // e.g. (0,0) invalid
		return sdkerrors.Wrapf(ErrInvalidPiecePosition, "out of board or out of playable cells")
	}

	return nil
}
