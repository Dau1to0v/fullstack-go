package handler

import (
	"fmt"
	"github.com/Dau1to0v/fullstack-go/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
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

	warehouseId, err := strconv.Atoi(c.Param("warehouse_id"))
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

func (h *Handler) getProductById(c *gin.Context) {
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

	var product models.Product
	product, err = h.services.Product.GetById(userId, productId)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, product)

}

func (h *Handler) updateProduct(c *gin.Context) {
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

	var input models.UpdateProductInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Выполняем обновление
	if err = h.services.Product.Update(userId, productId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Debugf("Fetching updated product for userId=%d, productId=%d", userId, productId)

	updatedProduct, err := h.services.Product.GetById(userId, productId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "could not retrieve updated product")
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"message":   "success",
		"warehouse": updatedProduct,
	})
}

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

func (h *Handler) searchProduct(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "user not authorized")
		return
	}

	warehouseId, err := strconv.Atoi(c.Query("warehouse_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid warehouse_id")
		return
	}

	text := c.Query("text")
	searchType := c.Query("type") // Может быть "name" или "category"

	products, err := h.services.Product.Search(userId, warehouseId, text, searchType)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *Handler) getExelFile(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "user not authorized")
		return
	}

	warehouseId, err := strconv.Atoi(c.Param("warehouse_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid warehouse ID"})
		return
	}

	// Получаем товары со склада
	products, err := h.services.Product.GetAll(userId, warehouseId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get products"})
		return
	}

	// Логируем количество полученных товаров
	logrus.Infof("Получено %d товаров для склада %d", len(products), warehouseId)

	f := excelize.NewFile()
	sheetName := "Products"
	f.SetSheetName("Sheet1", sheetName)

	// Устанавливаем активный лист
	index, err := f.GetSheetIndex(sheetName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get sheet index"})
		return
	}
	f.SetActiveSheet(index)

	// Заголовки таблицы
	headers := []string{"Название", "Количество", "Цена", "Категория"}
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		f.SetCellValue(sheetName, cell, header)
	}

	// Записываем данные о продуктах
	for i, product := range products {
		row := strconv.Itoa(i + 2) // Начинаем со 2-й строки
		f.SetCellValue(sheetName, "A"+row, product.Name)
		f.SetCellValue(sheetName, "B"+row, product.Quantity)
		f.SetCellValue(sheetName, "C"+row, product.Price)
		f.SetCellValue(sheetName, "D"+row, product.Category)
	}

	// Проверяем, записались ли данные
	logrus.Info("Файл успешно создан, отправляется пользователю")

	// Удаляем пустой стандартный лист
	f.DeleteSheet("Sheet1")

	// Устанавливаем заголовки для скачивания
	c.Header("Content-Disposition", "attachment; filename=products.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Access-Control-Expose-Headers", "Content-Disposition")

	// Записываем файл в HTTP-ответ
	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate file"})
	}
}
