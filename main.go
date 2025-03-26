package main

import (
	"go-final/controller"
)

func main() {
	// viper.SetConfigName("config")
	// viper.AddConfigPath(".")
	// err := viper.ReadInConfig()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(viper.Get("mysql.dsn"))
	// dsn := viper.GetString("mysql.dsn")

	// dialactor := mysql.Open(dsn)
	// db, err := gorm.Open(dialactor)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Connection successful")

	// customer := []model.Cart{}
	// result := db.Preload("CustomerData").Find(&customer)
	// if result.Error != nil {
	// 	panic(result.Error)
	// }
	// fmt.Println(customer)

	// Release mode
	//gin.SetMode(gin.ReleaseMode)
	controller.StartServer()
}
