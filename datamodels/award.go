package datamodels

type Award struct {
	ID   uint   `json:"aid" gorm:"primary_key"`
	Name string `json:"name"`
}
