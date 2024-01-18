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
	ID          uint    `json:"user_id" gorm:"primaryKey"`
	Login       string  `json:"login" gorm:"uniqueIndex;not null;type:TEXT;check:LENGTH(login) <= 255"`
	Password    string  `json:"password" gorm:"not null;type:TEXT;check:LENGTH(password) <= 255"`
	Balance     float64 `json:"balance" gorm:"type:NUMERIC;default:0"`
	IsModerator bool    `json:"is_moderator" gorm:"type:BOOLEAN"`
}

// Operation represents the 'operation' table
type Operation struct {
	ID           uint      `json:"operation_id" gorm:"primaryKey"`
	UserID       uint      `json:"user_id"`
	User         User      `json:"user" gorm:"foreignkey:UserID"`
	ModeratorID  uint      `json:"moderator_id" gorm:"default:NULL"`
	Moderator    User      `json:"moderator" gorm:"foreignkey:ModeratorID"`
	IsIncome     bool      `json:"is_income" gorm:"type:BOOLEAN"`
	Status       string    `json:"status" gorm:"type:TEXT"`
	CreatedAt    time.Time `json:"crated_at" gorm:"type:TIMESTAMPTZ;default:CURRENT_TIMESTAMP;not null"`
	FormationAt  time.Time `json:"formation_at" gorm:"type:TIMESTAMPTZ;default:CURRENT_TIMESTAMP"`
	CompletionAt time.Time `json:"comletion_at" gorm:"type:TIMESTAMPTZ"`
}

// Banknote represents the 'banknote' table
type Banknote struct {
	ID          uint    `json:"banknote_id" gorm:"primaryKey"`
	Nominal     float64 `json:"nominal" gorm:"type:NUMERIC"`
	Currency    string  `json:"currency" gorm:"type:VARCHAR(3)"`
	Description string  `json:"description" gorm:"type:TEXT"`
	ImageURL    string  `json:"image_url" gorm:"type:TEXT;not null;check:LENGTH(image_url) <= 500"`
	Status      string  `json:"strus" gorm:"type:VARCHAR(50)"`
}

type BanknoteList struct {
	DraftID   uint        `json:"draft_id"`
	Banknotes *[]Banknote `json:"banknotes_list"`
}

// OperationBanknote represents the 'operation_banknote' table
type OperationBanknote struct {
	ID          uint      `json:"opration_banknote_id" gorm:"primaryKey"`
	OperationID uint      `json:"operation_id" gorm:"not null"`
	BanknoteID  uint      `json:"banknote_id" gorm:"not null"`
	Operations  Operation `json:"opration" gorm:"foreignKey:OperationID"`
	Banknote    Banknote  `json:"application" gorm:"foreignKey:BanknoteID"`
	Quantity    int       `json:"quantity" gorm:"default:1;check:quantity >= 1"`
}
