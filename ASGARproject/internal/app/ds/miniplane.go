package ds

type Miniplane struct {
	ID       int    `gorm:"primaryKey"`
	UserID   uint   `gorm:"not null"`
	IsActive bool   `gorm:"type:boolean;default:true"`
	Total    string `gorm:"type:varchar(15)"`

	User Avius `gorm:"foreignKey:UserID"` // ← меняем Users на Avius
}
