package repository

import (
	"errors"
	"main/internal/app/ds"
	"main/internal/app/utils"
	"time"

	"gorm.io/gorm"
)

func (r *Repository) OperationByUserID(userID string) (*[]ds.OperationResponse, error) {
	var operations []ds.Operation
	var operationResponses = []ds.OperationResponse{}
	result := r.db.Preload("User").
		//Preload("OperationBanknotes.Operation.User").
		Preload("Moderator").
		Where("user_id = ? AND status != 'удален' AND status != 'черновик'", userID).
		Find(&operations)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}

	for _, op := range operations {
		operationResponse := ds.OperationResponse{
			ID:        op.ID,
			Name:      op.Name,
			UserName:  op.User.Name,
			UserLogin: op.User.Login,
			//UserRole:       op.User.Role,
			ModeratorName: op.Moderator.Name,
			Status:        op.Status,
			StatusCheck:   op.StatusCheck,
			//ModeratorRole:  op.Moderator.Role,
			ModeratorLogin:    op.Moderator.Login,
			CreatedAt:         op.CreatedAt,
			FormationAt:       op.FormationAt,
			CompletionAt:      op.CompletionAt,
			OperationBanknote: op.OperationBanknotes,
		}
		operationResponses = append(operationResponses, operationResponse)
	}

	return &operationResponses, result.Error
}

func (r *Repository) OperationByID(id uint) (*ds.OperationResponse, error) {
	operation := ds.Operation{}
	result := r.db.Preload("User").
		// Preload("OperationBanknotes.Operations").
		Preload("OperationBanknotes.Banknote").
		First(&operation, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}

	operationResponse := ds.OperationResponse{
		ID:        operation.ID,
		Name:      operation.Name,
		UserName:  operation.User.Name,
		UserLogin: operation.User.Login,
		//UserRole:       operation.User.Role,
		ModeratorName: operation.Moderator.Name,
		Status:        operation.Status,
		StatusCheck:   operation.StatusCheck,
		//ModeratorRole:  operation.Moderator.Role,
		ModeratorLogin:    operation.Moderator.Login,
		CreatedAt:         operation.CreatedAt,
		FormationAt:       operation.FormationAt,
		CompletionAt:      operation.CompletionAt,
		OperationBanknote: operation.OperationBanknotes,
	}
	return &operationResponse, result.Error
}

func (r *Repository) OperationModel(id uint) (*ds.Operation, error) {
	operation := ds.Operation{}

	result := r.db.Preload("User").
		//Preload("OperationBanknotes.Operations").
		Preload("OperationBanknotes.Banknote").
		First(&operation, id)
	return &operation, result.Error
}

func (r *Repository) OperationDraftId(userId uint) (uint, error) {
	var operation ds.Operation
	result := r.db.
		Where("status = ? AND user_id = ?", "черновик", userId).
		First(&operation)
	if result.RowsAffected == 0 {
		return 0, nil
	}
	return operation.ID, result.Error
}

func (r *Repository) GetOperationDraftID(creatorID uint) (uint, error) {
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
	//userInfo, err := GetUserInfo(r, creatorID)
	//if err != nil {
	//	return 0, err
	//}
	request := ds.Operation{
		UserID:      creatorID,
		Status:      "черновик",
		CreatedAt:   r.db.NowFunc(),
		ModeratorID: nil,
		//CreatorLogin: userInfo.Login,
	}

	if err := r.db.Create(&request).Error; err != nil {
		return 0, err
	}
	return request.ID, nil
}

func (r *Repository) GetOperationWithDataByID(requestID uint, userId uint, isAdmin bool) (ds.Operation, []ds.Banknote, error) {
	var OperationRequest ds.Operation
	var banknotes []ds.Banknote

	//ищем такую заявку
	result := r.db.First(&OperationRequest, "id =?", requestID)
	if result.Error != nil {
		r.logger.Error("error while getting monitoring request")
		return ds.Operation{}, nil, result.Error
	}
	if !isAdmin && OperationRequest.UserID == uint(userId) || isAdmin {
		res := r.db.
			Table("operation_banknotes").
			Select("banknotes.*").
			Where("status != ?", "удалён").
			Joins("JOIN banknotes ON operation_banknotes.\"BanknoteID\" = banknotes.id").
			Where("operation_banknotes.\"OperationID\" = ?", requestID).
			Find(&banknotes)
		if res.Error != nil {
			r.logger.Error("error while getting for tender request")
			return ds.Operation{}, nil, res.Error
		}
	} else {
		return ds.Operation{}, nil, errors.New("ошибка доступа к данной заявке")
	}

	return OperationRequest, banknotes, nil
}

func (r *Repository) OperationsList(statusID string, startDate time.Time, endDate time.Time) (*[]ds.OperationResponse, error) {
	var operations []ds.Operation
	operationResponses := []ds.OperationResponse{}
	if statusID == "" {
		result := r.db.
			Preload("User").
			Preload("Moderator").
			Where("status != 'удален' AND status != 'черновик' AND created_at BETWEEN ? AND ?", startDate, endDate).
			Order("id DESC").
			Find(&operations)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil, nil
			} else {
				return nil, result.Error
			}
		}

		for _, op := range operations {
			operationResponse := ds.OperationResponse{
				ID:        op.ID,
				Name:      op.Name,
				UserName:  op.User.Name,
				UserLogin: op.User.Login,
				//UserRole:       op.User.Role,
				ModeratorName: op.Moderator.Name,
				Status:        op.Status,
				StatusCheck:   op.StatusCheck,
				//ModeratorRole:  op.Moderator.Role,
				ModeratorLogin:    op.Moderator.Login,
				CreatedAt:         op.CreatedAt,
				FormationAt:       op.FormationAt,
				CompletionAt:      op.CompletionAt,
				OperationBanknote: op.OperationBanknotes,
			}
			operationResponses = append(operationResponses, operationResponse)
		}

		return &operationResponses, result.Error
	}

	result := r.db.
		Preload("User").
		Where("status = ? AND status != 'черновик' AND creation_date BETWEEN ? AND ?", statusID, startDate, endDate).
		Find(&operations)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}

	for _, op := range operations {
		operationResponse := ds.OperationResponse{
			ID:        op.ID,
			Name:      op.Name,
			UserName:  op.User.Name,
			UserLogin: op.User.Login,
			Status:    op.Status,
			//UserRole:       op.User.Role,
			ModeratorName: op.Moderator.Name,
			//ModeratorRole:  op.Moderator.Role,
			ModeratorLogin:    op.Moderator.Login,
			StatusCheck:       op.StatusCheck,
			CreatedAt:         op.CreatedAt,
			FormationAt:       op.FormationAt,
			CompletionAt:      op.CompletionAt,
			OperationBanknote: op.OperationBanknotes,
		}
		operationResponses = append(operationResponses, operationResponse)
	}

	return &operationResponses, result.Error
}

func (r *Repository) UpdateOperation(updatedOperation *ds.Operation) error {
	oldOperation := ds.Operation{}
	if result := r.db.First(&oldOperation, updatedOperation.ID); result.Error != nil {
		return result.Error
	}
	if updatedOperation.Name != "" {
		oldOperation.Name = updatedOperation.Name
	}
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

func (r *Repository) FormOperationRequestByIDAsynce(id uint, creatorID uint) (uint, error) {
	var req ds.Operation
	res := r.db.
		Where("id = ?", id).
		Where("user_id = ?", creatorID).
		//Where("status = ?", utils.Draft).
		Take(&req)

	if res.Error != nil {
		return 0, res.Error
	}
	if res.RowsAffected == 0 {
		return 0, errors.New("нет такой заявки")
	}

	req.StatusCheck = "В обработке"
	req.Status = "сформирован"
	req.FormationAt = time.Now()

	if err := r.db.Save(&req).Error; err != nil {
		return 0, err
	}

	return req.ID, nil
}

func (r *Repository) FormOperationRequestByID(creatorID uint) (uint, error) {
	var req ds.Operation
	res := r.db.
		//Where("id = ?", requestID).
		Where("user_id = ?", creatorID).
		Where("status = ?", utils.Draft).
		Take(&req)

	if res.Error != nil {
		return 0, res.Error
	}
	if res.RowsAffected == 0 {
		return 0, errors.New("нет такой заявки")
	}
	req.StatusCheck = "В обработке"
	req.Status = "сформирован"
	req.FormationAt = time.Now()
	if err := r.db.Save(&req).Error; err != nil {
		return 0, err
	}

	return req.ID, nil
}

func (r *Repository) GetOperationByUser(creatorID uint) (uint, error) {
	var req ds.Operation
	res := r.db.
		//Where("id = ?", requestID).
		Where("user_id = ?", creatorID).
		Where("status = ?", utils.Draft).
		Take(&req)

	if res.Error != nil {
		return 0, res.Error
	}
	if res.RowsAffected == 0 {
		return 0, errors.New("нет такой заявки")
	}

	return req.ID, nil
}

func (r *Repository) GetOperationByID(creatorID uint, id uint) (ds.Operation, error) {
	var req ds.Operation
	res := r.db.
		Where("id = ?", id).
		Where("user_id = ?", creatorID).
		Where("status = ?", utils.Draft).
		Take(&req)

	if res.Error != nil {
		return req, res.Error
	}
	if res.RowsAffected == 0 {
		return req, errors.New("нет такой заявки")
	}

	return req, nil
}

func (r *Repository) FinishRejectHelper(status string, requestID, moderatorID uint) error {
	//userInfo, err := GetUserInfo(r, moderatorID)
	//if err != nil {
	//	return err
	//}

	var req ds.Operation
	res := r.db.
		Where("id = ?", requestID).
		Where("status = ?", "сформирован").
		Take(&req)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("нет такой заявки")
	}

	req.ModeratorID = &moderatorID
	//req.ModeratorLogin = userInfo.Login
	req.Status = status

	req.CompletionAt = time.Now()

	if err := r.db.Save(&req).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteOperationByID(requestID uint) error { // ?
	var req ds.Operation
	if result := r.db.First(&req, requestID); result.Error != nil {
		return result.Error
	}

	req.Status = "удален"
	req.CompletionAt = time.Now()
	if err := r.db.Save(&req).Error; err != nil {
		return err
	}
	result := r.db.Save(&req)

	return result.Error
}

func (r *Repository) DeleteBanknoteFromRequest(id int) error {
	var dh ds.OperationBanknote
	if result := r.db.First(&dh, id); result.Error != nil {
		return result.Error
	}
	return r.db.Delete(&dh).Error
}

func (r *Repository) UpdateOperationBanknote(id uint, quantity int) error {
	var updateBanknote ds.OperationBanknote
	r.db.Where("id = ?", id).First(&updateBanknote)

	if updateBanknote.OperationID == 0 {
		return errors.New("нет такой заявки")
	}
	updateBanknote.Quantity = quantity

	if err := r.db.Save(&updateBanknote).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) SaveRequest(monitoringRequest ds.RequestAsyncService) error {
	var request ds.Operation
	err := r.db.First(&request, "id = ?", monitoringRequest.RequestId)
	if err.Error != nil {
		r.logger.Error("error while getting monitoring request")
		return err.Error
	}
	//request.CompletionDate = time.Now()
	request.StatusCheck = monitoringRequest.Status
	res := r.db.Save(&request)
	return res.Error
}
