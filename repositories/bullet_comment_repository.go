package repositories

import (
	"github.com/c-my/lottery_client_server/datamodels"
	"github.com/jinzhu/gorm"
)

type BulletCommentRepository interface {
	SelectAll() []datamodels.BulletComment
}

type BulletCommentSQLRepository struct {
	source *gorm.DB
}

func (r *BulletCommentSQLRepository) SelectAll() (bullets []datamodels.BulletComment) {
	r.source.Find(&bullets)
	return
}
