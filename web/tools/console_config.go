package tools

import (
	"encoding/json"
	"github.com/c-my/lottery_client_server/config"
	"github.com/c-my/lottery_client_server/web/logger"
	"io/ioutil"
	"os"
	"time"
)

var ConsoleConfig = Configure{}

func Run() {
	go ConsoleConfigLoop()
}

func ConsoleConfigLoop() {
	for {
		time.Sleep(1000)
		saveConfigure(config.ConfigureFile)
	}
}

func saveConfigure(fileName string) bool {
	jsonStr, _ := json.Marshal(ConsoleConfig)
	err := ioutil.WriteFile(fileName, jsonStr, os.ModeAppend)
	if err != nil {
		logger.Error.Println("failed to write config file:", err)
		return false
	}
	return true
}

func loadConfigure(fileName string) (config Configure) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		logger.Error.Println("failed to read config file:", err)
		return
	}
	json.Unmarshal(bytes, config)
	return
}

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

	IsActivityUnfinished bool `json:"is-activity-unfinished"`
}
