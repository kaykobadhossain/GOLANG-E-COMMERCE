package routes

import(
	"github.com/kaykobadhossain/e-commerce/controllers"
	"github.com/gin-gonic/gin"
)


func UserRoutes(router *gin.Engine) {
    router.POST("/users/signup", controllers.Signup())
    router.POST("/users/login", controllers.Login())
    router.GET("/users/productview", controllers.SearchProduct())
    router.GET("/users/search", controllers.SearchProductByQuery())
	router.POST("/admin/addproduct", controllers.ProductViewerAdmin())
}

// func userRoutes(incomingRoutes *gin.Engine){
// 	incomingRoutes.POST("/users/signup", controllers.Signup())
// 	incomingRoutes.POST("/users/login", controllers.Login())
// 	//incomingRoutes.POST("/admin/addproduct", controllers.ProductViewerAdmin())
// 	incomingRoutes.GET("/users/productview", controllers.SearchProduct())
// 	incomingRoutes.GET("/users/search", controllers.SearchProductByQuery())
// }
	