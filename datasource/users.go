package datasource

import "github.com/c-my/lottery_iris/datamodels"

// Users is a collection of user
var Users = map[int]datamodels.User{
	0: {
		ID:       20160001,
		Nickname: "dalao2",
		Avatar:   "/assets/avatar2.png",
	},

	1: {
		ID:       20160000,
		Nickname: "dalao",
		Avatar:   "/assets/avatar0.png",
	},
}
