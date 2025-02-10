package handler

import (
	"fmt"
	"github.com/Dau1to0v/fullstack-go/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (h *Handler) createWarehouse(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var input models.Warehouse
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Warehouse.Create(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	createdWarehouse := models.Warehouse{
		Id:       id,
		Name:     input.Name,
		Location: input.Location,
		UserId:   userId,
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"message":   "success",
		"warehouse": createdWarehouse,
	})

}

func (h *Handler) getAllWarehouse(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	warehouses, err := h.services.Warehouse.GetAll(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Не удалось получить склады")
		return
	}

	c.JSON(http.StatusOK, warehouses)
}

func (h *Handler) getWarehouseById(c *gin.Context) {}

func (h *Handler) updateWarehouse(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	warehouseId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}

	var input models.UpdateWarehouseInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Выполняем обновление
	if err = h.services.Warehouse.Update(userId, warehouseId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Debugf("Fetching updated warehouse for userId=%d, warehouseId=%d", userId, warehouseId)

	updatedWarehouse, err := h.services.Warehouse.GetById(userId, warehouseId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "could not retrieve updated warehouse")
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"message":   "success",
		"warehouse": updatedWarehouse,
	})
}

func (h *Handler) deleteWarehouse(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	warehouseId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}

	fmt.Println("Удаление склада userID:", userId, "warehouseID:", warehouseId)

	err = h.services.Warehouse.Delete(userId, warehouseId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "success",
	})
}

func (h *Handler) getWarehousesValue(c *gin.Context) {
	warehouseValues, err := h.services.Warehouse.CalculateWarehousesValue()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "could not calculate warehouses value")
		return
	}

	c.JSON(http.StatusOK, warehouseValues)
}
