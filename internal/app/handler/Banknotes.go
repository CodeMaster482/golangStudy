package handler

import (
	"errors"
	"main/internal/app/ds"
	"main/internal/app/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Services struct {
	Id          int
	Name        string
	Description string
	Img         string
}

/*var services = []Services{
	{Id: 0, Name: "credit", Description: "Lorem ipsum dolor sit amet, consectetur adipisicing elit", Img: "/static/credit.jpg"},
	{Id: 1, Name: "deposite", Description: "Lorem ipsum dolor sit amet, consectetur adipisicing elit", Img: "/static/deposite.jpg"},
	{Id: 2, Name: "transfer", Description: "Lorem ipsum dolor sit amet, consectetur adipisicing elit", Img: "/static/transfer.jpg"},
	{Id: 3, Name: "open account", Description: "Lorem ipsum dolor sit amet, consectetur adipisicing elit", Img: "/static/accaunt.jpg"},
	{Id: 4, Name: "exchange", Description: "Lorem ipsum dolor sit amet, consectetur adipisicing elit", Img: "/static/change.png"},
}*/

func (h *Handler) BanknotesList(ctx *gin.Context) {
	queryText, _ := ctx.GetQuery("banknote")
	banknotes, err := h.Repository.BanknotesList(queryText)
	if err != nil {
		h.errorHandler(ctx, http.StatusNoContent, err)
		return
	}

	draftID, err := h.Repository.GetOprationDraftID(creatorID) // creatorID(UserID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	banknotesList := ds.BanknoteList{
		DraftID:   draftID,
		Banknotes: banknotes,
	}

	h.successHandler(ctx, "banknotes", banknotesList)
}

func (h *Handler) BanknoteById(ctx *gin.Context) {
	//queryText, _ := ctx.GetQuery("banknote")

	id, err := strconv.ParseUint(ctx.Param("id")[:], 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
	}

	banknote, err := h.Repository.GetBanknoteById(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	h.successHandler(ctx, "company", banknote)
}

func (h *Handler) DeleteBanknote(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id")[:], 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if id == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("param `id` not found"))
		return
	}

	url := h.Repository.DeleteBanknoteImage(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	h.DeleteImage(ctx, utils.ExtractObjectNameFromUrl(url))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.DeleteBanknote(uint(id))
	if gorm.IsRecordNotFoundError(err) {
		h.errorHandler(ctx, http.StatusBadRequest, err)
	} else if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, "company deleted successfully")
}

func (h *Handler) AddBanknote(ctx *gin.Context) {
	var newService ds.Banknote

	if newService.ID != 0 {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("param `id` not found"))
		return
	}

	nominal, err := strconv.ParseFloat(ctx.Request.FormValue("nominal"), 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("неверное значение номинала"))
	}
	if nominal == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("имя компании не может быть пустой"))
		return
	}

	newService.Nominal = nominal

	newService.Description = ctx.Request.FormValue("description")
	if newService.Description == "" {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("описание не может быть пустой"))
		return
	}

	newService.Currency = ctx.Request.FormValue("currency")
	if newService.Currency == "" {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("описание не может быть пустой"))
		return
	}

	file, header, err := ctx.Request.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "ошибка при загрузке изображения"})
		return
	}

	if file != nil {
		// Изображение предоставлено, обрабатываем его
		newService.ImageURL, err = h.SaveImage(ctx.Request.Context(), file, header)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "ошибка при сохранении изображения"})
			return
		}
	} else {
		// Изображение не предоставлено, устанавливаем пустую строку или другое значение по умолчанию
		newService.ImageURL = ""
	}

	createID, err := h.Repository.AddBanknote(&newService)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	h.successAddHandler(ctx, "banknote_id", createID)
}

func (h *Handler) BanknoteUpdate(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id")[:], 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	file, header, err := ctx.Request.FormFile("image")
	// if err != nil {
	// 	h.errorHandler(ctx, http.StatusBadRequest, err)
	// 	return
	// }

	// var updatedBanknote ds.Banknote
	// if err := ctx.BindJSON(&updatedBanknote); err != nil {
	// 	h.errorHandler(ctx, http.StatusBadRequest, err)
	// 	return
	// }

	var updatedService ds.Banknote

	updatedService.ID = uint(id)

	if updatedService.ID == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("param `id` not found"))
	}

	nominal, err := strconv.ParseFloat(ctx.Request.FormValue("nominal"), 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("неверное значение номинала"))
	}

	updatedService.Nominal = nominal
	updatedService.Description = ctx.Request.FormValue("description")

	if header != nil && header.Size != 0 {
		if updatedService.ImageURL, err = h.SaveImage(ctx.Request.Context(), file, header); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}

		url := h.Repository.DeleteBanknoteImage(updatedService.ID)

		if err = h.DeleteImage(ctx, utils.ExtractObjectNameFromUrl(url)); err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
	}

	if _, err := h.Repository.UpdateBanknote(&updatedService); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	h.successHandler(ctx, "updated_company", gin.H{
		"id":           updatedService.ID,
		"company_name": updatedService.Nominal,
		"description":  updatedService.Description,
		"image_url":    updatedService.ImageURL,
		"status":       updatedService.Status,
	})
}

func (h *Handler) AddBanknoteToRequest(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)

	draftID, err := h.Repository.AddBanknoteToDraft(uint(id), creatorID)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"draftID": draftID,
	})
}
