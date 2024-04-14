package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"go-image-upload/models"
	"go-image-upload/util"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type App struct {
	DB        *gorm.DB
	CookieKey []byte // Key for AES encryption of JWT token
}

func NewApp(db *gorm.DB, cookieKey []byte) *App {
	return &App{DB: db, CookieKey: cookieKey}
}

// Define methods on the App struct to handle route logic
func (app *App) HandleIndex(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{})
}

// handle login
func (app *App) HandleLogin(c *fiber.Ctx) error {
	type LoginRequest struct {
		Username  string `form:"username"`
		Password  string `form:"password"`
		Email     string `form:"email"`
		CSRFToken string `form:"csrf_token"`
	}
	var loginReq LoginRequest
	if err := c.BodyParser(&loginReq); err != nil {
		return err
	}

	// Get CSRF token from request
	csrfToken := c.FormValue("csrf_token")

	// Verify CSRF token
	if csrfToken != util.GetCSRFCookie(c) {
		return fiber.NewError(fiber.StatusUnauthorized, "CSRF token mismatch")
	}

	var user models.User
	if err := app.DB.Where("username = ?", loginReq.Username).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).Redirect("/?message=Invalid%20username%20or%20password&type=error")
	}

	// Check if the password matches using bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).Redirect("/?message=Invalid%20username%20or%20password&type=error")
	}

	// If login successful, create JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["email"] = user.Email
	tokenString, err := token.SignedString(app.CookieKey)
	if err != nil {
		return err
	}

	// Encrypt token for storage in cookie
	encryptedToken, err := app.EncryptToken(tokenString)
	if err != nil {
		return err
	}

	// Set the token as a cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    encryptedToken,
		Expires:  time.Now().Add(5 * time.Hour),
		HTTPOnly: false,
		Secure:   true,
		Path:     "/",
	})

	return c.Redirect("/dashboard")
}

func (app *App) HandleDashboard(c *fiber.Ctx) error {

	type ImageData struct {
		ID      int
		URL     string
		Subject string
	}
	// Get CSRF token from request
	csrfToken := c.FormValue("csrf_token")

	// Verify CSRF token
	if csrfToken != util.GetCSRFCookie(c) {
		return fiber.NewError(fiber.StatusUnauthorized, "CSRF token mismatch")
	}

	// Get claims from Locals
	claims := c.Locals("user").(jwt.MapClaims)

	// Get username from claims
	username := claims["username"].(string)
	email := claims["email"].(string)

	// Query database to get images uploaded by the current user
	var images []models.Image
	if err := app.DB.Where("username = ?", username).Find(&images).Error; err != nil {
		return err
	}

	// Create a slice to store image data
	var imageData []ImageData

	// Iterate over the images and create image data
	for _, image := range images {
		imageData = append(imageData, ImageData{
			ID:      image.ID,
			URL:     image.ImageURL,
			Subject: image.Subject,
		})
	}

	// Render dashboard with username and image data
	return c.Render("dashboard", fiber.Map{
		"Username":  username,
		"Email":     email,
		"ImageData": imageData,
	})
}

func (app *App) HandleLogout(c *fiber.Ctx) error {
	// Clear JWT cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   true,
		Path:     "/",
	})

	// Redirect to login page
	return c.Redirect("/")
}

func (app *App) HandleSignup(c *fiber.Ctx) error {
	type SignupRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	var signupReq SignupRequest
	if err := c.BodyParser(&signupReq); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}

	// Check if user already exists
	var existingUser models.User
	if err := app.DB.Where("username = ?", signupReq.Username).First(&existingUser).Error; err == nil {
		// return fiber.NewError(fiber.StatusBadRequest, "Username already exists")
		return c.Status(fiber.StatusBadRequest).Redirect("/?message=Username%20already%20exists&type=error")
	}

	// Hash password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create new user entry in database
	newUser := models.User{
		Username: signupReq.Username,
		Password: string(hashedPassword),
		Email:    signupReq.Email,
	}
	if err := app.DB.Create(&newUser).Error; err != nil {
		return err
	}

	// Respond with success message
	return c.Status(fiber.StatusOK).Redirect("/?message=Signup%20successful")

}

func (app *App) HandleUpload(c *fiber.Ctx) error {
	// Get CSRF token from request
	csrfToken := c.FormValue("csrf_token")

	// Verify CSRF token
	if csrfToken != util.GetCSRFCookie(c) {
		return fiber.NewError(fiber.StatusUnauthorized, "CSRF token mismatch")
	}

	// Get claims from Locals
	claims := c.Locals("user").(jwt.MapClaims)

	// Get form values
	subject := c.FormValue("subject")

	// Get the file from the request
	file, err := c.FormFile("image")
	if err != nil {
		return err
	}

	// Validate image size
	if err := util.ValidateImageSize(file); err != nil {
		return err
	}

	// Validate image format
	if err := util.ValidateImageFormat(file); err != nil {
		return err
	}

	// Get username from claims
	username := claims["username"].(string)

	// Save file to server
	imagePath := "./uploads/" + file.Filename
	if err := c.SaveFile(file, imagePath); err != nil {
		return err
	}

	// Save image data to database
	image := models.Image{
		Username: username,
		Subject:  subject,
		ImageURL: imagePath,
	}
	if err := app.DB.Create(&image).Error; err != nil {
		return err
	}

	// Redirect to dashboard or another page
	return c.Redirect("/dashboard")
}

// Function to delete image from database and file directory
func (app *App) HandleDeleteImage(c *fiber.Ctx) error {
	// Get CSRF token from request
	csrfToken := c.FormValue("csrf_token")

	// Verify CSRF token
	if csrfToken != util.GetCSRFCookie(c) {
		return fiber.NewError(fiber.StatusUnauthorized, "CSRF token mismatch")
	}

	// Get image ID from form
	imageID := c.Params("id")

	// Convert imageID to integer
	id, err := strconv.ParseUint(imageID, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid image ID")
	}

	// Query database to find image by ID
	var image models.Image
	result := app.DB.Where("id = ?", id).First(&image)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusNotFound, "Image not found")
		}
		return result.Error
	}

	// Delete image record from database
	if err := app.DB.Delete(&image).Error; err != nil {
		return err
	}

	// Delete image file from file directory
	if err := os.Remove(image.ImageURL); err != nil {
		fmt.Println("Error:", err)
	}
	// Redirect to dashboard or another page
	return c.Redirect("/dashboard")
}

// Function to encrypt a JWT token using AES encryption
// Depends on handle login
func (app *App) EncryptToken(tokenString string) (string, error) {
	block, err := aes.NewCipher(app.CookieKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(tokenString), nil)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Function to decrypt a JWT token using AES encryption
// depends on middleware CheckJWT
func (app *App) DecryptToken(encryptedToken string) (string, error) {
	block, err := aes.NewCipher(app.CookieKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	decodedCiphertext, err := base64.URLEncoding.DecodeString(encryptedToken)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	if len(decodedCiphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}
	nonce, ciphertext := decodedCiphertext[:nonceSize], decodedCiphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
