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
	ID               uint    `json:"user_id" gorm:"primaryKey"`
	Login            string  `json:"login" gorm:"uniqueIndex;not null;type:TEXT;check:LENGTH(login) <= 255"`
	Password         string  `json:"password" gorm:"not null;type:TEXT;check:LENGTH(password) <= 255"`
	Name             string  `json:"name" gorm:"not null;type:TEXT;check:LENGTH(name) <= 255"`
	Balance          float64 `json:"balance" gorm:"type:NUMERIC;default:0"`
	Role             Role    `json:"role" sql:"type:string;`
	RegistrationDate time.Time
}

// Operation represents the 'operation' table
type Operation struct {
	ID                uint                `json:"operation_id" gorm:"primaryKey"`
	UserID            uint                `json:"user_id"`
	User              User                `json:"user" gorm:"foreignkey:UserID"`
	ModeratorID       *uint               `json:"moderator_id"`
	Moderator         User                `json:"moderator" gorm:"foreignkey:ModeratorID"`
	Name              string              `json:"is_income" gorm:"type:varchar(255)"`
	OperationBanknote []OperationBanknote `json:"operation_banknote" gorm:"foreignkey:OperationID"`
	Status            string              `json:"status" gorm:"type:TEXT"`
	StatusCheck       string              `json:"status_check"`
	CreatedAt         time.Time           `json:"crated_at" gorm:"type:TIMESTAMPTZ;default:CURRENT_TIMESTAMP;not null"`
	FormationAt       time.Time           `json:"formation_at" gorm:"type:TIMESTAMPTZ;default:CURRENT_TIMESTAMP"`
	CompletionAt      time.Time           `json:"comletion_at" gorm:"type:TIMESTAMPTZ"`
}

type OperationResponse struct {
	ID                uint                `json:"id"`
	Name              string              `json:"operation_name"`
	Status            string              `json:"status"`
	StatusCheck       string              `json:"status_check"`
	OperationBanknote []OperationBanknote `json:"operation_banknote"`
	UserName          string              `json:"user_name"`
	UserLogin         string              `json:"user_login"`
	ModeratorName     string              `json:"moderator_name"`
	ModeratorLogin    string              `json:"moderator_login"`
	FormationAt       time.Time           `json:"formation_at"`
	CreatedAt         time.Time           `json:"crated_at" `
	CompletionAt      time.Time           `json:"comletion_at"`
}

type NewStatus struct {
	Status      string `json:"status"`
	OperationID uint   `json:"operation_id"`
}

type OperationdetailsDetails struct {
	Banknote  *Banknote    `json:"banknote"`
	Operation *[]Operation `json:"operation"`
}

type RequestAsyncService struct {
	RequestId uint   `gorm:"primaryKey" json:"requestId"`
	Token     string `json:"server_Token"`
	Status    string `json:"status"`
}

type UpdateOperation struct {
	ID   uint   `json:"id"`
	Name string `json:"operation_name"`
}

// Banknote represents the 'banknote' table
type Banknote struct {
	ID          uint    `json:"banknote_id" gorm:"primaryKey"`
	Nominal     float64 `json:"nominal" gorm:"type:NUMERIC"`
	Currency    string  `json:"currency" gorm:"type:VARCHAR(3)"`
	Description string  `json:"description" gorm:"type:TEXT"`
	ImageURL    string  `json:"image_url" gorm:"type:TEXT;not null;check:LENGTH(image_url) <= 500"`
	Status      string  `json:"status" gorm:"type:VARCHAR(50)"`
}

type BanknoteList struct {
	DraftID   uint        `json:"draft_id"`
	Banknotes *[]Banknote `json:"banknotes_list"`
}

type AddToBanknoteID struct {
	BanknoteID uint `json:"banknote_id"`
	Quantity   int  `json:"quantity"`
}

// OperationBanknote represents the 'operation_banknote' table
type OperationBanknote struct {
	ID          uint      `json:"opration_banknote_id" gorm:"primaryKey"`
	OperationID uint      `json:"operation_id"`
	BanknoteID  uint      `json:"banknote_id"`
	Operations  Operation `json:"opration" gorm:"foreignKey:OperationID"`
	Banknote    Banknote  `json:"banknote" gorm:"foreignKey:BanknoteID"`
	Quantity    int       `json:"quantity" gorm:"type:int;default:1;check:quantity >= 1"`
}

type OperationBanknoteUpdate struct {
	ID       uint `json:"id" gorm:"primary_key"`
	Quantity int  `json:"quantity" gorm:"type:int;"`
}
