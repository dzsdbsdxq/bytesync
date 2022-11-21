package websocket

import (
	"EditSync/models"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

var (
	clientManager = NewClientManager() //管理者
	roomIds       []string
	serverIp      string
	serverPort    string
)

func GetRoomIds() []string {
	return roomIds
}

func GetServer() (server *models.Server) {
	server = models.NewServer(serverIp, serverPort)
	return
}

func IsLocal(server *models.Server) (isLocal bool) {
	if server.Ip == serverIp && server.Port == serverPort {
		isLocal = true
	}
	return
}
func SetRoomIds(roomId string) {
	roomIds = append(roomIds, roomId)
}

func InRoomIds(roomId string) (inRoomId bool) {

	for _, value := range roomIds {
		if value == roomId {
			inRoomId = true

			return
		}
	}

	return
}

func StartWebSocket() {
	http.HandleFunc("/acc", webSocketFunc)
	//添加处理程序
	go clientManager.start()
	wsBase := fmt.Sprintf("%s:%d", viper.GetString("socket.serverIp"), viper.GetInt("socket.serverPort"))
	fmt.Println("监听socket服务", wsBase)
	http.ListenAndServe(wsBase, nil)
}

func webSocketFunc(writer http.ResponseWriter, request *http.Request) {
	// 升级协议
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		fmt.Println("升级协议", "ua:", r.Header["User-Agent"], "referer:", r.Header["Referer"])
		return true
	}}).Upgrade(writer, request, nil)
	if err != nil {
		http.NotFound(writer, request)
		return
	}
	fmt.Println("webSocket 建立连接:", conn.RemoteAddr().String())
	//roomId := request.FormValue("roomId")
	//userId := request.FormValue("userId")
	//platform := request.FormValue("platform")

	currentTime := uint64(time.Now().Unix())
	client := NewClient(conn.RemoteAddr().String(), conn, currentTime)
	fmt.Println("webSocket client:", client)

	go client.read()
	go client.write()

	// 用户连接事件
	clientManager.Register <- client
}
