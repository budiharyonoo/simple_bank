package api

import (
	"database/sql"
	"errors"
	db "github.com/budiharyonoo/simple_bank/db/sqlc"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetListAccounts used for fetch query params
type GetListAccounts struct {
	Page  int32 `form:"page" binding:"required,min=1"`
	Limit int32 `form:"limit" binding:"required,min=1,max=50"`
}

func (server Server) getListAccounts(ctx *gin.Context) {
	var req GetListAccounts

	// Bind query params to the GetListAccounts struct
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Fetch number of total accounts rows
	totalRows, err := server.store.GetAccountsTotalRows(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Fetch list of accounts with pagination
	accounts, err := server.store.ListAccountsWithPagination(ctx, db.ListAccountsWithPaginationParams{
		Limit:  req.Limit,
		Offset: (req.Page - 1) * req.Limit,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, paginationResponse(totalRows, req.Page, req.Limit, accounts))
}

// CreateAccountRequest used for input payload request and validation
type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (server Server) createAccount(ctx *gin.Context) {
	var req CreateAccountRequest

	// Bind JSON payload to the CreateAccountRequest struct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, errorResponse(err))
		return
	}

	// Create Account to DB
	account, err := server.store.CreateAccount(ctx, db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, account)
}

// GetAccountByIdRequest used for fetch and validate the ID from URI
type GetAccountByIdRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server Server) getAccountById(ctx *gin.Context) {
	var req GetAccountByIdRequest

	// Bind URI to the GetAccountByIdRequest struct
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Fetch single account by Id
	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		// If account not exists
		if errors.Is(sql.ErrNoRows, err) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}
