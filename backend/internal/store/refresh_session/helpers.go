package refreshsession

import (
	"errors"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"gorm.io/gorm"
)

func ByID(id int) core.UpdateRefreshSessionParam {
	return func(tx *gorm.DB) error {
		if err := tx.
			Where("id=?", id).
			Delete(&core.RefreshSession{}).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return err
		}
		return nil
	}
}

func ByNothing() core.UpdateRefreshSessionParam {
	return func(*gorm.DB) error { return nil }
}
