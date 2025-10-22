package ds

type ASGARService struct {
	ID              int    `gorm:"primaryKey" json:"ID"`
	IsDelete        bool   `gorm:"type:boolean;default:false" json:"IsDelete"`
	Img             string `gorm:"type:varchar(100)" json:"Img"`
	Name            string `gorm:"type:varchar(25)" json:"Name"`
	Info            string `gorm:"type:varchar(100)" json:"Info"`
	Price           string `gorm:"type:varchar(15)" json:"Price"`
	FullDescription string `gorm:"type:text" json:"FullDescription"`
	Unit            string `gorm:"type:varchar(10)" json:"Unit"`
}
