package handlers

import (
	"time"

	"github.com/RathodViraj/sih-common/models"

	"github.com/gin-gonic/gin"
)

func SendResponse(c *gin.Context, response any, code int, message string) {
	res := models.Success{
		Response: response,
		Code:     code,
		Message:  message,
		Time:     time.Now(),
	}

	c.IndentedJSON(code, res)
}
