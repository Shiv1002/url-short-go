package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"

	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// URL Model: This defines the columns in your database table
type URL struct {
	ID        uint   `gorm:"primaryKey"`
	Code      string `gorm:"uniqueIndex"` // The short code (e.g., "aB12z")
	Original  string `gorm:"not null"`    // The long URL
	Username  string `gorm:"index;not null"`
	CreatedAt time.Time
}

var db *gorm.DB
var store = cookie.NewStore([]byte("very-secret-key-12345"))

const URL_CODE_LENGTH int = 6

func initDB() {
	var err error
	// This creates a file named "urls.db" in your folder automatically
	db, err = gorm.Open(sqlite.Open("urls.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// This creates the 'urls' table automatically based on our struct
	db.AutoMigrate(&URL{})
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateCode(codeLength int) string {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	b := make([]byte, codeLength)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

func main() {
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}
	// Set global options for every session created
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   true,
	})

	initDB()
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		// fmt.Printf("ClientIP: %s\n", c.ClientIP())
		var urls []URL
		// username, err := c.Cookie("uid")
		session, _ := store.Get(c.Request, "sid")

		if username, ok := session.Values["username"].(string); ok {

			db.Where("username = ?", username).Order("created_at desc").Find(&urls)

		} else {

			usernameCode := generateCode(20)
			session.Values["username"] = usernameCode
			session.Save(c.Request, c.Writer)
			// c.SetCookie("uid", usernameCode, 3600*24*15, "/", "", false, true)
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"urls": urls,
		})
	})

	// GET: Retrieve from SQLite and redirect
	router.GET("/get/:code", func(c *gin.Context) {
		// fmt.Printf("ClientIP: %s\n", c.ClientIP())
		code := c.Param("code")

		if len(code) != URL_CODE_LENGTH {
			c.JSON(400, gin.H{
				"error": "invalid code",
			})
			return
		}

		var urlEntry URL

		// Find the first record where code matches
		if err := db.Where("code = ?", code).First(&urlEntry).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
			return
		}

		c.Redirect(http.StatusMovedPermanently, urlEntry.Original)
	})

	// POST: Shorten a URL and save to SQLite
	router.POST("/shorten", func(c *gin.Context) {
		// longURL := c.PostForm("long_url")
		// username := c.PostForm("username")
		type ShortenRequest struct {
			LongURL string `form:"long_url" binding:"required,url"`
		}

		session, _ := store.Get(c.Request, "sid")

		username, ok := session.Values["username"].(string)

		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User"})
		}

		fmt.Println("username ", username)

		var req ShortenRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
			return
		}

		newURL := URL{
			Code:     generateCode(URL_CODE_LENGTH),
			Original: req.LongURL,
			Username: username,
		}

		// Save to the database
		result := db.Create(&newURL)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save URL"})
			return
		}

		c.Redirect(http.StatusSeeOther, "/")
	})

	router.Run(":" + PORT)
	fmt.Printf("Starting Server ar %s", PORT)
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin": // macOS
		err = exec.Command("open", url).Start()
	}

	if err != nil {
		fmt.Printf("Error opening browser: %v\n", err)
	}
}
