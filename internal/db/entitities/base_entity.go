package entitities

import (
	"time"

	"database/sql/driver"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return pq.StringArray(a).Value()
}

// Scan implements sql.Scanner interface
func (a *StringArray) Scan(value interface{}) error {
	var pqArray pq.StringArray
	if err := pqArray.Scan(value); err != nil {
		return err
	}
	*a = StringArray(pqArray)
	return nil
}

type BaseEntity struct {
	ID        string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (b *BaseEntity) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New().String()
	return
}
