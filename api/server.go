package api

import (
	"fmt"
	db "github.com/budiharyonoo/simple_bank/db/sqlc"
	"github.com/budiharyonoo/simple_bank/token"
	"github.com/budiharyonoo/simple_bank/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
)

// Server serves HTTP requests for our service
type Server struct {
	config     utils.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

// NewServer is a contructor of Server and setup router
func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker()
	if err != nil {
		return nil, fmt.Errorf("cannot init token maker %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	// Register custom validator to GIN
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validCurrency)
		if err != nil {
			log.Fatalln("Error register custom validator", err)
			return nil, err
		}
	}

	// === All Router ===
	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Server is up an running!"})
	})

	// Init auth middleware
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// Accounts Router
	authRoutes.GET("/v1/accounts/", server.getListAccounts)
	authRoutes.POST("/v1/accounts", server.createAccount)
	authRoutes.GET("/v1/accounts/:id", server.getAccountById)

	// Transfers Router
	authRoutes.POST("/v1/transfers", server.createTransfer)

	// Users Router
	router.POST("/v1/users", server.createUser)
	router.POST("/v1/users/login", server.loginUser)

	server.router = router

	return server, nil
}

// Start runs HTTP Server of specific port
// address is the server address
func (server Server) Start(address string) error {
	return server.router.Run(address)
}

// paginationResponse standarized response of pagination data
// gin.H is a shortcut for map[string]any
func paginationResponse(totalRows int64, page int32, limit int32, data any) gin.H {
	return gin.H{
		"metadata": gin.H{
			"total": totalRows,
			"page":  page,
			"limit": limit,
		},
		"data": data,
	}
}

// errorResponse serialize error before send to response client
// gin.H is a shortcut for map[string]any
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
