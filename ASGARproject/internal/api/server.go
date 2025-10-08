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

	r.LoadHTMLGlob("../../templates/*")
	r.Static("/static", "../../resources")

	r.GET("/ASGARcatalog", handler.GetAsgarList)
	r.GET("/SelectAviaSub/:id", handler.GetSelectAviaSub)
	r.GET("/miniplane", handler.GetMiniPlane)

	r.Run()
	log.Println("Server down")
}
