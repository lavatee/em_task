package endpoint

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lavatee/subs/internal/service"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Endpoint struct {
	services *service.Service
	logger   *logrus.Logger
}

func NewEndpoint(services *service.Service, logger *logrus.Logger) *Endpoint {
	return &Endpoint{
		services: services,
		logger:   logger,
	}
}

func (e *Endpoint) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	api := router.Group("/api/v1")
	{
		api.GET("/subscriptions", e.GetUserSubscriptions)
		api.POST("/subscriptions", e.CreateSubscription)
		api.GET("/subscriptions/:id", e.GetSubscription)
		api.PUT("/subscriptions/:id", e.UpdateSubscription)
		api.DELETE("/subscriptions/:id", e.DeleteSubscription)
		api.GET("/subscriptions/total", e.GetTotalCost)
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
