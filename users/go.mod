module gitlab.com/kafa1942/take10dashboard/users

go 1.13

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.6.2
	github.com/jinzhu/gorm v1.9.12
	gitlab.com/kafa1942/take10dashboard/common v0.0.0-20200425123241-974c90947003
	golang.org/x/crypto v0.0.0-20200423211502-4bdfaf469ed5
)

replace gitlab.com/kafa1942/take10dashboard/common => ../common
