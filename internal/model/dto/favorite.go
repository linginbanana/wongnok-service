package dto

type FavoriteRequest struct {
	FoodRecipeID uint `validate:"required"`
}

type FavoriteResponse struct {
	ID           uint   `json:"id"`
	FoodRecipeID uint   `json:"foodRecipeID"`
	UserID       string `json:"userID"`
}

type FavoritesResponse BaseListResponse[[]FavoriteResponse]