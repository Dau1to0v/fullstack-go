package handler

import (
	"github.com/Dau1to0v/fullstack-go/pkg/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Разрешить все источники
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // Кэширование preflight-запросов на 12 часов
	}))

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.signUp)
			auth.POST("/login", h.signIn)
			auth.GET("/getMe", h.userIdentity, h.getMe)
			auth.POST("/updateUser", h.userIdentity, h.updateUser)
		}

		protected := api.Group("/", h.userIdentity)
		{
			warehouse := protected.Group("/warehouse")
			{
				warehouse.POST("/create", h.createWarehouse)
				warehouse.GET("/getAll", h.getAllWarehouse)
				warehouse.GET("/", h.getWarehouseById)
				warehouse.POST("/update/:id", h.updateWarehouse)
				warehouse.POST("/delete/:id", h.deleteWarehouse)
			}

			product := protected.Group("/product")
			{
				product.POST("/create", h.createProduct)
				product.GET("/getAll/:warehouse_id", h.getAllProduct)
				product.GET("/get/:product_id", h.getProductById)
				product.POST("/update/:product_id", h.updateProduct)
				product.POST("/delete/:product_id", h.deleteProduct)

			}
		}
	}

	return router
}
