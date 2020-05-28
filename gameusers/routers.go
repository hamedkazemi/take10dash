package gameusers

import (
	"github.com/gin-gonic/gin"
)

func ConfigGinRouter(router *gin.RouterGroup) {
	// game users module group
	r := router.Group("/gameUsers")
	{
		configGinGameUsersRouter(r)
	}

}

func configGinGameUsersRouter(router gin.IRoutes) {
	// routes
	router.POST("/", GetAll)
	router.GET("/:id", Get)
	router.GET("/:id/profile", GetProfile)
	router.PATCH("/", Update)
	router.DELETE("/:id", Delete)
}
