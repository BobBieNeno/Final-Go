package model

import "time"

type Person struct {
	ID          int
	FirstName   string
	LastName    string
	Age         int `json:"-"`
	Email       string
	Password    string  `json:"-"`
	PostAddress Address `json:"Address"`
}
type User struct {
	CustomerID  int64     `json:"customer_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Address     string    `json:"address"`
	Password    string    `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
type Address struct {
	HouseNo string
	City    string
}