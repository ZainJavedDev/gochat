package routers

import (
	"chat-app/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/migrate", &controllers.MigrationController{})

	beego.Router("/signup", &controllers.SignupController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/refresh", &controllers.RefreshTokenController{})

	beego.Router("/chat", &controllers.ChatController{})

	beego.Router("/room/chat", &controllers.RoomChatController{})
	beego.Router("/room/create", &controllers.CreateRoomController{})
	beego.Router("/room/list", &controllers.ListRoomController{})
}
