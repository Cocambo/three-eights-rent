package handler

import (
	httpresponse "car-service/internal/response"

	"github.com/gin-gonic/gin"
)

func writeError(c *gin.Context, err error) {
	httpresponse.WriteError(c, err)
}

func writeSuccess(c *gin.Context, status int, data any) {
	httpresponse.WriteSuccess(c, status, data)
}
