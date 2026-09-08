package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a registered user
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Email     string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	Phone     string         `gorm:"size:20" json:"phone"`
	Role      string         `gorm:"size:20;default:'user'" json:"role"` // user or admin
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Orders    []Order        `json:"orders,omitempty"`
	Cart      *Cart          `json:"cart,omitempty"`
}

// Category for products
type Category struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;uniqueIndex;not null" json:"name"`
	Slug        string         `gorm:"size:100;uniqueIndex" json:"slug"`
	Description string         `gorm:"size:500" json:"description"`
	ImageURL    string         `gorm:"size:500" json:"image_url"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Products    []Product      `json:"products,omitempty"`
}

// Product
type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:200;not null" json:"name"`
	Slug        string         `gorm:"size:200;uniqueIndex" json:"slug"`
	Description string         `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"not null" json:"price"`
	CompareAt   float64        `json:"compare_at,omitempty"` // original price for discounts
	Stock       int            `gorm:"default:0" json:"stock"`
	ImageURL    string         `gorm:"size:500" json:"image_url"`
	CategoryID  uint           `json:"category_id"`
	Category    Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Rating      float64        `gorm:"default:0" json:"rating"`
	NumReviews  int            `gorm:"default:0" json:"num_reviews"`
	IsFeatured  bool           `gorm:"default:false" json:"is_featured"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// Cart for a user
type Cart struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"uniqueIndex;not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"-"`
	Items     []CartItem     `json:"items"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// CartItem
type CartItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CartID    uint      `gorm:"not null" json:"cart_id"`
	ProductID uint      `gorm:"not null" json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product"`
	Quantity  int       `gorm:"not null;default:1" json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Order
type Order struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	OrderNumber   string         `gorm:"size:50;uniqueIndex;not null" json:"order_number"`
	UserID        uint           `gorm:"not null" json:"user_id"`
	User          User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Status        string         `gorm:"size:30;default:'pending'" json:"status"`
	PaymentMethod string         `gorm:"size:50" json:"payment_method"`
	PaymentStatus string         `gorm:"size:30;default:'pending'" json:"payment_status"`
	Subtotal      float64        `json:"subtotal"`
	ShippingCost  float64        `json:"shipping_cost"`
	Total         float64        `json:"total"`
	CustomerName  string         `gorm:"size:100" json:"customer_name"`
	Email         string         `gorm:"size:100" json:"email"`
	Phone         string         `gorm:"size:20" json:"phone"`
	Address       string         `gorm:"size:300" json:"address"`
	City          string         `gorm:"size:100" json:"city"`
	State         string         `gorm:"size:100" json:"state"`
	PostalCode    string         `gorm:"size:20" json:"postal_code"`
	Items         []OrderItem    `json:"items"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// OrderItem
type OrderItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	OrderID   uint    `gorm:"not null" json:"order_id"`
	ProductID uint    `gorm:"not null" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Subtotal  float64 `json:"subtotal"`
}

// LoginRequest
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterRequest
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone"`
}

// AuthResponse
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// ProductFilter for search/filter
type ProductFilter struct {
	Search     string  `form:"search"`
	CategoryID uint    `form:"category_id"`
	MinPrice   float64 `form:"min_price"`
	MaxPrice   float64 `form:"max_price"`
	SortBy     string  `form:"sort_by"`
	Page       int     `form:"page"`
	Limit      int     `form:"limit"`
}

// CheckoutRequest
type CheckoutRequest struct {
	CustomerName  string `json:"customer_name" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	Phone         string `json:"phone" binding:"required"`
	Address       string `json:"address" binding:"required"`
	City          string `json:"city" binding:"required"`
	State         string `json:"state" binding:"required"`
	PostalCode    string `json:"postal_code" binding:"required"`
	PaymentMethod string `json:"payment_method" binding:"required"`
}

// AddToCartRequest
type AddToCartRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

// UpdateCartItemRequest
type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=0"`
}

// UpdateOrderStatusRequest
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// CreateProductRequest
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	CompareAt   float64 `json:"compare_at"`
	Stock       int     `json:"stock" binding:"gte=0"`
	ImageURL    string  `json:"image_url"`
	CategoryID  uint    `json:"category_id" binding:"required"`
	IsFeatured  bool    `json:"is_featured"`
	Rating      float64 `json:"rating"`
}

// UpdateProductRequest
type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CompareAt   float64 `json:"compare_at"`
	Stock       int     `json:"stock"`
	ImageURL    string  `json:"image_url"`
	CategoryID  uint    `json:"category_id"`
	IsFeatured  bool    `json:"is_featured"`
	IsActive    bool    `json:"is_active"`
	Rating      float64 `json:"rating"`
}
