package controllers

import (
	"chat-app/models"
	"chat-app/utils"

	"github.com/astaxie/beego"
)

type CreateRoomController struct {
	beego.Controller
}

type ListRoomController struct {
	beego.Controller
}

type Room struct {
	RoomName string `form:"roomname"`
}

func (c *CreateRoomController) Post() {

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

func (c *ListRoomController) Post() {

	var room Room
	err := c.ParseForm(&room)
	if err != nil {
		utils.CreateErrorResponse(&c.Controller, 500, err.Error())
	}

	tokenString := c.Ctx.Input.Header("Authorization")
	_, _, err = utils.Validate(tokenString)
	if err != nil {
		utils.CreateErrorResponse(&c.Controller, 400, "Invalid or expired token.")
	}

	db := utils.ConnectDB()
	defer db.Close()

	rooms := []models.Room{}
	result := db.Find(&rooms)

	c.Data["json"] = result
	c.ServeJSON()
}
