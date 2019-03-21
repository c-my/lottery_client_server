package datamodels

import "time"

// User is the data struct of user
type User struct {
	ID        uint   `json:"uid" gorm:"primary_key"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar" gorm:"default:'/assets/avatars/default.png'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
