package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"

	"github.com/hustlers/motivator-backend/pkg/response"
)

type AuthMiddleware struct {
	jwtSecret []byte
}

func NewAuth(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{jwtSecret: []byte(jwtSecret)}
}

func (a *AuthMiddleware) Required() fiber.Handler {
	return func(c fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return response.Unauthorized(c, "missing authorization header")
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			return response.Unauthorized(c, "invalid authorization format")
		}

		token, err := jwt.Parse(parts[1], func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return a.jwtSecret, nil
		})
		if err != nil || !token.Valid {
			return response.Unauthorized(c, "invalid or expired token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return response.Unauthorized(c, "invalid token claims")
		}

		userID, _ := claims["sub"].(string)
		if userID == "" {
			return response.Unauthorized(c, "missing user id in token")
		}

		email, _ := claims["email"].(string)

		c.Locals("userID", userID)
		c.Locals("email", email)

		return c.Next()
	}
}

func GetUserID(c fiber.Ctx) string {
	id, _ := c.Locals("userID").(string)
	return id
}

func GetEmail(c fiber.Ctx) string {
	email, _ := c.Locals("email").(string)
	return email
}
