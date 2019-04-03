package datamodels

import "time"

// User is the data struct of user
type User struct {
	ID        string `json:"uid" gorm:"primary_key"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar" gorm:"default:'/assets/avatars/default.png'"`
	Language  string `json:"language"`
	Country   string `json:"country"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Gender    uint   `json:"gender"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
