package util

import (
	"crypto/rand"

	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GenerateRandomKey(length int) ([]byte, error) {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// Function to get CSRF token from cookie
func GetCSRFCookie(c *fiber.Ctx) string {
	cookie := c.Cookies("csrf")
	return cookie
}

// Validate image size (less than 2MB)
func ValidateImageSize(file *multipart.FileHeader) error {
	if file.Size > 1<<20 { // 1MB
		return fiber.NewError(fiber.StatusBadRequest, "Image size must be less than 2MB")
	}
	return nil
}

// Validate image format (must be an image format)
func ValidateImageFormat(file *multipart.FileHeader) error {
	allowedFormats := []string{".jpg", ".jpeg", ".png", ".gif"}
	ext := filepath.Ext(file.Filename)
	for _, format := range allowedFormats {
		if strings.EqualFold(ext, format) {
			return nil
		}
	}
	return fiber.NewError(fiber.StatusBadRequest, "Invalid image format. Only JPG, JPEG, PNG, and GIF are allowed.")

}
