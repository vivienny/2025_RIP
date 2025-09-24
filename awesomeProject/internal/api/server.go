package api

import (
	"decode/internal/app/handler"
	"decode/internal/app/repository"
	"log"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория")
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()
	// добавляем наш html/шаблон
	r.LoadHTMLGlob("../../templates/*")
	r.Static("/static", "../../resources") // ← исправлено!
	// слева название папки, в которую выгрузится наша статика
	// справа путь к папке, в которой лежит статика

	r.GET("/hello", handler.GetOrders)
	r.GET("/order/:id", handler.GetOrder)
	r.GET("/cart", handler.GetCart) // ← ДОБАВЛЕН МАРШРУТ ДЛЯ КОРЗИНЫ

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	log.Println("Server down")
}
