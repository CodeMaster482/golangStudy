package handler

import (
	"fmt"
	"main/internal/ds"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
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

func (h *Handler) BanknotesList(context *gin.Context) {
	draftId, _ := h.Repository.GetOprationDraftID(creatorID)

	banknotesList, count, err := h.Repository.BanknoteList(draftId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrire to banknotes",
		})
	}

	query := context.DefaultQuery("banknote", "")
	if query != "" {
		result := []ds.Banknote{}

		for _, data := range *banknotesList {
			if strings.HasPrefix(fmt.Sprint(data.Nominal), query) {
				result = append(result, data)
			}
		}

		context.HTML(http.StatusOK, "index.html", gin.H{
			"services":  result,
			"banknote":  query,
			"BucketNum": count,
		})
	} else {
		context.HTML(http.StatusOK, "index.html", gin.H{
			"services":  banknotesList,
			"banknote":  query,
			"BucketNum": count,
		})
	}
}

func (h *Handler) AddBanknoteToRequest(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)

	_, err := h.Repository.AddBanknoteToDraft(uint(id), creatorID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	draftId, _ := h.Repository.GetOprationDraftID(creatorID)

	banknotesList, count, err := h.Repository.BanknoteList(draftId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrire to banknotes",
		})
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"services":  banknotesList,
		"BucketNum": count,
	})
}

func (h *Handler) BanknoteById(context *gin.Context) {
	idGet := context.Param("id")
	id, _ := strconv.Atoi(idGet)

	banknote, err := h.Repository.BanknoteById(id)
	if err != nil {
		context.String(http.StatusNotFound, "404 ---- Not Found")
	}

	context.HTML(http.StatusOK, "productPage.html", gin.H{
		"services": banknote,
	})
}

func (h *Handler) DeleteBanknote(context *gin.Context) {
	banknoteId := context.PostForm("banknote_id")
	err := h.Repository.DeleteBanknote(banknoteId)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Can't delete this banknote"})
	}
	context.Redirect(http.StatusFound, "/banknotes")
}
