package main

import (
	"github.com/c-my/lottery_client_server/repositories"
	"github.com/c-my/lottery_client_server/services"
	"github.com/c-my/lottery_client_server/web/controllers"
	"github.com/c-my/lottery_client_server/web/routers"
	"github.com/gorilla/mux"
	"github.com/kataras/iris/mvc"
	"log"
	"net/http"
	"time"
)

var userRepository = repositories.NewUserRepository()

const cloudWsServer string = "wss://sampling.alphamj.cn/ws"

//var app *iris.Application
var sendChan = make(chan []byte)

func main() {
	r := mux.NewRouter()
	routers.SetSubRouter("127.0.0.1:1923", r)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:1923",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
	//app = iris.New()
	//
	//
	//mvc.Configure(app.Party("get-awards"), awards)
	//app.Get("/", func(ctx iris.Context) {
	//	ctx.ServeFile("index.html", false)
	//})
	//
	//
	//setupWebsocket(app)
	//
	//wsc := websockets.NewWebsocketClient(cloudWsServer)
	//wsc.SetHandler(getWsCRecv)
	//wsc.Run()
	//
	//app.Run(
	//	iris.Addr(":1923"),
	//	iris.WithoutServerError(iris.ErrServerClosed),
	//	iris.WithOptimizations,
	//)

}

//func setupWebsocket(app *iris.Application) {
//	ws := websocket.New(websocket.Config{
//		ReadBufferSize:  1024,
//		WriteBufferSize: 1024,
//	})
//
//	ws.OnConnection(handleWebsocket)
//
//	app.Get("/ws", ws.Handler())
//
//	app.Any("/iris-ws.js", func(ctx iris.Context) {
//		ctx.Write(websocket.ClientSource)
//	})
//
//}

//
//func handleWebsocket(c websocket.Connection) {
//	c.OnMessage(func(data []byte) {
//		var msg message
//		json.Unmarshal(data, &msg)
//		fmt.Print("received from controller: ")
//		fmt.Println(msg)
//		switch msg["action"] {
//		case "stop-drawing":
//			luckyDog := userRepository.RandomSelect()
//
//			j, _ := addAction("who-is-lucky-dog", luckyDog)
//			fmt.Println(string(j))
//			c.To(websocket.All).EmitMessage(j)
//
//		// println("lucy dog is: ", luckyDog.ID)
//		case "append-user":
//			c.To(cloudWsServer).EmitMessage(data)
//			fmt.Println("append user:")
//			fmt.Println(msg)
//		case "start-drawing":
//			//var dmsg = drawMsg{dkind: "cube"}
//			//d, _ := addAction("start-drawing", dmsg)
//			fmt.Println("message delivered: " + string(data))
//			c.To(websocket.All).EmitMessage(data)
//			break
//		default:
//			c.To(websocket.Broadcast).EmitMessage(data)
//		}
//
//	})
//	go wsWriter(c)
//}

//func addAction(action string, content interface{}) ([]byte, error) {
//	m := map[string]interface{}{
//		"action":  action,
//		"content": content,
//	}
//	return json.Marshal(m)
//}

func users(app *mvc.Application) {
	repo := repositories.NewUserRepository()
	userService := services.NewUserService(repo)
	app.Register(userService)
	app.Handle(new(controllers.UserController))
}

//func wsWriter(c websocket.Connection) {
//	for {
//		var msg []byte
//		msg = <-sendChan
//		fmt.Println("message delivered to cloud:" + string(msg))
//		c.To(websocket.All).EmitMessage(msg)
//	}
//}

//func getWsCRecv(wsc *websockets.WebsocketClient, messageType int, p []byte) {
//	fmt.Println(string(p))
//	switch messageType {
//	case gwebsocket.TextMessage:
//		var msg message
//		json.Unmarshal(p, &msg)
//		fmt.Println("received from cloud: \n" + string(p))
//		switch msg["action"] {
//		case "append-user":
//			content := msg["content"]
//			var u datamodels.User
//			{
//			}
//			jcontent, _ := json.Marshal(content)
//			//fmt.Println("user to append: " + string(jcontent))
//			json.Unmarshal(jcontent, &u)
//			userRepository.Append(u)
//			sendChan <- p
//			break
//		case "send-danmu":
//			sendChan <- p
//			break
//		}
//		//fmt.Println(string(p))
//		break
//	}
//}
