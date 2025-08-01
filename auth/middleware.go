package auth

import (
	"golang-whatsapp-clone/config"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	jwtService *JWTService
	appConfig  *config.AppConfig
}

func NewAuthMiddleware(jwtService *JWTService, appConfig *config.AppConfig) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
		appConfig:  appConfig,
	}
}

// This middleware will try to attach the user data to the context request in case an authentication header exists
func (m *AuthMiddleware) AuthenticateUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var tokenString string

		// try to get token from Authorization header first (for mobile)
		authHeader := c.Get("Authorization")
		if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			// try to get token from cookie (web)
			tokenString = c.Cookies(m.appConfig.CookieName)
		}

		if tokenString == "" {
			// return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			// 	"error": "Missing authentication token",
			// })
			return c.Next()
		}

		// validate jwt
		claims, err := m.jwtService.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authentication token",
			})
		}

		// store user info in context for graphql resolvers
		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)

		// set headers for graphql handler
		c.Request().Header.Set("X-User-ID", claims.UserID)
		c.Request().Header.Set("X-User-Email", claims.Email)

		return c.Next()
	}
}
