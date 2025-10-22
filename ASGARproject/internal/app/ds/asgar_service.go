package ds

type ASGARService struct {
	ID              int    `gorm:"primaryKey"`
	IsDelete        bool   `gorm:"type:boolean;default:false"`
	Img             string `gorm:"type:varchar(100)"`
	Name            string `gorm:"type:varchar(25)"` // ← УБИРАЕМ NOT NULL
	Info            string `gorm:"type:varchar(100)"`
	Price           string `gorm:"type:varchar(15)"` // ← УБИРАЕМ NOT NULL
	FullDescription string `gorm:"type:text"`
	Unit            string `gorm:"type:varchar(10)"`
}
