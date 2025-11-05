package ds

type Subjserv struct {
	ID           uint   `gorm:"primaryKey"`
	FlightServID uint   `gorm:"not null;uniqueIndex:idx_flight_service"` // ← МЕНЯЕМ miniplane_id
	ServiceID    uint   `gorm:"not null;uniqueIndex:idx_flight_service"`
	Quantity     int    `gorm:"default:1"`
	Price        string `gorm:"type:varchar(15)"`

	FlightServ FlightServ   `gorm:"foreignKey:FlightServID"` // ← МЕНЯЕМ
	Service    ASGARService `gorm:"foreignKey:ServiceID"`
}
