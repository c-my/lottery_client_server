package main

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/c-my/lottery_client_server/datasource"
	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
)

func main() {
	app := iris.New()

	app.StaticWeb("/assets", "./assets")
	app.StaticWeb("/css", "./assets/css")
	app.StaticWeb("/js", "./assets/js")
	app.StaticWeb("/fonts", "./assets/fonts")
	app.Get("/get-exist-user", getExistUsers)

	app.Get("/", func(ctx iris.Context) {
		ctx.ServeFile("console.html", false)
	})

	app.Get("/screen", func(ctx iris.Context) {
		ctx.ServeFile("PrizeDraw.html", false)
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

type message map[string]interface{}

func handleWebsocket(c websocket.Connection) {
	c.OnMessage(func(data []byte) {
		var msg message
		json.Unmarshal(data, &msg)

		switch msg["action"] {
		case "stop-drawing":
			luckyNumber := rand.Intn(len(datasource.Users))
			j, _ := addAction("who-is-lucky-dog", datasource.Users[luckyNumber])
			fmt.Println(string(j))
			c.To(websocket.All).EmitMessage(j)

			println("lucy dog is: ", datasource.Users[luckyNumber].ID)
		default:
			c.To(websocket.Broadcast).EmitMessage(data)
		}

	})
}

func getExistUsers(ctx iris.Context) {
	ctx.JSON(datasource.Users)
}

func addAction(action string, content interface{}) ([]byte, error) {
	m := map[string]interface{}{
		"action":  action,
		"content": content,
	}
	return json.Marshal(m)
}
