package api

import (
	"database/sql"
	"errors"
	db "github.com/budiharyonoo/simple_bank/db/sqlc"
	"github.com/budiharyonoo/simple_bank/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
)

// GetListAccounts used for fetch query params
type GetListAccounts struct {
	Page  int32 `form:"page" binding:"required,min=1"`
	Limit int32 `form:"limit" binding:"required,min=1,max=50"`
}

func (server Server) getListAccounts(ctx *gin.Context) {
	var req GetListAccounts

	// Bind query params to the getListAccounts struct
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get auth payload data and convert as payload struct
	authPayload := ctx.MustGet(authPayloadKey).(*token.Payload)

	// Fetch number of total accounts rows
	totalRows, err := server.store.GetAccountsTotalRows(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Fetch list of accounts by owner with pagination
	accounts, err := server.store.ListAccountsWithPagination(ctx, db.ListAccountsWithPaginationParams{
		Owner:  authPayload.Username,
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
	Currency string `json:"currency" binding:"required,currency"`
}

func (server Server) createAccount(ctx *gin.Context) {
	var req CreateAccountRequest

	// Bind JSON payload to the CreateAccountRequest struct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, errorResponse(err))
		return
	}

	// Get auth payload data and convert as payload struct
	authPayload := ctx.MustGet(authPayloadKey).(*token.Payload)

	// Create Account to DB
	account, err := server.store.CreateAccount(ctx, db.CreateAccountParams{
		Owner:    authPayload.Username,
		Balance:  0,
		Currency: req.Currency,
	})

	if err != nil {
		var defaultStatusCode = http.StatusInternalServerError

		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code.Name() {
			case "unique_violation", "foreign_key_violation":
				defaultStatusCode = http.StatusForbidden
			}
		}

		ctx.JSON(defaultStatusCode, errorResponse(err))
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

	// Get auth payload data and convert as payload struct
	authPayload := ctx.MustGet(authPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}
