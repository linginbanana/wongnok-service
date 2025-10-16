/*
package favorite

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type Handler struct {
    db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
    return &Handler{db: db}
}

func (h *Handler) Add(c *gin.Context) {
    userID := c.Param("id")
    var req struct {
        RecipeID string `json:"recipeId"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    fav := Favorite{UserID: userID, RecipeID: req.RecipeID}
    if err := h.db.Create(&fav).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add favorite"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Added to favorites"})
}

func (h *Handler) Remove(c *gin.Context) {
    userID := c.Param("id")
    recipeID := c.Param("recipeId")

    if err := h.db.Where("user_id = ? AND recipe_id = ?", userID, recipeID).Delete(&Favorite{}).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove favorite"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Removed from favorites"})
}

func (h *Handler) Get(c *gin.Context) {
    userID := c.Param("id")
    var favorites []Favorite
    if err := h.db.Where("user_id = ?", userID).Find(&favorites).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch favorites"})
        return
    }

    c.JSON(http.StatusOK, favorites)
}
*/
package favorite

import (
	"errors"
	"net/http"
	"strconv"
	"wongnok/internal/global"
	"wongnok/internal/helper"
	"wongnok/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type IHandler interface {
	Get(ctx *gin.Context)
	GetByUser(ctx *gin.Context)
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Update(ctx *gin.Context)
}

type Handler struct {
	Service IService
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		Service: NewService(db),
	}
}

func (handler Handler) Get(ctx *gin.Context) {
	var id string

	pathParam := ctx.Param("id")
	if pathParam != "" {
		id = pathParam
	}
	favorite, err := handler.Service.Get(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "favorite not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, favorite.ToResponse())
}


func (handler Handler) GetByUser(ctx *gin.Context) {

	claims, err := helper.DecodeClaims(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	var foodRecipeQuery model.FoodRecipeQuery
	if err := ctx.ShouldBindQuery(&foodRecipeQuery); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	foodRecipes, total, err := handler.Service.GetByUser(foodRecipeQuery, claims)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "favorite not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, foodRecipes.ToResponse(total))
}


func (handler Handler) Create(ctx *gin.Context) {
	var id int
	pathParam := ctx.Param("id")
	if pathParam != "" {
		if parsed, err := strconv.Atoi(pathParam); err == nil && parsed > 0 {
			id = parsed
		}
	}

	claims, err := helper.DecodeClaims(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	favorite, err := handler.Service.Create(id, claims)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.As(err, &validator.ValidationErrors{}) {
			statusCode = http.StatusBadRequest
		}

		ctx.JSON(statusCode, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, favorite.ToResponse())
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
