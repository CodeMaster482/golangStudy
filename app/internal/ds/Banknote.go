package ds

import (
	"gorm.io/gorm"
)

/*
 	CREATE TABLE Banknotes (
     	banknote_id SERIAL PRIMARY KEY,
     	denomination INT NOT NULL,
     	currency VARCHAR(50) NOT NULL
 		);
*/

type Banknote struct {
	gorm.Model
	Id           uint   `json:"banknote_id" gorm:"primary_key"`
	Denomination int    `json:"denomination" gorm:"type:int"`
	Currency     string `json:"currency" gorm:"type:text"`
}
