package users

import (
	"github.com/gin-gonic/gin"
)

func ConfigGinRouter(router *gin.RouterGroup) {
	// group routers related to module
	r := router.Group("/users")
	{
		UsersRegister(r)
		UserRegister(r)
	}
}

func UsersRegister(router gin.IRoutes) {
	router.POST("/", UsersRegistration)
	router.POST("/login", UsersLogin)
}

func UserRegister(router gin.IRoutes) {
	// authentication check
	router.Use(AuthMiddleware(true))
	router.GET("/", UserRetrieve)
	router.PUT("/", UserUpdate)
}
