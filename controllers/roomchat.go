package controllers

import (
	"chat-app/models"
	"chat-app/utils"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

type Message struct {
	Receiver uint   `json:"receiver"`
	Value    string `json:"value"`
}

type RoomChatController struct {
	beego.Controller
}

var roomClients = make(map[*websocket.Conn]uint)

func (c *RoomChatController) Get() {
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
	_, _, err = utils.Validate(tokenString)
	if err != nil {
		errorResponse := "Invalid or expired token"
		utils.CreateErrorResponse(&c.Controller, 400, errorResponse)
	}

	db := utils.ConnectDB()
	defer db.Close()
	var roomFromDB models.Room
	roomId := c.Ctx.Input.Query("room")
	result := db.Where("id = ?", roomId).First(&roomFromDB)
	if result.Error != nil {
		errorMessage := "Unable to connect"
		utils.CreateErrorResponse(&c.Controller, 405, errorMessage)
	}
	roomIdUint, _ := strconv.ParseUint(roomId, 10, 0)
	var roomIdAsUint uint = uint(roomIdUint)
	clients[conn] = roomIdAsUint
	conn.WriteMessage(websocket.TextMessage, []byte("Hello Client!"))

	var message Message

	for {
		err := conn.ReadJSON(&message)
		if err != nil {
			removeClient(conn)
			return
		}

		for client, clientRoomId := range roomClients {
			if message.Receiver == clientRoomId {
				client.WriteMessage(websocket.TextMessage, []byte(message.Value))
			}
		}
	}
}
