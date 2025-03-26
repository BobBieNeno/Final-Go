package controller

import (
	"go-final/model"
	"net/http"
	"github.com/gin-gonic/gin"
)

// ฟังก์ชันในการเพิ่มสินค้าลงในรถเข็น
func AddToCart(c *gin.Context) {
	var input struct {
		Email     string  `json:"email"`
		CartName  string  `json:"cart_name"`
		ProductID int     `json:"product_id"`
		Quantity  int     `json:"quantity"`
		MinPrice  float64 `json:"min_price"`
		MaxPrice  float64 `json:"max_price"`
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

	// ค้นหาสินค้าในฐานข้อมูลที่ตรงกับช่วงราคาที่ลูกค้าต้องการ
	var products []model.Product
	result := db.Where("price BETWEEN ? AND ?", input.MinPrice, input.MaxPrice).Find(&products)
	if result.Error != nil || len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No products found in the specified price range"})
		return
	}

	// ค้นหาลูกค้าโดยใช้ email
	var customer model.Customer
	result = db.Where("email = ?", input.Email).First(&customer)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	// ค้นหารถเข็นที่ลูกค้าต้องการ
	var cart model.Cart
	result = db.Where("customer_id = ? AND cart_name = ?", customer.CustomerID, input.CartName).First(&cart)
	if result.Error != nil {
		// ถ้ารถเข็นไม่พบ ให้สร้างรถเข็นใหม่
		cart = model.Cart{
			CustomerID: customer.CustomerID,
			CartName:   input.CartName,
		}
		if err := db.Create(&cart).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
			return
		}
	}

	// ค้นหาสินค้าในรถเข็น
	var cartItem model.CartItem
	result = db.Where("cart_id = ? AND product_id = ?", cart.CartID, input.ProductID).First(&cartItem)
	if result.Error != nil {
		// ถ้าไม่พบสินค้าในรถเข็น ให้เพิ่มสินค้าใหม่
		cartItem = model.CartItem{
			CartID:    cart.CartID,
			ProductID: input.ProductID,
			Quantity:  input.Quantity,
		}
		if err := db.Create(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add product to cart"})
			return
		}
	} else {
		// ถ้ามีสินค้าในรถเข็นแล้ว ให้เพิ่มจำนวน
		cartItem.Quantity += input.Quantity
		if err := db.Save(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product quantity in cart"})
			return
		}
	}

	// ตอบกลับการเพิ่มสินค้าในรถเข็น
	c.JSON(http.StatusOK, gin.H{"message": "Product added to cart successfully"})
}
