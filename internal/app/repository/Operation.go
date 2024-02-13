package repository

import (
	"errors"
	"fmt"
	"main/internal/app/ds"
	"main/internal/app/utils"
	"time"

	"gorm.io/gorm"
)

func (r *Repository) GetOprationDraftID(creatorID uint) (uint, error) {
	var draftReq ds.Operation

	res := r.db.Where("user_id = ?", creatorID).Where("status = ?", utils.Draft).Take(&draftReq)
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

func (r *Repository) GetOperationWithDataByID(requestID uint) (ds.Operation, []ds.Banknote, error) {
	if requestID == 0 {
		return ds.Operation{}, nil, errors.New("record not found")
	}

	request := ds.Operation{ID: requestID}
	res := r.db.Take(&request)
	if err := res.Error; err != nil {
		return ds.Operation{}, nil, err
	}

	var dataService []ds.Banknote

	res = r.db.
		Table("operation_banknotes").
		Select("banknotes.*").
		Where("status != ?", "удалён").
		Joins("JOIN banknotes ON operation_banknotes.\"banknote_id\" = banknotes.id").
		Where("operation_banknotes.\"operation_id\" = ?", requestID).
		Find(&dataService)

	if err := res.Error; err != nil {
		return ds.Operation{}, nil, err
	}

	return request, dataService, nil
}

func (r *Repository) OperationList(status, start, end string) (*[]ds.Operation, error) {
	var operation []ds.Operation
	query := r.db.Where("status != ? AND status != ?", "удалён", "черновик")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if start != "" {
		query = query.Where("created_at >= ?", start)
	}

	if end != "" {
		query = query.Where("created_at <= ?", end)
	}
	result := query.Find(&operation)
	return &operation, result.Error
}

func (r *Repository) UpdateOperation(updatedOperation *ds.Operation) error {
	oldOperation := ds.Operation{}
	if result := r.db.First(&oldOperation, updatedOperation.ID); result.Error != nil {
		return result.Error
	}

	oldOperation.IsIncome = updatedOperation.IsIncome

	if updatedOperation.CreatedAt.String() != utils.EmptyDate {
		oldOperation.CreatedAt = updatedOperation.CreatedAt
	}
	if updatedOperation.CompletionAt.String() != utils.EmptyDate {
		oldOperation.CompletionAt = updatedOperation.CompletionAt
	}
	if updatedOperation.FormationAt.String() != utils.EmptyDate {
		oldOperation.FormationAt = updatedOperation.FormationAt
	}
	if updatedOperation.Status != "" {
		oldOperation.Status = updatedOperation.Status
	}

	*updatedOperation = oldOperation
	result := r.db.Save(updatedOperation)
	return result.Error
}

func (r *Repository) FormOperationRequestByID(requestID uint, creatorID uint) error {
	var req ds.Operation
	res := r.db.
		Where("id = ?", requestID).
		Where("user_id = ?", creatorID).
		Where("status = ?", utils.Draft).
		Take(&req)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("нет такой заявки")
	}

	req.Status = "сформирован"
	req.FormationAt = time.Now()
	req.ModeratorID = 1
	if err := r.db.Save(&req).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) RejectOperationRequestByID(requestID, moderatorID uint) error {
	return r.finishRejectHelper("отклонён", requestID, moderatorID)
}

func (r *Repository) FinishEncryptDecryptRequestByID(requestID, moderatorID uint) error {
	return r.finishRejectHelper("завершён", requestID, moderatorID)
}

func (r *Repository) finishRejectHelper(status string, requestID, moderatorID uint) error {
	var req ds.Operation
	res := r.db.
		Where("id = ?", requestID).
		Where("status = ?", "сформирован").
		Take(&req)

	if res.Error != nil {
		return errors.New("нет такой заявки")
	}
	if res.RowsAffected == 0 {
		return errors.New("нет такой заявки")
	}

	req.ModeratorID = moderatorID
	req.Status = status

	req.CompletionAt = time.Now()

	if err := r.db.Save(&req).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteOperationByID(requestID uint) error { // ?
	var req ds.Operation
	res := r.db.
		Where("id = ?", requestID). // ??
		Take(&req)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("нет такой заявки")
	}

	req.Status = "удалён"
	// delTime := time.Now()
	req.CompletionAt = time.Now()
	if err := r.db.Save(&req).Error; err != nil {
		return err
	}
	// .Update("deleted_at", delTime)
	if err := r.db.Model(&req).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteBanknoteFromRequest(deleteFromOperation ds.OperationBanknote) (ds.Operation, []ds.Banknote, error) {
	var deletedBanknoteOperation ds.OperationBanknote
	result := r.db.Where("\"BanknoteID\" = ? and \"OperationID\" = ?", deleteFromOperation.BanknoteID,
		deleteFromOperation.OperationID).Find(&deletedBanknoteOperation)
	if result.Error != nil {
		return ds.Operation{}, nil, result.Error
	}

	if result.RowsAffected == 0 {
		return ds.Operation{}, nil, fmt.Errorf("record not found")
	}
	if err := r.db.Delete(&deletedBanknoteOperation).Error; err != nil {
		return ds.Operation{}, nil, err
	}

	return r.GetOperationWithDataByID(deleteFromOperation.OperationID)
}

func (r *Repository) UpdateOperationBanknote(OperationID uint, BanknoteID uint, quantity int) error {
	var updateBanknote ds.OperationBanknote
	r.db.Where(" \"OperationID\" = ? and \"BanknoteID\" = ?", OperationID, BanknoteID).First(&updateBanknote)

	if updateBanknote.OperationID == 0 {
		return errors.New("нет такой заявки")
	}
	updateBanknote.Quantity = quantity

	if err := r.db.Save(&updateBanknote).Error; err != nil {
		return err
	}

	return nil
}
