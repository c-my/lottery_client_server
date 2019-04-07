package repositories

import (
	"github.com/c-my/lottery_client_server/datamodels"
	"github.com/c-my/lottery_client_server/datasource"
	"github.com/jinzhu/gorm"
	"sync"
)

type ActivityRepository interface {
	SelectAll() []datamodels.Activity
}

type activitySQLRepository struct {
	source *gorm.DB
	mu     sync.RWMutex
}

func (r activitySQLRepository) SelectAll() (activities []datamodels.Activity) {
	r.source.Find(&activities)
	return
}

func NewActivityRepository() ActivityRepository {
	db := datasource.DB
	if !db.HasTable(&datamodels.Activity{}) {
		db.CreateTable(&datamodels.Activity{})
	}
	return &activitySQLRepository{source: db}
}
