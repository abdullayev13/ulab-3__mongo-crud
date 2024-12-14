package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"my_app/internal/handlers"
	"my_app/internal/pkg/mongodb"
	"time"
)

func main() {
	err := mongodb.InitDB()
	if err != nil {
		panic(err)
	}
	defer mongodb.CloseDB()
	go mongodb.InitIndies()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	api := r.Group("/api", handlers.Auth)

	products := api.Group("/products")
	orders := api.Group("/orders")

	products.POST("/create", handlers.CreateProduct)
	products.GET("/get/:id", handlers.GetOneProduct)
	products.GET("/list", handlers.GetProducts) // search query
	products.PUT("/update/:id", handlers.UpdateProduct)
	products.DELETE("/delete/:id", handlers.DeleteProduct)

	orders.POST("/create", handlers.CreateOrder)
	orders.GET("/get/:id", handlers.GetOneOrder)
	orders.GET("/list", handlers.GetOrders)
	orders.DELETE("/delete/:id", handlers.DeleteOrder)

	log.Fatalln(r.Run(":8080"))

}
