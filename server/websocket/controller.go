package websocket

import (
	"bytesync/cache"
	"bytesync/common"
	"bytesync/models"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

// PingController ping
func PingController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {

	code = common.OK
	fmt.Println("webSocket_request ping接口", client.Addr, seq, message)

	data = "pong"

	return
}

// LoginController 用户登录
func LoginController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {

	code = common.OK
	currentTime := uint64(time.Now().Unix())

	request := &models.Login{}
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		fmt.Println("用户登录 解析数据失败", seq, err)

		return
	}

	fmt.Println("webSocket_request 用户登录", seq, "ServiceToken", request.Token, "UserId:", request.UserId)

	// TODO::进行用户权限认证，一般是客户端传入TOKEN，然后检验TOKEN是否合法，通过TOKEN解析出来用户ID
	// 本项目只是演示，所以直接过去客户端传入的用户ID
	if request.UserId == "" || len(request.UserId) >= 20 {
		code = common.UnauthorizedUserId
		fmt.Println("用户登录 非法的用户", seq, request.UserId)

		return
	}

	if !InRoomIds(request.RoomId) {
		code = common.Unauthorized
		fmt.Println("用户登录 不支持的平台", seq, request.RoomId)

		return
	}

	if client.IsLogin() {
		fmt.Println("用户登录 用户已经登录", client.RoomId, client.UserId, seq)
		code = common.OperationFailure

		return
	}

	client.Login(request.RoomId, request.UserId, currentTime)

	// 存储数据
	userOnline := models.UserLogin(serverIp, serverPort, request.RoomId, request.UserId, client.Addr, currentTime)
	err := cache.SetUserOnlineInfo(client.GetKey(), userOnline)
	if err != nil {
		code = common.ServerError
		fmt.Println("用户登录 SetUserOnlineInfo", seq, err)

		return
	}

	// 用户登录
	login := &login{
		RoomId: request.RoomId,
		UserId: request.UserId,
		Client: client,
	}
	clientManager.Login <- login

	fmt.Println("用户登录 成功", seq, client.Addr, request.UserId)

	return
}

// HeartbeatController 心跳接口
func HeartbeatController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {

	code = common.OK
	currentTime := uint64(time.Now().Unix())
	data = models.Json{
		"userId": "",
		"roomId": "",
	}

	request := &models.HeartBeat{}
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		fmt.Println("心跳接口 解析数据失败", seq, err)

		return
	}

	fmt.Println("webSocket_request 心跳接口", client.RoomId, client.UserId)

	if !client.IsLogin() {
		fmt.Println("心跳接口 用户未登录", client.RoomId, client.UserId, seq)
		code = common.NotLoggedIn

		return
	}

	userOnline, err := cache.GetUserOnlineInfo(client.GetKey())
	if err != nil {
		if err == redis.Nil {
			code = common.NotLoggedIn
			fmt.Println("心跳接口 用户未登录", seq, client.RoomId, client.UserId)

			return
		} else {
			code = common.ServerError
			fmt.Println("心跳接口 GetUserOnlineInfo", seq, client.RoomId, client.UserId, err)

			return
		}
	}

	client.Heartbeat(currentTime)
	userOnline.Heartbeat(currentTime)
	err = cache.SetUserOnlineInfo(client.GetKey(), userOnline)
	if err != nil {
		code = common.ServerError
		fmt.Println("心跳接口 SetUserOnlineInfo", seq, client.RoomId, client.UserId, err)

		return
	}

	data = models.Json{
		"userId": client.UserId,
		"roomId": client.RoomId,
	}

	return
}
