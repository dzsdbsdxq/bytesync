package models

type RoomOnline struct {
	RoomId     string //房间ID
	RoomTitle  string // 房间标题
	RoomDesc   string // 房间信息
	RoomOwner  string // 房主
	CreateTime uint32
}
