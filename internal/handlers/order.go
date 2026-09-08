package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/balaji/ecommerce/internal/database"
	"github.com/balaji/ecommerce/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderHandler struct{}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}

func generateOrderNumber() string {
	return fmt.Sprintf("ORD-%d-%d", time.Now().Unix(), time.Now().Nanosecond()%10000)
}

func (h *OrderHandler) Checkout(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req models.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate payment method
	if req.PaymentMethod != "cod" && req.PaymentMethod != "mock_online" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment method. Use 'cod' or 'mock_online'"})
		return
	}

	// Get cart
	var cart models.Cart
	if err := database.DB.Where("user_id = ?", userID).Preload("Items.Product").First(&cart).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}
	if len(cart.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}

	// Validate stock and calculate totals
	var subtotal float64
	for _, item := range cart.Items {
		if item.Quantity > item.Product.Stock {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":    fmt.Sprintf("Insufficient stock for %s", item.Product.Name),
				"available": item.Product.Stock,
			})
			return
		}
		subtotal += item.Product.Price * float64(item.Quantity)
	}

	shippingCost := 0.0
	if subtotal < 999 {
		shippingCost = 99.0 // Free shipping above 999
	}
	total := subtotal + shippingCost

	// Create order in transaction
	var order models.Order
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		order = models.Order{
			OrderNumber:   generateOrderNumber(),
			UserID:        userID.(uint),
			Status:        "pending",
			PaymentMethod: req.PaymentMethod,
			PaymentStatus: "pending",
			Subtotal:      subtotal,
			ShippingCost:  shippingCost,
			Total:         total,
			CustomerName:  req.CustomerName,
			Email:         req.Email,
			Phone:         req.Phone,
			Address:       req.Address,
			City:          req.City,
			State:         req.State,
			PostalCode:    req.PostalCode,
		}

		if req.PaymentMethod == "mock_online" {
			// Mock successful payment
			order.PaymentStatus = "paid"
			order.Status = "confirmed"
		}

		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		// Create order items and reduce stock
		for _, item := range cart.Items {
			orderItem := models.OrderItem{
				OrderID:   order.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Product.Price,
				Subtotal:  item.Product.Price * float64(item.Quantity),
			}
			if err := tx.Create(&orderItem).Error; err != nil {
				return err
			}

			// Reduce stock
			if err := tx.Model(&models.Product{}).Where("id = ?", item.ProductID).
				Update("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
				return err
			}
		}

		// Clear cart
		if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order: " + err.Error()})
		return
	}

	// Load full order
	database.DB.Preload("Items.Product").First(&order, order.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Order placed successfully",
		"order":   order,
	})
}

func (h *OrderHandler) MyOrders(c *gin.Context) {
	userID, _ := c.Get("userID")
	var orders []models.Order
	if err := database.DB.Where("user_id = ?", userID).
		Preload("Items.Product").
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var order models.Order
	query := database.DB.Preload("Items.Product").Preload("User")
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}
	if err := query.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

// Admin: List all orders
func (h *OrderHandler) AdminList(c *gin.Context) {
	var orders []models.Order
	if err := database.DB.Preload("Items.Product").Preload("User").
		Order("created_at DESC").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

// Admin: Update order status
func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var req models.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validStatuses := map[string]bool{
		"pending": true, "confirmed": true, "shipped": true,
		"delivered": true, "cancelled": true,
	}
	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// If cancelling, restore stock
	if req.Status == "cancelled" && order.Status != "cancelled" {
		var items []models.OrderItem
		database.DB.Where("order_id = ?", order.ID).Find(&items)
		for _, item := range items {
			database.DB.Model(&models.Product{}).Where("id = ?", item.ProductID).
				Update("stock", gorm.Expr("stock + ?", item.Quantity))
		}
	}

	order.Status = req.Status
	if err := database.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	database.DB.Preload("Items.Product").Preload("User").First(&order, order.ID)
	c.JSON(http.StatusOK, order)
}

// Admin: List users
func (h *OrderHandler) ListUsers(c *gin.Context) {
	var users []models.User
	if err := database.DB.Select("id, name, email, phone, role, created_at").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, users)
}
