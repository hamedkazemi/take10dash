module take10dashboard

go 1.13

require (
	github.com/denisenkom/go-mssqldb v0.0.0-20200206145737-bbfc9a55622e // indirect
	github.com/gin-gonic/gin v1.6.2
	github.com/guregu/null v3.4.0+incompatible
	github.com/jinzhu/gorm v1.9.12
	github.com/jinzhu/now v1.1.1 // indirect
	github.com/lib/pq v1.4.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.5.0
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14
	github.com/swaggo/gin-swagger v1.2.0
	gitlab.com/kafa1942/take10dashboard/common v0.0.0-20200425123241-974c90947003
	gitlab.com/kafa1942/take10dashboard/docs v0.0.0-20200425135542-8e92b6611041
	gitlab.com/kafa1942/take10dashboard/users v0.0.0-20200425123241-974c90947003
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/gin-gonic/gin.v1 v1.3.0
)

replace gitlab.com/kafa1942/take10dashboard/common => ./common

replace gitlab.com/kafa1942/take10dashboard/docs => ./docs

replace gitlab.com/kafa1942/take10dashboard/users => ./users
