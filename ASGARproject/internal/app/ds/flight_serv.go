package ds

import (
	"database/sql"
	"time"
)

type FlightServ struct {
	ID         uint
	Status     string    `gorm:"type:varchar(15);not null;default:'cart'"` // ← ДОБАВЛЯЕМ
	TotalPrice string    `gorm:"type:varchar(15)"`
	DateCreate time.Time `gorm:"not null"`
	DateUpdate time.Time
	DateFinish sql.NullTime
	UserID     uint `gorm:"not null"`

	User      Avius      `gorm:"foreignKey:UserID"`
	Subjservs []Subjserv `gorm:"foreignKey:FlightServID"` // ← ДОБАВЛЯЕМ
}
