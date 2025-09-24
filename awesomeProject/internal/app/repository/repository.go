package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
}

type Order struct {
	ID              int
	Title           string
	Description     string
	ImageURL        string
	Price           string
	Features        []string
	FullDescription string
	Unit            string
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

func (r *Repository) GetOrdersByTitle(title string) ([]Order, error) {
	orders, err := r.GetOrders()
	if err != nil {
		return []Order{}, err
	}

	var result []Order
	for _, order := range orders {
		if strings.Contains(strings.ToLower(order.Title), strings.ToLower(title)) {
			result = append(result, order)
		}
	}

	return result, nil
}

func (r *Repository) GetOrders() ([]Order, error) {
	orders := []Order{
		{
			ID:              1,
			Title:           "Заправочная мащина",
			Description:     "Машина для запрfвки воздушных судов",
			ImageURL:        "http://localhost:9000/test/zapravka.jpg",
			Price:           "45000 руб./сут",
			Features:        []string{"Защита до -50°C", "Быстрое действие", "Экологичность", "Легкое нанесение"},
			FullDescription: "Специализированный антифриз для обработки поверхностей воздушных судов. Обеспечивает надежную защиту от обледенения в любых погодных условиях.",
			Unit:            "руб./сут",
		},
		{
			ID:              2,
			Title:           "Авиационное топливо ",
			Description:     "Топливо JET A-1",
			ImageURL:        "http://localhost:9000/test//jeta1.jpg",
			Price:           "57677 руб./тонна",
			Features:        []string{"Высокая чистота", "Сертификация ISO", "Стабильное хранение", "Безопасность"},
			FullDescription: "Авиационное топливо JET A-1 соответствует международным стандартам качества. Используется для коммерческих авиаперевозок. Гарантия поставки в срок.",
			Unit:            "руб./час",
		},
		{
			ID:              3,
			Title:           "Ленточный погрузчик ",
			Description:     "Mulag Orbiter 12D/E",
			ImageURL:        "http://localhost:9000/test//lenta.jpg",
			Price:           "1250000 руб./час",
			Features:        []string{"Грузоподъемность 2т", "Электрический привод", "Плавный ход", "Низкий уровень шума"},
			FullDescription: "Немецкий ленточный погрузчик Mulag Orbiter 12D/E предназначен для эффективной погрузки багажа в самолеты. Надежность и долговечность гарантированы.",
			Unit:            "руб./сут.",
		},
		{
			ID:              4,
			Title:           "Буксировщик",
			Description:     "Тягач Challenger 550",
			ImageURL:        "http://localhost:9000/test/tyagach.jpg",
			Price:           "260000 руб./сут",
			Features:        []string{"Мощность 400 л.с.", "Гидравлическая система", "Точное управление", "Безопасность"},
			FullDescription: "Буксировщик для воздушных судов обеспечивает безопасное и точное перемещение самолетов по территории аэропорта. Соответствует всем требованиям безопасности.",
			Unit:            "руб./сут.",
		},

		{
			ID:              5,
			Title:           "Автолифт",
			Description:     "Амбулаторный автолифт DOLL X-PRM M",
			ImageURL:        "http://localhost:9000/test//lift.jpg",
			Price:           "150000 руб./сут",
			Features:        []string{"Оригинальное качество", "Сертификация", "Гарантия", "Быстрая поставка"},
			FullDescription: "Полный ассортимент оригинальных запчастей для различных моделей самолетов. Все детали проходят строгий контроль качества.",
			Unit:            "руб./сут.",
		},
		{
			ID:              6,
			Title:           "Антифриз",
			Description:     "MAXCool Hybrid",
			ImageURL:        "http://localhost:9000/test/noice.jpg",
			Price:           "68100 руб./бочка",
			Features:        []string{"Высокая вязкость", "Термостойкость", "Защита от износа", "Долгий срок службы"},
			FullDescription: "Высококачественное авиационное масло для турбовинтовых и турбореактивных двигателей. Обеспечивает надежную работу в экстремальных условиях.",
			Unit:            "руб./тонна",
		},
		{
			ID:              7,
			Title:           "Деайсер",
			Description:     "ELEPAHT BETA-15",
			ImageURL:        "http://localhost:9000/test/antifrizmash.jpg",
			Price:           "125000 руб./сут",
			Features:        []string{"Грузоподъемность 2т", "Электрический привод", "Плавный ход", "Низкий уровень шума"},
			FullDescription: "Немецкий ленточный погрузчик Mulag Orbiter 12D/E предназначен для эффективной погрузки багажа в самолеты. Надежность и долговечность гарантированы.",
			Unit:            "руб./сут",
		},
	}

	return orders, nil
}

func (r *Repository) GetOrder(id int) (Order, error) {
	orders, err := r.GetOrders()
	if err != nil {
		return Order{}, err
	}

	for _, order := range orders {
		if order.ID == id {
			return order, nil
		}
	}
	return Order{}, fmt.Errorf("заказ не найден")
}
