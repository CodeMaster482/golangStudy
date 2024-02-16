package handler

import (
	"errors"
	"main/internal/app/ds"
	"main/internal/app/utils"
	"mime/multipart"
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

// BanknotesList godoc
// @Summary      Banknotes List
// @Description  Banknotes List
// @Tags         Banknotes
// @Accept       json
// @Produce      json
// @Param        banknote_name query   string  false  "Query string to filter banknotes by nominal"
// @Success      200          {object}  ds.BanknoteList
// @Failure      500          {object}  error
// @Router       /api/banknote [get]
func (h *Handler) BanknotesList(ctx *gin.Context) {
	queryText, _ := ctx.GetQuery("banknote_name")
	banknotes, err := h.Repository.BanknotesList(queryText)
	if err != nil {
		h.errorHandler(ctx, http.StatusNoContent, err)
		return
	}
	userID, existsUser := ctx.Get("user_id")
	var draftIdRes uint = 0
	if existsUser {
		basketId, errBask := h.Repository.GetOperationDraftID(userID.(uint))
		if errBask != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, errBask)
			return
		}
		draftIdRes = basketId
	}
	//draftID, err := h.Repository.GetOperationDraftID(ctx.GetInt(userCtx)) // creatorID(UserID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	//banknotesList := ds.BanknotesList{
	//	DraftID:   draftIdRes,
	//	Banknotes: banknotes,
	//}

	ctx.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"banknotes": banknotes,
		"draft_id":  draftIdRes,
	})
}

// GetBanknotesById godoc
// @Summary      Banknotes By ID
// @Description  Banknote By ID
// @Tags         Banknotes
// @Accept       json
// @Produce      json
// @Param        id   path    int     true        "Banknotes ID"
// @Success      200          {object}  ds.Banknote
// @Failure      400          {object}  error
// @Failure      500          {object}  error
// @Router       /api/banknotes/{id} [get]
func (h *Handler) GetBanknotesById(ctx *gin.Context) {
	//queryText, _ := ctx.GetQuery("banknote")

	id, err := strconv.ParseUint(ctx.Param("id")[:], 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	banknote, err := h.Repository.GetBanknoteById(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	h.successHandler(ctx, "banknote", banknote)
}

// DeleteBanknote godoc
// @Summary      Delete banknote by ID
// @Description  Deletes a banknote with the given ID
// @Tags         Banknotes
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Banknote ID"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  error
// @Router       /api/banknotes [delete]
func (h *Handler) DeleteBanknote(ctx *gin.Context) {
	var request struct {
		ID string `json:"id"`
	}

	// id, err := strconv.ParseUint(ctx.Param("id")[:], 10, 64)
	// if err != nil {
	// 	h.errorHandler(ctx, http.StatusBadRequest, err)
	// 	return
	// }

	if err := ctx.BindJSON(&request); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	id, err := strconv.Atoi(request.ID)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if id == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("param `id` not found"))
		return
	}

	url := h.Repository.DeleteBanknoteImage(uint(id))

	err = h.DeleteImage(ctx, utils.ExtractObjectNameFromUrl(url))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.DeleteBanknote(uint(id))

	if gorm.IsRecordNotFoundError(err) {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	} else if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	h.successHandler(ctx, "deleted_id", id)
}

// AddBanknote godoc
// @Summary      Add new banknote
// @Description  Add a new banknote with image, nominal, currency
// @Tags         Banknotes
// @Accept       multipart/form-data
// @Produce      json
// @Param        image formData file true "Banknote image"
// @Param        nominal formData string true "Banknote nominal"
// @Param        description formData string false "Banknote description"
// @Param        currency formData string true "Banknote currency"
// @Success      201  {string}  map[string]any
// @Failure      400  {object}  map[string]any
// @Router       /api/banknotes [post]
func (h *Handler) AddBanknote(ctx *gin.Context) {
	var newBanknote ds.Banknote

	nominal, err := strconv.ParseFloat(ctx.Request.FormValue("nominal"), 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("неверное значение номинала"))
		return
	}
	if nominal == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("номинал не может быть пустым"))
		return
	}

	newBanknote.Nominal = nominal

	newBanknote.Currency = ctx.Request.FormValue("currency")
	if len(newBanknote.Currency) != 3 {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("невозможно добавить такой тип валюты"))
		return
	}

	newBanknote.Description = ctx.Request.FormValue("description")
	if newBanknote.Description == "" {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("описание не может быть пустой"))
		return
	}

	newBanknote.Status = ctx.Request.FormValue("status")

	file, header, err := ctx.Request.FormFile("image_url")
	if err != http.ErrMissingFile && err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "ошибка при загрузке изображения"})
		return
	}

	if newBanknote.ImageURL, err = h.SaveImage(ctx.Request.Context(), file, header); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "ошибка при сохранении изображения"})
		return
	}

	createID, err := h.Repository.AddBanknote(&newBanknote)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	h.successAddHandler(ctx, "banknote_id", createID)
}

// UpdateBanknote godoc
// @Summary      Update banknote by ID
// @Description  Updates a banknote with the given ID
// @Tags         Banknotes
// @Accept       multipart/form-data
// @Produce      json
// @Param        id          path        int     true        "ID"
// @Param        name        formData    string  false       "nominal"
// @Param        description formData    string  false       "description"
// @Param        IIN         formData    string  false       "currency"
// @Param        image       formData    file    false       "image"
// @Success      200         {object}    map[string]any
// @Failure      400         {object}    error
// @Router       /api/banknotes/ [put]
func (h *Handler) BanknoteUpdate(ctx *gin.Context) {
	var updatedBanknote ds.Banknote
	if err := ctx.BindJSON(&updatedBanknote); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if updatedBanknote.ImageURL != "" {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New(`image_url must be empty`))
		return
	}

	if updatedBanknote.Status != "действует" && updatedBanknote.Status != "удален" {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New(`status может быть только действует или удален`))
		return
	}

	if updatedBanknote.ID == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("param `id` not found"))
		return
	}

	if err := h.Repository.UpdateBanknote(&updatedBanknote); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	h.successHandler(ctx, "updated_banknote", gin.H{
		"id":               updatedBanknote.ID,
		"banknote_nominal": updatedBanknote.Nominal,
		"currency":         updatedBanknote.Currency,
		"description":      updatedBanknote.Description,
		"image_url":        updatedBanknote.ImageURL,
		"status":           updatedBanknote.Status,
	})
}

func (h *Handler) AddImage(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	banknoteId := ctx.Request.FormValue("banknote_id")

	if banknoteId == "" {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("param `id` not found"))
		return
	}
	if header == nil || header.Size == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("no file uploaded"))
		return
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	defer func(file multipart.File) {
		errLol := file.Close()
		if errLol != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, errLol)
			return
		}
	}(file)

	ID, _ := strconv.Atoi(banknoteId)
	url := h.Repository.DeleteBanknoteImage(uint(ID))

	if err = h.DeleteImage(ctx, utils.ExtractObjectNameFromUrl(url)); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var imageURL string
	if imageURL, err = h.SaveImage(ctx.Request.Context(), file, header); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	var updatedBanknote ds.Banknote
	updatedBanknote.ID = uint(ID)
	updatedBanknote.ImageURL = imageURL
	if err := h.Repository.UpdateBanknote(&updatedBanknote); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	h.successAddHandler(ctx, "image_url", imageURL)
}

// AddBanknoteToRequest godoc
// @Summary      Add banknote to request
// @Description  Adds a banknote to a operation request
// @Tags         Banknotes
// @Accept       json
// @Produce      json
// @Param request body ds.AddToBanknoteID true "Добавление банкноты"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  error
// @Router       /api/banknotes/request [post]
func (h *Handler) AddBanknoteToRequest(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		h.errorHandler(ctx, http.StatusUnauthorized, errors.New("user_id not found"))
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		h.errorHandler(ctx, http.StatusUnauthorized, errors.New("`user_id` must be uint number"))
		return
	}

	var request ds.AddToBanknoteID
	if err := ctx.BindJSON(&request); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	//request.UserID = userIDUint
	//id, err := strconv.ParseUint(ctx.Param("id")[:], 10, 64)
	//request.BanknoteID = uint(id)

	if request.BanknoteID == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "услуга не может быть пустой"})
		return
	}

	draftID, err := h.Repository.AddBanknoteToDraft(request.BanknoteID, userIDUint, request.Quantity)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	h.successHandler(ctx, "draft_id", draftID)
}
