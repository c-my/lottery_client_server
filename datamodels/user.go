package datamodels

//User is the data struct of user
type User struct {
	ID       uint   `json:"uid"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}
