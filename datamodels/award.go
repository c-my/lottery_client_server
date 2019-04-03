package datamodels

// Award is the data struct of award
type Award struct {
	ID    uint   `json:"aid" gorm:"primary_key"`
	Name  string `json:"name"`
	Pic   string `json:"picture"`
	Count uint   `json:"count"`
}
