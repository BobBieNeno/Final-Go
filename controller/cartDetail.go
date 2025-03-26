package controller

import (
	"go-final/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CartDetail(router *gin.Engine) {
	router.GET("/cartt", ping)
	router.POST("/cartdetail", GetCarts)
}

// ฟังก์ชันแสดงข้อมูลรถเข็นทั้งหมดของลูกค้า
func GetCarts(c *gin.Context) {
	var input struct {
		CustomerID int `json:"customer_id"`
	}

	// รับข้อมูลจาก body (CustomerID)
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// เชื่อมต่อกับฐานข้อมูล
	db, err := connectToDatabase()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ค้นหารถเข็นทั้งหมดที่เป็นของลูกค้า
	var carts []model.Cart
	result := db.Where("customer_id = ?", input.CustomerID).Find(&carts)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve carts"})
		return
	}

	// หากไม่มีรถเข็นของลูกค้า
	if len(carts) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No carts found for this customer"})
		return
	}

	// แสดงข้อมูลรายละเอียดของสินค้าภายในรถเข็นแต่ละคัน
	var cartDetails []gin.H
	for _, cart := range carts {
		// ค้นหาสินค้าที่อยู่ในรถเข็น
		var cartItems []model.CartItem
		result := db.Where("cart_id = ?", cart.CartID).Find(&cartItems)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart items"})
			return
		}

		// รายละเอียดของแต่ละสินค้าภายในรถเข็น
		var itemsDetails []gin.H
		for _, item := range cartItems {
			var product model.Product
			// ค้นหาสินค้า
			result := db.Where("id = ?", item.ProductID).First(&product)
			if result.Error != nil {
				continue // ถ้าสินค้าไม่พบให้ข้าม
			}

			// คำนวณราคาของแต่ละรายการ
			itemDetail := gin.H{
				"product_name": product.ProductName,
				"quantity":     item.Quantity,
				"price":        product.Price,
				"total_price":  float64(item.Quantity) * product.Price,
			}
			itemsDetails = append(itemsDetails, itemDetail)
		}

		// ข้อมูลรถเข็นแต่ละคัน
		cartDetail := gin.H{
			"cart_name": cart.CartName,
			"items":     itemsDetails,
		}
		cartDetails = append(cartDetails, cartDetail)
	}

	// ส่งข้อมูลทั้งหมดกลับไปยังผู้ใช้
	c.JSON(http.StatusOK, gin.H{
		"customer_id": input.CustomerID,
		"carts":       cartDetails,
	})
}
