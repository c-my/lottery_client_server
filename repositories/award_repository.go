package repositories

import (
	"github.com/c-my/lottery_client_server/datamodels"
	"github.com/c-my/lottery_client_server/datasource"
	"github.com/jinzhu/gorm"
)

type AwardSQLRepository struct {
	source gorm.DB
}

//NewAwardSQLRepository returns a award repository
func NewAwardSQLRepository() AwardSQLRepository {
	db := datasource.DB
	if !db.HasTable(&datamodels.Award{}) {
		db.CreateTable(&datamodels.Award{})
	}
	return AwardSQLRepository{}
}
