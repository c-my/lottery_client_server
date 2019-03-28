package datamodels

import "time"

type BulletComment struct {
	openid    string `json:"openid"`
	danmu     string `json:"danmu"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
