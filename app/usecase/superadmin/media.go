package superadmin_usecase

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

func (u *superadminAppUsecase) markMediaAsUnusedByIds(ctx context.Context, ids []string) {
	err := u.mongoDbRepo.UpdateManyMediaPartial(ctx, map[string]interface{}{
		"ids": ids,
	}, map[string]interface{}{
		"isUsed":    false,
		"updatedAt": time.Now(),
	})
	if err != nil {
		logrus.Errorf("Update media error %v", err)
	}
}
