package controller

import (
	"net/http"
	"time"
	"go-final/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)
func Cart(router *gin.Engine) {
	router.GET("/cart", ping)
	router.POST("/cart", AddToCart)
}
// AddToCart เพิ่มสินค้าลงในรถเข็น
func AddToCart(c *gin.Context) {
	var input struct {
		CustomerID int    `json:"customer_id"`
		CartName   string `json:"cart_name"`
		ProductID  int    `json:"product_id"`
		Quantity   int    `json:"quantity"`
	}

	// รับข้อมูลจาก body (JSON)
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

	// ค้นหารถเข็นตามชื่อและลูกค้า
	var cart model.Cart
	result := db.Where("customer_id = ? AND cart_name = ?", input.CustomerID, input.CartName).First(&cart)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// ถ้าไม่พบรถเข็น ให้สร้างใหม่
			cart = model.Cart{
				CustomerID: input.CustomerID,
				CartName:   input.CartName,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}
			if err := db.Create(&cart).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
	}

	// ค้นหาสินค้าในรถเข็น
	var cartItem model.CartItem
	result = db.Where("cart_id = ? AND product_id = ?", cart.CartID, input.ProductID).First(&cartItem)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// ถ้าไม่มีสินค้านี้ในรถเข็น ให้เพิ่มใหม่
			cartItem = model.CartItem{
				CartID:    cart.CartID,
				ProductID: input.ProductID,
				Quantity:  input.Quantity,
				UpdatedAt: time.Now(),
			}
			if err := db.Create(&cartItem).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add product to cart"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
	} else {
		// ถ้ามีสินค้านี้อยู่แล้ว ให้เพิ่มจำนวน
		cartItem.Quantity += input.Quantity
		cartItem.UpdatedAt = time.Now()
		if err := db.Save(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product quantity"})
			return
		}
	}

	// ตอบกลับสำเร็จ
	c.JSON(http.StatusOK, gin.H{"message": "Product added to cart successfully !!"})
}