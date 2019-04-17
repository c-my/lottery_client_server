package datamodels

type Activity struct {
	Id        int    `json:"id" gorm:"primary_key"`
	Name      string `json:"name"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}
