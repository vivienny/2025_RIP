package ds

type Miniplane struct {
	ID       int   `gorm:"primaryKey"`
	UserID   uint  `gorm:"not null"`
	IsActive bool  `gorm:"type:boolean;default:true"`
	User     Avius `gorm:"foreignKey:UserID"`
}
