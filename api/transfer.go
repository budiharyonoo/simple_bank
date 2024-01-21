package api

import (
	"database/sql"
	"errors"
	"fmt"
	db "github.com/budiharyonoo/simple_bank/db/sqlc"
	"github.com/budiharyonoo/simple_bank/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateTransferRequest used for input payload request and validation
type CreateTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Currency      string `json:"currency" binding:"required,currency"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
}

func (server Server) createTransfer(ctx *gin.Context) {
	var req CreateTransferRequest

	// Bind JSON payload to the CreateTransferRequest struct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, errorResponse(err))
		return
	}

	fromAccount, valid := server.validateAccount(ctx, req.FromAccountID, "FROM", req.Currency)
	if !valid {
		return
	}

	// Validate Owner and Auth Token Username must be the same
	authPayload := ctx.MustGet(authPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	_, valid = server.validateAccount(ctx, req.ToAccountID, "TO", req.Currency)
	if !valid {
		return
	}

	// Store Transfer Record
	transfer, err := server.store.TransferTx(ctx, db.TransferTxParams{
		FromAccountId: req.FromAccountID,
		ToAccountId:   req.ToAccountID,
		Amount:        req.Amount,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, transfer)
}

// validateAccount used for validate if the FROM & TO Account's exists and has the same currency
func (server Server) validateAccount(
	ctx *gin.Context,
	accountID int64,
	sourceAccount string,
	currency string,
) (db.Account, bool) {
	// Fetch single account by Id
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		// If account not exists
		if errors.Is(sql.ErrNoRows, err) {
			errMsg := fmt.Sprintf("%s account [%d] not found", sourceAccount, accountID)
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New(errMsg)))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		errMsg := fmt.Sprintf("The %s account [%d] currency is %s, not %s", sourceAccount, accountID, account.Currency, currency)
		err := errors.New(errMsg)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}
