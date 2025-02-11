package repositories

import (
	"github.com/c-my/lottery_client_server/datamodels"
	"github.com/c-my/lottery_client_server/datasource"
	"github.com/jinzhu/gorm"
)

type BulletCommentRepository interface {
	SelectAll() []datamodels.BulletComment
	Append(danmu datamodels.BulletComment) bool
}

type bulletCommentSQLRepository struct {
	source *gorm.DB
}

func (r *bulletCommentSQLRepository) SelectAll() (bullets []datamodels.BulletComment) {
	r.source.Find(&bullets)
	return
}

func (r *bulletCommentSQLRepository) Append(danmu datamodels.BulletComment) bool {
	r.source.Create(danmu)
	return true
}

func NewBulletCommentRepository() BulletCommentRepository {
	db := datasource.DB
	if (!db.HasTable(&datamodels.BulletComment{})) {
		db.CreateTable(&datamodels.BulletComment{})
	}
	return &bulletCommentSQLRepository{source: datasource.DB}
}
