package repository

import (
	"errors"
	"main/internal/app/ds"
	"main/internal/app/utils"
	"strconv"
	"strings"
	"time"
)

func (r *Repository) GetOpenBanknotes() (*[]ds.Banknote, error) {
	var tenders []ds.Banknote
	if err := r.db.Where("status = ?", "действует").Find(&tenders).Error; err != nil {
		return nil, err
	}
	return &tenders, nil
}

func (r *Repository) BanknotesList(name string) (*[]ds.Banknote, error) {
	//name = strings.ToLower(name)

	var serchval float64
	var banknotes []ds.Banknote

	if name != "" {
		serchval, _ = strconv.ParseFloat(name, 64)
		if err := r.db.Where("nominal = ? AND status != ?", serchval, "удален").Find(&banknotes).Error; err != nil {
			return nil, err
		}
	} else {
		if err := r.db.Where("status != ?", "удален").Find(&banknotes).Error; err != nil {
			return nil, err
		}
	}
	return &banknotes, nil
}

func (r *Repository) AddBanknote(banknote *ds.Banknote) (uint, error) {
	banknote.Status = "действует"
	//	fmt.Print(banknote.ID)
	result := r.db.Create(&banknote)
	return banknote.ID, result.Error
}

/*
func (r *Repository) SearchBanknotes(search string) (*[]ds.Banknote, error) {
	var banknotes []ds.Banknote
	r.db.Find(&banknotes)

	var filteredBanknotes []ds.Banknote
	for _, banknote := range banknotes {
		if strings.Contains(strings.ToLower(string(banknote.Amount)), strings.ToLower(search)) {
			filteredBanknotes = append(filteredBanknotes, banknote)
		}
	}

	return &filteredBanknotes, nil
}
*/

func (r *Repository) GetBanknoteById(id uint) (*ds.Banknote, error) {
	var company ds.Banknote
	if err := r.db.Where("status = ?", "действует").First(&company, id).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

// func (r *Repository) DeleteBanknote(id string) error {
// 	query := "UPDATE banknotes SET is_delete=t WHERE id=$1;"
// 	result := r.db.Exec(query, id)
// 	if result.Error != nil {
// 		return result.Error
// 	}
// 	return nil
// }

func (r *Repository) DeleteBanknote(id uint) error {
	banknote := ds.Banknote{}

	if err := r.db.First(&banknote, "id = ?", id).Error; err != nil {
		return err
	}

	if err := r.db.Model(&banknote).Update("status", "удален").Error; err != nil {
		return err
	}

	if err := r.db.Model(&banknote).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateBanknote(updatedBanknote *ds.Banknote) (*ds.Banknote, error) {
	var oldBanknote ds.Banknote

	if result := r.db.First(&oldBanknote, updatedBanknote.ID); result.Error != nil {
		return updatedBanknote, result.Error
	}

	if updatedBanknote.Nominal != 0 {
		oldBanknote.Nominal = updatedBanknote.Nominal
	}

	if updatedBanknote.ImageURL != "" {
		oldBanknote.ImageURL = updatedBanknote.ImageURL
	}

	if updatedBanknote.Description != "" {
		oldBanknote.Description = updatedBanknote.Description
	}

	if updatedBanknote.Currency != "" {
		oldBanknote.Currency = updatedBanknote.Currency
	}

	*updatedBanknote = oldBanknote
	result := r.db.Save(updatedBanknote)
	return updatedBanknote, result.Error
}

func (r *Repository) DeleteBanknoteImage(banknoteId uint) string {
	banknote := ds.Banknote{}

	r.db.First(&banknote, "id = ?", banknoteId)
	return banknote.ImageURL
}

func (r *Repository) AddBanknoteToDraft(dataID uint, creatorID uint) (uint, error) {
	// получаем услугу
	data, err := r.GetBanknoteById(dataID)
	if err != nil {
		return 0, err
	}

	if data == nil {
		return 0, errors.New("нет такой услуги")
	}
	if data.Status == "удален" {
		return 0, errors.New("услуга удалена")
	}

	// получаем черновик
	var draftReq ds.Operation
	res := r.db.Where("user_id = ?", creatorID).Where("status = ?", utils.Draft).Take(&draftReq)

	// создаем черновик, если его нет
	if res.RowsAffected == 0 {
		newDraftRequestID, err := r.CreateOperationDraft(creatorID)
		if err != nil {
			return 0, err
		}

		draftReq.ID = newDraftRequestID
	}

	// добавляем запись в мм
	requestToData := ds.OperationBanknote{
		OperationID: draftReq.ID,
		BanknoteID:  dataID,
		Quantity:    0,
	}

	err = r.db.Create(&requestToData).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return 0, errors.New("услуга уже существует в заявке")
		}

		return 0, err
	}

	return draftReq.ID, nil
}
