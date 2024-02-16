package handler

import (
	"errors"
	"fmt"
	"main/internal/app/ds"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func ParseDateString(dateString string) (time.Time, error) {
	format := "2006-01-02 15:04:05"
	re := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})`)
	matches := re.FindStringSubmatch(dateString)
	if len(matches) < 2 {
		return time.Time{}, nil
	}
	parsedTime, err := time.Parse(format, matches[1])
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}

// OperationList godoc
// @Summary      Get list of Operation requests
// @Description  Retrieves a list of Operation requests based on the provided parameters
// @Tags         Operations
// @Accept       json
// @Produce      json
// @Param        status_id      query  string    false  "Operation request status"
// @Param        start_date  query  string    false  "Start date in the format '2006-01-02T15:04:05Z'"
// @Param        end_date    query  string    false  "End date in the format '2006-01-02T15:04:05Z'"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success      200  {object}  []ds.Operation
// @Failure      400  {object}  error
// @Failure      500  {object}  error
// @Router       /api/operations [get]
func (h *Handler) OperationList(ctx *gin.Context) {
	userID, existsUser := ctx.Get("user_id")
	userRole, existsRole := ctx.Get("user_role")
	if !existsUser || !existsRole {
		h.errorHandler(ctx, http.StatusUnauthorized, errors.New("not fount `user_id` or `user_role`"))
		return
	}

	switch userRole {
	case ds.Buyer:
		h.operationByUserId(ctx, fmt.Sprintf("%d", userID))
		return
	default:
		break
	}

	queryStatus := ctx.Query("status_id")
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	if startDateStr == "" {
		startDateStr = "0001-01-01 00:00:00"
	}
	if endDateStr == "" {
		endDateStr = time.Now().Add(time.Hour * 24).String()
	}

	startDate, errStart := ParseDateString(startDateStr + " 00:00:00")
	endDate, errEnd := ParseDateString(endDateStr + " 00:00:00")
	h.Logger.Info(startDate, endDate)
	if errEnd != nil || errStart != nil {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("incorrect `start_date` or `end_date`"))
		return
	}

	operations, err := h.Repository.OperationsList(queryStatus, startDate, endDate)
	if err != nil {
		h.errorHandler(ctx, http.StatusNoContent, err)
		return
	}
	h.successHandler(ctx, "operations", operations)
}

func (h *Handler) operationByUserId(ctx *gin.Context, userID string) {
	operations, errDB := h.Repository.OperationByUserID(userID)
	if errDB != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, errDB)
		return
	}

	h.successHandler(ctx, "operations", operations)
}

// GetOperationById godoc
// @Summary      Get operation request by ID
// @Description  Retrieves a operation request with the given ID
// @Tags         Operations
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Operation Request ID"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success      200  {object}  nil
// @Failure      400  {object}  error
// @Router       /api/operations/{id} [get]
func (h *Handler) GetOperationById(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	//req, bank, err := h.Repository.GetOperationWithDataByID(uint(id))
	operation, err := h.Repository.OperationByID(uint(id))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	// c.JSON(http.StatusOK, gin.H{ "operation": req, "banknote":  bank})
	h.successHandler(c, "operation", operation)
}

// UpdateOperation godoc
// @Summary      Update Operation by admin
// @Description  Update Operation by admin
// @Tags         Operations
// @Accept       json
// @Produce      json
// @Param        input    body    ds.Operation  true    "updated Assembly"
// @Success      200          {object}  nil
// @Failure      400          {object}  error
// @Failure      500          {object}  error
// @Router       /api/operations [put]
func (h *Handler) UpdateOperation(ctx *gin.Context) {
	userID, existsUser := ctx.Get("user_id")
	userRole, existsRole := ctx.Get("user_role")
	if !existsUser || !existsRole {
		h.errorHandler(ctx, http.StatusUnauthorized, errors.New("not fount `user_id` or `user_role`"))
		return
	}

	var updatedOperation ds.UpdateOperation
	if err := ctx.BindJSON(&updatedOperation); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if updatedOperation.ID == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("id некоректен"))
		return
	}

	var updatedO ds.Operation
	updatedO.ID = updatedOperation.ID
	updatedO.Name = updatedOperation.Name

	operation, err := h.Repository.OperationModel(updatedO.ID)

	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, fmt.Errorf("hike with `id` = %d not found", operation.ID))
		return
	}

	if operation.UserID != userID && userRole == ds.Buyer {
		h.errorHandler(ctx, http.StatusForbidden, errors.New("you cannot change the hike if it's not yours"))
		return
	}

	if err := h.Repository.UpdateOperation(&updatedO); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	h.successHandler(ctx, "updated_operation", gin.H{
		"id":              updatedOperation.ID,
		"operation_name":  updatedOperation.Name,
		"creation_date":   operation.CreatedAt,
		"completion_date": operation.CompletionAt,
		"formation_date":  operation.FormationAt,
		"user_id":         operation.UserID,
		"status":          operation.Status,
	})
}

// func (h *Handler) CreateDraft(c *gin.Context) {
// 	draftID, err := h.Repository.CreateOperationDraft(creatorID)

// 	if err != nil {
// 		h.errorHandler(c, http.StatusInternalServerError, err)
// 	}

// 	c.JSON(http.StatusOK, gin.H{"draftID": draftID})
// }

// FormOperationsRequest godoc
// @Summary      Form Banknote by client
// @Description  Form Banknote by client
// @Tags         Operations
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Operation form ID"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success      200          {object}  nil
// @Failure      400          {object}  error
// @Failure      500          {object}  error
// @Router       /api/operations/form [put]
func (h *Handler) FormOperationRequest(c *gin.Context) {
	userID, existsUser := c.Get("user_id")
	if !existsUser {
		h.errorHandler(c, http.StatusUnauthorized, errors.New("not fount `user_id` or `user_role`"))
		return
	}

	_, err := h.Repository.FormOperationRequestByID(userID.(uint))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusOK)
}

// UpdateStatusOperationRequest godoc
// @Summary      Update operation request status by ID
// @Description  Updates the status of a operation request with the given ID on "завершен"/"отклонен"
// @Tags         Operation
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Request ID"
// @Param        input    body    ds.NewStatus  true    "update status"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  error
// @Router       /operations/updateStatus [put]
func (h *Handler) UpdateStatusOperationRequest(c *gin.Context) {
	var status ds.NewStatus
	if err := c.BindJSON(&status); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	userIDStr, existsUser := c.Get("user_id")
	if !existsUser {
		h.errorHandler(c, http.StatusUnauthorized, errors.New("not fount `user_id` or `user_role`"))
		return
	}
	userID := userIDStr.(uint)

	if status.Status != "отклонен" && status.Status != "завершен" {
		h.errorHandler(c, http.StatusBadRequest, errors.New("статус можно поменять только на 'отклонен' и 'завершен'"))
	}

	if err := h.Repository.FinishRejectHelper(status.Status, status.OperationID, userID); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusOK)
}

// func (h *Handler) FinishOperationRequest(c *gin.Context) {
// 	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

// 	if err := h.Repository.FinishEncryptDecryptRequestByID(uint(id), moderatorID); err != nil {
// 		h.errorHandler(c, http.StatusBadRequest, err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, "завершена")
// }

// DeleteBanknoteFromRequest godoc
// @Summary      Delete banknote from request
// @Description  Deletes a banknote from a request based on the user ID and banknote ID
// @Tags         Operation_Banknote
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "banknote ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  error
// @Router       /api/operation-request-banknote [delete]
func (h *Handler) DeleteBanknoteFromRequest(c *gin.Context) {
	var body struct {
		ID int `json:"id"`
	}

	if err := c.BindJSON(&body); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	if body.ID == 0 {
		h.errorHandler(c, http.StatusBadRequest, errors.New("param `id` not found"))
		return
	}

	err := h.Repository.DeleteBanknoteFromRequest(body.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	h.successHandler(c, "deleted_banknote_operation", body.ID)
}

// func (h *Handler) RejectOperationRequest(c *gin.Context) {
// 	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

// 	if err := h.Repository.RejectOperationRequestByID(uint(id), moderatorID); err != nil {
// 		h.errorHandler(c, http.StatusBadRequest, err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, "отклонена")
// }

// DeleteOperation godoc
// @Summary      Delete operation request by user ID
// @Description  Deletes a operation request for the given user ID
// @Tags         Operations
// @Accept       json
// @Produce      json
// @Param        user_id  path  int  true  "User ID"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  error
// @Router       /api/operations [delete]
func (h *Handler) DeleteOperation(c *gin.Context) {
	userID, existsUser := c.Get("user_id")
	userRole, existsRole := c.Get("user_role")
	if !existsUser || !existsRole {
		h.errorHandler(c, http.StatusUnauthorized, errors.New("not fount `user_id` or `user_role`"))
		return
	}

	//userId := c.GetInt(userCtx)
	var request struct {
		ID uint `json:"id"`
	}

	if err := c.BindJSON(&request); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	if request.ID == 0 {
		h.errorHandler(c, http.StatusBadRequest, errors.New("param `id` not found"))
		return
	}

	//userId := c.GetInt(userCtx)

	operation, err := h.Repository.OperationModel(request.ID)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, fmt.Errorf("operation with `id` = %d not found", operation.ID))
		return
	}

	if operation.UserID != userID && userRole == ds.Buyer {
		h.errorHandler(c, http.StatusForbidden, errors.New("you are not the creator. you can't delete a operation"))
		return
	}

	err = h.Repository.DeleteOperationByID(request.ID)
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	h.successHandler(c, "operation_id", request.ID)
}

// UpdateOperationBanknote godoc
// @Summary      Update money Operation Banknote
// @Description  Update money Operation Banknote by client
// @Tags         Operation_Banknote
// @Accept       json
// @Produce      json
// @Param        input    	  body    ds.OperationBanknote true    "Update quantity Operation Banknote"
// @Success      200          {object} map[string]string "update"
// @Failure      400          {object}  error
// @Failure      500          {object}  error
// @Router       /api/operation-request-banknote [put]
func (h *Handler) UpdateOperationBanknote(c *gin.Context) {
	//var operationBanknote ds.OperationBanknote
	var OperationBanknoteU ds.OperationBanknoteUpdate
	if err := c.BindJSON(&OperationBanknoteU); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	//if OperationBanknoteU.OperationID == 0 || OperationBanknoteU.BanknoteID == 0 {
	//	h.errorHandler(c, http.StatusBadRequest, errors.New("не верные id операции или банкноты"))
	//	return
	//}

	err := h.Repository.UpdateOperationBanknote(OperationBanknoteU.ID, OperationBanknoteU.Quantity)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, "update")
}
