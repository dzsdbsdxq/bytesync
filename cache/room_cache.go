package cache

import (
	"bytesync/middleware"
	"bytesync/models"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
)

const (
	roomOnlinePrefix    = "acc:room:online:" // 房间在线状态
	roomOnlineCacheTime = 24 * 60 * 60
)

func getRoomOnlineKey(roomId string) (key string) {
	key = fmt.Sprintf("%s%s", roomOnlinePrefix, roomId)
	return
}
func GetRoomOnlineInfo(roomId string) (roomOnline *models.RoomOnline, err error) {
	key := getRoomOnlineKey(roomId)
	redisClient := middleware.GetClient()
	data, err := redisClient.Get(key).Bytes()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("GetRoomOnlineInfo", roomId, err)

			return
		}
		fmt.Println("GetRoomOnlineInfo", roomId, err)
		return
	}

	roomOnline = &models.RoomOnline{}
	err = json.Unmarshal(data, roomOnline)
	if err != nil {
		fmt.Println("获取房间在线数据 json Unmarshal", roomId, err)

		return
	}
	return
}

// SetRoomInfo 设置房间信息
func SetRoomInfo(roomId string, roomOnline *models.RoomOnline) (err error) {
	key := getRoomOnlineKey(roomId)
	valueByte, err := json.Marshal(roomOnline)

	if err != nil {
		fmt.Println("设置房间在线数据 json Marshal", key, err)

		return
	}
	_, err = middleware.GetClient().Do("setEx", key, roomOnlineCacheTime, string(valueByte)).Result()
	if err != nil {
		fmt.Println("设置房间在线数据 ", key, err)

		return
	}
	return
}
