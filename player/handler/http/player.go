package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"quik/domain"
	"quik/internal/encryption"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PlayerHandler struct {
	PlayerService domain.PlayerService
	WalletService domain.WalletService
}

func NewPlayerHandler(router *gin.Engine, p domain.PlayerService, w domain.WalletService) {
	handler := &PlayerHandler{
		PlayerService: p,
		WalletService: w,
	}
	api := router.Group("/api/v1")
	api.POST("/players", handler.CreatePlayer)
	api.GET("/players/:id", handler.GetPlayerByID)
	api.PUT("/players/:id", handler.UpdatePlayerByID)
	api.DELETE("/players/:id", handler.DeletePlayerByID)
	api.POST("players/login", handler.Login)
}

var validate *validator.Validate

func isIDValid(ID string) error {
	id, err := strconv.ParseInt(ID, 10, 64)
	if err != nil || id < 1 {
		return errors.New("invalid id parameter")
	}
	return nil
}

func inputValidator(input interface{}) map[string]string {
	validate = validator.New()
	err := validate.Struct(input)
	if err != nil {
		validationErrors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			validationErrors[e.Field()] = fmt.Sprintf("%s failed validation", e.Field())
		}
		return validationErrors
	}
	return nil
}

func (p *PlayerHandler) CreatePlayer(c *gin.Context) {
	var input struct {
		Name     string `json:"name" validate:"gte=0,lte=500,required"`
		Email    string `json:"email" validate:"email,required"`
		Password string `json:"password" validate:"min=8,max=72,required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inputErr := inputValidator(input)
	if inputErr != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": inputErr})
		return
	}
	hashedPassword, err := encryption.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	input.Password = hashedPassword
	var ctx = context.TODO()
	var player domain.Player
	player.Name = input.Name
	player.Email = input.Email
	player.Password = input.Password
	err = p.PlayerService.Create(ctx, &player)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDuplicateRecord):
			c.JSON(http.StatusConflict, gin.H{"error": "Email is registered. Kindly login"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"payload": player})
}

func (p *PlayerHandler) GetPlayerByID(c *gin.Context) {
	id := c.Param("id")
	err := isIDValid(id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	var ctx = context.TODO()
	player, err := p.PlayerService.Get(ctx, id)
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
	c.JSON(http.StatusFound, gin.H{"payload": player})
}

func (p *PlayerHandler) UpdatePlayerByID(c *gin.Context) {
	id := c.Param("id")
	err := isIDValid(id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	var input struct {
		Name     string `json:"name" validate:"isdefault|gte=0,lte=500"`
		Email    string `json:"email" validate:"isdefault|email"`
		Password string `json:"password" validate:"isdefault|min=8,max=72"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inputErr := inputValidator(input)
	if inputErr != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": inputErr})
		return
	}
	var ctx = context.TODO()
	player, err := p.PlayerService.Get(ctx, id)
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
	var updatedPlayer domain.Player
	updatedPlayer.Email = input.Email
	updatedPlayer.Name = input.Name
	if input.Password != "" {
		hashedPassword, err := encryption.HashPassword(input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		updatedPlayer.Password = hashedPassword
	}
	err = p.PlayerService.Update(ctx, id, &player, updatedPlayer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, player)
}

func (p *PlayerHandler) DeletePlayerByID(c *gin.Context) {
	id := c.Param("id")
	err := isIDValid(id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	var ctx = context.TODO()
	var player domain.Player
	err = p.PlayerService.Delete(ctx, id, &player)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Player deleted successfully"})
}

func (p *PlayerHandler) Login(c *gin.Context) {
	var input struct {
		Password string `json:"password" validate:"min=8,max=72,required"`
		Email    string `json:"email" validate:"email,required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inputErr := inputValidator(input)
	if inputErr != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": inputErr})
		return
	}

	var ctx = context.TODO()
	var player domain.Player
	err := p.PlayerService.FindByEmail(ctx, input.Email, &player)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRecordNotFound):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ok, err := encryption.Matches(input.Password, player.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := encryption.CreateToken(player.Name, "player", player.ID, 2160)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	payload := map[string]interface{}{
		"token":  token,
		"player": player,
	}
	c.JSON(http.StatusOK, gin.H{"payload": payload})
}
