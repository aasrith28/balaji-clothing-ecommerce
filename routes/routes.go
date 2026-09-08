package routes

import (
	"github.com/balaji/ecommerce/internal/config"
	"github.com/balaji/ecommerce/internal/handlers"
	"github.com/balaji/ecommerce/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, cfg *config.Config) {
	authH := handlers.NewAuthHandler(cfg)
	productH := handlers.NewProductHandler()
	cartH := handlers.NewCartHandler()
	orderH := handlers.NewOrderHandler()

	// Public API
	api := r.Group("/api")
	{
		// Auth
		api.POST("/auth/register", authH.Register)
		api.POST("/auth/login", authH.Login)

		// Products (public)
		api.GET("/products", productH.List)
		api.GET("/products/featured", productH.Featured)
		api.GET("/products/:id", productH.Get)
		api.GET("/categories", productH.Categories)

		// Protected routes
		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware(cfg))
		{
			// Profile
			auth.GET("/auth/profile", authH.Profile)
			auth.PUT("/auth/profile", authH.UpdateProfile)

			// Cart
			auth.GET("/cart", cartH.GetCart)
			auth.POST("/cart/items", cartH.AddItem)
			auth.PUT("/cart/items/:id", cartH.UpdateItem)
			auth.DELETE("/cart/items/:id", cartH.RemoveItem)
			auth.DELETE("/cart", cartH.Clear)

			// Orders
			auth.POST("/checkout", orderH.Checkout)
			auth.GET("/orders", orderH.MyOrders)
			auth.GET("/orders/:id", orderH.GetOrder)

			// Admin routes
			admin := auth.Group("/admin")
			admin.Use(middleware.AdminMiddleware())
			{
				admin.GET("/products", productH.AdminList)
				admin.POST("/products", productH.Create)
				admin.PUT("/products/:id", productH.Update)
				admin.DELETE("/products/:id", productH.Delete)

				admin.GET("/orders", orderH.AdminList)
				admin.PUT("/orders/:id/status", orderH.UpdateStatus)

				admin.GET("/users", orderH.ListUsers)
			}
		}
	}

	// Serve frontend
	r.Static("/static", "./static")
	r.StaticFile("/", "./frontend/index.html")
	r.StaticFile("/shop", "./frontend/shop.html")
	r.StaticFile("/product", "./frontend/product.html")
	r.StaticFile("/cart", "./frontend/cart.html")
	r.StaticFile("/checkout", "./frontend/checkout.html")
	r.StaticFile("/login", "./frontend/login.html")
	r.StaticFile("/register", "./frontend/register.html")
	r.StaticFile("/orders", "./frontend/orders.html")
	r.StaticFile("/profile", "./frontend/profile.html")
	r.StaticFile("/admin", "./frontend/admin.html")
	r.StaticFile("/admin/products", "./frontend/admin-products.html")
	r.StaticFile("/admin/orders", "./frontend/admin-orders.html")
	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/index.html")
	})
}
