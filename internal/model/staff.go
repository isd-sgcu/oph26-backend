package model

type Staff struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	FirstName string `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName  string `gorm:"type:varchar(100);not null" json:"last_name"`
	Nickname  string `gorm:"type:varchar(100)" json:"nickname"`
	TelNumber string `gorm:"type:varchar(15);not null" json:"tel_number"`
	Faculty   string `gorm:"type:varchar(100);not null" json:"faculty"`
	Email     string `gorm:"type:varchar(100);not null;unique" json:"email"`
	Year      int    `gorm:"not null" json:"year"`
}
