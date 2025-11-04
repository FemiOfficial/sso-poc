package entitities
type Organization struct {
	BaseEntity
	Name string `gorm:"type:varchar(255); null" json:"name"`
	Domain string `gorm:"type:varchar(255); null" json:"domain"`
	Logo string `gorm:"type:varchar(255); null" json:"logo"`
	Description string `gorm:"type:varchar(255); null" json:"description"`
	Location string `gorm:"type:varchar(255); null" json:"location"`
	Industry string `gorm:"type:varchar(255); null" json:"industry"`
	Size int `gorm:"type:integer; null" json:"size"`
	Email string `gorm:"type:varchar(255); null; default:null" json:"email"`
}

func (Organization) TableName() string {
	return "organizations"
}
