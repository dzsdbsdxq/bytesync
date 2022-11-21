package router

import (
	"EditSync/controller"
	"EditSync/server/websocket"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(viper.GetString("app.runModel"))
	r.Use(httpCors())
	//r.POST("/api/getLoginQRCode", controller.GetLoginQRCode)
	//r.POST("/api/getWxOpenId", controller.GetWxOpenId)
	r.POST("/api/newPcConnect", controller.NewConnect)
	r.POST("/api/checkOwner", controller.CheckOwner)
	r.POST("/api/login", controller.Login)
	r.POST("/api/sendData", controller.SendData)
	r.POST("/api/logout", controller.Logout)
	r.POST("/api/saveDoc", controller.SaveDoc)
	return r

}

func httpCors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
		context.Next()
	}
}
func WebSocketInit() {
	websocket.Register("login", websocket.LoginController)
	websocket.Register("heartbeat", websocket.HeartbeatController)
	websocket.Register("ping", websocket.PingController)
}
