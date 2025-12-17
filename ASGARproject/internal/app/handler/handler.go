package handler

import (
	"decode/internal/app/config"
	"decode/internal/app/middleware"
	"decode/internal/app/redis"
	"decode/internal/app/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
	Redis      *redis.Client
	JWTConfig  config.JWTConfig
	Auth       *middleware.AuthMiddleware
}

func NewHandler(r *repository.Repository, redisClient *redis.Client, jwtConfig config.JWTConfig) *Handler {
	authMiddleware := middleware.NewAuthMiddleware(jwtConfig.Secret, redisClient)

	return &Handler{
		Repository: r,
		Redis:      redisClient,
		JWTConfig:  jwtConfig,
		Auth:       authMiddleware,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	// Public routes (no auth required)
	public := router.Group("/api")
	{
		public.POST("/auth/login", h.LoginUser)
		public.POST("/auth/register", h.RegisterUser)
		public.POST("/auth/logout", h.LogoutUser)
		public.GET("/services", h.GetAsgarServices)
		public.GET("/services/:id", h.GetAsgarService)
		public.GET("/uploads/:filename", h.ServeFile)
	}

	// User routes (require authentication)
	user := router.Group("/api")
	user.Use(h.Auth.AuthRequired())
	{
		user.GET("/users/profile", h.GetUserProfile)
		user.PUT("/users/profile", h.UpdateUserProfile)

		// Miniplane (cart)
		user.GET("/miniplane/icon", h.GetMiniplaneIcon)
		user.GET("/miniplane", h.GetMiniplane)
		user.POST("/services/:id/add-to-miniplane", h.AddServiceToMiniplane)

		// User's flight services
		user.GET("/flight-servs", h.GetUserFlightServs) // Используем новый метод
		user.GET("/flight-servs/:id", h.GetFlightServ)

		user.PUT("/flight-servs/:id", h.UpdateFlightServ)
		user.PUT("/flight-servs/:id/submit", h.SubmitFlightServ)
		user.DELETE("/flight-servs/:id", h.DeleteFlightServ)

		// Subjservs
		user.DELETE("/subjservs/:id", h.DeleteSubjserv)
		user.PUT("/subjservs/:id", h.UpdateSubjserv)
	}
	// Moderator routes (require moderator role)
	moderator := router.Group("/api")
	moderator.Use(h.Auth.AuthRequired(), h.Auth.RoleRequired("moderator", "admin"))
	{
		// All flight services (moderator can see all)
		moderator.GET("/moderator/flight-servs", h.GetAllFlightServs) // Используем метод для модераторов
		moderator.PUT("/moderator/flight-servs/:id/complete", h.CompleteFlightServ)

		// Service management
		moderator.POST("/services", h.CreateAsgarService)
		moderator.PUT("/services/:id", h.UpdateAsgarService)
		moderator.DELETE("/services/:id", h.DeleteAsgarService)
		moderator.POST("/services/:id/upload", h.UploadServiceImage)

		// User management
		moderator.GET("/moderator/users", h.GetAllUsers)
		moderator.PUT("/moderator/users/:id/role", h.UpdateUserRole)
	}

	// Static files
	router.Static("/uploads", "./uploads")
}

// CORSMiddleware добавляет CORS заголовки
func (h *Handler) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// LoggerMiddleware логирует запросы
func (h *Handler) LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		logrus.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   path,
			"ip":     c.ClientIP(),
		}).Info("Request")

		c.Next()
	}
}

func (h *Handler) ServeFile(c *gin.Context) {
	filename := c.Param("filename")
	c.File("./uploads/" + filename)
}

func (h *Handler) errorResponse(ctx *gin.Context, statusCode int, err error) {
	logrus.WithFields(logrus.Fields{
		"status": statusCode,
		"path":   ctx.Request.URL.Path,
		"method": ctx.Request.Method,
	}).Error(err.Error())

	ctx.JSON(statusCode, gin.H{
		"status": "error",
		"error":  err.Error(),
	})
}

// RegisterStatic регистрирует статические файлы и HTML шаблоны
func (h *Handler) RegisterStatic(router *gin.Engine) {
	// 1. Загружаем HTML шаблоны из папки templates
	router.LoadHTMLGlob("templates/*.html")

	// 2. Статика для ресурсов (ваши иконки, изображения)
	router.Static("/resources", "./resources")
	// Теперь доступны: /resources/img/cart-icon.png и т.д.

	// 3. Маршруты для всех ваших HTML страниц
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "asgar_catalog.html", nil)
	})

	router.GET("/catalog", func(c *gin.Context) {
		c.HTML(http.StatusOK, "asgar_catalog.html", nil)
	})

	router.GET("/miniplane", func(c *gin.Context) {
		c.HTML(http.StatusOK, "miniplane.html", nil)
	})

	router.GET("/product", func(c *gin.Context) {
		c.HTML(http.StatusOK, "selected_aviaproduct.html", nil)
	})

	logrus.Info("✅ Статика зарегистрирована: templates, uploads, resources")
}
