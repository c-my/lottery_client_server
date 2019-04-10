package repositories

import (
	"github.com/c-my/lottery_client_server/datamodels"
	"github.com/c-my/lottery_client_server/datasource"
	"github.com/jinzhu/gorm"
	"sync"
)

type ActivityRepository interface {
	SelectAll() []datamodels.Activity
	Append(activity datamodels.Activity) bool
}

type activitySQLRepository struct {
	source *gorm.DB
	mu     sync.RWMutex
}

func (r activitySQLRepository) SelectAll() (activities []datamodels.Activity) {
	r.source.Find(&activities)
	return
}

func (r *activitySQLRepository) Append(activity datamodels.Activity) bool {
	var act datamodels.Activity
	if err := r.source.Where("ID = ?", activity.Id).First(&act).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			r.source.Create(activity)
			return true
		}
	}
	return false
}

func NewActivityRepository() ActivityRepository {
	db := datasource.DB
	if !db.HasTable(&datamodels.Activity{}) {
		db.CreateTable(&datamodels.Activity{})
	}
	return &activitySQLRepository{source: db}
}
