package ds

type Subjserv struct {
	ID          uint   `gorm:"primaryKey"`
	MiniplaneID uint   `gorm:"not null;uniqueIndex:idx_miniplane_service"`
	ServiceID   uint   `gorm:"not null;uniqueIndex:idx_miniplane_service"`
	Quantity    int    `gorm:"default:1"`
	Price       string `gorm:"type:varchar(15)"`

	Miniplane Miniplane    `gorm:"foreignKey:MiniplaneID"`
	Service   ASGARService `gorm:"foreignKey:ServiceID"`
}
