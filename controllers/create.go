package controllers

import "github.com/astaxie/beego"

type CreateRoomController struct {
	beego.Controller
}

type Room struct {
	RoomName string `form:"roomname"`
}

func (c *CreateRoomController) Post() {

}
