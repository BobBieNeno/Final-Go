package controller

import "github.com/gin-gonic/gin"

func StartServer() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Api is now working",
		})
	})
	Login(router)
	Edit(router)
	Cart(router)
	CartDetail(router)
	router.Run()
}
