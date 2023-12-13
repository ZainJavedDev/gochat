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

type JoinRoomController struct {
	beego.Controller
}

type Message struct {
	Receiver uint   `json:"receiver"`
	Room     uint   `json:"room"`
	Value    string `json:"value"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients = make(map[*websocket.Conn]uint)
var roomClients = make(map[uint]uint) // key is client id and value is room id

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
	if err != nil {
		errorResponse := "Invalid or expired token"
		utils.CreateErrorResponse(&c.Controller, 400, errorResponse)
	}

	db := utils.ConnectDB()
	defer db.Close()

	var userFromDB models.User
	result := db.Where("id = ?", userId).First(&userFromDB)
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
		if message.Receiver != 0 {
			for clientConn, clientUserId := range clients {
				if message.Receiver == clientUserId {
					clientConn.WriteMessage(websocket.TextMessage, []byte(message.Value))
				}
			}
		} else {
			for roomClientUserId, roomId := range roomClients {
				if message.Room == roomId {
					for clientConn, clientUserId := range clients {
						if roomClientUserId == clientUserId {
							clientConn.WriteMessage(websocket.TextMessage, []byte(message.Value))
						}
					}
				}
			}
		}
	}
}

func removeClient(conn *websocket.Conn) {
	delete(clients, conn)
}

func (c *JoinRoomController) Post() {
	var room Room
	err := c.ParseForm(&room)
	if err != nil {
		utils.CreateErrorResponse(&c.Controller, 500, err.Error())
	}

	tokenString := c.Ctx.Input.Header("Authorization")
	userId, _, err := utils.Validate(tokenString)
	if err != nil {
		utils.CreateErrorResponse(&c.Controller, 400, "Invalid or expired token.")
	}

	db := utils.ConnectDB()
	defer db.Close()

	dbRoom := models.Room{}
	result := db.Where("name = ?", room.RoomName).First(&dbRoom)
	if result.Error != nil {
		utils.CreateErrorResponse(&c.Controller, 404, "No room with such name.")
	}
	roomClients[userId] = dbRoom.ID

	responseData := map[string]string{
		"message": "Room joined successfully!",
	}
	c.Data["json"] = responseData
	c.ServeJSON()
}
