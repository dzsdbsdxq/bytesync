package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"strings"
	"time"
)

type Resp struct {
	Code    uint32      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func RespFunc(c *gin.Context, httpCode int, code uint32, data interface{}) {

	c.JSON(httpCode, Resp{
		Code:    code,
		Message: GetErrorMessage(code, ""),
		Data: func(data interface{}) interface{} {
			if data == nil {
				return ""
			}
			return data
		}(data),
	})
}

func RandPass(lenNum int) string {
	var chars = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
	str := strings.Builder{}
	length := len(chars)
	rand.Seed(time.Now().UnixNano()) //重新播种，否则值不会变
	for i := 0; i < lenNum; i++ {
		str.WriteString(chars[rand.Intn(length)])

	}
	return str.String()
}

func GetUserIdByRandom() string {
	userId := uuid.Must(uuid.NewV4(), nil).String()
	h := md5.New()
	h.Write([]byte(userId + GetOrderIdTime()))
	return hex.EncodeToString(h.Sum(nil))[8:24]
}
func GetOrderIdTime() (orderId string) {
	currentTime := time.Now().Nanosecond()
	orderId = fmt.Sprintf("%d", currentTime)
	return
}
