package tools

import (
	"github.com/c-my/lottery_client_server/web/logger"
	"os/exec"
	"runtime"
)

func LaunchBrowser(url string) {
	switch runtime.GOOS {
	case "windows":
		err := exec.Command("cmd.exe", " /c start "+url).Start()
		if err != nil {
			logger.Warning.Println("failed to launch browser:", err)
		}
	case "linux":
		err := exec.Command("xdg-open", url).Start()
		if err != nil {
			logger.Warning.Println("failed to launch browser:", err)
		}
	case "darwin":
		err := exec.Command("open", url).Start()
		if err != nil {
			logger.Warning.Println("failed to launch browser:", err)
		}
	default:
		logger.Info.Println("failed to launch browser: unrecognized system")
	}
}
