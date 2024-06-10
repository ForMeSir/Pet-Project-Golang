package handler

import (
	"net/http"
	"pet/internal/model"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createItem(c *gin.Context) {
	var input model.Item
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token := c.GetHeader("token")

	parcedtoken, err := h.services.Authorization.ParseToken(token)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error()+" Невалидный токен")
		return
	}

	if parcedtoken.UserRole != "admin" {
		newErrorResponse(c, http.StatusBadRequest, "Нет доступа")
		return
	}

	if time.Now().Unix() >= parcedtoken.ExpiratedAt.Unix() {
		newErrorResponse(c, http.StatusBadRequest, "Обновите токен")
		return
	}

	id, err := h.services.CreateItem(input)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type Find struct{
	Title string `json:"title"`
}

func (h *Handler) findItem(c *gin.Context) {
   var input Find
	 if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	} 

	items,err := h.services.Item.FindItemByTitle(input.Title)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	} 
	
	c.JSON(http.StatusOK, map[string]interface{}{
		"items": items,
	})
}
