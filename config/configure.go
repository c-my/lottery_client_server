package config

const (
	CloudWsServer  string = "wss://sampling.alphamj.cn/ws"
	CloudLoginURL  string = "https://sampling.alphamj.cn/signin"
	CloudSignupURL string = "https://sampling.alphamj.cn/signup"
	LocalAddr      string = "127.0.0.1"
	LocalPort      string = "1923"
	LocalUrl       string = LocalAddr + ":" + LocalPort
	InitialPath    string = "login"
	InitialUrl     string = LocalUrl + "/" + InitialPath
	LaunchBrowser  bool   = true
)
