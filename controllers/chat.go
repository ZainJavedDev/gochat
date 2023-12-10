package controllers

import (
	"chat-app/models"
	"chat-app/utils"

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
	Receiver uint   `json:"receiver"`
	Value    string `json:"value"`
}

var clients = make(map[*websocket.Conn]uint)

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

	tokenString := c.Ctx.Input.Header("Authorization")
	userId, _, err := utils.Validate(tokenString)

	db := utils.ConnectDB()
	defer db.Close()

	var userFromDB models.User
	result := db.Where("user_id = ?", userId).First(&userFromDB)
	if result.Error != nil {
		errorMessage := "Unable to connect"
		utils.CreateErrorResponse(&c.Controller, 405, errorMessage)
	}

	clients[conn] = userFromDB.ID
	conn.WriteMessage(websocket.TextMessage, []byte("Hello Client!"))

	var message Message

	for {
		err := conn.ReadJSON(&message)
		if err != nil {
			removeClient(conn)
			return
		}

		for client, clientUserId := range clients {
			if message.Receiver == clientUserId {
				client.WriteMessage(websocket.TextMessage, []byte(message.Value))
			}
		}
	}
}

func removeClient(conn *websocket.Conn) {
	delete(clients, conn)
}
