package repository

import (
	"errors"
	"main/internal/ds"
	"strings"

	"gorm.io/gorm"
)

func (r *Repository) BanknoteList(id uint) (*[]ds.Banknote, int64, error) {
	var banknotes []ds.Banknote
	var count int64

	if err := r.db.Table("operation_banknotes").Where("operation_id = ?", id).Count(&count).Error; err != nil {
		return &banknotes, 0, err
	}

	if err := r.db.Where("status = ?", "действует").Find(&banknotes).Error; err != nil {
		return nil, 0, err
	}

	return &banknotes, count, nil
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
	res := r.db.Where("user_id = ?", creatorID).Where("status = ?", draftStatus).Take(&draftReq)

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

func (r *Repository) GetOprationDraftID(creatorID uint) (uint, error) {
	var draftReq ds.Operation

	res := r.db.Where("user_id = ?", creatorID).Where("status = ?", draftStatus).Take(&draftReq)
	if errors.Is(gorm.ErrRecordNotFound, res.Error) {
		return 0, nil
	}

	if res.Error != nil {
		return 0, res.Error
	}

	return draftReq.ID, nil
}

func (r *Repository) CreateOperationDraft(creatorID uint) (uint, error) {
	request := ds.Operation{
		// ModeratorID:  creatorID, // просто заглушка, потом придумаю, как сделать норм
		UserID:    creatorID,
		Status:    "черновик",
		CreatedAt: r.db.NowFunc(),
	}

	if err := r.db.Create(&request).Error; err != nil {
		return 0, err
	}
	return request.ID, nil
}

func (r *Repository) BanknoteById(id int) (*ds.Banknote, error) {
	var banknote ds.Banknote
	if err := r.db.First(&banknote, id).Error; err != nil {
		return nil, err
	}
	return &banknote, nil
}

func (r *Repository) DeleteBanknote(id string) error {
	query := "UPDATE banknotes SET is_delete=t WHERE id=$1;"
	result := r.db.Exec(query, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *Repository) GetBanknoteById(id uint) (*ds.Banknote, error) {
	var banknote ds.Banknote
	if err := r.db.Where("status = ?", "действует").First(&banknote, id).Error; err != nil {
		return nil, err
	}
	return &banknote, nil
}
