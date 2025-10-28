package entitities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseEntity struct {
	ID        string    `gorm:"primaryKey;type:varchar(255)"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (b *BaseEntity) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New().String()
	return
}