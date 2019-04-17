package datamodels

import "github.com/jinzhu/gorm"

type BulletComment struct {
	gorm.Model
	UID   string `json:"uid"`
	DanMu string `json:"danmu"`
}
