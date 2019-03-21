package main

import (
	"encoding/json"
	"fmt"

	"github.com/c-my/lottery_client_server/repositories"
	"github.com/c-my/lottery_client_server/services"
	"github.com/c-my/lottery_client_server/web/controllers"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/websocket"
)

var userRepository = repositories.NewUserRepository()

func main() {
	app := iris.New()

	app.StaticWeb("/assets", "./assets")
	app.StaticWeb("/css", "./assets/css")
	app.StaticWeb("/js", "./assets/js")
	app.StaticWeb("/fonts", "./assets/fonts")

	mvc.Configure(app.Party("/get-exist-user"), users)
	mvc.Configure(app.Party("get-awards"), awards)
	app.Get("/", func(ctx iris.Context) {
		ctx.ServeFile("console.html", false)
	})

	app.Get("/screen", func(ctx iris.Context) {
		ctx.ServeFile("PrizeDraw.html", false)
	})

	setupWebsocket(app)

	app.Run(
		iris.Addr(":8000"),
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

func handleWebsocket(c websocket.Connection) {
	c.OnMessage(func(data []byte) {
		var msg message
		json.Unmarshal(data, &msg)

		switch msg["action"] {
		case "stop-drawing":
			luckyDog := userRepository.RandomSelect()

			j, _ := addAction("who-is-lucky-dog", luckyDog)
			fmt.Println(string(j))
			c.To(websocket.All).EmitMessage(j)

			println("lucy dog is: ", luckyDog.ID)
		default:
			c.To(websocket.Broadcast).EmitMessage(data)
		}

	})
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
