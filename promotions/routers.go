package promotions

import (
	"github.com/gin-gonic/gin"
)

func ConfigGinRouter(router *gin.RouterGroup) {
	// game users module group
	r := router.Group("/promotions")
	{
		configGinPromotionsRouter(r)
	}

}

func configGinPromotionsRouter(router gin.IRoutes) {
	// routes
	router.POST("/", GetAll) // refactoring needed to make it standard rest
	router.PUT("/", Create)
	router.GET("/:id", Get)
	router.PATCH("/", Update)
	router.DELETE("/:id", Delete)
}
