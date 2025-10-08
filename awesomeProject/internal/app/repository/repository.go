package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
	FlightParam map[int]FlightParam
}

type ASGARserv struct {
	ID              int
	Title           string
	Description     string
	ImageURL        string
	Price           string
	FullDescription string
	Unit            string
}

type FlightParam struct {
	ID           int
	FlightNumber string
	Date         string
	Departure    string
	Arrival      string
	AircraftType string
	Services     []int
	TotalPrice   float64
	Status       string
}

func NewRepository() (*Repository, error) {
	return &Repository{
		FlightParam: make(map[int]FlightParam),
	}, nil
}

func (r *Repository) ToDoFlightInfo(app FlightParam) error {
	r.FlightParam[app.ID] = app
	return nil
}

func (r *Repository) GetFlightParam(id int) (FlightParam, error) {
	app, exists := r.FlightParam[id]
	if !exists {
		return FlightParam{}, fmt.Errorf("заявка не найдена")
	}
	return app, nil
}

func (r *Repository) GetFindAviaSub(title string) ([]ASGARserv, error) {
	aviationServices, err := r.GetAsgarList()
	if err != nil {
		return []ASGARserv{}, err
	}

	var result []ASGARserv
	for _, selsubavia := range aviationServices {
		if strings.Contains(strings.ToLower(selsubavia.Title), strings.ToLower(title)) {
			result = append(result, selsubavia)
		}
	}

	return result, nil
}

func (r *Repository) GetAsgarList() ([]ASGARserv, error) {
	aviationServices := []ASGARserv{
		{
			ID:              1,
			Title:           "Заправочная машина",
			Description:     "Машина для заправки воздушных судов",
			ImageURL:        "http://localhost:9000/test/zapravka.jpg",
			Price:           "45000 руб./сут",
			FullDescription: "Специализированный антифриз для обработки поверхностей воздушных судов. Обеспечивает надежную защиту от обледенения в любых погодных условиях. Характеристики: защита до -50°C, быстрое действие, экологичность, легкое нанесение.",
			Unit:            "руб./сут",
		},
		{
			ID:              2,
			Title:           "Авиационное топливо",
			Description:     "Топливо JET A-1",
			ImageURL:        "http://localhost:9000/test//jeta1.jpg",
			Price:           "57677 руб./тонна",
			FullDescription: "Высококачественное авиационное топливо для реактивных двигателей. JET A-1 соответствует всем международным стандартам качества и безопасности. Используется авиакомпаниями по всему миру для коммерческих и грузовых перевозок. Обеспечивает стабильную работу двигателей в различных климатических условиях. Поставляется в строгом соответствии с техническими требованиями. Характеристики: высокая чистота, сертификация ISO 9001, отличная термостабильность, антистатические свойства, низкая температура замерзания, увеличенный срок хранения.",
			Unit:            "руб./тонна",
		},
		{
			ID:              3,
			Title:           "Ленточный погрузчик",
			Description:     "Mulag Orbiter 12D/E",
			ImageURL:        "http://localhost:9000/test//lenta.jpg",
			Price:           "1250000 руб./час",
			FullDescription: "Немецкий ленточный погрузчик Mulag Orbiter 12D/E предназначен для эффективной погрузки багажа в самолеты. Надежность и долговечность гарантированы. Характеристики: грузоподъемность 2т, электрический привод, плавный ход, низкий уровень шума.",
			Unit:            "руб./час",
		},
		{
			ID:              4,
			Title:           "Буксировщик",
			Description:     "Тягач Challenger 550",
			ImageURL:        "http://localhost:9000/test/tyagach.jpg",
			Price:           "260000 руб./сут",
			FullDescription: "Буксировщик для воздушных судов обеспечивает безопасное и точное перемещение самолетов по территории аэропорта. Соответствует всем требованиям безопасности. Характеристики: мощность 400 л.с., гидравлическая система, точное управление, безопасность.",
			Unit:            "руб./сут",
		},
		{
			ID:              5,
			Title:           "Автолифт",
			Description:     "Амбулаторный автолифт DOLL X-PRM M",
			ImageURL:        "http://localhost:9000/test//lift.jpg",
			Price:           "150000 руб./сут",
			FullDescription: "Полный ассортимент оригинальных запчастей для различных моделей самолетов. Все детали проходят строгий контроль качества. Характеристики: оригинальное качество, сертификация, гарантия, быстрая поставка.",
			Unit:            "руб./сут",
		},
		{
			ID:              6,
			Title:           "Антифриз",
			Description:     "MAXCool Hybrid",
			ImageURL:        "http://localhost:9000/test/noice.jpg",
			Price:           "68100 руб./бочка",
			FullDescription: "Высококачественное авиационное масло для турбовинтовых и турбореактивных двигателей. Обеспечивает надежную работу в экстремальных условиях. Характеристики: высокая вязкость, термостойкость, защита от износа, долгий срок службы.",
			Unit:            "руб./бочка",
		},
		{
			ID:              7,
			Title:           "Деайсер",
			Description:     "ELEPAHT BETA-15",
			ImageURL:        "http://localhost:9000/test/antifrizmash.jpg",
			Price:           "125000 руб./сут",
			FullDescription: "Специализированная машина для удаления льда с воздушных судов. Обеспечивает безопасность полетов в зимних условиях. Характеристики: грузоподъемность 2т, электрический привод, плавный ход, низкий уровень шума.",
			Unit:            "руб./сут",
		},
	}

	return aviationServices, nil
}

func (r *Repository) GetSelectAviaSub(id int) (ASGARserv, error) {
	aviationServices, err := r.GetAsgarList()
	if err != nil {
		return ASGARserv{}, err
	}

	for _, selsubavia := range aviationServices {
		if selsubavia.ID == id {
			return selsubavia, nil
		}
	}
	return ASGARserv{}, fmt.Errorf("заказ не найден")
}
