package repositories

import (
	"sync"
	"time"

	"github.com/c-my/lottery_client_server/datamodels"
	"github.com/c-my/lottery_client_server/datasource"
	"github.com/jinzhu/gorm"
)

// UserRepository handles basic operations of user
type UserRepository interface {
	SelectByID(uid uint) (user datamodels.User, found bool)
	SelectAll() []datamodels.User
	RandomSelect() datamodels.User
	Append(user datamodels.User) bool
}

// UserSQLRepository handle users from database
type userSQLRepository struct {
	source *gorm.DB
	mu     sync.RWMutex
}

// SelectByID selects user by uid
func (r *userSQLRepository) SelectByID(uid uint) (user datamodels.User, found bool) {
	found = r.source.Where("ID = ?", uid).First(&user).RecordNotFound()
	return
}

// RandomSelect randomly select a user
func (r *userSQLRepository) RandomSelect() (user datamodels.User) {
	r.source.Order(gorm.Expr("random()")).First(&user)
	return
}

// SelectAll returns all users from database
func (r *userSQLRepository) SelectAll() (users []datamodels.User) {
	r.source.Find(&users)
	return
}

func (r *userSQLRepository) Append(user datamodels.User) bool {
	if !r.source.NewRecord(user) {
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		r.source.Create(user)
		return true
	}
	return false
}

// NewUserRepository is
func NewUserRepository() UserRepository {
	db := datasource.DB
	if !db.HasTable(&datamodels.User{}) {
		db.CreateTable(&datamodels.User{})
	}
	return &userSQLRepository{source: db}
}
