package handler

import (
	"github.com/Dau1to0v/fullstack-go/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) signUp(c *gin.Context) {
	var user models.User

	logrus.Info("Получен запрос на регистрацию")

	if err := c.BindJSON(&user); err != nil {
		logrus.Error("Ошибка при разборе JSON: ", err)
		newErrorResponse(c, http.StatusBadRequest, "Invalid input data")
		return
	}

	logrus.Info("Данные пользователя успешно получены: ", user)

	id, err := h.services.Authorization.CreateUser(user)
	if err != nil {
		logrus.Error("Ошибка при создании пользователя: ", err)
		newErrorResponse(c, http.StatusInternalServerError, "Could not create user")
		return
	}

	logrus.Info("Пользователь успешно создан с ID: ", id)

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var user signInInput

	if err := c.BindJSON(&user); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Authorization.GenerateToken(user.Username, user.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
