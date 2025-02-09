package handler

import (
	"fmt"
	"github.com/Dau1to0v/fullstack-go/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) createProduct(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var input models.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Ошибка биндинга JSON: %v", err)
		newErrorResponse(c, http.StatusBadRequest, "invalid product data")
		return
	}

	log.Printf("Создание продукта: %+v", input)

	id, err := h.services.Product.Create(userId, int(input.WarehouseId), input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	createdProduct := models.Product{
		Id:          id,
		Name:        input.Name,
		Quantity:    input.Quantity,
		Price:       input.Price,
		Category:    input.Category,
		UserId:      userId,
		Description: input.Description,
		Image:       input.Image,
		WarehouseId: input.WarehouseId,
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "success",
		"product": createdProduct,
	})
}

func (h *Handler) getAllProduct(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	warehouseId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid warehouse id")
		return
	}

	products, err := h.services.Product.GetAll(userId, warehouseId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *Handler) getProductById(c *gin.Context) {}

func (h *Handler) updateProduct(c *gin.Context) {}

func (h *Handler) deleteProduct(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "user not authorized")
		return
	}

	productId, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid product id")
		return
	}

	fmt.Println("Удаление товара userID:", userId, "productID:", productId)

	err = h.services.Product.Delete(userId, productId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "success",
	})
}
