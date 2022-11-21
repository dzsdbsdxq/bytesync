package controller

import (
	"EditSync/cache"
	"EditSync/common"
	"EditSync/middleware"
	"EditSync/models"
	"EditSync/server/websocket"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
	"time"
)

// NewConnect 创建连接，返回房间号
func NewConnect(c *gin.Context) {

	roomTitle := ""
	roomDesc := ""
	roomCTime := uint32(time.Now().Nanosecond())
	//自动生成临时唯一房间号
	roomId := common.RandPass(6)

	//自动创建用户ID
	userId := common.GetUserIdByRandom()

	//检查redis是否安装启动
	_, err := middleware.GetClient().Ping().Result()
	if err != nil {
		common.RespFunc(c, http.StatusBadRequest, common.ServerError, "")
		return
	}

	//将房间号注册到系统中
	roomOnline := &models.RoomOnline{
		RoomId:     roomId,
		RoomTitle:  roomTitle,
		RoomDesc:   roomDesc,
		RoomOwner:  userId,
		CreateTime: roomCTime,
	}
	cache.SetRoomInfo(roomId, roomOnline)

	websocket.SetRoomIds(roomId)

	common.RespFunc(c, http.StatusOK, common.OK, models.Json{
		"roomId":     roomOnline.RoomId,
		"roomOwner":  roomOnline.RoomOwner,
		"userId":     userId,
		"roomTitle":  roomOnline.RoomTitle,
		"roomDesc":   roomOnline.RoomDesc,
		"createTime": roomCTime,
	})
}

func Login(c *gin.Context) {
	roomId := com.StrTo(c.PostForm("roomId")).String()
	//获取房间信息
	roomOnline, err := cache.GetRoomOnlineInfo(roomId)
	if err != nil {
		common.RespFunc(c, http.StatusBadRequest, common.RoomInfoNotExist, models.Json{})
		return
	}
	//检查redis是否安装启动
	_, err = middleware.GetClient().Ping().Result()
	if err != nil {
		common.RespFunc(c, http.StatusBadRequest, common.ServerError, "")
		return
	}

	//自动创建用户ID
	userId := common.GetUserIdByRandom()

	//返回房间的信息和用户ID
	common.RespFunc(c, http.StatusOK, common.OK, models.Json{
		"roomId":    roomOnline.RoomId,
		"roomOwner": roomOnline.RoomOwner,
		"userId":    userId,
		"roomTitle": roomOnline.RoomTitle,
		"roomDesc":  roomOnline.RoomDesc,
	})

}

func CheckOwner(c *gin.Context) {
	roomId := com.StrTo(c.PostForm("roomId")).String()
	userId := com.StrTo(c.PostForm("userId")).String()

	//获取房间信息
	roomOnline, err := cache.GetRoomOnlineInfo(roomId)
	if err != nil {
		common.RespFunc(c, http.StatusBadRequest, common.RoomInfoNotExist, models.Json{})
		return
	}
	if roomOnline.RoomOwner == userId {
		common.RespFunc(c, http.StatusOK, common.OK, models.Json{"owner": 1})
		return
	}
	common.RespFunc(c, http.StatusOK, common.OK, models.Json{"owner": 0})
}
func SendData(c *gin.Context) {
	//获取到roomId、userId、content
	roomId := com.StrTo(c.PostForm("roomId")).String()
	userId := com.StrTo(c.PostForm("userId")).String()
	msgId := com.StrTo(c.PostForm("msgId")).String()
	content := com.StrTo(c.PostForm("content")).String()

	//将内容传送到在线用户
	fmt.Println("http_request 给全体用户发送消息", roomId, userId, msgId, content)
	data := make(map[string]interface{})
	if cache.SeqDuplicates(msgId) {
		fmt.Println("给用户发送消息 重复提交:", msgId)
		common.RespFunc(c, http.StatusBadRequest, common.OK, data)
		return
	}
	sendResults, err := websocket.SendUserMessageAll(roomId, userId, msgId, models.MessageCmdMsg, content)
	if err != nil {
		data["sendResultsErr"] = err.Error()

	}
	data["sendResults"] = sendResults
	common.RespFunc(c, http.StatusOK, common.OK, data)
}

// Logout 退出登录
func Logout(c *gin.Context) {
	roomId := com.StrTo(c.PostForm("roomId")).String()
	userId := com.StrTo(c.PostForm("userId")).String()
	//获取房间信息
	roomOnline, err := cache.GetRoomOnlineInfo(roomId)
	if err != nil {
		common.RespFunc(c, http.StatusBadRequest, common.RoomInfoNotExist, models.Json{})
		return
	}
	if roomOnline.RoomOwner == userId {

		return
	} else {

	}
}
