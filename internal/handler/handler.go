package handler

import (
	"pet/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{services: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.POST("/sign-up", h.signUp)
	router.GET("/sign-in", h.signIn)
	router.GET("/refresh", h.refresh)
	shop := router.Group("/shop")
	{
		shop.POST("/item", h.createItem)
		shop.GET("/item", h.findItem)
		shop.DELETE("/deleteitem", h.deleteItem)
		shop.PATCH("/updateitem", h.updateItem)
	}
	return router
}
