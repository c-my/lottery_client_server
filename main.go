package main

import (
	"encoding/json"
	"fmt"
	"github.com/c-my/lottery_client_server/datamodels"
	"github.com/c-my/lottery_client_server/repositories"
	"github.com/c-my/lottery_client_server/services"
	"github.com/c-my/lottery_client_server/web/controllers"
	"github.com/c-my/lottery_client_server/web/websockets"
	gwebsocket "github.com/gorilla/websocket"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/websocket"
)

var userRepository = repositories.NewUserRepository()

const cloudWsServer string = "wss://sampling.alphamj.cn/ws"

var app *iris.Application
var sendChan = make(chan []byte)

func main() {
	app = iris.New()

	app.StaticWeb("/assets", "./assets")
	app.StaticWeb("/css", "./assets/css")
	app.StaticWeb("/js", "./assets/js")
	app.StaticWeb("/fonts", "./assets/fonts")
	app.StaticWeb("/img", "./assets/img")

	mvc.Configure(app.Party("/get-exist-user"), users)
	mvc.Configure(app.Party("get-awards"), awards)
	app.Get("/", func(ctx iris.Context) {
		ctx.ServeFile("index.html", false)
	})

	app.Get("/index", func(ctx iris.Context) {
		ctx.ServeFile("index.html", false)
	})

	app.Get("/console", func(ctx iris.Context) {
		ctx.ServeFile("console.html", false)
	})

	app.Get("/screen", func(ctx iris.Context) {
		ctx.ServeFile("PrizeDraw.html", false)
	})

	setupWebsocket(app)

	wsc := websockets.NewWebsocketClient(cloudWsServer)
	wsc.SetHandler(getWsCRecv)
	wsc.Run()

	app.Run(
		iris.Addr(":1923"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)

}

func setupWebsocket(app *iris.Application) {
	ws := websocket.New(websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	})

	ws.OnConnection(handleWebsocket)

	app.Get("/ws", ws.Handler())

	app.Any("/iris-ws.js", func(ctx iris.Context) {
		ctx.Write(websocket.ClientSource)
	})

}

type message map[string]interface{}
type drawMsg struct {
	dkind string
}

func handleWebsocket(c websocket.Connection) {
	c.OnMessage(func(data []byte) {
		var msg message
		json.Unmarshal(data, &msg)
		fmt.Print("received from controller: ")
		fmt.Println(msg)
		switch msg["action"] {
		case "stop-drawing":
			luckyDog := userRepository.RandomSelect()

			j, _ := addAction("who-is-lucky-dog", luckyDog)
			fmt.Println(string(j))
			c.To(websocket.All).EmitMessage(j)

		// println("lucy dog is: ", luckyDog.ID)
		case "append-user":
			c.To(cloudWsServer).EmitMessage(data)
			fmt.Println("append user:")
			fmt.Println(msg)
		case "start-draw":
			var dmsg = drawMsg{dkind: "cube"}
			d, _ := addAction("start-draw", dmsg)
			fmt.Println("message delivered: " + string(d))
			c.To(websocket.All).EmitMessage(d)
			break
		default:
			c.To(websocket.Broadcast).EmitMessage(data)
		}

	})
	go wsWriter(c)
}

func addAction(action string, content interface{}) ([]byte, error) {
	m := map[string]interface{}{
		"action":  action,
		"content": content,
	}
	return json.Marshal(m)
}

func users(app *mvc.Application) {
	repo := repositories.NewUserRepository()
	userService := services.NewUserService(repo)
	app.Register(userService)
	app.Handle(new(controllers.UserController))
}

func awards(app *mvc.Application) {
	repo := repositories.NewAwardSQLRepository()
	awardService := services.NewAwardService(repo)
	app.Register(awardService)
	app.Handle(new(controllers.AwardController))
}

func wsWriter(c websocket.Connection) {
	for {
		var msg []byte
		msg = <-sendChan
		fmt.Println("message delivered to cloud:" + string(msg))
		c.To(websocket.All).EmitMessage(msg)
	}
}

func getWsCRecv(wsc *websockets.WebsocketClient, messageType int, p []byte) {
	fmt.Println(string(p))
	switch messageType {
	case gwebsocket.TextMessage:
		var msg message
		json.Unmarshal(p, &msg)
		fmt.Println("received from cloud: \n" + string(p))
		switch msg["action"] {
		case "append-user":
			content := msg["content"]
			var u datamodels.User
			{
			}
			jcontent, _ := json.Marshal(content)
			//fmt.Println("user to append: " + string(jcontent))
			json.Unmarshal(jcontent, &u)
			userRepository.Append(u)
			sendChan <- p
			break
		case "send-danmu":
			sendChan <- p
			break
		}
		//fmt.Println(string(p))
		break
	}
}
