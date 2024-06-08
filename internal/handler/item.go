package handler

import (
	"net/http"
	"pet/internal/model"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createItem(c *gin.Context) {
	var input model.Item
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
}
