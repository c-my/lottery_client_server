package repositories

import (
	"sync"

	"github.com/c-my/lottery_client_server/datamodels"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

//UserSQLRepository handle users from database
type UserSQLRepository struct {
	source *gorm.DB
	mu     sync.RWMutex
}

//SelectByID selects user by uid
func (r *UserSQLRepository) SelectByID(uid uint) (user datamodels.User, found bool) {
	found = r.source.Where("ID = ?", uid).First(&user).RecordNotFound()
	return
}

//RandomSelect randomly select a user
func (r *UserSQLRepository) RandomSelect() (user datamodels.User) {
	r.source.Order(gorm.Expr("random()")).First(&user)
	return
}

//SelectAll returns all users from database
func (r *UserSQLRepository) SelectAll() (users []datamodels.User) {
	r.source.Find(&users)
	return
}

//NewUserRepository is
func NewUserRepository() UserSQLRepository {
	var db *gorm.DB
	// var err error
	db, _ = gorm.Open("sqlite3", "./datasource/user.db")
	if (!db.HasTable(&datamodels.User{})) {
		db.CreateTable(&datamodels.User{})
	}
	return UserSQLRepository{source: db}
}
