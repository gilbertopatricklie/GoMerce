package router

import (
	"gomerce/internal/handler"
	"gomerce/internal/middleware"
	"log/slog"

	"github.com/gin-gonic/gin"
)

//This function returns engine so later wherever this func is imported, using engine method is possible
func SetupRouter(logger *slog.Logger) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	//Custom middleware from internal/middleware/request_id.go
	router.Use(middleware.RequestID())
	//custom middleware from internal/middleware/access_log.go
	//AccessLog after RequestID is because AccessLog requires request_id, generated from RequestID
	router.Use(middleware.AccessLog(logger))

	healthHandler := handler.NewHealthHandler()
	productHandler := handler.NewProductHandler()

	router.GET("/health", healthHandler.Check)
	router.POST("/products", productHandler.Create)

	return router
}