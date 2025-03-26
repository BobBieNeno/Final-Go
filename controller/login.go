package controller

import (
	"fmt"
	"go-final/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Login(router *gin.Engine) {
	router.GET("/login", ping)
	router.POST("/login", LoginUser)
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "login!!",
	})
}

// ฟังก์ชันในการเชื่อมต่อกับฐานข้อมูล
func connectToDatabase() (*gorm.DB, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	dsn := viper.GetString("mysql.dsn")
	dialactor := mysql.Open(dsn)
	db, err := gorm.Open(dialactor, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}
	return db, nil
}

func LoginUser(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// รับ JSON จาก body
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	// input.Email = "somchai@example.com"
	// input.Password = "password123"

	// เชื่อมต่อกับฐานข้อมูล
	db, err := connectToDatabase()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ค้นหาข้อมูลลูกค้าตาม email
	var customerData model.Customer
	if err := db.Where("email = ?", input.Email).First(&customerData).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// // ตรวจสอบรหัสผ่าน (เปรียบเทียบแบบปกติ ไม่ใช้ Hashing)
	// if customerData.Password != input.Password {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
	// 	return
	// }

	err = bcrypt.CompareHashAndPassword([]byte(customerData.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// ล็อกอินสำเร็จ
	c.JSON(http.StatusOK, gin.H{
		"customer_id":  customerData.CustomerID,
		"first_name":   customerData.FirstName,
		"last_name":    customerData.LastName,
		"email":        customerData.Email,
		"phone_number": customerData.PhoneNumber,
		"address":      customerData.Address,
		"created_at":   customerData.CreatedAt,
		"updated_at":   customerData.UpdatedAt,
	})
}

// ฟังก์ชันสำหรับการเข้ารหัสรหัสผ่าน
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
