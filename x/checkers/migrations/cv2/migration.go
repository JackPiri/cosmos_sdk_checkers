package cv2

import (
	"github.com/alice/checkers/x/checkers/keeper"
	cv2Keeper "github.com/alice/checkers/x/checkers/migrations/cv2/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func PerformMigration(ctx sdk.Context, k keeper.Keeper, storedGameChunk uint64) error {
	ctx.Logger().Info("Start to compute checkers games to player info calculation...")
	err := cv2Keeper.MapStoredGamesReduceToPlayerInfo(ctx, k, storedGameChunk)
	if err != nil {
		ctx.Logger().Error("Checkers games to player info computation ended with error")
	} else {
		ctx.Logger().Info("Checkers games to player info computation done")
	}
	return err
}
