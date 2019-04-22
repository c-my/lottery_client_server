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
		SaveConfigure(config.ConfigureFile)
	}
}

func SaveConfigure(fileName string) bool {
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
	ActivityName string `json:"activity_name"`

	// drawing settings
	DrawModeChosen   string `json:"draw_mode_chosen"`
	RewardItemsNames string `json:"reward_items_names"`
	PrizeNames       string `json:"prize_names"`

	// methods and themes
	LotteryStype  string `json:"lottery_style"`
	TpoicColor    string `json:"topic_color"`
	LotteryMusic  string `json:"lottery_music"`
	WinPrizeMusic string `json:"win_prize_music"`
	GetPrizeMusic string `json:"get_prize_music"`

	// bullet settings
	BulletFontSize          string `json:"bullet_font_size"`
	BulletTransparentDegree string `json:"bullet_transparent_degree"`
	BulletFont              string `json:"bullet_font"`
	BulletColor             string `json:"bullet_color"`
	BulletVelocity          string `json:"bullet_velocity"`
	BulletLocation          string `json:"bullet_location"`
	BulletEnable            string `json:"bullet_enable"`
	BulletCheckEnable       string `json:"bullet_check_enable"`

	// drawing status bar on the left
	RewardUsersList string `json:"reward_users_list"`
	RewardRemain    string `json:"reward_remain"`

	IsActivityUnfinished bool `json:"is-activity-unfinished"`
}
