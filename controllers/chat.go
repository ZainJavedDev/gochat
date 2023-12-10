package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

type ChatController struct {
	beego.Controller
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message struct {
	Receiver string `json:"receiver"`
	Value    string `json:"value"`
}

var clients = make(map[*websocket.Conn]string)

func (c *ChatController) Get() {
	defer c.ServeJSON()
	conn, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
	if err != nil {
		beego.Error("Error upgrading to WebSocket:", err)
		return
	}
	defer func() {
		removeClient(conn)
		conn.Close()
	}()

	username := c.GetString("username")
	clients[conn] = username
	conn.WriteMessage(websocket.TextMessage, []byte("Hello Client!"))

	var message Message

	for {
		err := conn.ReadJSON(&message)
		if err != nil {
			removeClient(conn)
			return
		}

		for client, clientUsername := range clients {
			if message.Receiver == clientUsername {
				client.WriteMessage(websocket.TextMessage, []byte(message.Value))
			}
		}
	}
}

func removeClient(conn *websocket.Conn) {
	delete(clients, conn)
}
