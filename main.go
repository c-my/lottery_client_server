package main

import (
	"encoding/json"
	"math/rand"

	"github.com/c-my/lottery/datasource"
	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
)

func main() {
	app := iris.New()

	app.StaticWeb("/assets", "./assets")
	app.StaticWeb("/css", "./assets/css")
	app.StaticWeb("/js", "./assets/js")
	app.Get("/get-exist-user", getExistUsers)

	app.Get("/", func(ctx iris.Context) {
		ctx.ServeFile("console.html", false)
	})

	setupWebsocket(app)

	app.Run(iris.Addr(":8000"))
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

func handleWebsocket(c websocket.Connection) {
	c.On("stop-drawing", func() {
		luckyNumber := rand.Intn(len(datasource.Users))
		j, _ := json.Marshal(datasource.Users[luckyNumber])
		c.To(websocket.All).Emit("who-is-lucky-dog", j)

		println("lucy dog is: ", datasource.Users[luckyNumber].ID)

	})
}

func getExistUsers(ctx iris.Context) {
	ctx.JSON(datasource.Users)
}
