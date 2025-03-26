package controller

import (
	"go-final/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Edit(router *gin.Engine) {
	router.GET("/Edit", ping)
	router.PUT("/Edit", ChangePassword)
}

// ฟังก์ชันสำหรับการเปลี่ยนรหัสผ่าน (เปรียบเทียบรหัสผ่านเก่า)
func ChangePassword(c *gin.Context) {
	var input struct {
		Email       string `json:"email"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	// รับข้อมูล JSON จาก body
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

	// ค้นหาผู้ใช้ตามอีเมล
	var customer model.Customer
	result := db.Where("email = ?", input.Email).First(&customer)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// เปรียบเทียบรหัสผ่านเก่าที่กรอกกับรหัสผ่านที่แฮชในฐานข้อมูล
	err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(input.OldPassword))
	if err != nil {
		// หากรหัสผ่านเก่าไม่ถูกต้อง
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Old password is incorrect"})
		return
	}

	// ถ้ารหัสผ่านเก่าถูกต้อง สามารถทำการเปลี่ยนรหัสผ่านใหม่ได้
	hashedNewPassword, err := HashPassword(input.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing new password"})
		return
	}

	// อัพเดตฐานข้อมูลด้วยรหัสผ่านใหม่ที่แฮชแล้ว และอัพเดตเวลา UpdatedAt
	customer.Password = hashedNewPassword
	customer.UpdatedAt = time.Now() // อัพเดตเวลาปัจจุบัน

	// บันทึกข้อมูลลงในฐานข้อมูล
	if err := db.Save(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user data"})
		return
	}

	// ตอบกลับว่าเปลี่ยนรหัสผ่านสำเร็จ
	c.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully",
		"update":  customer.UpdatedAt,
	})
}
