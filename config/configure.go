package config

const (
	CloudWsServer         string = "wss://sampling.alphamj.cn/ws"
	CloudURL              string = "sampling.alphamj.cn"
	CloudLoginURL         string = "https://" + CloudURL + "/signin"
	CloudSignupURL        string = "https://" + CloudURL + "/signup"
	CloudAppendActivities string = "https://" + CloudURL + "/append-activity"
	CloudGetACtivities    string = "https://" + CloudURL + "/get-activities"
	LocalAddr             string = "127.0.0.1"
	LocalPort             string = "1923"
	LocalUrl              string = LocalAddr + ":" + LocalPort
	InitialPath           string = "login"
	InitialUrl            string = LocalUrl + "/" + InitialPath
	LaunchBrowser         bool   = true
)
