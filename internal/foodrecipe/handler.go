package foodrecipe

import (
	"net/http"
	"strconv"
	"wongnok/internal/global"
	"wongnok/internal/helper"
	"wongnok/internal/model"
	"wongnok/internal/model/dto"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type IHandler interface {
	Create(ctx *gin.Context)
	Get(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type Handler struct {
	Service IService
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		Service: NewService(db),
	}
}


func (handler Handler) Create(ctx *gin.Context) {
	claims, err := helper.DecodeClaims(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	var request dto.FoodRecipeRequest

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	recipe, err := handler.Service.Create(request, claims)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.As(err, &validator.ValidationErrors{}) {
			statusCode = http.StatusBadRequest
		}

		if errors.Is(err, global.ErrForbidden) {
			statusCode = http.StatusForbidden
		}

		ctx.JSON(statusCode, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, recipe.ToResponse())
}


func (handler Handler) Get(ctx *gin.Context) {
	var foodRecipeQuery model.FoodRecipeQuery
	if err := ctx.ShouldBindQuery(&foodRecipeQuery); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	recipes, total, err := handler.Service.Get(foodRecipeQuery)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, recipes.ToResponse(total))
}


func (handler Handler) GetByID(ctx *gin.Context) {
	var id int

	pathParam := ctx.Param("id")
	if pathParam != "" {
		if parsed, err := strconv.Atoi(pathParam); err == nil && parsed > 0 {
			id = parsed
		}
	}

	recipe, err := handler.Service.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Recipe not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, recipe.ToResponse())
}


func (handler Handler) Update(ctx *gin.Context) {
	claims, err := helper.DecodeClaims(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	var request dto.FoodRecipeRequest
	var id int

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	pathParam := ctx.Param("id")
	if pathParam != "" {
		if parsed, err := strconv.Atoi(pathParam); err == nil && parsed > 0 {
			id = parsed
		}
	}

	recipe, err := handler.Service.Update(request, id, claims)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.As(err, &validator.ValidationErrors{}) {
			statusCode = http.StatusBadRequest
		}

		if errors.Is(err, global.ErrForbidden) {
			statusCode = http.StatusForbidden
		}

		ctx.JSON(statusCode, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, recipe.ToResponse())
}


func (handler Handler) Delete(ctx *gin.Context) {
	claims, err := helper.DecodeClaims(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	var id int

	pathParam := ctx.Param("id")
	if pathParam != "" {
		if parsed, err := strconv.Atoi(pathParam); err == nil && parsed > 0 {
			id = parsed
		}
	}

	if err := handler.Service.Delete(id, claims); err != nil {
		statusCode := http.StatusInternalServerError

		if errors.Is(err, global.ErrForbidden) {
			statusCode = http.StatusForbidden
		}

		ctx.JSON(statusCode, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Recipe deleted successfully"})
}
