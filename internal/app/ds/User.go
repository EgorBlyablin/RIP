package ds

type User struct {
	ID          uint
	Login       string `gorm:"type:varchar(30);unique;not null"`
	Password    string `gorm:"type:varchar(100);not null"`
	IsModerator bool   `gorm:"default:false;not null"`

	GenerationRequests []GenerationRequest `gorm:"foreignKey:CreatedByID"`
}
