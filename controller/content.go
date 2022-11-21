package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

func SaveDoc(c *gin.Context) {
	content := com.StrTo(c.PostForm("content")).String()

	fmt.Println(content)
}
