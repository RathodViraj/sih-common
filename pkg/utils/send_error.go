package handlers

import (
	"time"

	"github.com/RathodViraj/sih-common/models"

	"github.com/gin-gonic/gin"
)

func SendError(c *gin.Context, data any, code int, message string) {
	err := models.Error{
		Data:    data,
		Code:    code,
		Message: message,
		Time:    time.Now(),
	}
	c.IndentedJSON(code, err)
}
