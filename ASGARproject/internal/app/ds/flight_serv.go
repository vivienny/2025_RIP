package ds

import (
	"database/sql"
	"time"
)

type FlightServ struct {
	ID         uint      `gorm:"primaryKey"`
	Status     string    `gorm:"type:varchar(15);not null"` // "pending", "confirmed", "completed"
	TotalPrice string    `gorm:"type:varchar(15)"`
	DateCreate time.Time `gorm:"not null"`
	DateUpdate time.Time
	DateFinish sql.NullTime
	UserID     uint `gorm:"not null"`

	User Avius `gorm:"foreignKey:UserID"`
}
