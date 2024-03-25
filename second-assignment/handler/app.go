package handler

import (
	"second-assignment/docs"
	"second-assignment/infra/config"
	"second-assignment/infra/database"
	"second-assignment/pkg/errs"
	"second-assignment/repository/item_repository/item_pg"
	"second-assignment/repository/order_repository/order_pg"
	"second-assignment/service/order_service"

	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		ctx.Next()
	}
}

func OrderAuthorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		orderId := ctx.Param("orderId")

		if orderId == "2" {
			forbiddenAccessErr := errs.NewUnauthorizedError("forbidden")
			ctx.AbortWithStatusJSON(forbiddenAccessErr.Status(), forbiddenAccessErr)
			return
		}

		ctx.Next()

	}
}

func StartApp() {
	config.LoadAppConfig()
	database.InitiliazeDatabase()

	db := database.GetDatabaseInstance()

	orderRepo := order_pg.NewRepository(db)

	itemRepo := item_pg.NewRepository(db)

	orderService := order_service.NewService(orderRepo, itemRepo)

	orderHandler := NewOrderHandler(orderService)

	r := gin.Default()

	docs.SwaggerInfo.Title = "H8 Assignment 2"
	docs.SwaggerInfo.Description = "Ini adalah tugas ke 2 dari kelas Kominfo"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.Schemes = []string{"http"}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.POST("/orders", orderHandler.CreateOrder)

	r.GET("/orders", Middleware(), orderHandler.GetOrders)

	r.PUT("/orders/:orderId", OrderAuthorization(), orderHandler.UpdateOrder)

	r.DELETE("/orders/:orderId", OrderAuthorization(), orderHandler.DeleteOrder)

	r.Run(":8080")
}
