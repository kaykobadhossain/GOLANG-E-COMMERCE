package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kaykobadhossain/e-commerce/controllers"
	"github.com/kaykobadhossain/e-commerce/database"
	"github.com/kaykobadhossain/e-commerce/middlewares"
	"github.com/kaykobadhossain/e-commerce/routes"
)

func main(){
	port := os.Getenv("PORT")
	if port ==""{
		port="8000"
	}

	app:= controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	router := gin.new()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.use(middlewares.Authentication())

	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem",app.RemoveItem())
	router.GET("/cartcheckout",app.BuyFromCart())
	router.GET("/instantbuy",app.InsatantBuy())

	log.Fatal(router.Run(":"+port))
}