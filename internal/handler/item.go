package handler

import (
	"net/http"
	"pet/internal/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ИЗМЕНИТЬ В FIND ITEM  выходные данные на массив
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
		newErrorResponse(c, 403, "Нет доступа")
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

type Find struct {
	Title  string `json:"title"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

func (h *Handler) findItem(c *gin.Context) {
	var input Find
	var err error
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	limit, ok := c.GetQuery("limit")
	if ok {
		input.Limit, err = strconv.Atoi(limit)
		if err != nil {
			newErrorResponse(c, 400, "Некорректные данные")
			return
		}
	}

	offset, ok := c.GetQuery("offset")
	if ok {
		input.Offset, err = strconv.Atoi(offset)
		if err != nil {
			newErrorResponse(c, 400, "Некорректные данные")
			return
		}
	}

	items, err := h.services.Item.FindItemByTitle(input.Title, input.Limit, input.Offset)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"items": items,
	})
}

type DeleteItem struct {
	Id uuid.UUID `json:"id"`
}

func (h *Handler) deleteItem(c *gin.Context) {
	var input DeleteItem
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
		newErrorResponse(c, 403, "Нет доступа")
		return
	}

	if time.Now().Unix() >= parcedtoken.ExpiratedAt.Unix() {
		newErrorResponse(c, http.StatusBadRequest, "Обновите токен")
		return
	}

	err = h.services.Item.DeleteItem(input.Id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status" : "ok",
	})
}

func (h *Handler) updateItem(c *gin.Context) {
	
}