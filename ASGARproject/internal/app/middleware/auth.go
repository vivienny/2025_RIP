package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"decode/internal/app/config"
	"decode/internal/app/ds"
	"decode/internal/app/redis"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
)

type AuthMiddleware struct {
	jwtSecret string
	redis     *redis.Client
}

func NewAuthMiddleware(jwtSecret string, redisClient *redis.Client) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
		redis:     redisClient,
	}
}

// AuthRequired middleware проверяет JWT токен
func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Требуется токен авторизации",
			})
			return
		}

		// Проверяем, не в черном списке ли токен (если Redis доступен)
		if m.redis != nil {
			inBlacklist, err := m.redis.IsInBlacklist(c.Request.Context(), tokenString)
			if err != nil {
				logrus.Errorf("Ошибка проверки черного списка: %v", err)
			}
			if inBlacklist {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Токен отозван",
				})
				return
			}
		}

		// Парсим JWT токен
		token, err := jwt.ParseWithClaims(tokenString, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
			}
			return []byte(m.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Неверный или просроченный токен",
			})
			return
		}

		claims, ok := token.Claims.(*ds.JWTClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Неверные данные токена",
			})
			return
		}

		// Сохраняем информацию о пользователе в контексте
		c.Set("user_id", claims.UserID)
		c.Set("user_login", claims.Login)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// RoleRequired middleware проверяет роль пользователя
func (m *AuthMiddleware) RoleRequired(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Роль пользователя не найдена",
			})
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Неверный тип роли пользователя",
			})
			return
		}

		// Проверяем, есть ли у пользователя нужная роль
		hasRole := false
		for _, requiredRole := range roles {
			if roleStr == requiredRole {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": fmt.Sprintf("Требуемая роль: %v", roles),
			})
			return
		}

		c.Next()
	}
}

// GenerateJWTToken генерирует новый JWT токен
func GenerateJWTToken(user *ds.Avius, cfg config.JWTConfig) (string, time.Time, error) {
	expirationTime := time.Now().Add(cfg.ExpiresIn)

	claims := &ds.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "asgar-avia",
		},
		UserID: user.ID,
		Login:  user.Login,
		Role:   user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

func extractToken(c *gin.Context) string {
	bearerToken := c.GetHeader(AuthorizationHeader)
	if bearerToken == "" {
		// Пробуем получить из куки
		token, err := c.Cookie("token")
		if err == nil && token != "" {
			return token
		}
		return ""
	}

	if strings.HasPrefix(bearerToken, BearerPrefix) {
		return strings.TrimPrefix(bearerToken, BearerPrefix)
	}

	return bearerToken
}

// GetUserFromContext извлекает пользователя из контекста
func GetUserFromContext(c *gin.Context) (*ds.AuthUser, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return nil, fmt.Errorf("пользователь не аутентифицирован")
	}

	userLogin, _ := c.Get("user_login")
	userRole, _ := c.Get("user_role")

	id, ok := userID.(uint)
	if !ok {
		return nil, fmt.Errorf("неверный тип ID пользователя")
	}

	login, _ := userLogin.(string)
	role, _ := userRole.(string)

	return &ds.AuthUser{
		ID:    id,
		Login: login,
		Role:  role,
	}, nil
}
