package middleware

import (
	"go-image-upload/handlers"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-jwt/jwt"
)

func CheckJWT(app *handlers.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Read JWT token from cookie
		cookie := c.Cookies("jwt")
		if cookie == "" {
			return fiber.NewError(fiber.StatusBadRequest, "Token is missing")
		}

		// Decrypt JWT token
		tokenString, err := app.DecryptToken(cookie)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Failed to decrypt token")
		}

		// Parse JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return app.CookieKey, nil
		})
		if err != nil || !token.Valid {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid token")
		}

		// Store claims in Locals for further processing
		c.Locals("user", token.Claims)

		// Continue to next handler
		return c.Next()
	}
}

func CSRFMiddleware() fiber.Handler {
	return csrf.New(csrf.Config{
		CookieSameSite: "Strict",
		Expiration:     5 * time.Hour,
		KeyLookup:      "header:X-Csrf-Token",
	})
}

func LoginLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        5,
		Expiration: 5 * time.Minute,
	})
}

func LoggingMiddleware() fiber.Handler {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	logger := logger.New(logger.Config{
		Format:     "${time} - ${ip} - ${status} - ${method} - ${path} - ${error}\n",
		TimeFormat: "2006-01-02 15:04:05",
		Output:     file,
	})

	return logger
}
