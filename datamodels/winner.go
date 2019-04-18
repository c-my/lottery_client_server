package datamodels

type Winner struct {
	UID string `json:"uid" gorm:"primary_key"`
	AID string `json:"aid" gorm:"primary_key"`
}
