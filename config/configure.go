package config

const (
	CloudWsServer string = "wss://sampling.alphamj.cn/ws"
	CloudLoginURL string = "https://sampling.alphamj.cn/signin"
	LocalAddr     string = "127.0.0.1"
	LocalPort     string = "1923"
	LocalUrl      string = LocalAddr + ":" + LocalPort
	InitialPath   string = "login"
	InitialUrl    string = LocalUrl + "/" + InitialPath
)
