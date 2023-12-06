package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

type MainController struct {
	beego.Controller
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients = make(map[*websocket.Conn]bool)

func (c *MainController) Get() {
	defer c.ServeJSON()
	conn, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
	if err != nil {
		beego.Error("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()
	clients[conn] = true
	conn.WriteMessage(websocket.TextMessage, []byte("Hello Client!"))
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			beego.Error("Error reading from WebSocket:", err)
			return
		}
		beego.Info("Message received from client:", string(msg))
		for client := range clients {
			if client == conn {
				continue
			}
			err = client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				beego.Error("Error writing to WebSocket:", err)
				return
			}
		}
	}
}
