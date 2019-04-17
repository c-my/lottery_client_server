package datamodels

type Configure struct {
	Online   string `json:"online"`
	UserType string `json:"userType"`

	// basic information
	ActivityName string `json:"activity-name"`

	// drawing settings
	DrawModeChosen   string `json:"draw_mode-chosen"`
	RewardItemsNames string `json:"reward-items-names"`
	PrizeNames       string `json:"prize-names"`

	// methods and themes
	LotteryStype  string `json:"lottery-style"`
	TpoicColor    string `json:"tpoic-color"`
	LotteryMusic  string `json:"lottery-music"`
	WinPrizeMusic string `json:"win-prize-music"`
	GetPrizeMusic string `json:"get-prize-music"`

	// bullet settings
	BulletFontSize          int    `json:"bullet-font-size"`
	BulletTransparentDegree int    `json:"bullet-transparent-degree"`
	BulletFont              string `json:"bullet-font"`
	BulletColor             string `json:"bullet-color"`
	BulletVelocity          string `json:"bullet-velocity"`
	BulletLocation          string `json:"bullet-location"`
	BulletEnable            string `json:"bullet-enable"`
	BulletCheckEnable       string `json:"bullet-check-enable"`

	// drawing status bar on the left
	RewardUsersList string `json:"reward-users-list"`
	RewardRemain    string `json:"reward-remain"`
}
