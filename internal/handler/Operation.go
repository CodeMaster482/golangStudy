package handler

import (
	"errors"
	"main/internal/ds"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) OperationList(ctx *gin.Context) {
	queryStatus, _ := ctx.GetQuery("status")

	queryStart, _ := ctx.GetQuery("start")

	queryEnd, _ := ctx.GetQuery("end")

	operations, err := h.Repository.OperationList(queryStatus, queryStart, queryEnd)

	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"operations": operations})
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
