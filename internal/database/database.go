package database

import (
	"fmt"
	"log"
	"time"

	"github.com/balaji/ecommerce/internal/config"
	"github.com/balaji/ecommerce/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) error {
	var err error
	DB, err = gorm.Open(postgres.Open(cfg.DBDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connected successfully")
	return nil
}

func Migrate() error {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.Cart{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
	)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	log.Println("Database migrated successfully")
	return nil
}

func Seed() error {
	var count int64
	DB.Model(&models.User{}).Where("role = ?", "admin").Count(&count)
	if count > 0 {
		log.Println("Database already seeded, skipping...")
		return nil
	}

	adminPass, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := models.User{
		Name:     "Admin User",
		Email:    "admin@balaji.com",
		Password: string(adminPass),
		Phone:    "9999999999",
		Role:     "admin",
	}
	if err := DB.Create(&admin).Error; err != nil {
		return err
	}

	userPass, _ := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
	user := models.User{
		Name:     "John Doe",
		Email:    "user@balaji.com",
		Password: string(userPass),
		Phone:    "9876543210",
		Role:     "user",
	}
	if err := DB.Create(&user).Error; err != nil {
		return err
	}

	categories := []models.Category{
		{Name: "Men", Slug: "men", Description: "Men's Clothing Collection"},
		{Name: "Women", Slug: "women", Description: "Women's Clothing Collection"},
		{Name: "Kids",Slug: "kids", Description: "Kids Clothing Collection"},
		{Name: "Accessories",Slug: "accessories", Description: "Fashion Accessories"},
		{Name: "Footwear",Slug: "footwear", Description: "Shoes and Footwear"},
	}
	for i := range categories {
		if err := DB.Create(&categories[i]).Error; err != nil {
			return err
		}
	}

	products := []models.Product{
		{Name: "Classic Formal Shirt Black", Slug: "classic-formal-shirt-black", Description: "Premium cotton formal shirt for office wear.", Price: 1299, CompareAt: 1799, Stock: 50, ImageURL: "https://images.unsplash.com/photo-1596755094514-f87e34085b2c?w=400", CategoryID: 1, Rating: 4.5, NumReviews: 128, IsFeatured: true},
		{Name: "Casual Denim Shirt Blue",Slug: "casual-denim-shirt-blue", Description: "Comfortable denim shirt perfect for casual outings.", Price: 999, CompareAt: 1499, Stock: 40, ImageURL: "https://images.unsplash.com/photo-1603252109303-2751441dd157?w=400", CategoryID: 1, Rating: 4.2, NumReviews: 89, IsFeatured: true},
		{Name: "Slim Fit Chinos Beige",Slug: "slim-fit-chinos-beige", Description: "Modern slim fit chinos in beige.", Price: 1499, CompareAt: 1999, Stock: 35, ImageURL: "https://images.unsplash.com/photo-1473966968600-fa801b869a1a?w=400", CategoryID: 1, Rating: 4.3, NumReviews: 67},
		{Name: "Cotton Polo T-Shirt Navy",Slug: "cotton-polo-navy", Description: "Classic navy polo t-shirt made from pure cotton.", Price: 799, CompareAt: 999, Stock: 80, ImageURL: "https://images.unsplash.com/photo-1586790170083-2f9ceadc732d?w=400", CategoryID: 1, Rating: 4.6, NumReviews: 210, IsFeatured: true},
		{Name: "Leather Jacket Brown",Slug: "leather-jacket-brown", Description: "Genuine leather jacket with classic biker style.", Price: 4999, CompareAt: 6999, Stock: 15, ImageURL: "https://images.unsplash.com/photo-1551028719-00167b16eac5?w=400", CategoryID: 1, Rating: 4.8, NumReviews: 45},
		{Name: "Cargo Pants Olive",Slug: "cargo-pants-olive", Description: "Multi-pocket cargo pants in olive green.", Price: 1699, CompareAt: 2199, Stock: 30, ImageURL: "https://images.unsplash.com/photo-1624378439575-d8705ad7ae80?w=400", CategoryID: 1, Rating: 4.1, NumReviews: 54},
		{Name: "Floral Summer Dress",Slug: "floral-summer-dress", Description: "Light and breezy floral print summer dress.", Price: 1899, CompareAt: 2499, Stock: 25, ImageURL: "https://images.unsplash.com/photo-1572804013309-59a88b7e92f1?w=400", CategoryID: 2, Rating: 4.7, NumReviews: 156, IsFeatured: true},
		{Name: "Blue Lehenga Set",Slug: "blue-lehenga-set", Description: "Elegant blue lehenga with intricate embroidery.", Price: 4999, CompareAt: 6999, Stock: 12, ImageURL: "https://images.unsplash.com/photo-1610030469983-98e550d6193c?w=400", CategoryID: 2, Rating: 4.9, NumReviews: 78, IsFeatured: true},
		{Name: "Designer Kurti Pink",Slug: "designer-kurti-pink", Description: "Beautiful pink designer kurti with modern prints.", Price: 1299, CompareAt: 1799, Stock: 40, ImageURL: "https://images.unsplash.com/photo-1583391733956-3750e0ff4e8b?w=400", CategoryID: 2, Rating: 4.4, NumReviews: 92},
		{Name: "High-Waist Jeans Black",Slug: "high-waist-jeans-black", Description: "Stretchable high-waist black jeans with perfect fit.", Price: 1599, CompareAt: 1999, Stock: 45, ImageURL: "https://images.unsplash.com/photo-1541099649105-f69ad21f3246?w=400", CategoryID: 2, Rating: 4.5, NumReviews: 134},
		{Name: "Silk Saree Red",Slug: "silk-saree-red", Description: "Traditional pure silk saree in vibrant red.", Price: 5999, CompareAt: 7999, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1610030469983-98e550d6193c?w=400", CategoryID: 2, Rating: 4.8, NumReviews: 61, IsFeatured: true},
		{Name: "Crop Top White",Slug: "crop-top-white", Description: "Trendy white crop top for casual and party wear.", Price: 699, CompareAt: 999, Stock: 60, ImageURL: "https://images.unsplash.com/photo-1564257631407-4deb1f99d992?w=400", CategoryID: 2, Rating: 4.2, NumReviews: 88},
		{Name: "Kids Graphic T-Shirt",Slug: "kids-graphic-tshirt", Description: "Fun graphic print t-shirt for kids.", Price: 499, CompareAt: 699, Stock: 70, ImageURL: "https://images.unsplash.com/photo-1519238262714-846ea7ef2a2c?w=400", CategoryID: 3, Rating: 4.3, NumReviews: 45},
		{Name: "Boys Formal Suit",Slug: "boys-formal-suit", Description: "Complete formal suit set for boys.", Price: 2499, CompareAt: 3499, Stock: 20, ImageURL: "https://images.unsplash.com/photo-1503919545889-aef636e10ad4?w=400", CategoryID: 3, Rating: 4.6, NumReviews: 32},
		{Name: "Girls Party Frock",Slug: "girls-party-frock", Description: "Beautiful party frock with lace detailing.", Price: 1499, CompareAt: 1999, Stock: 25, ImageURL: "https://images.unsplash.com/photo-1518831959646-742c3a14ebf7?w=400", CategoryID: 3, Rating: 4.7, NumReviews: 41, IsFeatured: true},
		{Name: "Leather Wallet Brown",Slug: "leather-wallet-brown", Description: "Genuine leather bifold wallet.", Price: 899, CompareAt: 1299, Stock: 55, ImageURL: "https://images.unsplash.com/photo-1627123424574-724758594e93?w=400", CategoryID: 4, Rating: 4.4, NumReviews: 76},
		{Name: "Classic Sunglasses",Slug: "classic-sunglasses", Description: "UV protection classic aviator sunglasses.", Price: 1299, CompareAt: 1799, Stock: 40, ImageURL: "https://images.unsplash.com/photo-1572635196237-14b3f281503f?w=400", CategoryID: 4, Rating: 4.5, NumReviews: 112, IsFeatured: true},
		{Name: "Canvas Belt Black",Slug: "canvas-belt-black", Description: "Durable canvas belt with metal buckle.", Price: 499, CompareAt: 699, Stock: 80, ImageURL: "https://images.unsplash.com/photo-1553062407-98eeb64c6a62?w=400", CategoryID: 4, Rating: 4.1, NumReviews: 38},
		{Name: "Running Sneakers White",Slug: "running-sneakers-white", Description: "Lightweight running sneakers with cushioned sole.", Price: 2499, CompareAt: 3499, Stock: 35, ImageURL: "https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=400", CategoryID: 5, Rating: 4.6, NumReviews: 189, IsFeatured: true},
		{Name: "Formal Leather Shoes",Slug: "formal-leather-shoes", Description: "Classic black formal leather shoes.", Price: 2999, CompareAt: 3999, Stock: 28, ImageURL: "https://images.unsplash.com/photo-1614252235816-8e4a2a2d1f4b?w=400", CategoryID: 5, Rating: 4.4, NumReviews: 67},
		{Name: "Casual Loafers Brown",Slug: "casual-loafers-brown", Description: "Comfortable brown loafers for everyday wear.", Price: 1899, CompareAt: 2499, Stock: 32, ImageURL: "https://images.unsplash.com/photo-1533867617858-e7b97e060509?w=400", CategoryID: 5, Rating: 4.3, NumReviews: 54},
		{Name: "Sports Sandals",Slug: "sports-sandals", Description: "Durable sports sandals for outdoor activities.", Price: 999, CompareAt: 1399, Stock: 45, ImageURL: "https://images.unsplash.com/photo-1603487742131-4160ec2625da?w=400", CategoryID: 5, Rating: 4.2, NumReviews: 43},
	}

	for i := range products {
		if err := DB.Create(&products[i]).Error; err != nil {
			return err
		}
	}

	log.Println("Database seeded successfully with sample data")
	log.Println("Admin: admin@balaji.com / admin123")
	log.Println("User:  user@balaji.com / user123")
	return nil
}
