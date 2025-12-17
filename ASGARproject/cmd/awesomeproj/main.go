package main

import (
	"context"
	"decode/internal/app/config"
	"decode/internal/app/dsn"
	"decode/internal/app/handler"
	"decode/internal/app/redis"
	"decode/internal/app/repository"
	"decode/internal/pkg"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "decode/docs" // импорт сгенерированной документации

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title ASGAR Avia API
// @version 1.0
// @description API для системы заказа авиауслуг ASGAR Avia
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.asgaravia.com/support
// @contact.email support@asgaravia.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token

func main() {
	// ==================== НАСТРОЙКА ЛОГГЕРА ====================
	// Устанавливаем формат логов в JSON для удобства
	logrus.SetFormatter(&logrus.JSONFormatter{})
	// Логи выводим в стандартный вывод
	logrus.SetOutput(os.Stdout)
	// Уровень логирования - Info
	logrus.SetLevel(logrus.InfoLevel)

	logrus.Info("🚀 Запуск приложения ASGAR Avia")

	// ==================== ЗАГРУЗКА КОНФИГУРАЦИИ ====================
	// Загружаем конфигурацию из файлов и переменных окружения
	conf, err := config.NewConfig()
	if err != nil {
		// Если не удалось загрузить конфиг - завершаем с ошибкой
		logrus.Fatalf("❌ Ошибка загрузки конфигурации: %v", err)
	}
	logrus.Info("✅ Конфигурация успешно загружена")

	// ==================== ИНИЦИАЛИЗАЦИЯ REDIS (НОВОЕ ДЛЯ ЛР4) ====================
	// Создаем контекст для работы с Redis
	ctx := context.Background()
	var redisClient *redis.Client

	// Проверяем, настроен ли Redis в конфигурации
	if conf.Redis.Host != "" {
		// Пытаемся подключиться к Redis
		redisClient, err = redis.NewClient(ctx, conf.Redis)
		if err != nil {
			// Если не удалось - работаем без Redis (предупреждение)
			logrus.Warnf("⚠️ Не удалось подключиться к Redis: %v. Работаем без Redis.", err)
			redisClient = nil
		} else {
			// Закрываем соединение при завершении программы
			defer func() {
				if redisClient != nil {
					redisClient.Close()
				}
			}()
			logrus.Info("✅ Успешное подключение к Redis")
		}
	} else {
		logrus.Info("ℹ️ Redis не настроен, работаем без Redis")
	}

	// ==================== ИНИЦИАЛИЗАЦИЯ БАЗЫ ДАННЫХ ====================
	// Получаем строку подключения к PostgreSQL
	postgresString := dsn.FromEnv()
	fmt.Println("🔗 Строка подключения к БД:", postgresString)

	// Создаем репозиторий для работы с БД
	rep, errRep := repository.New(postgresString)
	if errRep != nil {
		// Если не удалось подключиться к БД - завершаем работу
		logrus.Fatalf("❌ Ошибка инициализации репозитория: %v", errRep)
	}
	logrus.Info("✅ Успешное подключение к базе данных")

	// ==================== ИНИЦИАЛИЗАЦИЯ ОБРАБОТЧИКОВ ====================
	// Создаем обработчики HTTP-запросов
	// ВАЖНО: нужно обновить NewHandler чтобы он принимал redisClient и JWT config
	hand := handler.NewHandler(rep, redisClient, conf.JWT)
	// Временная заглушка - позже нужно будет обновить конструктор
	logrus.Info("✅ Обработчики инициализированы")

	// ==================== НАСТРОЙКА GIN ФРЕЙМВОРКА ====================
	// Создаем роутер Gin
	router := gin.Default()

	// ==================== SWAGGER ДОКУМЕНТАЦИЯ ====================
	// ДОБАВЛЯЕМ ПЕРЕД регистрацией других маршрутов
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ==================== РЕГИСТРАЦИЯ МАРШРУТОВ ====================
	// Регистрируем все обработчики запросов
	hand.RegisterHandler(router)
	logrus.Info("✅ Маршруты зарегистрированы")

	// ==================== СОЗДАНИЕ И ЗАПУСК ПРИЛОЖЕНИЯ ====================
	// Создаем основное приложение
	application := pkg.NewApp(conf, router, hand)

	// ==================== НАСТРОЙКА ГРАЦИОЗНОГО ЗАВЕРШЕНИЯ ====================
	// Создаем канал для получения сигналов от системы
	quit := make(chan os.Signal, 1)
	// Настраиваем перехват сигналов Ctrl+C и системного завершения
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем сервер в отдельной горутине
	go func() {
		logrus.Infof("🌐 Сервер запускается на %s:%d", conf.ServiceHost, conf.ServicePort)
		logrus.Info("📝 Swagger будет доступен по адресу: http://localhost:8080/swagger/index.html")
		logrus.Info("🛑 Для остановки нажмите Ctrl+C")

		// Запускаем веб-сервер (этот метод уже существует в вашем коде)
		application.RunApp()
	}()

	// ==================== ОЖИДАНИЕ СИГНАЛА ЗАВЕРШЕНИЯ ====================
	// Ждем сигнала завершения
	<-quit
	logrus.Info("🔄 Получен сигнал завершения, останавливаем сервер...")

	// Даем время на завершение работы
	time.Sleep(2 * time.Second)
	logrus.Info("✅ Сервер успешно остановлен")
}
