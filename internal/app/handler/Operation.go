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
// @Success      200  {object}  []ds.Operation
// @Failure      400  {object}  error
// @Failure      500  {object}  error
// @Router       /api/operation [get]
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

	tenders, err := h.Repository.OperationList(queryStatus, startDate, endDate)

	if err != nil {
		h.errorHandler(ctx, http.StatusNoContent, err)
		return
	}
	h.successHandler(ctx, "tenders", tenders)
}

func (h *Handler) operationByUserId(ctx *gin.Context, userID string) {
	tenders, errDB := h.Repository.OperationByUserID(userID)
	if errDB != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, errDB)
		return
	}

	h.successHandler(ctx, "tenders", tenders)
}

func (h *Handler) GetOperationById(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	req, bank, err := h.Repository.GetOperationWithDataByID(uint(id))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"operation": req,
		"banknote":  bank,
	})
}

func (h *Handler) UpdateOperation(ctx *gin.Context) {
	var updatedOperation ds.Operation
	if err := ctx.BindJSON(&updatedOperation); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if updatedOperation.ID == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("id некоректен"))
		return
	}
	if err := h.Repository.UpdateOperation(&updatedOperation); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, "")
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
// @Router       /operation/updateStatus [put]
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

func (h *Handler) CreateDraft(c *gin.Context) {
	draftID, err := h.Repository.CreateOperationDraft(creatorID)

	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, gin.H{"draftID": draftID})
}

func (h *Handler) FormOperationRequest(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	err := h.Repository.FormOperationRequestByID(uint(id), creatorID)
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	req, com, err := h.Repository.GetOperationWithDataByID(uint(id))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"operation": req,
		"banknotes": com,
	})
}

func (h *Handler) RejectOperationRequest(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.Repository.RejectOperationRequestByID(uint(id), moderatorID); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, "отклонена")
}

func (h *Handler) FinishOperationRequest(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.Repository.FinishEncryptDecryptRequestByID(uint(id), moderatorID); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, "завершена")
}
func (h *Handler) DeleteBanknoteFromRequest(c *gin.Context) {
	var deleteFromOperation ds.OperationBanknote
	if err := c.BindJSON(&deleteFromOperation); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}
	if deleteFromOperation.OperationID <= 0 {
		h.errorHandler(c, http.StatusBadRequest, errors.New("id не найден"))
		return
	}

	if deleteFromOperation.BanknoteID <= 0 {
		h.errorHandler(c, http.StatusBadRequest, errors.New("id не найден"))
		return
	}

	request, banknotes, err := h.Repository.DeleteBanknoteFromRequest(deleteFromOperation)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "банкнота удалена из заявки", "banknotes": banknotes, "monitoring-request": request})
}

func (h *Handler) DeleteOperation(c *gin.Context) {
	//ModeratorID тут проверка, что это модератор -> 5 lab
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	err := h.Repository.DeleteOperationByID(uint(id))
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, "deleted")
}

func (h *Handler) UpdateOperationBanknote(c *gin.Context) {
	var OperationBanknote ds.OperationBanknote
	if err := c.BindJSON(&OperationBanknote); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	if OperationBanknote.OperationID == 0 || OperationBanknote.BanknoteID == 0 {
		h.errorHandler(c, http.StatusBadRequest, errors.New("не верные id тендера или кампапии"))
		return
	}

	err := h.Repository.UpdateOperationBanknote(OperationBanknote.OperationID, OperationBanknote.BanknoteID, OperationBanknote.Quantity)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, "update")
}
