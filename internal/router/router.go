package router

import (
	"log/slog"
	"time"

	"gomerce/internal/handler"
	"gomerce/internal/middleware"
	"gomerce/internal/response"
	"gomerce/internal/security"
	"gomerce/internal/service"

	"github.com/gin-gonic/gin"
)

//This function returns engine so later wherever this func is imported, using engine method is possible
func SetupRouter(logger *slog.Logger, productService *service.ProductService,
				 categoryService *service.CategoryService, userService *service.UserService,
				 tokenManager *security.TokenManager, refreshTTL time.Duration) *gin.Engine {

	router := gin.New()

	router.HandleMethodNotAllowed = true

	router.Use(gin.Recovery())
	//Custom middleware from internal/middleware/request_id.go
	router.Use(middleware.RequestID())
	//custom middleware from internal/middleware/access_log.go
	//AccessLog after RequestID is because AccessLog requires request_id, generated from RequestID
	router.Use(middleware.AccessLog(logger))
	//Handles if no Route requested exists
	router.NoRoute(func(c *gin.Context) {
		response.Fail(c, 404, "ROUTE_NOT_FOUND", "Requested Endpoint Does Not Exist", nil)
	})
	//Handles if Method not allowed on certain endpoint
	router.NoMethod(func(c *gin.Context) {
		response.Fail(c, 405, "METHOD_NOT_ALLOWED", "HTTP Method is Not Allowed for This Endpoint", nil)
	})

	//Routes Handlers
	healthHandler := handler.NewHealthHandler()
	productHandler := handler.NewProductHandler(productService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	authHandler := handler.NewAuthHandler(userService, tokenManager, refreshTTL)

	//Health Endpoint
	router.GET("/health", healthHandler.Check)
	//Product Endpoint
	router.GET("/products", productHandler.GetAll)
	router.GET("/products/:id", productHandler.GetByID)
	//Admin products exclusive Endpoint
	adminProducts := router.Group("/products")
	adminProducts.Use(middleware.AuthRequired(tokenManager), middleware.RequireRole("admin"))
	adminProducts.POST("", productHandler.Create)
	adminProducts.PUT("/:id", productHandler.Update)
	adminProducts.DELETE("/:id", productHandler.Delete)
	//Category Endpoint
	router.GET("/categories", categoryHandler.GetAll)
	//Admin categories exclusive Endpoint
	adminCategories := router.Group("/categories")
	adminCategories.Use(middleware.AuthRequired(tokenManager), middleware.RequireRole("admin"))
	router.POST("", categoryHandler.Create)
	router.PUT("/:id", categoryHandler.Update)
	router.DELETE("/:id", categoryHandler.Delete)
	//User Endpoint
	router.POST("/auth/signup", authHandler.SignUp)
	router.POST("/auth/signin", authHandler.SignIn)
	router.POST("/auth/refresh", authHandler.Refresh)

	return router
}