package handler

import "github.com/gin-gonic/gin"

type HttpHandler interface {
	RegisterRoutes(router *gin.RouterGroup)
}
