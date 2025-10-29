package entitities
type Organization struct {
	BaseEntity
	Name string `gorm:"type:varchar(255); null"`
	Domain string `gorm:"type:varchar(255); null"`
	Logo string `gorm:"type:varchar(255); null"`
	Description string `gorm:"type:varchar(255); null"`
	Location string `gorm:"type:varchar(255); null"`
	Industry string `gorm:"type:varchar(255); null"`
	Size int `gorm:"type:integer; null"`
	Email string `gorm:"type:varchar(255); null; default:null"`
}

func (Organization) TableName() string {
	return "organizations"
}
