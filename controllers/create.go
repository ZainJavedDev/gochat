package controllers

import (
	"chat-app/models"
	"chat-app/utils"

	"github.com/astaxie/beego"
)

type CreateRoomController struct {
	beego.Controller
}

type Room struct {
	RoomName string `form:"roomname"`
}

func (c *CreateRoomController) Post() {
	tokenString := c.Ctx.Input.Header("Authorization")
	userId, _, err := utils.Validate(tokenString)
	if err != nil {
		utils.CreateErrorResponse(&c.Controller, 400, "Invalid or expired token.")
	}
	db := utils.ConnectDB()
	defer db.Close()

	var room Room
	result := db.Create(&models.Room{UserId: userId, Name: room.RoomName})
	if result.Error != nil {
		utils.CreateErrorResponse(&c.Controller, 400, "Unable to create room!")
	}

	responseData := map[string]interface{}{
		"message": "Room created successfully!",
	}

	c.Data["json"] = responseData
	c.ServeJSON()
}
