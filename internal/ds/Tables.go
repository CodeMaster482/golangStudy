package ds

import (
	"time"
)

/*
Банкомат.
Услуги - различные виды купюр,
заявки - операции внесения/снятия наличных
*/

// User represents the 'user' table
type User struct {
	ID          uint    `gorm:"primaryKey"`
	Login       string  `gorm:"uniqueIndex;not null;type:TEXT;check:LENGTH(login) <= 255"`
	Password    string  `gorm:"not null;type:TEXT;check:LENGTH(password) <= 255"`
	Balance     float64 `gorm:"type:NUMERIC;default:0"`
	IsModerator bool    `gorm:"type:BOOLEAN"`
}

// Operation represents the 'operation' table
type Operation struct {
	ID           uint       `gorm:"primaryKey"`
	UserID       uint       `gorm:"not null"`
	IsIncome     bool       `gorm:"type:BOOLEAN"`
	Status       string     `gorm:"type:TEXT"`
	CreatedAt    time.Time  `gorm:"type:TIMESTAMPTZ;default:CURRENT_TIMESTAMP;not null"`
	FormationAt  *time.Time `gorm:"type:TIMESTAMPTZ;default:CURRENT_TIMESTAMP"`
	CompletionAt *time.Time `gorm:"type:TIMESTAMPTZ"`
}

// Banknote represents the 'banknote' table
type Banknote struct {
	ID          uint    `gorm:"primaryKey"`
	Nominal     float64 `gorm:"type:NUMERIC"`
	Description string  `gorm:"type:TEXT"`
	ImageURL    string  `gorm:"type:TEXT;not null;check:LENGTH(image_url) <= 500"`
	IsDelete    bool    `gorm:"type:BOOLEAN"`
}

// OperationBanknote represents the 'operation_banknote' table
type OperationBanknote struct {
	ID          uint `gorm:"primaryKey"`
	OperationID uint `gorm:"not null"`
	BanknoteID  uint `gorm:"not null"`
	Quantity    int  `gorm:"default:1;check:quantity >= 1"`
}
