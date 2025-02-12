package models

import (
	"errors"
	"time"
)

type User struct {
	Id        int       `json:"_id" db:"id"`
	Username  string    `json:"username" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	Password  string    `json:"password" binding:"required"`
	CreatedAt time.Time `json:"-" db:"created_at"`
}

type Warehouse struct {
	Id        int       `json:"_id" db:"id"`
	Name      string    `json:"name" db:"name" binding:"required"`
	Location  string    `json:"location" db:"location" binding:"required"`
	UserId    int       `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"-" db:"created_at"`
}

type Product struct {
	Id          int         `json:"_id" db:"id"`
	Name        string      `json:"name" binding:"required" db:"name"`
	Quantity    int         `json:"quantity" binding:"required" db:"quantity"`
	Price       float64     `json:"price" binding:"required" db:"price"`
	Category    string      `json:"category" binding:"required" db:"category"`
	Description string      `json:"description" binding:"required" db:"description"`
	Image       string      `json:"image" binding:"required" db:"image"`
	UserId      int         `json:"user_id" db:"user_id"`
	WarehouseId WarehouseId `json:"warehouse_id" db:"warehouse_id"`
	CreatedAt   time.Time   `json:"-" db:"created_at"`
}

type UpdateWarehouseInput struct {
	Id       int     `json:"-"`
	Name     *string `json:"name" db:"name"`
	Location *string `json:"location" db:"location"`
}

func (i UpdateWarehouseInput) Validate() error {
	if i.Name == nil && i.Location == nil {
		return errors.New("update structures has no values")
	}

	return nil
}

type UpdateProductInput struct {
	Name        *string  `json:"name" db:"name"`
	Quantity    *int     `json:"quantity" db:"quantity"`
	Price       *float64 `json:"price" db:"price"`
	Category    *string  `json:"category" db:"category"`
	Description *string  `json:"description" db:"description"`
	Image       *string  `json:"image" db:"image"`
}

func (i UpdateProductInput) Validate() error {
	if i.Name == nil && i.Quantity == nil && i.Price == nil && i.Category == nil && i.Description == nil && i.Image == nil {
		return errors.New("update structures has no values")
	}

	return nil
}

type UpdateUserInput struct {
	Username *string `json:"username" binding:"required"`
	Email    *string `json:"email" binding:"required,email"`
}

func (i UpdateUserInput) Validate() error {
	if i.Username == nil && i.Email == nil {
		return errors.New("update structures has no values")
	}

	return nil
}

type WarehouseNetWorth struct {
	Id        int     `json:"_id" db:"warehouse_id"`
	Location  string  `json:"location" db:"warehouse_location"`
	Warehouse string  `json:"warehouse" db:"warehouse_name"`
	NetWorth  float64 `json:"netWorth" db:"net_worth"`
}
