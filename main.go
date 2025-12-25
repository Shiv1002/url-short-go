package main

import (
	"fmt"
	"math/rand"

	// "net/http"
	"time"
	// "github.com/gin-gonic/gin"
	// "gorm.io/driver/sqlite"
	// "gorm.io/gorm"
)

func generateCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	fmt.Println("this is ", time.Now())

	source := rand.NewSource(time.Now().UnixNano())

	r := rand.New(source)
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

func main() {
	fmt.Println(generateCode())
}

// // URL Model: This defines the columns in your database table
// type URL struct {
// 	ID        uint   `gorm:"primaryKey"`
// 	Code      string `gorm:"uniqueIndex"` // The short code (e.g., "aB12z")
// 	Original  string `gorm:"not null"`    // The long URL
// 	CreatedAt time.Time
// }

// var db *gorm.DB

// func initDB() {
// 	var err error
// 	// This creates a file named "urls.db" in your folder automatically
// 	db, err = gorm.Open(sqlite.Open("urls.db"), &gorm.Config{})
// 	if err != nil {
// 		panic("Failed to connect to database")
// 	}

// 	// This creates the 'urls' table automatically based on our struct
// 	db.AutoMigrate(&URL{})
// }

// func generateCode() string {
// 	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
// 	source := rand.NewSource(time.Now().UnixNano())
// 	r := rand.New(source)
// 	b := make([]byte, 6)
// 	for i := range b {
// 		b[i] = charset[r.Intn(len(charset))]
// 	}
// 	return string(b)
// }

// func main() {
// 	initDB()
// 	r := gin.Default()

// 	// POST: Shorten a URL and save to SQLite
// 	r.POST("/shorten", func(c *gin.Context) {
// 		var input struct {
// 			LongURL string `json:"long_url" binding:"required"`
// 		}

// 		if err := c.ShouldBindJSON(&input); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
// 			return
// 		}

// 		newURL := URL{
// 			Code:     generateCode(),
// 			Original: input.LongURL,
// 		}

// 		// Save to the database
// 		result := db.Create(&newURL)
// 		if result.Error != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save URL"})
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{"short_url": "http://localhost:8080/" + newURL.Code})
// 	})

// 	// GET: Retrieve from SQLite and redirect
// 	r.GET("/:code", func(c *gin.Context) {
// 		code := c.Param("code")
// 		var urlEntry URL

// 		// Find the first record where code matches
// 		if err := db.Where("code = ?", code).First(&urlEntry).Error; err != nil {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
// 			return
// 		}

// 		c.Redirect(http.StatusMovedPermanently, urlEntry.Original)
// 	})

// 	r.Run(":8080")
// }
