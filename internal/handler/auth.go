package handler

import (
	"net/http"
	"pet/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type signUpInput struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signUp(c *gin.Context) {
	var input signUpInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(input.Name, input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.services.Authorization.GetUser(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Неверный логин или пароль")
		return
	}

	logrus.Info(user)

	ref, ac, err := h.services.GenerateToken(user.ID, user.Role)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"Accesstoken":       ac,
		"AccessExpiratedAt": time.Now().Add(service.AcTimeLive),
		"Refreshtoken":      ref,
		"RefExpiratedAt":    time.Now().Add(service.RefTimeLive),
	})
}

type refreshInput struct {
	Refreshtoken string `json:"Refreshtoken"`
	Accesstoken  string `json:"Accesstoken"`
}

//Доделать refresh

func (h *Handler) refresh(c *gin.Context) {
	var input refreshInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	Accesstoken, err := h.services.Authorization.Refresh(input.Refreshtoken, input.Accesstoken)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"Accesstoken":       Accesstoken,
		"AccessExpiratedAt": time.Now().Add(service.AcTimeLive),
	})

}
