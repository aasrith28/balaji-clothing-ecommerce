package handlers

import (
	"net/http"
	"strconv"

	"github.com/balaji/ecommerce/internal/database"
	"github.com/balaji/ecommerce/internal/models"
	"github.com/gin-gonic/gin"
)

type CartHandler struct{}

func NewCartHandler() *CartHandler {
	return &CartHandler{}
}

func getOrCreateCart(userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := database.DB.Where("user_id = ?", userID).Preload("Items.Product").First(&cart).Error
	if err != nil {
		// Create new cart
		cart = models.Cart{UserID: userID}
		if err := database.DB.Create(&cart).Error; err != nil {
			return nil, err
		}
		database.DB.Preload("Items.Product").First(&cart, cart.ID)
	}
	return &cart, nil
}

func (h *CartHandler) GetCart(c *gin.Context) {
	userID, _ := c.Get("userID")
	cart, err := getOrCreateCart(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cart"})
		return
	}

	var subtotal float64
	var itemCount int
	for _, item := range cart.Items {
		subtotal += item.Product.Price * float64(item.Quantity)
		itemCount += item.Quantity
	}

	c.JSON(http.StatusOK, gin.H{
		"cart":       cart,
		"subtotal":   subtotal,
		"item_count": itemCount,
	})
}

func (h *CartHandler) AddItem(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req models.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check product
	var product models.Product
	if err := database.DB.First(&product, req.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	if product.Stock < req.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock", "available": product.Stock})
		return
	}

	cart, err := getOrCreateCart(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cart"})
		return
	}

	// Check if item already in cart
	var existing models.CartItem
	err = database.DB.Where("cart_id = ? AND product_id = ?", cart.ID, req.ProductID).First(&existing).Error
	if err == nil {
		// Update quantity
		newQty := existing.Quantity + req.Quantity
		if newQty > product.Stock {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock", "available": product.Stock})
			return
		}
		existing.Quantity = newQty
		database.DB.Save(&existing)
	} else {
		item := models.CartItem{
			CartID:    cart.ID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}
		if err := database.DB.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item"})
			return
		}
	}

	// Return updated cart
	cart, _ = getOrCreateCart(userID.(uint))
	var subtotal float64
	var itemCount int
	for _, item := range cart.Items {
		subtotal += item.Product.Price * float64(item.Quantity)
		itemCount += item.Quantity
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Item added to cart",
		"cart":       cart,
		"subtotal":   subtotal,
		"item_count": itemCount,
	})
}

func (h *CartHandler) UpdateItem(c *gin.Context) {
	userID, _ := c.Get("userID")
	itemID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var req models.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := getOrCreateCart(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cart"})
		return
	}

	var item models.CartItem
	if err := database.DB.Where("id = ? AND cart_id = ?", itemID, cart.ID).Preload("Product").First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}

	if req.Quantity == 0 {
		database.DB.Delete(&item)
	} else {
		if req.Quantity > item.Product.Stock {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock", "available": item.Product.Stock})
			return
		}
		item.Quantity = req.Quantity
		database.DB.Save(&item)
	}

	cart, _ = getOrCreateCart(userID.(uint))
	var subtotal float64
	var itemCount int
	for _, i := range cart.Items {
		subtotal += i.Product.Price * float64(i.Quantity)
		itemCount += i.Quantity
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Cart updated",
		"cart":       cart,
		"subtotal":   subtotal,
		"item_count": itemCount,
	})
}

func (h *CartHandler) RemoveItem(c *gin.Context) {
	userID, _ := c.Get("userID")
	itemID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	cart, err := getOrCreateCart(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cart"})
		return
	}

	result := database.DB.Where("id = ? AND cart_id = ?", itemID, cart.ID).Delete(&models.CartItem{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}

	cart, _ = getOrCreateCart(userID.(uint))
	var subtotal float64
	var itemCount int
	for _, i := range cart.Items {
		subtotal += i.Product.Price * float64(i.Quantity)
		itemCount += i.Quantity
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Item removed",
		"cart":       cart,
		"subtotal":   subtotal,
		"item_count": itemCount,
	})
}

func (h *CartHandler) Clear(c *gin.Context) {
	userID, _ := c.Get("userID")
	cart, err := getOrCreateCart(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cart"})
		return
	}
	database.DB.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{})
	c.JSON(http.StatusOK, gin.H{"message": "Cart cleared"})
}
