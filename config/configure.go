package config

const (
	CloudWsServer         string = "wss://sampling.alphamj.cn/ws"
	CloudURL              string = "sampling.alphamj.cn"
	CloudLoginURL         string = "https://" + CloudURL + "/signin"
	CloudSignupURL        string = "https://" + CloudURL + "/signup"
	CloudAppendActivities string = "https://" + CloudURL + "/append-activity"
	CloudGetACtivities    string = "https://" + CloudURL + "/get-activities"
	CloudConfirmCode      string = "BLXDNZ"
	LocalAddr             string = "0.0.0.0"
	LocalPort             string = "1923"
	LocalUrl              string = LocalAddr + ":" + LocalPort
	InitialPath           string = "login"
	InitialUrl            string = "127.0.0.1" + ":" + LocalPort + "/" + InitialPath
	LaunchBrowser         bool   = true
	ConfigureFile         string = "console_configure.json"
)
