package http

import (
	"context"
	"errors"
	"net/http"
	"quik/domain"
	"quik/wallet/handler/middleware"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	WalletService domain.WalletService
}

func NewWalletHandler(router *gin.Engine, ws domain.WalletService) {
	handler := &WalletHandler{
		WalletService: ws,
	}

	api := router.Group("/api/v1")
	api.POST("/wallets", middleware.AuthPlayer(), handler.CreateWallet)
	api.GET("/wallets/:wallet_id/balance", middleware.AuthPlayer(), handler.GetWalletBalance)
	api.POST("/wallets/:wallet_id/credit", middleware.AuthPlayer(), handler.CreditWallet)
	api.POST("/wallets/:wallet_id/debit", middleware.AuthPlayer(), handler.DebitWallet)
}

func isValidInteger(value string) bool {
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil || intValue < 1 {
		return false
	}
	return true
}

func (w *WalletHandler) CreateWallet(c *gin.Context) {
	var ctx = context.TODO()
	var wallet domain.Wallet

	id, _ := c.Get("playerId")
	playerId, _ := id.(int)
	wallet.PlayerID = playerId
	err := w.WalletService.Create(ctx, &wallet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"payload": wallet})
}

func (w *WalletHandler) CreditWallet(c *gin.Context) {
	var input struct {
		Amount string `json:"amount" `
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	walletId := c.Param("wallet_id")
	if !isValidInteger(walletId) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid amount"})
		return
	}
	var ctx = context.TODO()
	err := w.WalletService.Credit(ctx, walletId, input.Amount)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		case errors.Is(err, domain.ErrInvalidAmount):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "wallet credited"})
}

func (w *WalletHandler) DebitWallet(c *gin.Context) {
	var input struct {
		Amount string `json:"amount" validate:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	walletId := c.Param("wallet_id")
	if !isValidInteger(walletId) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid amount"})
		return
	}
	var ctx = context.TODO()
	err := w.WalletService.Debit(ctx, walletId, input.Amount)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		case errors.Is(err, domain.ErrInsufficientFunds):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		case errors.Is(err, domain.ErrInvalidAmount):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "wallet debited"})
}

func (w *WalletHandler) GetWalletBalance(c *gin.Context) {
	walletId := c.Param("wallet_id")
	if !isValidInteger(walletId) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid amount"})
		return
	}
	var ctx = context.TODO()
	wallet, err := w.WalletService.Get(ctx, walletId)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	payload := map[string]interface{}{
		"balance": wallet.Balance,
	}
	c.JSON(http.StatusOK, gin.H{"payload": payload})
}
